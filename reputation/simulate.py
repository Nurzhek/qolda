"""Behaviour simulator for the reputation model.

Users are borrowers with a hidden profile (reliable, careless, malicious, colluder). Each step one
rental happens between a random borrower and a random lender; the outcome is drawn from the
borrower's profile, the lender rates the borrower, and a dispute is opened for damage or loss.
Colluders also run fake rentals with each other that always end well and get 5 stars.

ponytail: only the borrower side is modelled (deposits protect lenders), and behaviour does not react
to the deposit size. Add strategic users (build reputation, then defect) when the baseline
comparison is settled.
"""

from collections import defaultdict
from dataclasses import dataclass, field

import numpy as np

from model import (RECIPROCAL_WEIGHT, Campus, History, Vote, history_from_votes, neutral_score,
                   recommended_deposit, trust_score)

PROFILES = {
    # per-rental probabilities of a late return, damage, and never returning the item
    "reliable": {"late": 0.05, "damage": 0.01, "lost": 0.002},
    "careless": {"late": 0.30, "damage": 0.08, "lost": 0.010},
    "malicious": {"late": 0.30, "damage": 0.15, "lost": 0.150},
    "colluder": {"late": 0.30, "damage": 0.15, "lost": 0.150},
}
RISKY = {"malicious", "colluder"}

DEFAULTS = {
    "users": 300,
    "transactions": 4000,
    "horizon_days": 240,  # two semesters
    "mix": {"reliable": 0.60, "careless": 0.25, "malicious": 0.10, "colluder": 0.05},
    "inflation": 0.3,  # chance a lender gives 5 stars whatever happened (social pressure)
    "collusion_intensity": 0.5,  # fake colluder rentals per real rental
    "collusion_one_way": False,  # a smarter ring: fake rentals only go one way, so no pair is mutual
    "moderator_accuracy": 0.9,  # chance a dispute is decided correctly
    "false_dispute_rate": 0.02,  # lender opens a dispute although the rental went fine
    "item_value_median": 15_000,  # KZT
    "item_value_sigma": 0.8,
    "base_deposit_share": 0.5,  # owner's base deposit as a share of the item value
    "damage_cost_share": 0.4,  # repair cost as a share of the item value
    "new_user_max_rentals": 3,  # "few reviews" subset for the few-review-bias metric
}


@dataclass
class Pair:
    """Everything between one borrower and one lender."""

    rating_sum: float = 0.0
    rating_count: int = 0
    returns: int = 0
    on_time: int = 0
    closed: int = 0
    disputes_lost: int = 0


@dataclass
class User:
    id: int
    profile: str
    join_day: float
    pairs: dict = field(default_factory=lambda: defaultdict(Pair))  # lender id -> Pair
    real_rentals: int = 0

    @property
    def risk(self) -> float:
        """Expected loss per rental as a share of the item value: the ground truth."""
        p = PROFILES[self.profile]
        return p["lost"] + 0.4 * p["damage"] + 0.05 * p["late"]


def votes(u: User) -> list[Vote]:
    return [Vote(rater=lid, rating=p.rating_sum / p.rating_count, on_time=p.on_time / p.returns,
                 clean=1 - p.disputes_lost / p.closed) for lid, p in u.pairs.items()]


def age(u: User, day: float) -> float:
    return max(day - u.join_day, 0.0)


def history(u: User, day: float, pair_cap: bool, weights: dict[int, float] | None = None) -> History:
    """Aggregate a user's rentals: one vote per counterparty (pair_cap) or every rental counted."""
    if pair_cap:
        return history_from_votes(votes(u), age(u, day), weights)
    h = History(age_days=age(u, day))
    for p in u.pairs.values():
        h.rating_sum += p.rating_sum
        h.rating_count += p.rating_count
        h.returns += p.returns
        h.on_time += p.on_time
        h.closed += p.closed
        h.disputes_lost += p.disputes_lost
    return h


def simulate(params: dict, seed: int) -> dict:
    cfg = {**DEFAULTS, **params}
    rng = np.random.default_rng(seed)

    names = list(cfg["mix"])
    weights = np.array([cfg["mix"][n] for n in names], dtype=float)
    profiles = rng.choice(names, size=cfg["users"], p=weights / weights.sum())
    join_days = rng.uniform(0, cfg["horizon_days"] * 0.6, size=cfg["users"])
    users = [User(i, str(profiles[i]), float(join_days[i])) for i in range(cfg["users"])]
    colluders = [u for u in users if u.profile == "colluder"]

    # running campus-wide totals: rating sum/count, rentals due/on time, closed/lost disputes
    totals = {"rating_sum": 0.0, "rating_count": 0, "returns": 0, "on_time": 0, "closed": 0, "lost": 0}
    policies = ("fixed", "dynamic_avg", "dynamic_ours")
    loss = dict.fromkeys(policies, 0.0)
    deposits = {p: {"honest": [], "risky": []} for p in policies}

    def campus() -> Campus:
        c = Campus()
        if totals["rating_count"]:
            c.rating = totals["rating_sum"] / totals["rating_count"]
        if totals["returns"]:
            c.on_time = totals["on_time"] / totals["returns"]
        if totals["closed"]:
            c.clean = 1 - totals["lost"] / totals["closed"]
        return c

    def record(pair: Pair, stars: float, on_time: bool, lost_dispute: bool):
        pair.rating_sum += stars
        pair.rating_count += 1
        pair.returns += 1
        pair.on_time += int(on_time)
        pair.closed += 1
        pair.disputes_lost += int(lost_dispute)
        for k, v in (("rating_sum", stars), ("rating_count", 1), ("returns", 1), ("on_time", int(on_time)),
                     ("closed", 1), ("lost", int(lost_dispute))):
            totals[k] += v

    for t in range(cfg["transactions"]):
        day = cfg["horizon_days"] * (t + 1) / cfg["transactions"]
        active = [u for u in users if u.join_day <= day]
        if len(active) < 2:
            continue
        i, j = rng.choice(len(active), size=2, replace=False)
        borrower, lender = active[i], active[j]

        # deposits are computed before the outcome, from the history so far
        value = float(np.clip(rng.lognormal(np.log(cfg["item_value_median"]), cfg["item_value_sigma"]), 2_000, 300_000))
        base = value * cfg["base_deposit_share"]
        c = campus()
        # same as model.reciprocity_weights, without building the whole vote graph per rental
        mutual = {lid: RECIPROCAL_WEIGHT if borrower.id in users[lid].pairs else 1.0 for lid in borrower.pairs}
        ours = trust_score(history(borrower, day, pair_cap=True, weights=mutual), c).score
        raw = history(borrower, day, pair_cap=False)
        # baseline policy, anchored the same way: average rating vs the campus mean, 1 star = ±50%
        avg_factor = 1 - 0.5 * (raw.rating_sum / raw.rating_count - c.rating) if raw.rating_count else 1.0
        held = {
            "fixed": base,
            "dynamic_avg": base * min(max(avg_factor, 0.5), 1.5),
            "dynamic_ours": recommended_deposit(base, ours, neutral_score(c)),
        }

        # outcome
        p = PROFILES[borrower.profile]
        lost = rng.random() < p["lost"]
        damaged = not lost and rng.random() < p["damage"]
        late = not lost and rng.random() < p["late"]
        cost = value if lost else value * cfg["damage_cost_share"] if damaged else 0.0

        group = "risky" if borrower.profile in RISKY else "honest"
        for policy, d in held.items():
            loss[policy] += max(cost - d, 0.0)
            deposits[policy][group].append(d)

        # record history; every rental is due back, an item that never comes back is the worst kind of late
        borrower.real_rentals += 1
        at_fault = lost or damaged
        lost_dispute = False
        if at_fault or rng.random() < cfg["false_dispute_rate"]:
            correct = rng.random() < cfg["moderator_accuracy"]
            lost_dispute = at_fault == correct

        if rng.random() < cfg["inflation"]:
            stars = 5.0
        elif lost:
            stars = 1.0
        elif damaged:
            stars = float(rng.choice([1, 2]))
        elif late:
            stars = float(rng.choice([3, 4]))
        else:
            stars = 5.0 if rng.random() < 0.8 else 4.0
        record(borrower.pairs[lender.id], stars, on_time=not lost and not late, lost_dispute=lost_dispute)

        # colluders pad each other's history with fake rentals that always end well
        if len(colluders) >= 2 and rng.random() < cfg["collusion_intensity"]:
            a, b = rng.choice(len(colluders), size=2, replace=False)
            if cfg["collusion_one_way"]:
                a, b = min(a, b), max(a, b)  # only ever borrows from a higher-numbered member
            record(colluders[a].pairs[colluders[b].id], 5.0, on_time=True, lost_dispute=False)

    return {"users": users, "campus": campus(), "day": cfg["horizon_days"], "cfg": cfg,
            "loss": loss, "deposits": deposits}

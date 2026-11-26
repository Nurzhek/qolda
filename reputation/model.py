"""Trust score and dynamic deposit: the reputation model prototype.

Pure functions over a user's aggregated history, so the same formula can be ported to the Go
backend as-is. Every component is in [0, 1] and is returned for explainability (FR10).

Each behavioural component is a Bayesian estimate shrunk toward the campus average: a user with no
history is an "average student" on every component (the neutral prior from the design document),
and a few transactions move them only a little (few-review bias).

Collusion resistance has two layers:
1. Each distinct counterparty is one Vote (the mean of all rentals between the two users), so fake
   rentals inside a ring don't add up.
2. A mutual vote (the user has also rated the voter: they rented to each other) counts as
   RECIPROCAL_WEIGHT of a vote. Rings that rent back and forth lose most of their self-vouching, while
   honest students on a large campus rarely have mutual pairs. Weighting votes by the voter's own trust
   (EigenTrust) was tried and rejected: a careless borrower can still be an honest rater, so it threw
   away good evidence (see README).
"""

from dataclasses import dataclass

# weights of the score components; sum to 1
W_RATING, W_ON_TIME, W_CLEAN, W_AGE = 0.45, 0.25, 0.20, 0.10
PRIOR_WEIGHT = 5  # how many "virtual" campus-average transactions every user starts with
AGE_FULL_DAYS = 365  # account age stops adding trust after a year

# deposit = base × factor; factor is 1 for a new user and moves by 0.5 per DEPOSIT_SPAN points
DEPOSIT_SPAN = 20
DEPOSIT_MIN_FACTOR, DEPOSIT_MAX_FACTOR = 0.5, 1.5
DEPOSIT_ROUND = 100  # KZT

RECIPROCAL_WEIGHT = 0.2  # a mutual vote counts as a fifth of a vote


@dataclass
class History:
    """What the platform knows about a user as a borrower."""

    rating_sum: float = 0.0
    rating_count: float = 0.0
    returns: float = 0.0  # rentals that were due back
    on_time: float = 0.0  # of those, returned on time (an unreturned item is not on time)
    closed: float = 0.0  # closed orders (any outcome)
    disputes_lost: float = 0.0  # disputes where the user was found at fault
    age_days: float = 0.0


@dataclass
class Vote:
    """What one counterparty's rentals say about a user, averaged over all rentals between the two."""

    rater: int
    rating: float  # mean stars
    on_time: float  # share returned on time
    clean: float  # share of closed orders without a lost dispute


@dataclass
class Campus:
    """Campus-wide averages the components are shrunk toward."""

    rating: float = 4.5
    on_time: float = 0.9
    clean: float = 0.95  # share of closed orders without a lost dispute


@dataclass
class Score:
    score: float  # 0..100
    rating: float
    on_time: float
    clean: float
    age: float


def trust_score(h: History, campus: Campus) -> Score:
    k = PRIOR_WEIGHT
    bayes_rating = (campus.rating * k + h.rating_sum) / (k + h.rating_count)
    rating = (bayes_rating - 1) / 4
    on_time = (campus.on_time * k + h.on_time) / (k + h.returns)
    clean = (campus.clean * k + h.closed - h.disputes_lost) / (k + h.closed)
    age = min(h.age_days / AGE_FULL_DAYS, 1.0)

    score = 100 * (W_RATING * rating + W_ON_TIME * on_time + W_CLEAN * clean + W_AGE * age)
    return Score(min(max(score, 0.0), 100.0), rating, on_time, clean, age)


def neutral_score(campus: Campus) -> float:
    """Score of a user who just joined: the deposit anchor."""
    return trust_score(History(), campus).score


def deposit_factor(score: float, neutral: float) -> float:
    """1 for a new user; lower for a better-than-new score, higher for a worse one, within limits."""
    factor = 1 - 0.5 * (score - neutral) / DEPOSIT_SPAN
    return min(max(factor, DEPOSIT_MIN_FACTOR), DEPOSIT_MAX_FACTOR)


def recommended_deposit(base_deposit: float, score: float, neutral: float) -> int:
    raw = base_deposit * deposit_factor(score, neutral)
    return int(round(raw / DEPOSIT_ROUND) * DEPOSIT_ROUND)


def history_from_votes(votes: list[Vote], age_days: float, weights: dict[int, float] | None = None) -> History:
    """One vote per counterparty; weights (rater -> weight, default 1) scale individual votes."""
    h = History(age_days=age_days)
    for v in votes:
        w = weights.get(v.rater, 1.0) if weights else 1.0
        h.rating_sum += w * v.rating
        h.rating_count += w
        h.returns += w
        h.on_time += w * v.on_time
        h.closed += w
        h.disputes_lost += w * (1 - v.clean)
    return h


def reciprocity_weights(user: int, votes: dict[int, list[Vote]]) -> dict[int, float]:
    """rater -> weight for the votes about `user`: RECIPROCAL_WEIGHT if `user` has also rated that rater."""
    return {v.rater: RECIPROCAL_WEIGHT if any(x.rater == user for x in votes.get(v.rater, [])) else 1.0
            for v in votes.get(user, [])}


def trust_scores(votes: dict[int, list[Vote]], ages: dict[int, float], campus: Campus) -> dict[int, float]:
    """Scores of all users. votes[u] are the votes about u; ages[u] is u's account age in days."""
    return {u: trust_score(history_from_votes(votes.get(u, []), age, reciprocity_weights(u, votes)), campus).score
            for u, age in ages.items()}

"""Write golden cases for the Go port of the model (backend/internal/domain/trust).

The Go test reads this file and must reproduce every number, so the formula can't drift between the
prototype and the backend. CI regenerates it and fails if it changed without being committed.

    python golden.py
"""

import json
from dataclasses import asdict
from pathlib import Path

import numpy as np

from model import Campus, History, Vote, deposit_factor, neutral_score, recommended_deposit, trust_score, trust_scores

OUT = Path(__file__).parent.parent / "backend/internal/domain/trust/service/testdata/golden.json"


def main():
    rng = np.random.default_rng(2026)
    campuses = [Campus(), Campus(rating=4.7, on_time=0.85, clean=0.97), Campus(rating=3.9, on_time=0.7, clean=0.8)]

    histories = [History(), History(age_days=900),
                 History(rating_sum=5, rating_count=1, returns=1, on_time=1, closed=1),
                 History(rating_sum=192, rating_count=40, returns=40, on_time=40, closed=40, age_days=200)]
    for _ in range(20):
        n = float(rng.integers(0, 40))
        histories.append(History(
            rating_sum=round(float(rng.uniform(1, 5)) * n, 3), rating_count=n, returns=n,
            on_time=float(rng.integers(0, n + 1)), closed=n, disputes_lost=float(rng.integers(0, n // 4 + 1)),
            age_days=float(rng.integers(0, 800))))

    scores = []
    for c in campuses:
        for h in histories:
            scores.append({"history": asdict(h), "campus": asdict(c), "expected": asdict(trust_score(h, c)),
                           "neutral": neutral_score(c)})

    deposits = []
    neutral = neutral_score(Campus())
    for base in (0.0, 5_000.0, 7_350.0, 20_000.0):
        for delta in (-40, -20, -7.3, 0, 3, 11.1, 20, 55):
            score = neutral + delta
            deposits.append({"base": base, "score": score, "neutral": neutral,
                             "factor": deposit_factor(score, neutral),
                             "deposit": recommended_deposit(base, score, neutral)})

    # a small ring: 1 and 2 vouch for each other, 3 and 4 are real lenders
    v = lambda rater, rating, on_time, clean: Vote(rater, rating, on_time, clean)
    votes = {1: [v(2, 5, 1, 1), v(3, 2, 0.5, 0.5), v(4, 1, 0, 0)], 2: [v(1, 5, 1, 1), v(3, 5, 1, 1)],
             3: [v(4, 4.5, 1, 1)], 4: []}
    ages = {1: 40.0, 2: 40.0, 3: 300.0, 4: 10.0, 5: 0.0}
    graph = {"votes": {str(u): [asdict(x) for x in vs] for u, vs in votes.items()},
             "ages": {str(u): a for u, a in ages.items()}, "campus": asdict(Campus()),
             "expected": {str(u): s for u, s in trust_scores(votes, ages, Campus()).items()}}

    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text(json.dumps({"scores": scores, "deposits": deposits, "graph": graph}, indent=1) + "\n")
    print(f"wrote {OUT}")


if __name__ == "__main__":
    main()

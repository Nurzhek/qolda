import json
import subprocess
import sys

from model import (Campus, History, Vote, deposit_factor, neutral_score, recommended_deposit, reciprocity_weights,
                   trust_score, trust_scores)

CAMPUS = Campus(rating=4.6, on_time=0.9, clean=0.95)


def s(h: History) -> float:
    return trust_score(h, CAMPUS).score


def test_neutral_start_is_the_average_student():
    new = neutral_score(CAMPUS)
    average = History(rating_sum=4.6 * 40, rating_count=40, returns=40, on_time=36, closed=40, disputes_lost=2)
    assert abs(s(average) - new) < 0.5, "no history = average on every component"
    worst = History(rating_sum=50, rating_count=50, returns=50, closed=50, disputes_lost=50)
    best = History(rating_sum=500, rating_count=100, returns=100, on_time=100, closed=100, age_days=900)
    assert 0 <= s(worst) < new < s(best) <= 100


def test_few_review_bias():
    # the example from the design document
    one_five = History(rating_sum=5, rating_count=1, returns=1, on_time=1, closed=1)
    forty = History(rating_sum=40 * 4.8, rating_count=40, returns=40, on_time=40, closed=40)
    assert s(one_five) < s(forty)


def test_bad_behaviour_lowers_score():
    good = History(rating_sum=50, rating_count=10, returns=10, on_time=10, closed=10, age_days=100)
    late = History(rating_sum=50, rating_count=10, returns=10, on_time=3, closed=10, age_days=100)
    disputes = History(rating_sum=50, rating_count=10, returns=10, on_time=10, closed=10, disputes_lost=4, age_days=100)
    assert s(late) < s(good) and s(disputes) < s(good)


def test_deposit():
    n = neutral_score(CAMPUS)
    assert deposit_factor(n, n) == 1.0
    assert deposit_factor(n + 20, n) == 0.5 and deposit_factor(n - 20, n) == 1.5
    assert deposit_factor(n + 99, n) == 0.5 and deposit_factor(n - 99, n) == 1.5
    assert recommended_deposit(10_000, n, n) == 10_000
    deposits = [recommended_deposit(20_000, v, n) for v in range(0, 101, 5)]
    assert deposits == sorted(deposits, reverse=True), "a higher score never means a higher deposit"


def test_mutual_votes_count_less():
    perfect = lambda rater: Vote(rater=rater, rating=5, on_time=1, clean=1)
    bad = lambda rater: Vote(rater=rater, rating=1, on_time=0, clean=0)
    # users 1 and 2 vouch for each other; user 3 is a real lender who had a bad rental with user 1
    votes = {1: [perfect(2), bad(3)], 2: [perfect(1)]}
    assert reciprocity_weights(1, votes) == {2: 0.2, 3: 1.0}
    ages = {1: 30, 2: 30, 3: 30}
    one_way = trust_scores({1: [perfect(2), bad(3)]}, ages, CAMPUS)[1]
    assert trust_scores(votes, ages, CAMPUS)[1] < one_way, "the ring's mutual vote is discounted"


def test_experiment_is_reproducible(tmp_path):
    cfg = {"name": "tiny", "seeds": [1], "users": 40, "transactions": 300, "scenarios": {"base": {}}}
    path = tmp_path / "tiny.json"
    path.write_text(json.dumps(cfg))
    runs = []
    for i in range(2):
        out = tmp_path / f"out{i}"
        subprocess.run([sys.executable, "experiment.py", str(path), "--out", str(out), "--no-plots"], check=True)
        runs.append((out / "metrics.json").read_text())
    assert runs[0] == runs[1], "same config + seed must give identical results"

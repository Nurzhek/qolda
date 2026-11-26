# Reputation model prototype

Fall 2026 deliverable: a working prototype of the trust score, a behaviour simulator, and a first
comparison with the average-rating baselines (Requirements and Design, "Reputation engine").

```
model.py       trust score + dynamic deposit (pure functions)
simulate.py    behaviour simulator: reliable, careless, malicious and colluding users
experiment.py  runs a config over all seeds, computes metrics, writes results/<name>/
golden.py      golden cases for the Go port (backend/internal/domain/trust)
configs/       versioned experiment configs
results/       saved results: metrics.json, summary.md, plots
```

```bash
pip install -r requirements.txt
python -m pytest                              # model properties + reproducibility
python experiment.py configs/fall2026.json    # ~5 s on a laptop → results/fall2026-v1/
python golden.py                              # after any change to model.py
```

**Go port.** The formula lives in the backend as `backend/internal/domain/trust` (Go). Its test
reproduces every number in `golden.py`'s output, so the prototype and the backend can't drift apart.
Change `model.py`, run `python golden.py`, then update the Go code until `go test ./internal/domain/trust/...`
passes. CI fails if the golden file is stale.

Every run is reproducible: fixed seeds from the config, and `metrics.json` records the config and a
hash of the code (`code_version`) so a number in the thesis can be traced to the exact model.

## Model

For each user (as a borrower) four components in [0, 1]:

| component | weight | estimate |
|---|---|---|
| rating | 0.45 | Bayesian average of received stars, shrunk toward the campus mean with 5 virtual ratings, mapped 1–5 → 0–1 |
| on-time returns | 0.25 | share of rentals returned on time (an unreturned item counts as not on time), shrunk toward the campus rate |
| clean record | 0.20 | share of closed orders without a lost dispute, shrunk toward the campus rate |
| account age | 0.10 | days since sign-up / 365, capped at 1 |

`score = 100 × Σ weight × component`, bounded 0–100. With no history every behavioural component equals
the campus average, so a new user starts from a neutral prior (an "average student"), not 0 or 100, and a
handful of transactions moves them only a little.

**Collusion resistance**, two layers:

1. *Per-counterparty cap.* Each distinct counterparty contributes one vote (the mean of everything
   between the two users), so a ring gains at most one vote per member, however many fake rentals it
   runs. `ours_nocap` in the results is the model without the cap.
2. *Mutual-vote discount.* If the user has also rated the voter (they rented to each other), that vote
   counts as 0.2 of a vote. Rings that rent back and forth lose most of their self-vouching; honest
   students on a campus of hundreds rarely have mutual pairs. `ours_no_mutual` is the model without it.

Tried and rejected: weighting each vote by the voter's own trust score (EigenTrust-style). It made
every metric worse, with or without rescaling the weights, because a careless or even malicious
*borrower* can still be an honest *rater*: down-weighting them threw away good evidence.

**Dynamic deposit.** `deposit = base × factor`, `factor = 1 − 0.5 × (score − neutral) / 20`, clamped to
[0.5, 1.5] and rounded to 100 KZT. A new user pays exactly the owner's base deposit; each point above a
new user's score lowers the deposit by 2.5%, each point below raises it by 2.5%. This is the
explanation a user sees: "your score is 88, a new member's is 83, so your deposit is 12.5% lower".

## Simulator

300 users over two semesters (240 days), 4,000 rentals. Hidden profiles and per-rental probabilities:

| profile | share | late | damage | never returned |
|---|---|---|---|---|
| reliable | 60% | 5% | 1% | 0.2% |
| careless | 25% | 30% | 8% | 1% |
| malicious | 10% | 30% | 15% | 15% |
| colluder | 5% | 30% | 15% | 15% |

The lender rates the borrower after each rental (5 stars when fine, 3–4 when late, 1–2 for damage, 1 when
lost); with probability `inflation` the lender gives 5 stars regardless (social pressure, as on Airbnb).
Damage or loss opens a dispute; moderators decide correctly 90% of the time, and 2% of fine rentals get
a false dispute. Colluders additionally run fake rentals with each other (0.5 per real rental) that
always end on time with 5 stars; with `collusion_one_way` a ring only rents in one direction, so none of
its pairs are mutual (the best attack on the mutual-vote discount).

Ground truth per user: expected loss share per rental (`lost + 0.4·damage + 0.05·late`). Risky users =
malicious + colluders.

**Baselines:** `avg`, the plain average rating (no ratings → campus mean); `avg_min3`, the average only
after 3+ reviews.

**Learned baselines** (`logreg`, `boosting`): logistic regression and gradient boosting
(scikit-learn) over what the platform can observe — the trust-score components, number of distinct
counterparties, share of mutual votes, raw average rating, number of rentals — trained on 10 separate
seeds per scenario (`train_seeds`, never evaluated) with the simulator's ground-truth labels. They are
not candidates for production (no labels exist on a real campus yet, and FR10 requires an explainable
score); they measure how much the transparent formula leaves on the table.

**Metrics:** Spearman correlation between score and true reliability (all users, and users with ≤3
rentals); average precision for finding risky users by lowest score; precision/recall when moderators
review the lowest-scored 15%; lenders' uncovered loss (`max(cost − deposit, 0)`) and mean deposit held
under fixed vs dynamic deposits.

## Results (`fall2026-v1`, 10 seeds)

Base scenario (30% rating inflation, 5% colluders):

| method | Spearman | Spearman, ≤3 rentals | avg precision | recall @15% |
|---|---|---|---|---|
| avg | 0.50 ± 0.02 | 0.28 ± 0.15 | 0.45 ± 0.06 | 0.48 ± 0.06 |
| avg_min3 | 0.50 ± 0.03 | 0.17 ± 0.28 | 0.46 ± 0.07 | 0.48 ± 0.06 |
| ours_nocap | 0.58 ± 0.03 | 0.28 ± 0.22 | 0.52 ± 0.07 | 0.51 ± 0.06 |
| ours_no_mutual | 0.76 ± 0.02 | 0.36 ± 0.19 | 0.60 ± 0.07 | 0.54 ± 0.06 |
| **ours** | **0.79 ± 0.01** | **0.42 ± 0.14** | **0.71 ± 0.06** | **0.65 ± 0.06** |

Spearman under stress:

| scenario | avg | ours_no_mutual | ours |
|---|---|---|---|
| no colluders | 0.66 | 0.78 | 0.77 |
| 15% colluders, renting back and forth | 0.41 | 0.60 | **0.70** |
| 15% colluders, one-way only (best attack on the discount) | 0.44 | 0.67 | 0.66 |
| 90% rating inflation | 0.24 | 0.77 | **0.79** |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.05 ± 0.18 | 10,228 | 10,514 |
| dynamic, average rating | 0.92 ± 0.15 | 10,345 | 12,172 |
| dynamic, ours | 0.92 ± 0.15 | 9,808 | 12,104 |

Full tables: `results/fall2026-v1/summary.md`. Plots: `sweep_inflation.png`, `sweep_colluders.png`,
`sweep_one_way_colluders.png`, `deposits_base.png`.

**What this shows**

- The trust score ranks users by true reliability far better than the average rating (0.79 vs 0.50)
  and finds risky users better (average precision 0.71 vs 0.45; moderators reviewing the lowest 15%
  catch 65% of risky users vs 48%).
- It holds under rating inflation: at 90% inflation the average rating falls to 0.24 while the trust
  score stays at 0.79, because on-time returns and dispute outcomes can't be inflated by politeness.
- Both collusion layers pull their weight. The cap alone takes 5% collusion from 0.58 to 0.76, and the
  mutual-vote discount takes 15% collusion from 0.60 to 0.70. Neither costs anything without collusion.
- A ring that avoids mutual pairs defeats the discount, but it has to give up half of its fake votes to
  do so, so its best attack (0.66) is still weaker than the old worst case (0.60).
- Dynamic deposits cut lenders' uncovered loss by about 12% vs a fixed deposit. The average-rating policy
  does the same, but only the trust score does it while also lowering honest users' deposits (−4%).

**Learned baselines: what a model trained on ground truth adds**

| scenario | ours: Spearman / AP | logreg: Spearman / AP | boosting: Spearman / AP |
|---|---|---|---|
| base | 0.79 / 0.71 | 0.80 / **0.88** | 0.76 / 0.85 |
| no colluders | 0.77 / 0.70 | 0.78 / 0.76 | 0.71 / 0.74 |
| 15% colluders | 0.70 / 0.60 | 0.82 / **0.95** | 0.83 / 0.95 |
| 90% rating inflation | 0.79 / 0.69 | 0.81 / 0.90 | 0.73 / 0.88 |

- Boosting is never better than logistic regression and ranks users worse: a non-linear black box buys
  nothing here, so the production model should stay linear and explainable.
- Logistic regression ranks about as well as our formula but finds risky users much better (AP 0.88 vs
  0.71). Its coefficients (base scenario, standardized, + = riskier) say why: many rentals (+2.14) from
  few distinct counterparties (−1.02), i.e. renting again and again from the same people, the signature
  of a ring, is the strongest signal; then disputes lost (clean −1.47) and on-time returns (−0.68);
  the rating matters least (−0.26). Our formula has no "repeat counterparty" signal and gives the rating
  the largest weight (0.45) and disputes only 0.20.
- Caveat: the regression learned the simulator's own rules (a ring runs 0.5 fake rentals per real
  one). A real ring could keep its volume low, so the repeat-counterparty signal must be confirmed on
  pilot data before it changes the model.
- Next step: add a transparent "rentals per counterparty" component and move weight from the rating to
  disputes, then re-run; the formula stays explainable and portable to Go.

**What it doesn't show yet (honestly)**

- One-way rings remain the open problem (0.66 at 15%). Next candidates: discount votes from users whose
  own history comes mostly from one tight cluster, or require payment through Kaspi for a rental to count
  (fake rentals then cost real money).
- The few-review metric (≤3 rentals) is noisy at this population size (std ≈ 0.15–0.3).
- The deposit effect is modest because simulated users don't react to the deposit size (no deterrence)
  and the base deposit is only half the item value.
- Only the borrower side is modelled, and there are no strategic users (build a reputation, then defect).
- Weights (0.45/0.25/0.20/0.10), the prior weight (5), the deposit span (20) and the mutual-vote weight
  (0.2) are first guesses, not tuned. Tune them on held-out seeds, then confirm on pilot data.

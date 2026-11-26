# fall2026-v1

First comparison of the trust score with the average-rating baselines: base campus, rating inflation sweep, collusion sweep. Defaults: simulate.DEFAULTS.

Mean ± std over 10 seeds. Flag budget: lowest 15% of users.

## base

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.497 ± 0.024 | 0.278 ± 0.152 | 0.454 ± 0.058 | 0.478 ± 0.102 | 0.480 ± 0.055 |
| avg_min3 | 0.498 ± 0.028 | 0.168 ± 0.281 | 0.458 ± 0.066 | 0.476 ± 0.109 | 0.477 ± 0.059 |
| ours_nocap | 0.577 ± 0.032 | 0.282 ± 0.222 | 0.519 ± 0.073 | 0.507 ± 0.101 | 0.510 ± 0.063 |
| ours_no_mutual | 0.758 ± 0.019 | 0.359 ± 0.187 | 0.595 ± 0.072 | 0.536 ± 0.086 | 0.542 ± 0.056 |
| ours | 0.791 ± 0.011 | 0.418 ± 0.144 | 0.707 ± 0.056 | 0.636 ± 0.078 | 0.646 ± 0.059 |
| logreg | 0.800 ± 0.016 | 0.674 ± 0.133 | 0.879 ± 0.041 | 0.771 ± 0.079 | 0.784 ± 0.060 |
| boosting | 0.758 ± 0.019 | 0.647 ± 0.241 | 0.853 ± 0.037 | 0.731 ± 0.066 | 0.746 ± 0.073 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.05 ± 0.18 | 10,228 | 10,514 |
| dynamic_avg | 0.92 ± 0.15 | 10,345 | 12,172 |
| dynamic_ours | 0.92 ± 0.15 | 9,808 | 12,104 |

Logistic regression coefficients (standardized features, + = riskier): rating -0.26, on_time -0.68, clean -1.47, age -0.10, votes -1.02, mutual_share +0.30, raw_avg_rating +0.03, rentals +2.14

## inflation=0.0

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.525 ± 0.031 | 0.453 ± 0.256 | 0.471 ± 0.082 | 0.504 ± 0.100 | 0.507 ± 0.058 |
| avg_min3 | 0.523 ± 0.033 | 0.315 ± 0.349 | 0.477 ± 0.080 | 0.504 ± 0.100 | 0.507 ± 0.058 |
| ours_nocap | 0.555 ± 0.040 | 0.346 ± 0.427 | 0.507 ± 0.092 | 0.484 ± 0.101 | 0.487 ± 0.066 |
| ours_no_mutual | 0.754 ± 0.028 | 0.410 ± 0.350 | 0.593 ± 0.085 | 0.511 ± 0.080 | 0.517 ± 0.052 |
| ours | 0.789 ± 0.018 | 0.498 ± 0.287 | 0.700 ± 0.069 | 0.620 ± 0.076 | 0.630 ± 0.064 |
| logreg | 0.777 ± 0.014 | 0.582 ± 0.209 | 0.868 ± 0.033 | 0.769 ± 0.082 | 0.781 ± 0.046 |
| boosting | 0.724 ± 0.024 | 0.256 ± 0.389 | 0.856 ± 0.047 | 0.751 ± 0.091 | 0.763 ± 0.058 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 0.99 ± 0.18 | 10,240 | 10,249 |
| dynamic_avg | 0.81 ± 0.14 | 10,408 | 12,132 |
| dynamic_ours | 0.84 ± 0.14 | 9,844 | 12,028 |

Logistic regression coefficients (standardized features, + = riskier): rating -1.12, on_time -0.27, clean -1.04, age -0.06, votes -1.19, mutual_share +0.64, raw_avg_rating -0.12, rentals +2.08

## inflation=0.3

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.497 ± 0.024 | 0.278 ± 0.152 | 0.454 ± 0.058 | 0.478 ± 0.102 | 0.480 ± 0.055 |
| avg_min3 | 0.498 ± 0.028 | 0.168 ± 0.281 | 0.458 ± 0.066 | 0.476 ± 0.109 | 0.477 ± 0.059 |
| ours_nocap | 0.577 ± 0.032 | 0.282 ± 0.222 | 0.519 ± 0.073 | 0.507 ± 0.101 | 0.510 ± 0.063 |
| ours_no_mutual | 0.758 ± 0.019 | 0.359 ± 0.187 | 0.595 ± 0.072 | 0.536 ± 0.086 | 0.542 ± 0.056 |
| ours | 0.791 ± 0.011 | 0.418 ± 0.144 | 0.707 ± 0.056 | 0.636 ± 0.078 | 0.646 ± 0.059 |
| logreg | 0.800 ± 0.016 | 0.674 ± 0.133 | 0.879 ± 0.041 | 0.771 ± 0.079 | 0.784 ± 0.060 |
| boosting | 0.758 ± 0.019 | 0.647 ± 0.241 | 0.853 ± 0.037 | 0.731 ± 0.066 | 0.746 ± 0.073 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.05 ± 0.18 | 10,228 | 10,514 |
| dynamic_avg | 0.92 ± 0.15 | 10,345 | 12,172 |
| dynamic_ours | 0.92 ± 0.15 | 9,808 | 12,104 |

Logistic regression coefficients (standardized features, + = riskier): rating -0.26, on_time -0.68, clean -1.47, age -0.10, votes -1.02, mutual_share +0.30, raw_avg_rating +0.03, rentals +2.14

## inflation=0.6

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.399 ± 0.047 | 0.124 ± 0.287 | 0.377 ± 0.052 | 0.411 ± 0.060 | 0.418 ± 0.053 |
| avg_min3 | 0.401 ± 0.046 | -0.013 ± 0.242 | 0.379 ± 0.052 | 0.409 ± 0.060 | 0.415 ± 0.052 |
| ours_nocap | 0.572 ± 0.038 | 0.360 ± 0.324 | 0.498 ± 0.061 | 0.489 ± 0.064 | 0.496 ± 0.044 |
| ours_no_mutual | 0.741 ± 0.031 | 0.416 ± 0.291 | 0.571 ± 0.062 | 0.507 ± 0.059 | 0.516 ± 0.055 |
| ours | 0.773 ± 0.022 | 0.552 ± 0.224 | 0.658 ± 0.037 | 0.600 ± 0.034 | 0.613 ± 0.055 |
| logreg | 0.805 ± 0.017 | 0.639 ± 0.144 | 0.875 ± 0.029 | 0.769 ± 0.066 | 0.784 ± 0.064 |
| boosting | 0.769 ± 0.029 | 0.555 ± 0.274 | 0.877 ± 0.037 | 0.782 ± 0.076 | 0.795 ± 0.049 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.07 ± 0.13 | 10,274 | 10,530 |
| dynamic_avg | 0.96 ± 0.12 | 10,317 | 11,634 |
| dynamic_ours | 0.96 ± 0.11 | 9,816 | 11,726 |

Logistic regression coefficients (standardized features, + = riskier): rating +0.26, on_time -0.97, clean -1.44, age -0.26, votes -0.93, mutual_share -0.09, raw_avg_rating -0.42, rentals +2.85

## inflation=0.9

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.237 ± 0.051 | 0.104 ± 0.252 | 0.275 ± 0.052 | 0.267 ± 0.055 | 0.270 ± 0.046 |
| avg_min3 | 0.233 ± 0.052 | 0.006 ± 0.284 | 0.273 ± 0.050 | 0.264 ± 0.057 | 0.267 ± 0.046 |
| ours_nocap | 0.611 ± 0.043 | 0.190 ± 0.506 | 0.519 ± 0.053 | 0.500 ± 0.082 | 0.507 ± 0.064 |
| ours_no_mutual | 0.767 ± 0.029 | 0.222 ± 0.529 | 0.603 ± 0.053 | 0.533 ± 0.072 | 0.542 ± 0.064 |
| ours | 0.793 ± 0.021 | 0.278 ± 0.546 | 0.694 ± 0.032 | 0.638 ± 0.050 | 0.651 ± 0.063 |
| logreg | 0.813 ± 0.022 | 0.430 ± 0.436 | 0.900 ± 0.028 | 0.793 ± 0.060 | 0.808 ± 0.057 |
| boosting | 0.729 ± 0.034 | 0.340 ± 0.466 | 0.877 ± 0.033 | 0.773 ± 0.062 | 0.788 ± 0.065 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.08 ± 0.21 | 10,365 | 10,341 |
| dynamic_avg | 1.05 ± 0.21 | 10,361 | 10,682 |
| dynamic_ours | 1.01 ± 0.20 | 9,879 | 11,261 |

Logistic regression coefficients (standardized features, + = riskier): rating -0.11, on_time -0.90, clean -1.39, age -0.10, votes -1.04, mutual_share +0.22, raw_avg_rating +0.01, rentals +2.36

## colluders=0.0

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.659 ± 0.031 | 0.403 ± 0.348 | 0.595 ± 0.069 | 0.471 ± 0.069 | 0.729 ± 0.075 |
| avg_min3 | 0.664 ± 0.030 | 0.468 ± 0.187 | 0.602 ± 0.075 | 0.464 ± 0.072 | 0.718 ± 0.078 |
| ours_nocap | 0.779 ± 0.025 | 0.505 ± 0.328 | 0.712 ± 0.057 | 0.502 ± 0.054 | 0.782 ± 0.081 |
| ours_no_mutual | 0.776 ± 0.026 | 0.502 ± 0.327 | 0.709 ± 0.060 | 0.496 ± 0.050 | 0.772 ± 0.084 |
| ours | 0.773 ± 0.026 | 0.502 ± 0.327 | 0.696 ± 0.061 | 0.491 ± 0.050 | 0.764 ± 0.077 |
| logreg | 0.777 ± 0.024 | 0.503 ± 0.329 | 0.762 ± 0.067 | 0.529 ± 0.074 | 0.819 ± 0.079 |
| boosting | 0.713 ± 0.027 | 0.368 ± 0.370 | 0.740 ± 0.057 | 0.518 ± 0.076 | 0.801 ± 0.079 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 0.74 ± 0.15 | 10,318 | 10,187 |
| dynamic_avg | 0.61 ± 0.12 | 9,875 | 12,584 |
| dynamic_ours | 0.66 ± 0.12 | 9,554 | 11,797 |

Logistic regression coefficients (standardized features, + = riskier): rating +0.01, on_time -0.89, clean -1.42, age -0.00, votes -0.08, mutual_share -0.08, raw_avg_rating -0.31, rentals -0.12

## colluders=0.05

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.497 ± 0.024 | 0.278 ± 0.152 | 0.454 ± 0.058 | 0.478 ± 0.102 | 0.480 ± 0.055 |
| avg_min3 | 0.498 ± 0.028 | 0.168 ± 0.281 | 0.458 ± 0.066 | 0.476 ± 0.109 | 0.477 ± 0.059 |
| ours_nocap | 0.577 ± 0.032 | 0.282 ± 0.222 | 0.519 ± 0.073 | 0.507 ± 0.101 | 0.510 ± 0.063 |
| ours_no_mutual | 0.758 ± 0.019 | 0.359 ± 0.187 | 0.595 ± 0.072 | 0.536 ± 0.086 | 0.542 ± 0.056 |
| ours | 0.791 ± 0.011 | 0.418 ± 0.144 | 0.707 ± 0.056 | 0.636 ± 0.078 | 0.646 ± 0.059 |
| logreg | 0.800 ± 0.016 | 0.674 ± 0.133 | 0.879 ± 0.041 | 0.771 ± 0.079 | 0.784 ± 0.060 |
| boosting | 0.758 ± 0.019 | 0.647 ± 0.241 | 0.853 ± 0.037 | 0.731 ± 0.066 | 0.746 ± 0.073 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.05 ± 0.18 | 10,228 | 10,514 |
| dynamic_avg | 0.92 ± 0.15 | 10,345 | 12,172 |
| dynamic_ours | 0.92 ± 0.15 | 9,808 | 12,104 |

Logistic regression coefficients (standardized features, + = riskier): rating -0.26, on_time -0.68, clean -1.47, age -0.10, votes -1.02, mutual_share +0.30, raw_avg_rating +0.03, rentals +2.14

## colluders=0.1

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.433 ± 0.036 | 0.303 ± 0.260 | 0.455 ± 0.065 | 0.527 ± 0.084 | 0.389 ± 0.064 |
| avg_min3 | 0.435 ± 0.034 | 0.068 ± 0.360 | 0.455 ± 0.061 | 0.522 ± 0.080 | 0.386 ± 0.062 |
| ours_nocap | 0.516 ± 0.026 | 0.084 ± 0.313 | 0.505 ± 0.067 | 0.569 ± 0.089 | 0.418 ± 0.061 |
| ours_no_mutual | 0.676 ± 0.043 | 0.099 ± 0.312 | 0.562 ± 0.071 | 0.573 ± 0.090 | 0.422 ± 0.064 |
| ours | 0.775 ± 0.038 | 0.240 ± 0.241 | 0.683 ± 0.082 | 0.682 ± 0.082 | 0.503 ± 0.067 |
| logreg | 0.821 ± 0.016 | 0.460 ± 0.358 | 0.932 ± 0.013 | 0.951 ± 0.041 | 0.702 ± 0.064 |
| boosting | 0.785 ± 0.014 | 0.486 ± 0.296 | 0.924 ± 0.017 | 0.949 ± 0.034 | 0.700 ± 0.054 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.47 ± 0.13 | 10,302 | 10,446 |
| dynamic_avg | 1.32 ± 0.12 | 10,329 | 11,660 |
| dynamic_ours | 1.34 ± 0.13 | 9,801 | 11,505 |

Logistic regression coefficients (standardized features, + = riskier): rating -0.31, on_time -0.74, clean -1.36, age -0.77, votes -0.75, mutual_share +1.03, raw_avg_rating -0.02, rentals +2.14

## colluders=0.15

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.407 ± 0.027 | 0.040 ± 0.321 | 0.445 ± 0.045 | 0.527 ± 0.074 | 0.308 ± 0.025 |
| avg_min3 | 0.407 ± 0.030 | -0.120 ± 0.338 | 0.440 ± 0.047 | 0.527 ± 0.081 | 0.307 ± 0.024 |
| ours_nocap | 0.495 ± 0.017 | -0.059 ± 0.319 | 0.493 ± 0.032 | 0.556 ± 0.064 | 0.326 ± 0.025 |
| ours_no_mutual | 0.596 ± 0.030 | -0.055 ± 0.311 | 0.525 ± 0.031 | 0.562 ± 0.065 | 0.329 ± 0.019 |
| ours | 0.702 ± 0.040 | 0.105 ± 0.240 | 0.597 ± 0.033 | 0.598 ± 0.058 | 0.351 ± 0.027 |
| logreg | 0.817 ± 0.015 | 0.606 ± 0.167 | 0.948 ± 0.017 | 0.993 ± 0.010 | 0.588 ± 0.065 |
| boosting | 0.832 ± 0.013 | 0.506 ± 0.265 | 0.953 ± 0.017 | 1.000 ± 0.000 | 0.592 ± 0.064 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.82 ± 0.28 | 10,272 | 10,476 |
| dynamic_avg | 1.67 ± 0.24 | 10,178 | 11,416 |
| dynamic_ours | 1.72 ± 0.26 | 9,697 | 11,064 |

Logistic regression coefficients (standardized features, + = riskier): rating +0.21, on_time -0.80, clean -1.58, age -1.24, votes +0.82, mutual_share +0.94, raw_avg_rating -0.31, rentals +1.27

## one_way_colluders=0.05

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.529 ± 0.020 | 0.278 ± 0.152 | 0.478 ± 0.056 | 0.502 ± 0.100 | 0.505 ± 0.053 |
| avg_min3 | 0.530 ± 0.024 | 0.168 ± 0.281 | 0.483 ± 0.063 | 0.498 ± 0.102 | 0.500 ± 0.055 |
| ours_nocap | 0.618 ± 0.025 | 0.317 ± 0.184 | 0.546 ± 0.064 | 0.527 ± 0.096 | 0.531 ± 0.061 |
| ours_no_mutual | 0.780 ± 0.017 | 0.404 ± 0.155 | 0.673 ± 0.061 | 0.618 ± 0.084 | 0.627 ± 0.060 |
| ours | 0.775 ± 0.015 | 0.399 ± 0.159 | 0.664 ± 0.062 | 0.604 ± 0.078 | 0.613 ± 0.050 |
| logreg | 0.798 ± 0.017 | 0.686 ± 0.137 | 0.872 ± 0.042 | 0.760 ± 0.079 | 0.773 ± 0.061 |
| boosting | 0.743 ± 0.018 | 0.631 ± 0.192 | 0.846 ± 0.036 | 0.724 ± 0.072 | 0.737 ± 0.061 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.05 ± 0.18 | 10,228 | 10,514 |
| dynamic_avg | 0.91 ± 0.16 | 10,345 | 12,326 |
| dynamic_ours | 0.93 ± 0.16 | 9,808 | 11,963 |

Logistic regression coefficients (standardized features, + = riskier): rating -0.34, on_time -0.65, clean -1.44, age -0.04, votes -1.14, mutual_share +0.09, raw_avg_rating +0.12, rentals +3.07

## one_way_colluders=0.1

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.468 ± 0.037 | 0.303 ± 0.260 | 0.487 ± 0.059 | 0.558 ± 0.075 | 0.411 ± 0.056 |
| avg_min3 | 0.471 ± 0.035 | 0.068 ± 0.360 | 0.487 ± 0.055 | 0.556 ± 0.071 | 0.410 ± 0.055 |
| ours_nocap | 0.557 ± 0.026 | 0.084 ± 0.313 | 0.545 ± 0.057 | 0.604 ± 0.073 | 0.444 ± 0.046 |
| ours_no_mutual | 0.740 ± 0.039 | 0.143 ± 0.254 | 0.654 ± 0.064 | 0.660 ± 0.067 | 0.486 ± 0.051 |
| ours | 0.732 ± 0.039 | 0.146 ± 0.253 | 0.640 ± 0.068 | 0.651 ± 0.072 | 0.479 ± 0.050 |
| logreg | 0.826 ± 0.016 | 0.473 ± 0.362 | 0.930 ± 0.015 | 0.944 ± 0.046 | 0.697 ± 0.058 |
| boosting | 0.779 ± 0.016 | 0.490 ± 0.186 | 0.915 ± 0.015 | 0.938 ± 0.042 | 0.691 ± 0.053 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.47 ± 0.13 | 10,302 | 10,446 |
| dynamic_avg | 1.30 ± 0.11 | 10,329 | 11,826 |
| dynamic_ours | 1.35 ± 0.12 | 9,801 | 11,432 |

Logistic regression coefficients (standardized features, + = riskier): rating -0.53, on_time -0.67, clean -1.35, age -0.45, votes -1.88, mutual_share +0.09, raw_avg_rating +0.25, rentals +4.80

## one_way_colluders=0.15

| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |
|---|---|---|---|---|---|
| avg | 0.444 ± 0.022 | 0.067 ± 0.332 | 0.483 ± 0.047 | 0.580 ± 0.070 | 0.341 ± 0.036 |
| avg_min3 | 0.445 ± 0.024 | -0.081 ± 0.355 | 0.478 ± 0.048 | 0.576 ± 0.078 | 0.337 ± 0.035 |
| ours_nocap | 0.529 ± 0.018 | -0.006 ± 0.259 | 0.538 ± 0.033 | 0.607 ± 0.076 | 0.356 ± 0.033 |
| ours_no_mutual | 0.669 ± 0.029 | 0.011 ± 0.246 | 0.607 ± 0.029 | 0.647 ± 0.075 | 0.380 ± 0.036 |
| ours | 0.655 ± 0.031 | 0.018 ± 0.257 | 0.592 ± 0.032 | 0.629 ± 0.072 | 0.370 ± 0.043 |
| logreg | 0.833 ± 0.012 | 0.611 ± 0.139 | 0.935 ± 0.014 | 0.982 ± 0.013 | 0.582 ± 0.065 |
| boosting | 0.840 ± 0.013 | 0.533 ± 0.208 | 0.941 ± 0.018 | 0.987 ± 0.015 | 0.584 ± 0.063 |

| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |
|---|---|---|---|
| fixed | 1.82 ± 0.28 | 10,272 | 10,476 |
| dynamic_avg | 1.63 ± 0.23 | 10,178 | 11,632 |
| dynamic_ours | 1.72 ± 0.25 | 9,697 | 11,109 |

Logistic regression coefficients (standardized features, + = riskier): rating -0.12, on_time -0.79, clean -1.57, age -0.86, votes -1.72, mutual_share +0.00, raw_avg_rating +0.14, rentals +4.81

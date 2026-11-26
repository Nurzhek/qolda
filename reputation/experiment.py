"""Run a versioned experiment config: every scenario × every seed, compare the trust score with the
baselines, save metrics.json, summary.md and plots under results/<name>/.

    python experiment.py configs/fall2026.json
"""

import argparse
import hashlib
import json
from concurrent.futures import ProcessPoolExecutor
from pathlib import Path

import numpy as np
import pandas as pd

from model import trust_score, trust_scores
from simulate import RISKY, age, history, simulate, votes

HERE = Path(__file__).parent
FLAG_BUDGET = 0.15  # moderators can review the lowest-scored 15% of users


# --- methods: a score per user, lower means riskier


def method_scores(run: dict) -> dict[str, np.ndarray]:
    users, day, campus = run["users"], run["day"], run["campus"]
    mean = campus.rating
    weighted = trust_scores({u.id: votes(u) for u in users}, {u.id: age(u, day) for u in users}, campus)
    raw = [history(u, day, pair_cap=False) for u in users]
    avg = [h.rating_sum / h.rating_count if h.rating_count else mean for h in raw]
    return {
        # plain average rating; no ratings yet -> campus mean
        "avg": np.array(avg),
        # average only once a user has 3+ reviews
        "avg_min3": np.array([a if h.rating_count >= 3 else mean for a, h in zip(avg, raw)]),
        # our model without the per-counterparty cap (ablation)
        "ours_nocap": np.array([trust_score(h, campus).score for h in raw]),
        # one vote per counterparty, no discount for mutual votes (ablation)
        "ours_no_mutual": np.array([trust_score(history(u, day, pair_cap=True), campus).score for u in users]),
        # the model: one vote per counterparty, mutual votes discounted
        "ours": np.array([weighted[u.id] for u in users]),
    }


# --- metrics


def spearman(a: np.ndarray, b: np.ndarray) -> float:
    if len(a) < 3 or np.ptp(a) == 0 or np.ptp(b) == 0:
        return 0.0
    # Spearman = Pearson on average ranks; done by hand so pandas doesn't need scipy
    return float(pd.Series(a).rank().corr(pd.Series(b).rank()))


def average_precision(scores: np.ndarray, risky: np.ndarray) -> float:
    """Area under the precision-recall curve for flagging low scores; ties share a threshold."""
    if not risky.any():
        return float("nan")
    ap, prev_recall = 0.0, 0.0
    for thr in np.unique(scores):
        flagged = scores <= thr
        tp = (flagged & risky).sum()
        recall = tp / risky.sum()
        ap += (recall - prev_recall) * tp / flagged.sum()
        prev_recall = recall
    return float(ap)


def at_budget(scores: np.ndarray, risky: np.ndarray) -> tuple[float, float]:
    k = max(int(round(len(scores) * FLAG_BUDGET)), 1)
    flagged = np.zeros(len(scores), dtype=bool)
    flagged[np.argsort(scores, kind="stable")[:k]] = True  # ties broken by user id
    tp = (flagged & risky).sum()
    return float(tp / k), float(tp / risky.sum()) if risky.any() else float("nan")


def evaluate(run: dict) -> dict:
    users = run["users"]
    reliability = -np.array([u.risk for u in users])
    risky = np.array([u.profile in RISKY for u in users])
    new = np.array([1 <= u.real_rentals <= run["cfg"]["new_user_max_rentals"] for u in users])

    result = {}
    for name, s in method_scores(run).items():
        precision, recall = at_budget(s, risky)
        result[name] = {
            "spearman": spearman(s, reliability),
            "spearman_new_users": spearman(s[new], reliability[new]),
            "average_precision": average_precision(s, risky),
            "precision_at_budget": precision,
            "recall_at_budget": recall,
        }
    for policy, groups in run["deposits"].items():
        result[f"deposit:{policy}"] = {
            "uncovered_loss_mln_kzt": run["loss"][policy] / 1e6,
            "mean_deposit_honest": float(np.mean(groups["honest"])),
            "mean_deposit_risky": float(np.mean(groups["risky"])) if groups["risky"] else float("nan"),
        }
    return result


def run_one(job: tuple[dict, int]) -> dict:
    params, seed = job
    return evaluate(simulate(params, seed))


# --- experiment


def aggregate(per_seed: list[dict]) -> dict:
    out = {}
    for row in per_seed[0]:
        out[row] = {}
        for metric in per_seed[0][row]:
            vals = np.array([r[row][metric] for r in per_seed], dtype=float)
            out[row][metric] = {"mean": round(float(np.nanmean(vals)), 6), "std": round(float(np.nanstd(vals)), 6)}
    return out


def code_version() -> str:
    h = hashlib.sha256()
    for f in ("model.py", "simulate.py", "experiment.py"):
        h.update((HERE / f).read_bytes())
    return h.hexdigest()[:12]


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("config")
    ap.add_argument("--out", help="default: results/<config name>")
    ap.add_argument("--no-plots", action="store_true")
    args = ap.parse_args()

    cfg = json.loads(Path(args.config).read_text())
    shared = {k: v for k, v in cfg.items() if k not in ("name", "description", "seeds", "scenarios")}
    out = Path(args.out) if args.out else HERE / "results" / cfg["name"]
    out.mkdir(parents=True, exist_ok=True)

    jobs = [({**shared, **params}, seed) for params in cfg["scenarios"].values() for seed in cfg["seeds"]]
    with ProcessPoolExecutor() as pool:
        results = list(pool.map(run_one, jobs))

    n = len(cfg["seeds"])
    scenarios = {name: aggregate(results[i * n:(i + 1) * n]) for i, name in enumerate(cfg["scenarios"])}
    metrics = {"config": cfg, "code_version": code_version(), "scenarios": scenarios}
    (out / "metrics.json").write_text(json.dumps(metrics, indent=2, sort_keys=True) + "\n")
    (out / "summary.md").write_text(summary(cfg, scenarios))
    if not args.no_plots:
        plots(scenarios, out)
    print(f"wrote {out}")


def summary(cfg: dict, scenarios: dict) -> str:
    methods = ("avg", "avg_min3", "ours_nocap", "ours_no_mutual", "ours")
    lines = [f"# {cfg['name']}", "", cfg.get("description", ""), "",
             f"Mean ± std over {len(cfg['seeds'])} seeds. Flag budget: lowest {FLAG_BUDGET:.0%} of users.", ""]
    for name, rows in scenarios.items():
        lines += [f"## {name}", "",
                  "| method | Spearman | Spearman (≤3 rentals) | avg precision | precision@budget | recall@budget |",
                  "|---|---|---|---|---|---|"]
        for m in methods:
            r = rows[m]
            cell = lambda k: f"{r[k]['mean']:.3f} ± {r[k]['std']:.3f}"
            lines.append(f"| {m} | {cell('spearman')} | {cell('spearman_new_users')} | {cell('average_precision')} "
                         f"| {cell('precision_at_budget')} | {cell('recall_at_budget')} |")
        lines += ["", "| deposit policy | uncovered loss, mln KZT | mean deposit, honest | mean deposit, risky |",
                  "|---|---|---|---|"]
        for p in ("fixed", "dynamic_avg", "dynamic_ours"):
            r = rows[f"deposit:{p}"]
            lines.append(f"| {p} | {r['uncovered_loss_mln_kzt']['mean']:.2f} ± {r['uncovered_loss_mln_kzt']['std']:.2f} "
                         f"| {r['mean_deposit_honest']['mean']:,.0f} | {r['mean_deposit_risky']['mean']:,.0f} |")
        lines.append("")
    return "\n".join(lines)


def plots(scenarios: dict, out: Path):
    import matplotlib

    matplotlib.use("Agg")
    import matplotlib.pyplot as plt

    # sweeps are scenarios named "<param>=<value>"
    sweeps: dict[str, list] = {}
    for name in scenarios:
        if "=" in name:
            param, value = name.split("=", 1)
            sweeps.setdefault(param, []).append((float(value), name))

    for param, points in sweeps.items():
        points.sort()
        fig, axes = plt.subplots(1, 2, figsize=(10, 4))
        for ax, metric, title in ((axes[0], "spearman", "Spearman vs true reliability"),
                                  (axes[1], "average_precision", "Average precision, risky users")):
            for m in ("avg", "avg_min3", "ours_nocap", "ours_no_mutual", "ours"):
                means = [scenarios[n][m][metric]["mean"] for _, n in points]
                stds = [scenarios[n][m][metric]["std"] for _, n in points]
                ax.errorbar([v for v, _ in points], means, yerr=stds, marker="o", capsize=3, label=m)
            ax.set_xlabel(param)
            ax.set_title(title)
            ax.grid(alpha=0.3)
        axes[0].legend()
        fig.tight_layout()
        fig.savefig(out / f"sweep_{param}.png", dpi=120)
        plt.close(fig)

    if "base" in scenarios:
        policies = ("fixed", "dynamic_avg", "dynamic_ours")
        rows = [scenarios["base"][f"deposit:{p}"] for p in policies]
        fig, axes = plt.subplots(1, 2, figsize=(10, 4))
        axes[0].bar(policies, [r["uncovered_loss_mln_kzt"]["mean"] for r in rows],
                    yerr=[r["uncovered_loss_mln_kzt"]["std"] for r in rows], capsize=3)
        axes[0].set_title("Lenders' uncovered loss, mln KZT")
        x = np.arange(len(policies))
        axes[1].bar(x - 0.2, [r["mean_deposit_honest"]["mean"] for r in rows], 0.4, label="honest")
        axes[1].bar(x + 0.2, [r["mean_deposit_risky"]["mean"] for r in rows], 0.4, label="risky")
        axes[1].set_xticks(x, policies)
        axes[1].set_title("Mean deposit held, KZT")
        axes[1].legend()
        fig.tight_layout()
        fig.savefig(out / "deposits_base.png", dpi=120)
        plt.close(fig)


if __name__ == "__main__":
    main()

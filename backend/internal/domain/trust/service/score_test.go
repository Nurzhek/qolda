package service

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"qolda/internal/domain/trust/model"
)

const tolerance = 1e-9

// json shapes of reputation/golden.py (the domain model has no serialization tags)
type goldenHistory struct {
	RatingSum    float64 `json:"rating_sum"`
	RatingCount  float64 `json:"rating_count"`
	Returns      float64 `json:"returns"`
	OnTime       float64 `json:"on_time"`
	Closed       float64 `json:"closed"`
	DisputesLost float64 `json:"disputes_lost"`
	AgeDays      float64 `json:"age_days"`
}

type goldenCampus struct {
	Rating float64 `json:"rating"`
	OnTime float64 `json:"on_time"`
	Clean  float64 `json:"clean"`
}

type goldenVote struct {
	Rater  int64   `json:"rater"`
	Rating float64 `json:"rating"`
	OnTime float64 `json:"on_time"`
	Clean  float64 `json:"clean"`
}

type golden struct {
	Scores []struct {
		History  goldenHistory `json:"history"`
		Campus   goldenCampus  `json:"campus"`
		Neutral  float64       `json:"neutral"`
		Expected struct {
			Score  float64 `json:"score"`
			Rating float64 `json:"rating"`
			OnTime float64 `json:"on_time"`
			Clean  float64 `json:"clean"`
			Age    float64 `json:"age"`
		} `json:"expected"`
	} `json:"scores"`
	Deposits []struct {
		Base    float64 `json:"base"`
		Score   float64 `json:"score"`
		Neutral float64 `json:"neutral"`
		Factor  float64 `json:"factor"`
		Deposit int64   `json:"deposit"`
	} `json:"deposits"`
	Graph struct {
		Votes    map[string][]goldenVote `json:"votes"`
		Ages     map[string]float64      `json:"ages"`
		Campus   goldenCampus            `json:"campus"`
		Expected map[string]float64      `json:"expected"`
	} `json:"graph"`
}

func (c goldenCampus) domain() model.Campus {
	return model.Campus{Rating: c.Rating, OnTime: c.OnTime, Clean: c.Clean}
}

func loadGolden(t *testing.T) golden {
	t.Helper()
	data, err := os.ReadFile("testdata/golden.json")
	require.NoError(t, err, "regenerate with: cd reputation && python golden.py")
	var g golden
	require.NoError(t, json.Unmarshal(data, &g))
	require.NotEmpty(t, g.Scores)
	return g
}

// TestMatchesPrototype — the Go port reproduces the Python prototype on every golden case
func TestMatchesPrototype(t *testing.T) {
	g := loadGolden(t)

	for i, c := range g.Scores {
		h := model.History(c.History)
		got := TrustScore(h, c.Campus.domain())
		assert.InDelta(t, c.Expected.Score, got.Score, tolerance, "score, case %d", i)
		assert.InDelta(t, c.Expected.Rating, got.Rating, tolerance, "rating, case %d", i)
		assert.InDelta(t, c.Expected.OnTime, got.OnTime, tolerance, "on_time, case %d", i)
		assert.InDelta(t, c.Expected.Clean, got.Clean, tolerance, "clean, case %d", i)
		assert.InDelta(t, c.Expected.Age, got.Age, tolerance, "age, case %d", i)
		assert.InDelta(t, c.Neutral, NeutralScore(c.Campus.domain()), tolerance, "neutral, case %d", i)
	}

	for i, d := range g.Deposits {
		assert.InDelta(t, d.Factor, DepositFactor(d.Score, d.Neutral), tolerance, "factor, case %d", i)
		assert.Equal(t, d.Deposit, RecommendedDeposit(d.Base, d.Score, d.Neutral), "deposit, case %d", i)
	}

	votes := map[int64][]model.Vote{}
	for u, vs := range g.Graph.Votes {
		id, err := strconv.ParseInt(u, 10, 64)
		require.NoError(t, err)
		for _, v := range vs {
			votes[id] = append(votes[id], model.Vote(v))
		}
	}
	ages := map[int64]float64{}
	for u, a := range g.Graph.Ages {
		id, err := strconv.ParseInt(u, 10, 64)
		require.NoError(t, err)
		ages[id] = a
	}
	got := TrustScores(votes, ages, g.Graph.Campus.domain())
	require.Len(t, got, len(g.Graph.Expected))
	for u, want := range g.Graph.Expected {
		id, _ := strconv.ParseInt(u, 10, 64)
		assert.InDelta(t, want, got[id], tolerance, "graph score, user %s", u)
	}
}

func TestProperties(t *testing.T) {
	c := model.DefaultCampus()
	neutral := NeutralScore(c)

	// the design document's example: one 5-star review does not beat forty 4.8-star reviews
	one := TrustScore(model.History{RatingSum: 5, RatingCount: 1, Returns: 1, OnTime: 1, Closed: 1}, c)
	forty := TrustScore(model.History{RatingSum: 192, RatingCount: 40, Returns: 40, OnTime: 40, Closed: 40}, c)
	assert.Less(t, one.Score, forty.Score)

	assert.Equal(t, 1.0, DepositFactor(neutral, neutral))
	assert.Equal(t, int64(10_000), RecommendedDeposit(10_000, neutral, neutral))
	prev := int64(1 << 62)
	for s := 0.0; s <= 100; s += 5 {
		d := RecommendedDeposit(20_000, s, neutral)
		assert.LessOrEqual(t, d, prev, "a higher score never means a higher deposit")
		prev = d
	}

	// a mutual vote is discounted
	votes := map[int64][]model.Vote{1: {{Rater: 2, Rating: 5, OnTime: 1, Clean: 1}}, 2: {{Rater: 1, Rating: 5, OnTime: 1, Clean: 1}}}
	assert.Equal(t, map[int64]float64{2: ReciprocalWeight}, ReciprocityWeights(1, votes))
}

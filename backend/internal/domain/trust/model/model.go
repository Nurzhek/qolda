package model

// History — what the platform knows about a user as a borrower.
type History struct {
	RatingSum    float64
	RatingCount  float64
	Returns      float64 // rentals that were due back
	OnTime       float64 // of those, returned on time (an unreturned item is not on time)
	Closed       float64 // closed orders (any outcome)
	DisputesLost float64 // disputes where the user was found at fault
	AgeDays      float64
}

// Campus — campus-wide averages every component is shrunk toward.
type Campus struct {
	Rating float64
	OnTime float64
	Clean  float64 // share of closed orders without a lost dispute
}

// DefaultCampus is used until the campus has its own statistics.
func DefaultCampus() Campus {
	return Campus{Rating: 4.5, OnTime: 0.9, Clean: 0.95}
}

// Score — the trust score (0..100) and its components (0..1), shown to the user as the explanation.
type Score struct {
	Score  float64
	Rating float64
	OnTime float64
	Clean  float64
	Age    float64
}

// Vote — what one counterparty's rentals say about a user, averaged over all rentals between the two.
type Vote struct {
	Rater  int64
	Rating float64 // mean stars
	OnTime float64 // share returned on time
	Clean  float64 // share of closed orders without a lost dispute
}

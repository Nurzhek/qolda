package constant

import "slices"

const (
	ServiceName = "qolda"

	MaxPageSize = 1000
)

// Version, Commit, BuiltAt are set at build time (Makefile, ldflags): what exactly is running,
// version, full commit SHA, build time (logged on start).
var (
	Version = "dev"
	Commit  = ""
	BuiltAt = ""
)

// roles
const (
	RoleStudent   = "student"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

// listing types
const (
	ListingTypeRent     = "rent"
	ListingTypeExchange = "exchange"
	ListingTypeSale     = "sale"
)

// listing statuses
const (
	ListingStatusActive   = "active"
	ListingStatusPaused   = "paused"
	ListingStatusArchived = "archived"
)

var ListingTypes = []string{ListingTypeRent, ListingTypeExchange, ListingTypeSale}

var ListingStatuses = []string{ListingStatusActive, ListingStatusPaused, ListingStatusArchived}

// ponytail: categories are hardcoded; move them to a table when they need editing without a release
var Categories = []string{
	"electronics", "books", "study", "sports_outdoor", "travel",
	"furniture", "kitchen", "clothing", "music", "other",
}

func IsCategory(v string) bool { return slices.Contains(Categories, v) }

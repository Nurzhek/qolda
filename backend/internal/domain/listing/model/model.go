package model

import (
	"time"

	commonModel "qolda/internal/domain/common/model"
)

// Main is a listing. Availability dates are calendar dates, "YYYY-MM-DD"; empty means not set.
type Main struct {
	Id            int64
	CreatedAt     time.Time
	ModifiedAt    time.Time
	OwnerId       int64
	Type          string
	Status        string
	Title         string
	Description   string
	Category      string
	Campus        string
	Photos        []string
	Price         int64
	Deposit       int64
	Wanted        string
	AvailableFrom string
	AvailableTo   string

	// filled by the usecase for responses, untouched by the repo
	Owner *Owner
}

// Owner is a short owner card in responses
type Owner struct {
	Id     int64
	Name   string
	Campus string
}

type Edit struct {
	CreatedAt     *time.Time
	ModifiedAt    *time.Time
	OwnerId       *int64
	Type          *string
	Status        *string
	Title         *string
	Description   *string
	Category      *string
	Campus        *string
	Photos        *[]string
	Price         *int64
	Deposit       *int64
	Wanted        *string
	AvailableFrom *string
	AvailableTo   *string
}

type ListReq struct {
	commonModel.ListParams

	Query    *string
	Type     *string
	Category *string
	Campus   *string
	Status   *string
	OwnerId  *int64
	PriceMin *int64
	PriceMax *int64
	FreeFrom *string // the availability window covers [FreeFrom, FreeTo]
	FreeTo   *string
}

// Apply returns the listing state after an edit (to check invariants before writing)
func (m Main) Apply(e *Edit) Main {
	set := func(dst *string, src *string) {
		if src != nil {
			*dst = *src
		}
	}
	set(&m.Type, e.Type)
	set(&m.Status, e.Status)
	set(&m.Title, e.Title)
	set(&m.Description, e.Description)
	set(&m.Category, e.Category)
	set(&m.Campus, e.Campus)
	set(&m.Wanted, e.Wanted)
	set(&m.AvailableFrom, e.AvailableFrom)
	set(&m.AvailableTo, e.AvailableTo)
	if e.OwnerId != nil {
		m.OwnerId = *e.OwnerId
	}
	if e.Photos != nil {
		m.Photos = *e.Photos
	}
	if e.Price != nil {
		m.Price = *e.Price
	}
	if e.Deposit != nil {
		m.Deposit = *e.Deposit
	}
	return m
}

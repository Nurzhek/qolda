package model

import (
	"time"

	domainModel "qolda/internal/domain/listing/model"
)

type Select struct {
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
}

func (m *Select) ListColumnMap() map[string]any {
	return map[string]any{
		"id":          &m.Id,
		"created_at":  &m.CreatedAt,
		"modified_at": &m.ModifiedAt,
		"owner_id":    &m.OwnerId,
		"type":        &m.Type,
		"status":      &m.Status,
		"title":       &m.Title,
		"description": &m.Description,
		"category":    &m.Category,
		"campus":      &m.Campus,
		"photos":      &m.Photos,
		"price":       &m.Price,
		"deposit":     &m.Deposit,
		"wanted":      &m.Wanted,
		// date → "YYYY-MM-DD", NULL → ""
		`coalesce(to_char(available_from, 'YYYY-MM-DD'), '')`: &m.AvailableFrom,
		`coalesce(to_char(available_to, 'YYYY-MM-DD'), '')`:   &m.AvailableTo,
	}
}

func (m *Select) PKColumnMap() map[string]any {
	return map[string]any{"id": m.Id}
}

func (m *Select) DefaultSortColumns() []string {
	return []string{"created_at DESC", "id DESC"}
}

// DTO

func EncodeSelect(v *Select, _ int) *domainModel.Main {
	return &domainModel.Main{
		Id:            v.Id,
		CreatedAt:     v.CreatedAt,
		ModifiedAt:    v.ModifiedAt,
		OwnerId:       v.OwnerId,
		Type:          v.Type,
		Status:        v.Status,
		Title:         v.Title,
		Description:   v.Description,
		Category:      v.Category,
		Campus:        v.Campus,
		Photos:        v.Photos,
		Price:         v.Price,
		Deposit:       v.Deposit,
		Wanted:        v.Wanted,
		AvailableFrom: v.AvailableFrom,
		AvailableTo:   v.AvailableTo,
	}
}

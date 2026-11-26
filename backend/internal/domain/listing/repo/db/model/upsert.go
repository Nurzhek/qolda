package model

import (
	"time"

	domainModel "qolda/internal/domain/listing/model"
)

type Upsert struct {
	PKId  int64
	NewId int64

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

func (m *Upsert) CreateColumnMap() map[string]any {
	result := make(map[string]any, 15)
	if m.CreatedAt != nil {
		result["created_at"] = *m.CreatedAt
	}
	if m.ModifiedAt != nil {
		result["modified_at"] = *m.ModifiedAt
	}
	if m.OwnerId != nil {
		result["owner_id"] = *m.OwnerId
	}
	if m.Type != nil {
		result["type"] = *m.Type
	}
	if m.Status != nil {
		result["status"] = *m.Status
	}
	if m.Title != nil {
		result["title"] = *m.Title
	}
	if m.Description != nil {
		result["description"] = *m.Description
	}
	if m.Category != nil {
		result["category"] = *m.Category
	}
	if m.Campus != nil {
		result["campus"] = *m.Campus
	}
	if m.Photos != nil {
		// pgx writes a nil slice as NULL, the column is not null
		result["photos"] = append([]string{}, *m.Photos...)
	}
	if m.Price != nil {
		result["price"] = *m.Price
	}
	if m.Deposit != nil {
		result["deposit"] = *m.Deposit
	}
	if m.Wanted != nil {
		result["wanted"] = *m.Wanted
	}
	if m.AvailableFrom != nil {
		result["available_from"] = nullIfEmpty(*m.AvailableFrom)
	}
	if m.AvailableTo != nil {
		result["available_to"] = nullIfEmpty(*m.AvailableTo)
	}
	return result
}

func (m *Upsert) UpdateColumnMap() map[string]any {
	return m.CreateColumnMap()
}

func (m *Upsert) PKColumnMap() map[string]any {
	return map[string]any{"id": m.PKId}
}

func (m *Upsert) ReturningColumnMap() map[string]any {
	return map[string]any{"id": &m.NewId}
}

// DTO

func DecodeUpsert(v *domainModel.Edit) *Upsert {
	return &Upsert{
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

// nullIfEmpty writes an empty date as NULL
func nullIfEmpty(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

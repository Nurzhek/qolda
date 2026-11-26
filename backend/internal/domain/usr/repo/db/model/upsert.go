package model

import (
	"time"

	domainModel "qolda/internal/domain/usr/model"
)

type Upsert struct {
	PKId  int64
	NewId int64

	CreatedAt  *time.Time
	ModifiedAt *time.Time
	Email      *string
	Name       *string
	Campus     *string
	Role       *string
	Blocked    *bool
	Telegram   *string
}

func (m *Upsert) CreateColumnMap() map[string]any {
	result := make(map[string]any, 8)
	if m.CreatedAt != nil {
		result["created_at"] = *m.CreatedAt
	}
	if m.ModifiedAt != nil {
		result["modified_at"] = *m.ModifiedAt
	}
	if m.Email != nil {
		result["email"] = *m.Email
	}
	if m.Name != nil {
		result["name"] = *m.Name
	}
	if m.Campus != nil {
		result["campus"] = *m.Campus
	}
	if m.Role != nil {
		result["role"] = *m.Role
	}
	if m.Blocked != nil {
		result["blocked"] = *m.Blocked
	}
	if m.Telegram != nil {
		result["telegram"] = *m.Telegram
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
		CreatedAt:  v.CreatedAt,
		ModifiedAt: v.ModifiedAt,
		Email:      v.Email,
		Name:       v.Name,
		Campus:     v.Campus,
		Role:       v.Role,
		Blocked:    v.Blocked,
		Telegram:   v.Telegram,
	}
}

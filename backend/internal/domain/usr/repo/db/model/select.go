package model

import (
	"time"

	domainModel "qolda/internal/domain/usr/model"
)

type Select struct {
	Id         int64
	CreatedAt  time.Time
	ModifiedAt time.Time
	Email      string
	Name       string
	Campus     string
	Role       string
	Blocked    bool
	Telegram   string
}

func (m *Select) ListColumnMap() map[string]any {
	return map[string]any{
		"id":          &m.Id,
		"created_at":  &m.CreatedAt,
		"modified_at": &m.ModifiedAt,
		"email":       &m.Email,
		"name":        &m.Name,
		"campus":      &m.Campus,
		"role":        &m.Role,
		"blocked":     &m.Blocked,
		"telegram":    &m.Telegram,
	}
}

func (m *Select) PKColumnMap() map[string]any {
	return map[string]any{"id": m.Id}
}

func (m *Select) DefaultSortColumns() []string {
	return []string{"created_at DESC"}
}

// DTO

func EncodeSelect(v *Select, _ int) *domainModel.Main {
	return &domainModel.Main{
		Id:         v.Id,
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

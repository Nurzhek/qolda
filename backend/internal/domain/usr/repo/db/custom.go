package db

import (
	"context"
	"fmt"

	"qolda/internal/domain/usr/model"
	repoModel "qolda/internal/domain/usr/repo/db/model"
)

var allowedSortFields = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"name":       "name",
}

func (r *Repo) getConditions(pars *model.ListReq) (map[string]any, map[string][]any) {
	conditions := make(map[string]any, 2)
	conditionExps := make(map[string][]any, 2)

	if pars == nil {
		return conditions, conditionExps
	}

	if pars.IDs != nil {
		conditions["id"] = *pars.IDs
	}

	return conditions, conditionExps
}

// GetByEmail: email is unique, so we use Get with email as the "PK"
func (r *Repo) GetByEmail(ctx context.Context, email string) (*model.Main, bool, error) {
	m := &byEmail{key: email}
	found, err := r.ModelStore.Get(ctx, m)
	if err != nil {
		return nil, false, fmt.Errorf("ModelStore.Get: %w", err)
	}
	if !found {
		return nil, false, nil
	}
	return repoModel.EncodeSelect(&m.Select, 0), true, nil
}

type byEmail struct {
	repoModel.Select
	key string
}

func (m *byEmail) PKColumnMap() map[string]any {
	return map[string]any{"email": m.key}
}

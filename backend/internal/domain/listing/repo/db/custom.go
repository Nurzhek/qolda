package db

import (
	"strings"

	"qolda/internal/domain/listing/model"
)

var allowedSortFields = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"price":      "price",
}

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

func (r *Repo) getConditions(pars *model.ListReq) (map[string]any, map[string][]any) {
	conditions := make(map[string]any, 6)
	conditionExps := make(map[string][]any, 6)

	if pars == nil {
		return conditions, conditionExps
	}

	if pars.Type != nil {
		conditions["type"] = *pars.Type
	}
	if pars.Category != nil {
		conditions["category"] = *pars.Category
	}
	if pars.Campus != nil {
		conditions["campus"] = *pars.Campus
	}
	if pars.Status != nil {
		conditions["status"] = *pars.Status
	}
	if pars.OwnerId != nil {
		conditions["owner_id"] = *pars.OwnerId
	}
	if pars.PriceMin != nil {
		conditionExps["price >= ?"] = []any{*pars.PriceMin}
	}
	if pars.PriceMax != nil {
		conditionExps["price <= ?"] = []any{*pars.PriceMax}
	}
	// ponytail: ILIKE on title/description is a full scan; pg_trgm/tsvector when the catalogue grows
	if pars.Query != nil {
		q := "%" + likeEscaper.Replace(*pars.Query) + "%"
		conditionExps["(title ILIKE ? OR description ILIKE ?)"] = []any{q, q}
	}
	// a single date means "available on that day"
	if pars.FreeFrom != nil || pars.FreeTo != nil {
		from, to := pars.FreeFrom, pars.FreeTo
		if from == nil {
			from = to
		}
		if to == nil {
			to = from
		}
		conditionExps["available_from <= ?::date"] = []any{*from}
		conditionExps["available_to >= ?::date"] = []any{*to}
	}

	return conditions, conditionExps
}

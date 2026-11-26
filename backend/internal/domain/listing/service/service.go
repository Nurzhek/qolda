package service

import (
	"context"
	"fmt"
	"time"

	"qolda/internal/constant"
	"qolda/internal/domain/listing/model"
	"qolda/internal/errs"
)

type Service struct {
	repoDb RepoDbI
}

func New(repoDb RepoDbI) *Service {
	return &Service{repoDb: repoDb}
}

func (s *Service) List(ctx context.Context, pars *model.ListReq) ([]*model.Main, int64, error) {
	items, tCount, err := s.repoDb.List(ctx, pars)
	if err != nil {
		return nil, 0, fmt.Errorf("repoDb.List: %w", err)
	}
	return items, tCount, nil
}

func (s *Service) Get(ctx context.Context, id int64, errNE bool) (*model.Main, bool, error) {
	result, found, err := s.repoDb.Get(ctx, id)
	if err != nil {
		return nil, false, fmt.Errorf("repoDb.Get: %w", err)
	}
	if !found {
		if errNE {
			return nil, false, errs.ObjectNotFound
		}
		return nil, false, nil
	}
	return result, true, nil
}

func (s *Service) Create(ctx context.Context, obj *model.Edit) (int64, error) {
	obj.ModifiedAt = nil // set by the DB default
	normalize(obj)

	newId, err := s.repoDb.Create(ctx, obj)
	if err != nil {
		return 0, fmt.Errorf("repoDb.Create: %w", err)
	}
	return newId, nil
}

func (s *Service) Update(ctx context.Context, id int64, obj *model.Edit) error {
	obj.ModifiedAt = new(time.Now())
	normalize(obj)

	if err := s.repoDb.Update(ctx, id, obj); err != nil {
		return fmt.Errorf("repoDb.Update: %w", err)
	}
	return nil
}

// normalize enforces the type invariant: fields of other types are zeroed so no leftovers stay in the DB
// (e.g. a deposit on a sale after the form changed). The type is known only on create.
func normalize(obj *model.Edit) {
	if obj.Type == nil {
		return
	}
	switch *obj.Type {
	case constant.ListingTypeRent:
		obj.Wanted = new("")
	case constant.ListingTypeSale:
		obj.Deposit, obj.Wanted = new(int64(0)), new("")
		obj.AvailableFrom, obj.AvailableTo = new(""), new("")
	case constant.ListingTypeExchange:
		obj.Price, obj.Deposit = new(int64(0)), new(int64(0))
		obj.AvailableFrom, obj.AvailableTo = new(""), new("")
	}
}

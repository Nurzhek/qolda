package model

import (
	"time"

	commonModel "qolda/internal/domain/common/model"
)

type Main struct {
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

type Edit struct {
	CreatedAt  *time.Time
	ModifiedAt *time.Time
	Email      *string
	Name       *string
	Campus     *string
	Role       *string
	Blocked    *bool
	Telegram   *string
}

type ListReq struct {
	commonModel.ListParams

	IDs *[]int64
}

package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"qolda/internal/config"
	"qolda/internal/constant"
	proto "qolda/pkg/proto/qolda_v1"
)

// Dict serves dictionaries and public config; there is no logic, so no usecase
type Dict struct {
	proto.UnsafeDictServiceServer
	googleEnabled bool
}

func NewDict(googleEnabled bool) *Dict {
	return &Dict{googleEnabled: googleEnabled}
}

func (h *Dict) Get(context.Context, *emptypb.Empty) (*proto.DictRep, error) {
	return &proto.DictRep{
		Categories:     constant.Categories,
		Campuses:       config.Conf.Campuses,
		AllowedDomains: config.Conf.AllowedDomains,
		GoogleLogin:    h.googleEnabled,
		DevLogin:       config.Conf.AuthDevLogin,
	}, nil
}

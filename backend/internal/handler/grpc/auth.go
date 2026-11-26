package grpc

import (
	"context"

	"qolda/internal/handler/grpc/dto"
	usecase "qolda/internal/usecase/auth"
	proto "qolda/pkg/proto/qolda_v1"
)

type Auth struct {
	proto.UnsafeAuthServiceServer
	usecase *usecase.Usecase
}

func NewAuth(uc *usecase.Usecase) *Auth {
	return &Auth{usecase: uc}
}

func (h *Auth) DevLogin(ctx context.Context, req *proto.AuthDevLoginReq) (*proto.AuthTokenRep, error) {
	tokens, err := h.usecase.DevLogin(ctx, req.GetEmail(), req.GetName())
	if err != nil {
		return nil, err
	}
	return dto.EncodeAuthTokenRep(tokens), nil
}

func (h *Auth) Refresh(ctx context.Context, req *proto.AuthRefreshReq) (*proto.AuthTokenRep, error) {
	tokens, err := h.usecase.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	return dto.EncodeAuthTokenRep(tokens), nil
}

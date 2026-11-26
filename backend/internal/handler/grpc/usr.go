package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"qolda/internal/handler/grpc/dto"
	usecase "qolda/internal/usecase/usr"
	proto "qolda/pkg/proto/qolda_v1"
)

type Usr struct {
	proto.UnsafeUsrServiceServer
	usecase *usecase.Usecase
}

func NewUsr(uc *usecase.Usecase) *Usr {
	return &Usr{usecase: uc}
}

func (h *Usr) Get(ctx context.Context, req *proto.UsrGetReq) (*proto.UsrMain, error) {
	item, err := h.usecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return dto.EncodeUsrMain(item, 0), nil
}

func (h *Usr) GetMe(ctx context.Context, _ *emptypb.Empty) (*proto.UsrMain, error) {
	item, err := h.usecase.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	return dto.EncodeUsrMain(item, 0), nil
}

func (h *Usr) UpdateMe(ctx context.Context, req *proto.UsrUpdateMeReq) (*emptypb.Empty, error) {
	if err := h.usecase.UpdateMe(ctx, dto.DecodeUsrUpdateMeReq(req)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

package grpc

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"

	"qolda/internal/handler/grpc/dto"
	usecase "qolda/internal/usecase/listing"
	"qolda/pkg/proto/common"
	proto "qolda/pkg/proto/qolda_v1"
)

type Listing struct {
	proto.UnsafeListingServiceServer
	usecase *usecase.Usecase
}

func NewListing(uc *usecase.Usecase) *Listing {
	return &Listing{usecase: uc}
}

func (h *Listing) List(ctx context.Context, req *proto.ListingListReq) (*proto.ListingListRep, error) {
	items, tCount, err := h.usecase.List(ctx, dto.DecodeListingListReq(req), req.GetMine())
	if err != nil {
		return nil, err
	}
	return &proto.ListingListRep{
		PaginationInfo: &common.PaginationInfoSt{
			Page:       req.GetListParams().GetPage(),
			PageSize:   req.GetListParams().GetPageSize(),
			TotalCount: tCount,
		},
		Results: lo.Map(items, dto.EncodeListingMain),
	}, nil
}

func (h *Listing) Get(ctx context.Context, req *proto.ListingGetReq) (*proto.ListingMain, error) {
	item, err := h.usecase.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return dto.EncodeListingMain(item, 0), nil
}

func (h *Listing) Create(ctx context.Context, req *proto.ListingCreateReq) (*proto.ListingCreateRep, error) {
	newId, err := h.usecase.Create(ctx, dto.DecodeListingCreateReq(req))
	if err != nil {
		return nil, err
	}
	return &proto.ListingCreateRep{Id: newId}, nil
}

func (h *Listing) Update(ctx context.Context, req *proto.ListingUpdateReq) (*emptypb.Empty, error) {
	if err := h.usecase.Update(ctx, req.GetId(), dto.DecodeListingUpdateReq(req)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Listing) SetStatus(ctx context.Context, req *proto.ListingSetStatusReq) (*emptypb.Empty, error) {
	if err := h.usecase.SetStatus(ctx, req.GetId(), req.GetStatus()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

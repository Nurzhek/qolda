package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domainModel "qolda/internal/domain/listing/model"
	proto "qolda/pkg/proto/qolda_v1"
)

// domain → proto

func EncodeListingMain(v *domainModel.Main, _ int) *proto.ListingMain {
	if v == nil {
		return nil
	}
	result := &proto.ListingMain{
		Id:            v.Id,
		CreatedAt:     timestamppb.New(v.CreatedAt),
		ModifiedAt:    timestamppb.New(v.ModifiedAt),
		OwnerId:       v.OwnerId,
		Type:          v.Type,
		Status:        v.Status,
		Title:         v.Title,
		Description:   v.Description,
		Category:      v.Category,
		Campus:        v.Campus,
		Photos:        v.Photos,
		Price:         int32(v.Price), // amounts are capped in the usecase (maxAmount)
		Deposit:       int32(v.Deposit),
		Wanted:        v.Wanted,
		AvailableFrom: v.AvailableFrom,
		AvailableTo:   v.AvailableTo,
	}
	if v.Owner != nil {
		result.Owner = &proto.ListingOwner{Id: v.Owner.Id, Name: v.Owner.Name, Campus: v.Owner.Campus}
	}
	return result
}

// proto → domain

func DecodeListingListReq(v *proto.ListingListReq) *domainModel.ListReq {
	return &domainModel.ListReq{
		ListParams: DecodeListParams(v.GetListParams()),
		Query:      v.Query,
		Type:       v.Type,
		Category:   v.Category,
		Campus:     v.Campus,
		Status:     v.Status,
		OwnerId:    v.OwnerId,
		PriceMin:   i64(v.PriceMin),
		PriceMax:   i64(v.PriceMax),
		FreeFrom:   v.FreeFrom,
		FreeTo:     v.FreeTo,
	}
}

func DecodeListingCreateReq(v *proto.ListingCreateReq) *domainModel.Edit {
	return &domainModel.Edit{
		Type:          v.Type,
		Title:         v.Title,
		Description:   v.Description,
		Category:      v.Category,
		Campus:        v.Campus,
		Photos:        new(v.GetPhotos()),
		Price:         i64(v.Price),
		Deposit:       i64(v.Deposit),
		Wanted:        v.Wanted,
		AvailableFrom: v.AvailableFrom,
		AvailableTo:   v.AvailableTo,
	}
}

// photos are always replaced wholesale (repeated has no presence)
func DecodeListingUpdateReq(v *proto.ListingUpdateReq) *domainModel.Edit {
	return &domainModel.Edit{
		Title:         v.Title,
		Description:   v.Description,
		Category:      v.Category,
		Campus:        v.Campus,
		Photos:        new(v.GetPhotos()),
		Price:         i64(v.Price),
		Deposit:       i64(v.Deposit),
		Wanted:        v.Wanted,
		AvailableFrom: v.AvailableFrom,
		AvailableTo:   v.AvailableTo,
	}
}

func i64(v *int32) *int64 {
	if v == nil {
		return nil
	}
	return new(int64(*v))
}

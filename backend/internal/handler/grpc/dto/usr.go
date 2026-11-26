package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	domainModel "qolda/internal/domain/usr/model"
	proto "qolda/pkg/proto/qolda_v1"
)

func EncodeUsrMain(v *domainModel.Main, _ int) *proto.UsrMain {
	if v == nil {
		return nil
	}
	return &proto.UsrMain{
		Id:        v.Id,
		CreatedAt: timestamppb.New(v.CreatedAt),
		Email:     v.Email,
		Name:      v.Name,
		Campus:    v.Campus,
		Role:      v.Role,
		Telegram:  v.Telegram,
	}
}

func DecodeUsrUpdateMeReq(v *proto.UsrUpdateMeReq) *domainModel.Edit {
	return &domainModel.Edit{
		Name:     v.Name,
		Campus:   v.Campus,
		Telegram: v.Telegram,
	}
}

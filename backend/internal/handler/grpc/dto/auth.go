package dto

import (
	sessionModel "qolda/internal/domain/session/model"
	proto "qolda/pkg/proto/qolda_v1"
)

func EncodeAuthTokenRep(v *sessionModel.Tokens) *proto.AuthTokenRep {
	return &proto.AuthTokenRep{AccessToken: v.Access, RefreshToken: v.Refresh}
}

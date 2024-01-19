package model

import "folkatech-customerIdentity/src/pkg/common"

type AuthLogin struct {
	ID             int64  `json:"id" bson:"_id"`
	Username       string `json:"username" bson:"username" validate:"required"`
	AccountNumber  string `json:"account_number" bson:"account_number"`
	EmailAddress   string `json:"email_address" bson:"email_address" validate:"required"`
	IdentityNumber string `json:"identity_number" bson:"identity_number"`
	Password       string `json:"password,omitempty" bson:"password" validate:"required"`
	common.BaseTime
}

type AuthLoginResponse struct {
	Tokens JwtModel  `json:"token"`
	User   AuthLogin `json:"user"`
}

type JwtModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

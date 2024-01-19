package model

import "folkatech-customerIdentity/src/pkg/common"

type User struct {
	ID             int64  `json:"id" bson:"_id"`
	Username       string `json:"username" bson:"username"`
	AccountNumber  string `json:"account_number" bson:"account_number"`
	EmailAddress   string `json:"email_address" bson:"email_address"`
	IdentityNumber string `json:"identity_number" bson:"identity_number"`
	Password       string `json:"password,omitempty" bson:"password"`
	common.BaseTime
}

type UserResp struct {
	Username       string `json:"username" bson:"username"`
	AccountNumber  string `json:"account_number" bson:"account_number"`
	EmailAddress   string `json:"email_address" bson:"email_address"`
	IdentityNumber string `json:"identity_number" bson:"identity_number"`
	common.BaseTime
}

type FilterUser struct {
	ID             int    `json:"id" bson:"_id"`
	AccountNumber  string `json:"account_number" bson:"account_number"`
	IdentityNumber string `json:"identity_number" bson:"identity_number"`
	Search         string `json:"search" bson:"-"`
}

type ListUsersResponse struct {
	User      []UserResp `json:"user" bson:"user"`
	TotalPage int64      `json:"total_page" bson:"total_page"`
	TotalItem int64      `json:"total_item" bson:"total_item"`
}

type Counter struct {
	ID    string `bson:"_id"`
	Value int64  `bson:"value"`
}

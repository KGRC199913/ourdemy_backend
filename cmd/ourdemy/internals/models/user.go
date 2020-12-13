package models

import "github.com/kamva/mgm/v3"

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Fullname         string `json:"fullname" bson:"fullname"`
	Email            string `json:"email" bson:"email"`
	HPassword        string `json:"hpass" bson:"hpass"`
	CurOtp           string `json:"-" bson:"otp"`
	RecoverCode      string `json:"recover" bson:"recover"`
	RefreshToken     string `json:"rf" bson:"rf"`
}

func NewUser(fullname string, email string, password string) *User {
	return &User{
		Fullname:     fullname,
		Email:        email,
		HPassword:    password, //TODO: hash
		CurOtp:       "123",    //TODO: gen
		RecoverCode:  "123",    //TODO: gen
		RefreshToken: "123",    //TODO: gen
	}
}

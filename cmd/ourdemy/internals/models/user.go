package models

import (
	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	field.DefaultField `bson:",inline"`
	Fullname           string `json:"fullname" bson:"fullname"`
	Email              string `json:"email" bson:"email"`
	HPassword          string `json:"hpass" bson:"hpass"`
	CurOtp             string `json:"-" bson:"otp"`
	RecoverCode        string `json:"recover" bson:"recover"`
	RefreshToken       string `json:"rf" bson:"rf"`
}

func (User) collName() string {
	return "users"
}

func NewUser(fullname string, email string, password string) *User {
	return &User{
		Fullname:     fullname,
		Email:        email,
		HPassword:    password,
		CurOtp:       "123",
		RecoverCode:  "123",
		RefreshToken: "123",
	}
}

func (u *User) Save() error {
	_, err := db.Collection(u.collName()).InsertOne(ctx, u)
	return err
}

func (u *User) FindByName(name string) error {
	return db.Collection(u.collName()).Find(ctx, bson.M{"fullname": name}).One(u)
}

//Hooks
func (u *User) BeforeInsert() error {
	hashed, err := scrypt.GenerateFromPassword([]byte(u.HPassword), scrypt.DefaultParams)
	if err != nil {
		return err
	}
	u.HPassword = string(hashed)
	return nil
}

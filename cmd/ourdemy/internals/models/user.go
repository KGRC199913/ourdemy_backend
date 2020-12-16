package models

import (
	"errors"
	"fmt"
	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type User struct {
	field.DefaultField     `json:"-" bson:",inline"`
	Fullname               string    `json:"fullname" bson:"fullname" binding:"required"`
	Username               string    `json:"username" bson:"username" binding:"required"`
	Email                  string    `json:"email" bson:"email" binding:"required"`
	HPassword              string    `json:"pass" bson:"hpass" binding:"required"`
	CurOtp                 string    `json:"-" bson:"otp"`
	CurOtpExpiredTime      time.Time `json:"-" bson:"otp_exp"`
	RecoverCode            string    `json:"-" bson:"recover"`
	RecoverCodeExpiredTime time.Time `json:"-" bson:"rec_exp"`
	RefreshToken           string    `json:"-" bson:"rf"`
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

func (u *User) FindByUsername(username string) error {
	return db.Collection(u.collName()).Find(ctx, bson.M{"username": username}).One(u)
}

func (u *User) FindByEmail(email string) error {
	return db.Collection(u.collName()).Find(ctx, bson.M{"email": email}).One(u)
}

func (u *User) ConfirmOtp(username string, otp string) error {
	if err := u.FindByUsername(username); err != nil {
		return err
	}

	if u.CurOtp != otp {
		return errors.New("otp not matched")
	}

	if u.CurOtpExpiredTime.Before(time.Now()) {
		return errors.New("otp expired")
	}

	fmt.Println(u.Id)
	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"otp":     "",
			"otp_exp": time.Now(),
		},
	})
}

//Hooks
func (u *User) BeforeInsert() error {
	dupUser := &User{}

	if err := dupUser.FindByUsername(u.Username); err == nil {
		return errors.New("username is already existed")
	}

	if err := dupUser.FindByEmail(u.Email); err == nil {
		return errors.New("user's email is already existed")
	}

	hashed, err := scrypt.GenerateFromPassword([]byte(u.HPassword), scrypt.DefaultParams)
	if err != nil {
		return err
	}
	u.HPassword = string(hashed)
	return nil
}

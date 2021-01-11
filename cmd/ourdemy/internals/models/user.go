package models

import (
	"errors"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/qiniu/qmgo/field"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type User struct {
	field.DefaultField     `json:",inline" bson:",inline"`
	Fullname               string    `json:"fullname" bson:"fullname" binding:"required"`
	Username               string    `json:"username" bson:"username" binding:"required"`
	Email                  string    `json:"email" bson:"email" binding:"required"`
	HPassword              string    `json:"pass" bson:"hpass" binding:"required"`
	CurOtp                 string    `json:"-" bson:"otp"`
	LastOtpUpdated         time.Time `json:"-" bson:"last_otp_updated"`
	CurOtpExpiredTime      time.Time `json:"-" bson:"otp_exp"`
	RecoverCode            string    `json:"-" bson:"recover"`
	RecoverCodeExpiredTime time.Time `json:"-" bson:"rec_exp"`
	RefreshToken           string    `json:"-" bson:"rf"`
	IsLec                  bool      `json:"isLec" bson:"is_lec"`
	IsBanned               bool      `json:"is_banned" bson:"is_banned"`
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

func AllUser() ([]User, error) {
	var res []User
	err := db.Collection(User{}.collName()).Find(ctx, bson.M{}).All(&res)
	if res == nil {
		res = []User{}
	}
	return res, err
}

func (u *User) Ban() error {
	u.IsBanned = true
	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"is_banned": true,
		},
	})
}

func (u *User) Unban() error {
	u.IsBanned = false
	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"is_banned": false,
		},
	})
}

func (u *User) FindById(oid primitive.ObjectID) error {
	return db.Collection(u.collName()).Find(ctx, bson.M{"_id": oid}).One(u)
}

func (u *User) FindByUsername(username string) error {
	return db.Collection(u.collName()).Find(ctx, bson.M{"username": username}).One(u)
}

func (u *User) FindByEmail(email string) error {
	return db.Collection(u.collName()).Find(ctx, bson.M{"email": email}).One(u)
}

func (u *User) FindByIdAndRfToken(oid primitive.ObjectID, rf string) error {
	return db.Collection(u.collName()).Find(ctx, bson.M{"_id": oid, "rf": rf}).One(u)
}

func (u *User) GenerateRfToken() error {
	u.RefreshToken = randstr.Hex(8)
	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"rf": u.RefreshToken,
		},
	})
}

func (u *User) UpdateTokens() (*string, *string, error) {
	newAccessToken, err := ultis.CreateToken(u.Id, u.IsLec)
	if err != nil {
		return nil, nil, err
	}

	rfToken := randstr.Hex(8)

	return &newAccessToken, &rfToken, db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"rf": rfToken,
		},
	})
}

func (u *User) GenerateOtp() error {
	if u.CurOtp != "" {
		return errors.New("otp is already generated")
	}

	u.CurOtp = strings.ToUpper(randstr.Hex(4))
	u.LastOtpUpdated = time.Now()
	u.CurOtpExpiredTime = time.Now().Add(time.Minute * 30)

	return nil
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

	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"otp":     "",
			"otp_exp": time.Now(),
		},
	})
}

func (u *User) UpdateOtp(username string) (*string, error) {
	if err := u.FindByUsername(username); err != nil {
		return nil, err
	}

	//fmt.Println(u.LastOtpUpdated)
	//fmt.Println(time.Now())
	if u.LastOtpUpdated.Add(time.Second * 30).After(time.Now()) {
		return nil, errors.New("otp request too frequent")
	}

	newOtp := strings.ToUpper(randstr.Hex(4))

	return &newOtp, db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"otp":              newOtp,
			"last_otp_updated": time.Now(),
			"otp_exp":          time.Now().Add(time.Minute * 30),
		},
	})
}

func (u *User) GenerateRecoverCode() error {
	u.RecoverCode = strings.ToUpper(randstr.Hex(4))
	u.RecoverCodeExpiredTime = time.Now().AddDate(1, 0, 0)

	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"recover": u.RecoverCode,
			"rec_exp": u.RecoverCodeExpiredTime,
		},
	})
}

func (u *User) ConfirmResetPassword(recoverCode string, newPassword string) error {
	if recoverCode != u.RecoverCode {
		return errors.New("recover code not match")
	}

	if u.RecoverCodeExpiredTime.Before(time.Now()) {
		return errors.New("recover code expired")
	}

	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"hpass":   newPassword,
			"recover": "",
			"rec_exp": time.Now(),
		},
	})
}

func (u *User) UpdateProfile(newFullname string, newEmail string) error {
	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"fullname": newFullname,
			"email":    newEmail,
		},
	})
}

func (u *User) UpdatePassword(newPassword string) error {
	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"hpass": newPassword,
		},
	})
}

func (u *User) UpdateLecturerStatus(isLec bool) error {
	u.IsLec = isLec
	return db.Collection(u.collName()).UpdateOne(ctx, bson.M{
		"_id": u.Id,
	}, bson.M{
		"$set": bson.M{
			"is_lec": isLec,
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

func (u *User) BeforeUpdate() error {
	hashed, err := scrypt.GenerateFromPassword([]byte(u.HPassword), scrypt.DefaultParams)
	if err != nil {
		return err
	}
	u.HPassword = string(hashed)
	return nil
}

func (u *User) AfterInsert() error {
	wl := WatchList{
		UserId:    u.Id,
		CoursesId: []primitive.ObjectID{},
	}

	return wl.Save()
}

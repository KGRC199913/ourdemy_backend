package models

import (
	"errors"
	scrypt "github.com/elithrar/simple-scrypt"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lecturer struct {
	field.DefaultField `json:"-" bson:",inline"`
	Username           string `json:"username" bson:"username" binding:"required"`
	Email              string `json:"email" bson:"email" binding:"required"`
	HPassword          string `json:"pass" bson:"pass" binding:"required"`
}

func (Lecturer) collName() string {
	return "lecturers"
}

func NewLecturer(username string, email string, password string) *Lecturer {
	return &Lecturer{
		Username:  username,
		Email:     email,
		HPassword: password,
	}
}

func (lec *Lecturer) Save() error {
	_, err := db.Collection(lec.collName()).InsertOne(ctx, lec)
	return err
}

func (lec *Lecturer) FindByName(name string) error {
	return db.Collection(lec.collName()).Find(ctx, bson.M{"fullname": name}).One(lec)
}

func (lec *Lecturer) FindById(oid primitive.ObjectID) error {
	return db.Collection(lec.collName()).Find(ctx, bson.M{"_id": oid}).One(lec)
}

//Hooks
func (lec *Lecturer) BeforeInsert() error {
	dupLec := &Lecturer{}
	if err := dupLec.FindByName(lec.Username); err != nil {
		return errors.New("lecturer username is already existed")
	}

	hashed, err := scrypt.GenerateFromPassword([]byte(lec.HPassword), scrypt.DefaultParams)
	if err != nil {
		return err
	}
	lec.HPassword = string(hashed)
	return nil
}

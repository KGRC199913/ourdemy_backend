package models

import (
	"errors"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Approve struct {
	field.DefaultField `json:"-" bson:",inline"`
	LecId              primitive.ObjectID `json:"lid" bson:"lid"`
}

func (Approve) collName() string {
	return "approves"
}

func (a *Approve) Save() error {
	_, err := db.Collection(a.collName()).InsertOne(ctx, a)
	return err
}

func GetAllApprovingLecturers() ([]Approve, error) {
	var apprs []Approve
	err := db.Collection(Approve{}.collName()).Find(ctx, bson.M{}).All(&apprs)

	if err != nil {
		return nil, err
	}
	return apprs, nil
}

func (appr *Approve) FindById(oid primitive.ObjectID) error {
	return db.Collection(appr.collName()).Find(ctx, bson.M{"_id": oid}).One(appr)
}

func (appr *Approve) FindByLecId(oid primitive.ObjectID) error {
	return db.Collection(appr.collName()).Find(ctx, bson.M{"lid": oid}).One(appr)
}

func (appr *Approve) Remove() error {
	return db.Collection(appr.collName()).Remove(ctx, bson.M{"_id": appr.Id})
}

//Hooks
func (appr *Approve) BeforeInsert() error {
	targetUser := &User{}

	if err := targetUser.FindById(appr.Id); err == nil {
		return errors.New("is already a lecturer")
	}

	return nil
}

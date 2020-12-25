package models

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Approve struct {
	field.DefaultField `json:"-" bson:",inline"`
	LecId              primitive.ObjectID `json:"lid" bson:"username"`
}

func (Approve) collName() string {
	return "approves"
}

func GetAllApprovingLecturers() (apprs []Approve, err error) {
	err = db.Collection(Category{}.collName()).Find(ctx, bson.M{}).All(&apprs)
	if err != nil {
		return nil, err
	}
	return apprs, nil
}

func (appr *Approve) FindById(oid primitive.ObjectID) error {
	return db.Collection(appr.collName()).Find(ctx, bson.M{"_id": oid}).One(appr)
}

func (appr *Approve) Remove() error {
	return db.Collection(appr.collName()).Remove(ctx, bson.M{"_id": appr.Id})
}

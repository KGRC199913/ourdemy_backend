package models

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimeMark struct {
	field.DefaultField `bson:",inline"`
	UserId             primitive.ObjectID `json:"uid" bson:"uid"`
	VideoId            primitive.ObjectID `json:"vid" bson:"vid"`
	CurTime            int64              `json:"cur_time" bson:"cur_time"`
}

func (TimeMark) collName() string {
	return "timemarks"
}

func (tm *TimeMark) FindByVideoId(vid primitive.ObjectID) error {
	return db.Collection(tm.collName()).Find(ctx, bson.M{"vid": vid}).One(tm)
}

func (tm *TimeMark) Save() error {
	_, err := db.Collection(tm.collName()).InsertOne(ctx, tm)
	return err
}

func (tm *TimeMark) UpdateTimeMark(newCurTime int64) error {
	return db.Collection(tm.collName()).UpdateOne(ctx, bson.M{
		"_id": tm.Id,
	}, bson.M{
		"$set": bson.M{
			"cur_time": newCurTime,
		},
	})
}

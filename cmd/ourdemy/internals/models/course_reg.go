package models

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseRegister struct {
	field.DefaultField `bson:",inline"`
	UserId             primitive.ObjectID   `json:"uid" bson:"uid"`
	CourseId           []primitive.ObjectID `json:"cid" bson:"cid"`
}

type TimeMark struct {
	field.DefaultField `bson:",inline"`
	UserId             primitive.ObjectID `json:"uid" bson:"uid"`
	VideoId            primitive.ObjectID `json:"vid" bson:"vid"`
	CurTime            int64              `json:"cur_time" bson:"cur_time"`
}

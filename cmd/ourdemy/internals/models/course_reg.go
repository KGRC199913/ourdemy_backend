package models

import "github.com/qiniu/qmgo/field"

type CourseRegister struct {
	field.DefaultField `bson:",inline"`
	UserId             string `json:"uid" bson:"uid"`
	CourseId           string `json:"cid" bson:"cid"`
}

type TimeMark struct {
	field.DefaultField `bson:",inline"`
	UserId             string `json:"uid" bson:"uid"`
	VideoId            string `json:"vid" bson:"vid"`
	CurTime            int64  `json:"cur_time" bson:"cur_time"`
}

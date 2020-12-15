package models

import "github.com/qiniu/qmgo/field"

type Review struct {
	field.DefaultField `bson:",inline"`
	UserId             string  `json:"uid" bson:"uid"`
	CourseId           string  `json:"cid" bson:"cid"`
	Content            string  `json:"content" bson:"content"`
	Score              float32 `json:"score" bson:"score"`
}

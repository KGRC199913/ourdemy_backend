package models

import "github.com/qiniu/qmgo/field"

type WatchList struct {
	field.DefaultField `bson:",inline"`
	UserId             string   `json:"uid" bson:"uid"`
	CoursesId          []string `json:"cids" bson:"cids"`
}

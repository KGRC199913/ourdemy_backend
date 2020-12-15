package models

import "github.com/qiniu/qmgo/field"

type Lecturer struct {
	field.DefaultField `bson:",inline"`
	Username           string `json:"username" bson:"username"`
	HPassword          string `json:"hpass" bson:"hpass"`
}

package models

import (
	"github.com/qiniu/qmgo/field"
)

type Category struct {
	field.DefaultField `bson:",inline"`
	Name               string `json:"name" bson:"name"`
}

type SubCategory struct {
	field.DefaultField `bson:",inline"`
	Name               string `json:"name" bson:"name"`
	ParentCategoryId   string `json:"parentCategoryId" bson:"parentCategoryId"`
}

func NewCategory(name string) (*Category, error) {
	return &Category{
		Name: name,
	}, nil
}

func NewSubCategory(name string, parentCatId string) (*SubCategory, error) {

	return &SubCategory{
		Name:             name,
		ParentCategoryId: parentCatId,
	}, nil
}

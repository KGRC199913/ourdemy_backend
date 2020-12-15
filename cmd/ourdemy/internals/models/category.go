package models

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	field.DefaultField `bson:",inline"`
	Name               string `json:"name" bson:"name"`
}

func (Category) collName() string {
	return "categories"
}

type SubCategory struct {
	field.DefaultField `bson:",inline"`
	Name               string `json:"name" bson:"name"`
	ParentCategoryId   string `json:"parentCategoryId" bson:"parentCategoryId"`
}

func (SubCategory) collName() string {
	return "subcategories"
}

func NewCategory(name string) (*Category, error) {
	return &Category{
		Name: name,
	}, nil
}

func (cat *Category) Save() error {
	_, err := db.Collection(cat.collName()).InsertOne(ctx, cat)
	return err
}

func FindCategoryByName(name string) (cat *Category, err error) {
	cat = &Category{}
	err = db.Collection(cat.collName()).Find(ctx, bson.M{"name": name}).One(cat)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func FindCategoryById(id string) (cat *Category, err error) {
	cat = &Category{}
	var oid primitive.ObjectID
	oid, err = primitive.ObjectIDFromHex(id)
	err = db.Collection(cat.collName()).Find(ctx, bson.M{"_id": oid}).One(cat)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func GetAllCategory() (cats []Category, err error) {
	err = db.Collection(Category{}.collName()).Find(ctx, bson.M{}).All(&cats)
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func NewSubCategory(name string, parentCatId string) (*SubCategory, error) {
	return &SubCategory{
		Name:             name,
		ParentCategoryId: parentCatId,
	}, nil
}

func (subcat *SubCategory) Save() error {
	_, err := db.Collection(subcat.collName()).InsertOne(ctx, subcat)
	return err
}

//Hooks
func (subcat *SubCategory) BeforeInsert() error {
	_, err := FindCategoryById(subcat.ParentCategoryId)
	if err != nil {
		return err
	}
	return nil
}

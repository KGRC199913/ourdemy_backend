package models

import (
	"errors"
	"fmt"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	field.DefaultField `bson:",inline" json:"-"`
	Name               string `json:"name" bson:"name" binding:"required"`
}

func (Category) collName() string {
	return "categories"
}

type SubCategory struct {
	field.DefaultField `bson:",inline" json:"-"`
	Name               string             `json:"name" bson:"name"`
	ParentCategoryId   primitive.ObjectID `json:"parentCategoryId" bson:"parentCategoryId"`
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

func FindCategoryById(oid primitive.ObjectID) (cat *Category, err error) {
	cat = &Category{}
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

// HOOKS
func (cat *Category) BeforeInsert() error {
	_, err := FindCategoryByName(cat.Name)
	if err == nil {
		return errors.New(fmt.Sprintf("duplicate category with name: %s", cat.Name))
	}
	return nil
}

// SUBCAT

func CreateSubCategory(name string, catName string) (*SubCategory, error) {
	cat, err := FindCategoryByName(catName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("category with name: %s not found", catName))
	}

	return &SubCategory{
		Name:             name,
		ParentCategoryId: cat.Id,
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
		return errors.New(fmt.Sprintf("parent category not found with id: %s", subcat.ParentCategoryId))
	}
	return nil
}

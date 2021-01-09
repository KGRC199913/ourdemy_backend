package models

import (
	"errors"
	"fmt"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	field.DefaultField `json:",inline" bson:",inline" json:"-"`
	Name               string `json:"name" bson:"name" binding:"required"`
}

func (Category) collName() string {
	return "categories"
}

type SubCategory struct {
	field.DefaultField `bson:",inline" json:",inline"`
	Name               string             `json:"name" bson:"name"`
	ParentCategoryId   primitive.ObjectID `json:"parentCategoryId" bson:"parentCategoryId"`
}

func (SubCategory) collName() string {
	return "subcategories"
}

func (cat *Category) Save() error {
	_, err := db.Collection(cat.collName()).InsertOne(ctx, cat)
	return err
}

func (cat *Category) FindCategoryByName(name string) error {
	return db.Collection(cat.collName()).Find(ctx, bson.M{"name": name}).One(cat)
}

func (cat *Category) FindCategoryById(oid primitive.ObjectID) error {
	return db.Collection(cat.collName()).Find(ctx, bson.M{"_id": oid}).One(cat)
}

func (cat *Category) UpdateName(name string) error {
	cat.Name = name
	return db.Collection(cat.collName()).UpdateOne(ctx, bson.M{
		"_id": cat.Id,
	}, bson.M{
		"$set": bson.M{
			"name": cat.Name,
		},
	})
}

func (cat *Category) Remove() error {
	return db.Collection(cat.collName()).RemoveId(ctx, cat.Id)
}

func (subcat *SubCategory) Remove() error {
	return db.Collection(subcat.collName()).RemoveId(ctx, subcat.Id)
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
	var dupCat Category
	err := dupCat.FindCategoryByName(cat.Name)
	if err == nil {
		return errors.New(fmt.Sprintf("duplicate category with name: %s", cat.Name))
	}
	return nil
}

// SUBCAT

func CreateSubCategory(name string, catName string) (*SubCategory, error) {
	var cat Category
	err := cat.FindCategoryByName(catName)
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

func (subcat *SubCategory) FindByName(name string) error {
	return db.Collection(subcat.collName()).Find(ctx, bson.M{"name": name}).One(subcat)
}

func FindSubcatsByCatName(name string) ([]SubCategory, error) {
	var cat Category
	if err := cat.FindCategoryByName(name); err != nil {
		return nil, err
	}

	var subcat []SubCategory
	err := db.Collection(SubCategory{}.collName()).Find(ctx, bson.M{"parentCategoryId": cat.Id}).All(&subcat)
	return subcat, err
}

func (subcat *SubCategory) FindSubCategoryById(oid primitive.ObjectID) error {
	return db.Collection(subcat.collName()).Find(ctx, bson.M{"_id": oid}).One(subcat)
}

func FindByParentCategoryId(ParentCatId primitive.ObjectID) ([]SubCategory, error) {
	var subcats []SubCategory
	err := db.Collection(SubCategory{}.collName()).Find(ctx, bson.M{"parentCategoryId": ParentCatId}).All(&subcats)
	if err != nil {
		return nil, err
	}
	if subcats == nil {
		subcats = []SubCategory{}
	}
	return subcats, nil
}

//Hooks
func (subcat *SubCategory) BeforeInsert() error {
	var cat Category
	err := cat.FindCategoryById(subcat.ParentCategoryId)
	if err != nil {
		return errors.New(fmt.Sprintf("parent category not found with id: %s", subcat.ParentCategoryId))
	}
	var subcatDup SubCategory
	if err := subcatDup.FindByName(subcat.Name); err == nil {
		return errors.New(fmt.Sprintf("duplicate subcat: %s ", subcat.Name))
	}

	return nil
}

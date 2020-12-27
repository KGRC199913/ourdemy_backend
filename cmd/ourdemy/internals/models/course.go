package models

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Course struct {
	field.DefaultField `bson:",inline"`
	LecId              primitive.ObjectID `json:"lid" bson:"lid" binding:"required"`
	CatId              primitive.ObjectID `json:"cat_id" bson:"cat_id" binding:"required"`
	Ava                string             `json:"ava" bson:"ava"`
	Name               string             `json:"name" bson:"name"  binding:"required"`
	ShortDesc          string             `json:"short_desc" bson:"short_desc" binding:"required"`
	FullDesc           string             `json:"full_desc" bson:"full_desc"  binding:"required"`
	Fee                float64            `json:"fee" bson:"fee" binding:"required"`
	Discount           float64            `json:"discount" bson:"discount" binding:"required"`
	ChapterCount       int                `json:"chapter_count" bson:"chapter_count" binding:"required"`
	IsDone             bool               `json:"is_done" bson:"is_done"`
	RegCount           int                `json:"reg_count" bson:"reg_count"`
}

type CourseChapter struct {
	field.DefaultField `bson:",inline"`
	CourseId           string `json:"cid" bson:"cid"`
	Title              string `json:"title" bson:"title"`
}

type Video struct {
	field.DefaultField `bson:",inline"`
	ChapterId          string `json:"chap_id" bson:"chap_id"`
	CourseId           string `json:"cid" bson:"cid"`
	Path               string `json:"path" bson:"path"`
	Title              string `json:"title" bson:"title"`
	Previewable        bool   `json:"previewable" bson:"previewable"`
}

type SimpleCourse struct {
	Id           string  `json:"id"`
	Title        string  `json:"title"`
	CategoryId   string  `json:"cid"`
	Category     string  `json:"category"`
	LecturerId   string  `json:"lid"`
	Lecturer     string  `json:"lecturer"`
	ReviewScore  float32 `json:"review_score"`
	Ava          string  `json:"ava"`
	CurrentPrice float64 `json:"current_price"`
}

func (Course) collName() string {
	return "courses"
}

func CreateCourseTextIndexModels() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{Key: "name", Value: bsonx.String("text")}},
		},
	}
}

func (c *Course) Save() error {
	_, err := db.Collection(c.collName()).InsertOne(ctx, c)
	return err
}

func (c *Course) FindById(oid primitive.ObjectID) error {
	return db.Collection(c.collName()).Find(ctx, bson.M{"_id": oid}).One(c)
}

func FindByLecId(lid primitive.ObjectID) ([]Course, error) {
	var res []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{"lid": lid}).All(res)
	return res, err
}

func (c *Course) FindByCatId(cid primitive.ObjectID) ([]Course, error) {
	var res []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{"cat_id": cid}).All(res)
	return res, err
}

func (c *Course) ConvertToSimpleCourse() (*SimpleCourse, error) {
	var simple SimpleCourse
	simple.Id = c.Id.String()
	simple.Title = c.Name
	simple.Ava = c.Ava
	simple.CurrentPrice = c.Fee * c.Discount
	simple.CategoryId = c.CatId.String()
	category := Category{}
	if err := category.FindCategoryById(c.CatId); err != nil {
		return nil, err
	}
	simple.Category = category.Name

	var lecturer User
	if err := lecturer.FindById(c.LecId); err != nil {
		return nil, err
	}
	simple.LecturerId = lecturer.Fullname

	var err error
	simple.ReviewScore, err = CalcAvgScore(c.CatId)
	if err != nil {
		return nil, err
	}

	return &simple, nil
}

// Hooks
func (c *Course) BeforeInsert() error {
	return nil
}

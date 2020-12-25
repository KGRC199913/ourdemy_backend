package models

import (
	"github.com/qiniu/qmgo/field"
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
	Discount           float32            `json:"discount" bson:"discount" binding:"required"`
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

//func NewCourse(lecId string, catId string,
//	ava string, name string,
//	shortDes string, fullDesc string,
//	fee float64,
//	chapterCount int) *Course {
//	return &Course{
//		LecId:        lecId,
//		CatId:        catId,
//		AvaUrl:       ava,
//		Name:         name,
//		ShortDesc:    shortDes,
//		FullDesc:     fullDesc,
//		Fee:          fee,
//		Discount:     0,
//		ChapterCount: chapterCount,
//		IsDone:       false,
//		RegCount:     0,
//	}
//}

func (c *Course) Save() error {
	_, err := db.Collection(c.collName()).InsertOne(ctx, c)
	return err
}

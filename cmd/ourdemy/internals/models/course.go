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
	Previewable        bool   `json:"previewable" bson:"previewable"`
}

type Video struct {
	field.DefaultField `bson:",inline"`
	ChapterId          string `json:"chap_id" bson:"chap_id"`
	CourseId           string `json:"cid" bson:"cid"`
	Path               string `json:"path" bson:"path"`
	Title              string `json:"title" bson:"title"`
}

type simpleCourse struct {
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

type fullCourse struct {
	Id              string          `json:"id"`
	Title           string          `json:"title"`
	CategoryId      string          `json:"cid"`
	Category        string          `json:"category"`
	LecturerId      string          `json:"lid"`
	Lecturer        string          `json:"lecturer"`
	ReviewScore     float32         `json:"review_score"`
	Ava             string          `json:"ava"`
	Fee             float64         `json:"fee" bson:"fee"`
	Discount        float64         `json:"discount" bson:"discount"`
	ShortDesc       string          `json:"short_desc" bson:"short_desc"`
	FullDesc        string          `json:"full_desc" bson:"full_desc"`
	IsDone          bool            `json:"is_done"`
	Chapters        []CourseChapter `json:"ch"`
	PreviewChapters []CourseChapter `json:"pch"`
}

func (Course) collName() string {
	return "courses"
}

func (CourseChapter) collName() string {
	return "course_chapters"
}

func (Video) collName() string {
	return "videos"
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

func (c *Course) ConvertToSimpleCourse() (*simpleCourse, error) {
	category := Category{}
	if err := category.FindCategoryById(c.CatId); err != nil {
		return nil, err
	}

	var lecturer User
	if err := lecturer.FindById(c.LecId); err != nil {
		return nil, err
	}

	var err error
	reviewScore, err := CalcAvgScore(c.CatId)
	if err != nil {
		return nil, err
	}

	return &simpleCourse{
		Id:           c.Id.String(),
		Title:        c.Name,
		CategoryId:   c.CatId.String(),
		Category:     category.Name,
		LecturerId:   c.LecId.String(),
		Lecturer:     lecturer.Fullname,
		ReviewScore:  reviewScore,
		Ava:          c.Ava,
		CurrentPrice: c.Fee * c.Discount,
	}, nil
}

func (c *Course) ConvertToFullCourse() (*fullCourse, error) {
	category := Category{}
	if err := category.FindCategoryById(c.CatId); err != nil {
		return nil, err
	}

	var lecturer User
	if err := lecturer.FindById(c.LecId); err != nil {
		return nil, err
	}

	var err error
	reviewScore, err := CalcAvgScore(c.CatId)
	if err != nil {
		return nil, err
	}

	chapters, err := getAllChapterByCourseId(c.Id)
	if err != nil {
		return nil, err
	}

	var previewableCc []CourseChapter
	for _, cc := range chapters {
		if cc.Previewable {
			previewableCc = append(previewableCc, cc)
		}
	}

	return &fullCourse{
		Id:              c.Id.String(),
		Title:           c.Name,
		CategoryId:      c.CatId.String(),
		Category:        category.Name,
		LecturerId:      c.LecId.String(),
		Lecturer:        lecturer.Fullname,
		ReviewScore:     reviewScore,
		Ava:             c.Ava,
		Fee:             c.Fee,
		Discount:        c.Discount,
		ShortDesc:       c.ShortDesc,
		FullDesc:        c.FullDesc,
		IsDone:          c.IsDone,
		Chapters:        chapters,
		PreviewChapters: previewableCc,
	}, nil
}

func getAllChapterByCourseId(cid primitive.ObjectID) (cc []CourseChapter, err error) {
	err = db.Collection(CourseChapter{}.collName()).Find(ctx, bson.M{"cid": cid}).All(&cc)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

// Hooks
func (c *Course) BeforeInsert() error {
	return nil
}

func (c *Course) AfterInsert() error {
	rgC := regCourse{
		CourseId: c.Id,
		JoinInfo: []courseJoinInfo{},
	}
	return rgC.Save()
}

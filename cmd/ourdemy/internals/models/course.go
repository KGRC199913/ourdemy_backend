package models

import "github.com/qiniu/qmgo/field"

type Course struct {
	field.DefaultField `bson:",inline"`
	LecId              string  `json:"lid" bson:"lid"`
	CatId              string  `json:"cat_id" bson:"cat_id"`
	AvaUrl             string  `json:"ava_url" bson:"ava_url"`
	ShortDesc          string  `json:"short_desc" bson:"short_desc"`
	FullDesc           string  `json:"full_desc" bson:"full_desc"`
	Fee                float64 `json:"fee" bson:"fee"`
	Discount           float32 `json:"discount" bson:"discount"`
	ChapterCount       int     `json:"chapter_count" bson:"chapter_count"`
	IsDone             bool    `json:"is_done" bson:"is_done"`
	RegCount           int     `json:"reg_count" bson:"reg_count"`
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

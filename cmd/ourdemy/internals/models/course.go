package models

import (
	"errors"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"os"
	"sort"
)

type Course struct {
	field.DefaultField `bson:",inline"`
	LecId              primitive.ObjectID `json:"lid" bson:"lid" binding:"required"`
	CatId              primitive.ObjectID `json:"cat_id" bson:"cat_id" binding:"required"`
	Ava                string             `json:"ava" bson:"ava"`
	Name               string             `json:"name" bson:"name"  binding:"required"`
	ShortDesc          string             `json:"short_desc" bson:"short_desc"`
	FullDesc           string             `json:"full_desc" bson:"full_desc"`
	Fee                float64            `json:"fee" bson:"fee" binding:"required"`
	Discount           float64            `json:"discount" bson:"discount" binding:"required"`
	ChapterCount       int                `json:"chapter_count" bson:"chapter_count"`
	IsDone             bool               `json:"is_done" bson:"is_done"`
	RegCount           int                `json:"reg_count" bson:"reg_count"`
	WatchCount         int                `json:"watch_count" bson:"watch_count"`
	Disabled           bool               `json:"disabled" bson:"disabled"`
}

type CourseChapter struct {
	field.DefaultField `bson:",inline" json:",inline"`
	CourseId           primitive.ObjectID `json:"cid" bson:"cid" binding:"required"`
	Title              string             `json:"title" bson:"title" binding:"required"`
	Previewable        *bool              `json:"previewable" bson:"previewable" binding:"required"`
	Videos             []VideoMetadata    `json:"videos" bson:"-" binding:"-"`
}

type VideoMetadata struct {
	field.DefaultField `bson:",inline" json:",inline"`
	ChapterId          primitive.ObjectID `json:"chap_id" bson:"chap_id"`
	CourseId           primitive.ObjectID `json:"cid" bson:"cid"`
	Title              string             `json:"title" bson:"title"`
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
	IsDiscount   bool    `json:"is_discount"`
	IsDone       bool    `json:"is_done"`
	ReviewCount  int64   `json:"review_count"`
}

type FullCourse struct {
	Id            string          `json:"id"`
	Title         string          `json:"title"`
	CategoryId    string          `json:"cid"`
	Category      string          `json:"category"`
	LecturerId    string          `json:"lid"`
	Lecturer      string          `json:"lecturer"`
	LecturerEmail string          `json:"lecturer_email"`
	ReviewScore   float32         `json:"review_score"`
	Ava           string          `json:"ava"`
	Fee           float64         `json:"fee" bson:"fee"`
	Discount      float64         `json:"discount" bson:"discount"`
	ShortDesc     string          `json:"short_desc" bson:"short_desc"`
	FullDesc      string          `json:"full_desc" bson:"full_desc"`
	IsDone        bool            `json:"is_done"`
	Chapters      []CourseChapter `json:"chapters"`
	ReviewCount   int64           `json:"review_count"`
	Disabled      bool            `json:"disabled"`
}

func (Course) collName() string {
	return "courses"
}

func (CourseChapter) collName() string {
	return "course_chapters"
}

func (VideoMetadata) collName() string {
	return "video_metadata"
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

func (cc *CourseChapter) Save() error {
	_, err := db.Collection(cc.collName()).InsertOne(ctx, cc)
	return err
}

func (vm *VideoMetadata) Save() error {
	_, err := db.Collection(vm.collName()).InsertOne(ctx, vm)
	return err
}

func (vm *VideoMetadata) Remove() error {
	return db.Collection(vm.collName()).Remove(ctx, bson.M{
		"_id": vm.Id,
	})
}

func (c *Course) FindById(oid primitive.ObjectID) error {
	return db.Collection(c.collName()).Find(ctx, bson.M{"_id": oid}).One(c)
}

func (cc *CourseChapter) FindById(ccid primitive.ObjectID) error {
	return db.Collection(cc.collName()).Find(ctx, bson.M{"_id": ccid}).One(cc)
}

func (vm *VideoMetadata) FindById(vmid primitive.ObjectID) error {
	return db.Collection(vm.collName()).Find(ctx, bson.M{"_id": vmid}).One(vm)
}

func FindAllChapterByCourseId(cid primitive.ObjectID) ([]CourseChapter, error) {
	var ccs []CourseChapter
	err := db.Collection(CourseChapter{}.collName()).Find(ctx, bson.M{"cid": cid}).All(&ccs)
	return ccs, err
}

func FindByLecId(lid primitive.ObjectID) ([]Course, error) {
	var res []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{"lid": lid}).All(&res)
	return res, err
}

func GetAllCourse(excludeDisabled bool) ([]Course, error) {
	var res []Course
	var err error
	if excludeDisabled {
		err = db.Collection(Course{}.collName()).Find(ctx, bson.M{
			"disabled": bson.M{
				"$ne": true,
			},
		}).All(&res)
	} else {
		err = db.Collection(Course{}.collName()).Find(ctx, bson.M{}).All(&res)
	}

	if res == nil {
		res = []Course{}
	}
	return res, err
}

func GetAllCourseAsSimple() ([]SimpleCourse, error) {
	c, err := GetAllCourse(true)

	if err != nil {
		return nil, err
	}

	if c == nil {
		c = []Course{}
	}

	var res []SimpleCourse
	for _, item := range c {
		f, err := item.ConvertToSimpleCourse()
		if err != nil {
			return nil, err
		}

		res = append(res, *f)
	}

	if res == nil {
		res = []SimpleCourse{}
	}

	return res, err
}

func GetAllCourseBySubcatIdAsSimple(subcatId primitive.ObjectID) ([]SimpleCourse, error) {
	var res []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{"cat_id": subcatId, "disabled": bson.M{
		"$ne": true,
	}}).All(&res)

	if res == nil {
		res = []Course{}
	}

	var returnRes []SimpleCourse
	for _, course := range res {
		f, err := course.ConvertToSimpleCourse()
		if err != nil {
			return nil, err
		}

		returnRes = append(returnRes, *f)
	}

	if returnRes == nil {
		returnRes = []SimpleCourse{}
	}

	return returnRes, err
}

func GetAllCourseByCatIdAsSimple(catId primitive.ObjectID) ([]SimpleCourse, error) {
	var res []Course
	subCats, err := FindByParentCategoryId(catId)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, subCat := range subCats {
		err := db.Collection(Course{}.collName()).Find(ctx,
			bson.M{"cat_id": subCat.Id,
				"disabled": bson.M{
					"$ne": true,
				}}).All(&courses)
		if err != nil {
			return nil, err
		}
		res = append(res, courses...)
	}

	if res == nil {
		res = []Course{}
	}

	var returnRes []SimpleCourse
	for _, course := range res {
		f, err := course.ConvertToSimpleCourse()
		if err != nil {
			return nil, err
		}

		returnRes = append(returnRes, *f)
	}

	if returnRes == nil {
		returnRes = []SimpleCourse{}
	}

	return returnRes, err
}

func GetAllCourseAsFull() ([]FullCourse, error) {
	c, err := GetAllCourse(false)

	if err != nil {
		return nil, err
	}

	if c == nil {
		c = []Course{}
	}

	var res []FullCourse
	for _, item := range c {
		f, err := item.ConvertToFullCourse()
		if err != nil {
			return nil, err
		}

		res = append(res, *f)
	}

	if res == nil {
		res = []FullCourse{}
	}

	return res, err
}

func GetTop10NewestCourse() ([]SimpleCourse, error) {
	var courses []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{
		"disabled": bson.M{
			"$ne": true,
		},
	}).Sort("-createAt").Limit(10).All(&courses)

	if courses == nil {
		courses = []Course{}
	}

	var res []SimpleCourse
	for _, course := range courses {
		simpleCourse, _ := course.ConvertToSimpleCourse()
		res = append(res, *simpleCourse)
	}

	if res == nil {
		res = []SimpleCourse{}
	}

	return res, err
}

func GetTop10MostWatchCourse() ([]SimpleCourse, error) {
	var courses []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{
		"disabled": bson.M{
			"$ne": true,
		},
	}).Sort("-watch_count").Limit(10).All(&courses)
	if courses == nil {
		courses = []Course{}
	}
	var res []SimpleCourse
	for _, course := range courses {
		simpleCourse, _ := course.ConvertToSimpleCourse()
		res = append(res, *simpleCourse)
	}

	if res == nil {
		res = []SimpleCourse{}
	}

	return res, err
}

func GetTop4HighlightCourse() ([]SimpleCourse, error) {
	var courses []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{
		"disabled": bson.M{
			"$ne": true,
		},
	}).All(&courses)

	var res []SimpleCourse
	for _, course := range courses {
		simpleCourse, _ := course.ConvertToSimpleCourse()
		res = append(res, *simpleCourse)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ReviewScore > res[j].ReviewScore
	})

	if res == nil {
		res = []SimpleCourse{}
	}

	end := 4
	if len(res) < 4 {
		end = len(res)
	}

	return res[0:end], err
}

func GetTop10MostRegisterCourse() ([]Course, error) {
	var courses []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{
		"disabled": bson.M{
			"$ne": true,
		},
	}).Sort("-reg_count").Limit(10).All(&courses)
	if courses == nil {
		courses = []Course{}
	}

	return courses, err
}

func Get5RandomCourseBySubcat(cid primitive.ObjectID) ([]Course, error) {
	var courses []Course
	err := db.Collection(Course{}.collName()).Aggregate(ctx, mongo.Pipeline{
		{{"$match", bson.D{{"cat_id", cid}}}},
		{{"$sample", bson.D{{"size", 5}}}},
	}).All(&courses)

	if courses == nil {
		courses = []Course{}
	}

	return courses, err
}

func (c *Course) UpdateCourseDesc(short string, full string) error {
	c.ShortDesc = short
	c.FullDesc = full
	return db.Collection(c.collName()).UpdateOne(ctx, bson.M{
		"_id": c.Id,
	}, bson.M{
		"$set": bson.M{
			"short_desc": short,
			"full_desc":  full,
		},
	})
}

func (c *Course) UpdateDisableStatus(disabled bool) error {
	c.Disabled = disabled
	return db.Collection(c.collName()).UpdateOne(ctx, bson.M{
		"_id": c.Id,
	}, bson.M{
		"$set": bson.M{
			"disabled": disabled,
		},
	})
}

func (c *Course) UpdateDiscount(discount float64) error {
	c.Discount = discount
	return db.Collection(c.collName()).UpdateOne(ctx, bson.M{
		"_id": c.Id,
	}, bson.M{
		"$set": bson.M{
			"discount": discount,
		},
	})
}

func (c *Course) UpdateCourseStatus(isDone bool) error {
	return db.Collection(c.collName()).UpdateOne(ctx, bson.M{
		"_id": c.Id,
	}, bson.M{
		"$set": bson.M{
			"is_done": isDone,
		},
	})
}

func (c *Course) UpdateChapterCount(count int) error {
	return db.Collection(c.collName()).UpdateOne(ctx, bson.M{
		"_id": c.Id,
	}, bson.M{
		"$set": bson.M{
			"chapter_count": count,
		},
	})
}

func (c *Course) UpdateWatchCount() error {
	return db.Collection(c.collName()).UpdateOne(ctx, bson.M{
		"_id": c.Id,
	}, bson.M{
		"$set": bson.M{
			"watch_count": c.WatchCount + 1,
		},
	})
}

func SearchByKeyword(keyword string, limit int64, offset int64) ([]Course, error) {
	var res []Course
	err := db.Collection(Course{}.collName()).Aggregate(ctx, mongo.Pipeline{
		bson.D{{"$match",
			[]bson.E{
				{"$text", bson.D{{"$search", keyword}}}}}},
	}).All(&res)

	if res == nil {
		res = []Course{}
	}

	var filteredRes []Course
	for _, course := range res {
		if !course.Disabled {
			filteredRes = append(filteredRes, course)
		}
	}

	if filteredRes == nil {
		filteredRes = []Course{}
	}

	return paginateCourse(filteredRes, offset, limit), err
}

func FindByCatId(cid primitive.ObjectID, limit int64, offset int64) ([]Course, error) {
	var res []Course
	subCats, err := FindByParentCategoryId(cid)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, subCat := range subCats {
		err := db.Collection(Course{}.collName()).Find(ctx, bson.M{"cat_id": subCat.Id}).Skip(offset).Limit(limit).All(&courses)
		if err != nil {
			return nil, err
		}
		res = append(res, courses...)
	}
	return res, err
}

func FindAllCourseByCatId(cid primitive.ObjectID) ([]Course, error) {
	var res []Course
	subCats, err := FindByParentCategoryId(cid)
	if err != nil {
		return nil, err
	}

	var courses []Course
	for _, subCat := range subCats {
		err := db.Collection(Course{}.collName()).Find(ctx, bson.M{"cat_id": subCat.Id}).All(&courses)
		if err != nil {
			return nil, err
		}
		res = append(res, courses...)
	}
	return res, err
}

func FindAllCourseBySubcatId(scid primitive.ObjectID) ([]Course, error) {
	var res []Course

	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{"cat_id": scid}).All(&res)
	return res, err
}

func FindBySubcatId(subcatId primitive.ObjectID, limit int64, offset int64) ([]Course, error) {
	var res []Course
	err := db.Collection(Course{}.collName()).Find(ctx, bson.M{"cat_id": subcatId}).Skip(offset).Limit(limit).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindAllVideoMetadataByChapterId(ccid primitive.ObjectID) ([]VideoMetadata, error) {
	var res []VideoMetadata
	err := db.Collection(VideoMetadata{}.collName()).Find(ctx, bson.M{"chap_id": ccid}).All(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Course) ConvertToSimpleCourse() (*SimpleCourse, error) {
	subcat := SubCategory{}
	if err := subcat.FindSubCategoryById(c.CatId); err != nil {
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

	rvc, err := CountReviewByCourseId(c.Id)

	return &SimpleCourse{
		Id:           c.Id.String(),
		Title:        c.Name,
		CategoryId:   c.CatId.String(),
		Category:     subcat.Name,
		LecturerId:   c.LecId.String(),
		Lecturer:     lecturer.Fullname,
		ReviewScore:  reviewScore,
		Ava:          c.Ava,
		CurrentPrice: c.Fee * (1.0 - c.Discount),
		IsDone:       c.IsDone,
		IsDiscount:   c.Discount < 0.0001,
		ReviewCount:  rvc,
	}, nil
}

func (c *Course) ConvertToFullCourse() (*FullCourse, error) {
	category := SubCategory{}
	if err := category.FindSubCategoryById(c.CatId); err != nil {
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

	if chapters == nil {
		chapters = []CourseChapter{}
	}

	for index, _ := range chapters {
		chapters[index].Videos, err = FindAllVideoMetadataByChapterId(chapters[index].Id)
		if err != nil {
			return nil, err
		}
		if chapters[index].Videos == nil {
			chapters[index].Videos = []VideoMetadata{}
		}
	}

	rvc, err := CountReviewByCourseId(c.Id)

	return &FullCourse{
		Id:            c.Id.String(),
		Title:         c.Name,
		CategoryId:    c.CatId.String(),
		Category:      category.Name,
		LecturerId:    c.LecId.String(),
		Lecturer:      lecturer.Fullname,
		LecturerEmail: lecturer.Email,
		ReviewScore:   reviewScore,
		Ava:           c.Ava,
		Fee:           c.Fee,
		Discount:      c.Discount,
		ShortDesc:     c.ShortDesc,
		FullDesc:      c.FullDesc,
		IsDone:        c.IsDone,
		Chapters:      chapters,
		ReviewCount:   rvc,
		Disabled:      c.Disabled,
	}, nil
}

func getAllChapterByCourseId(cid primitive.ObjectID) (cc []CourseChapter, err error) {
	err = db.Collection(CourseChapter{}.collName()).Find(ctx, bson.M{"cid": cid}).All(&cc)

	if err != nil {
		return nil, err
	}
	return cc, nil
}

func (c *Course) ForceRemove() error {
	chaps, err := FindAllChapterByCourseId(c.Id)
	if err != nil {
		return err
	}

	if chaps == nil {
		chaps = []CourseChapter{}
	}

	for _, chap := range chaps {
		vids, err := FindAllVideoMetadataByChapterId(chap.Id)
		if err != nil {
			return err
		}

		if vids == nil {
			vids = []VideoMetadata{}
		}

		for _, vid := range vids {
			if err := os.Remove("vid/" + vid.Id.Hex()); err != nil {
				return err
			}

			if err := vid.Remove(); err != nil {
				return err
			}
		}

		if err = chap.Remove(); err != nil {
			return err
		}
	}

	return c.Remove()
}

func (c *Course) Remove() error {
	return db.Collection(c.collName()).Remove(ctx, bson.M{
		"_id": c.Id,
	})
}

func (vm *VideoMetadata) UpdateVideoTitle(title string) error {
	vm.Title = title
	return db.Collection(vm.collName()).UpdateOne(ctx, bson.M{
		"_id": vm.Id,
	}, bson.M{
		"$set": bson.M{
			"title": vm.Title,
		},
	})
}

func (cc *CourseChapter) Remove() error {
	return db.Collection(cc.collName()).Remove(ctx, bson.M{
		"_id": cc.Id,
	})
}

// Hooks c
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

// Hooks cc
func (cc *CourseChapter) BeforeInsert() error {
	var course Course
	if err := course.FindById(cc.CourseId); err != nil {
		return errors.New("course not exist")
	}

	if err := course.UpdateChapterCount(course.ChapterCount + 1); err != nil {
		return errors.New("update chapter count failed: " + err.Error())
	}

	return nil
}

func (cc *CourseChapter) AfterRemove() error {
	var course Course
	if err := course.FindById(cc.CourseId); err != nil {
		return errors.New("course not exist")
	}

	if err := course.UpdateChapterCount(course.ChapterCount - 1); err != nil {
		return errors.New("update chapter count failed: " + err.Error())
	}

	return nil
}

func paginateCourse(x []Course, offset int64, limit int64) []Course {
	if offset > int64(len(x)) {
		offset = int64(len(x))
	}

	end := offset + limit
	if end > int64(len(x)) {
		end = int64(len(x))
	}

	return x[offset:end]
}

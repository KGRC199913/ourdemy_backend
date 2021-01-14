package models

import (
	"errors"
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type regCourse struct {
	field.DefaultField `bson:",inline"`
	CourseId           primitive.ObjectID `json:"cid" bson:"cid"`
	JoinInfo           []courseJoinInfo   `json:"join_info" bson:"join_info"`
}

type courseJoinInfo struct {
	JoinDate time.Time          `json:"join_date" bson:"join_date"`
	UserId   primitive.ObjectID `json:"uid" bson:"uid"`
}

type XSimpleCourse struct {
	CourseId   primitive.ObjectID `json:"cid"`
	CourseName string             `json:"course_name"`
}

func (regCourse) collName() string {
	return "reg_course"
}

func (rgC *regCourse) Save() error {
	_, err := db.Collection(rgC.collName()).InsertOne(ctx, rgC)
	return err
}

func AddUserToCourseInfo(uid primitive.ObjectID, cid primitive.ObjectID) error {
	var rgC regCourse
	if err := rgC.FindByCourseId(cid); err != nil {
		return errors.New("reg course info not found")
	}

	index := rgcIndexOfUid(uid, rgC.JoinInfo)
	if index != -1 {
		return errors.New("already joined")
	}

	rgC.JoinInfo = append(rgC.JoinInfo, courseJoinInfo{
		JoinDate: time.Now(),
		UserId:   uid,
	})

	return db.Collection(rgC.collName()).UpdateOne(ctx, bson.M{
		"_id": rgC.Id,
	}, bson.M{
		"$set": bson.M{
			"join_info": rgC.JoinInfo,
		},
	})
}

func (rgC *regCourse) RemoveUserFromCourseInfo(uid primitive.ObjectID) error {
	index := rgcIndexOfUid(uid, rgC.JoinInfo)
	rgC.JoinInfo = rgcRemoveUidFromIndex(rgC.JoinInfo, index)

	return db.Collection(rgC.collName()).UpdateOne(ctx, bson.M{
		"_id": rgC.Id,
	}, bson.M{
		"$set": bson.M{
			"join_info": rgC.JoinInfo,
		},
	})
}

func (rgC *regCourse) FindByCourseId(cid primitive.ObjectID) error {
	return db.Collection(rgC.collName()).Find(ctx, bson.M{"cid": cid}).One(rgC)
}

func IsUserJoined(cid primitive.ObjectID, uid primitive.ObjectID) bool {
	var rgC regCourse
	if err := rgC.FindByCourseId(cid); err != nil {
		return false
	}

	for _, joinInfo := range rgC.JoinInfo {
		if joinInfo.UserId == uid {
			return true
		}
	}

	return false
}

func GetRegByUid(uid primitive.ObjectID) ([]XSimpleCourse, error) {
	var data []regCourse
	err := db.Collection(regCourse{}.collName()).Find(ctx, bson.M{}).All(&data)
	if err != nil {
		return []XSimpleCourse{}, err
	}

	var res []XSimpleCourse
	for _, r := range data {
		for _, join := range r.JoinInfo {
			if join.UserId == uid {
				var course Course
				if err := course.FindById(r.CourseId); err != nil {
					return nil, err
				}

				res = append(res, XSimpleCourse{
					CourseId:   r.CourseId,
					CourseName: course.Name,
				})

				break
			}
		}
	}

	if res == nil {
		res = []XSimpleCourse{}
	}

	return res, nil
}

// Hooks
func (rgC *regCourse) BeforeInsert() error {
	rgC.JoinInfo = rgcUnique(rgC.JoinInfo)
	return nil
}

func rgcUnique(data []courseJoinInfo) []courseJoinInfo {
	keys := make(map[primitive.ObjectID]bool)
	var list []courseJoinInfo
	for _, entry := range data {
		if _, value := keys[entry.UserId]; !value {
			keys[entry.UserId] = true
			list = append(list, entry)
		}
	}
	return list
}

func rgcRemoveUidFromIndex(s []courseJoinInfo, index int) []courseJoinInfo {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func rgcIndexOfUid(uid primitive.ObjectID, data []courseJoinInfo) (index int) {
	for index, val := range data {
		if val.UserId == uid {
			return index
		}
	}
	return -1 //not found.
}

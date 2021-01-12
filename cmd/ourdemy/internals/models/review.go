package models

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"time"
)

type Review struct {
	field.DefaultField `bson:",inline"`
	UserId             primitive.ObjectID `json:"uid" bson:"uid"`
	CourseId           primitive.ObjectID `json:"cid" bson:"cid"`
	Content            string             `json:"content" bson:"content"`
	Score              float32            `json:"score" bson:"score"`
}

type DisplayableReview struct {
	Id       primitive.ObjectID `json:"id"`
	Content  string             `json:"content"`
	Score    float32            `json:"score"`
	Username string             `json:"username"`
	Time     time.Time          `json:"time"`
}

func (Review) collName() string {
	return "reviews"
}

func (r *Review) Save() error {
	_, err := db.Collection(r.collName()).InsertOne(ctx, r)
	return err
}

func (r *Review) FindById(oid primitive.ObjectID) error {
	return db.Collection(r.collName()).Find(ctx, bson.M{"_id": oid}).One(r)
}

func (r *Review) FindByUid(uid primitive.ObjectID) error {
	return db.Collection(r.collName()).Find(ctx, bson.M{"uid": uid}).One(r)
}

func (r *Review) FindByUidAndCid(uid primitive.ObjectID, cid primitive.ObjectID) error {
	return db.Collection(r.collName()).Find(ctx, bson.M{"uid": uid, "cid": cid}).One(r)
}

func FindByCourseId(cid primitive.ObjectID) ([]Review, error) {
	var revs []Review
	err := db.Collection(Review{}.collName()).Find(ctx, bson.M{"cid": cid}).All(&revs)
	if err != nil {
		return nil, err
	}

	if revs == nil {
		revs = []Review{}
	}

	return revs, nil
}

func (rv *Review) ConvertToDisplayableReview() (*DisplayableReview, error) {
	var u User
	if err := u.FindById(rv.UserId); err != nil {
		return nil, err
	}

	return &DisplayableReview{
		Id:       rv.Id,
		Content:  rv.Content,
		Score:    rv.Score,
		Username: u.Username,
		Time:     rv.CreateAt,
	}, nil
}

func CalcAvgScore(cid primitive.ObjectID) (float32, error) {
	reviews, err := FindByCourseId(cid)
	if err != nil {
		return 0.0, err
	}

	var totalScore float32
	for _, review := range reviews {
		totalScore += review.Score
	}
	avg := totalScore / float32(len(reviews))
	if math.IsNaN(float64(avg)) {
		avg = 5.0
	}
	return avg, nil
}

func (r *Review) UpdateReview(newContent string, newScore float32) error {
	return db.Collection(r.collName()).UpdateOne(ctx, bson.M{
		"_id": r.Id,
	}, bson.M{
		"$set": bson.M{
			"content": newContent,
			"score":   newScore,
		},
	})
}

func (r *Review) DeleteReview(rid primitive.ObjectID) error {
	return db.Collection(r.collName()).Remove(ctx, bson.M{
		"_id": rid,
	})
}

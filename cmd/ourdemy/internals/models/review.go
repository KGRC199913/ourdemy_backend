package models

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	field.DefaultField `bson:",inline"`
	UserId             primitive.ObjectID `json:"uid" bson:"uid"`
	CourseId           primitive.ObjectID `json:"cid" bson:"cid"`
	Content            string             `json:"content" bson:"content"`
	Score              float32            `json:"score" bson:"score"`
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

func FindByCourseId(cid primitive.ObjectID) (revs []Review, err error) {
	err = db.Collection(Review{}.collName()).Find(ctx, bson.M{"cid": cid}).All(&revs)
	if err != nil {
		return nil, err
	}
	return revs, nil
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

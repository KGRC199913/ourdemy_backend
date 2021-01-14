package models

import (
	"github.com/qiniu/qmgo/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WatchList struct {
	field.DefaultField `bson:",inline"`
	UserId             primitive.ObjectID   `json:"uid" bson:"uid"`
	CoursesId          []primitive.ObjectID `json:"cids" bson:"cids"`
}

func (WatchList) collName() string {
	return "watch_list"
}

func (wl *WatchList) Save() error {
	_, err := db.Collection(wl.collName()).InsertOne(ctx, wl)
	return err
}

func (wl *WatchList) FindByUid(uid primitive.ObjectID) error {
	return db.Collection(wl.collName()).Find(ctx, bson.M{"uid": uid}).One(wl)
}

func appendIfMissing(slice []primitive.ObjectID, i primitive.ObjectID) []primitive.ObjectID {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func (wl *WatchList) AddCourseToWatchList(cid primitive.ObjectID) error {
	wl.CoursesId = appendIfMissing(wl.CoursesId, cid)

	return db.Collection(wl.collName()).UpdateOne(ctx, bson.M{
		"_id": wl.Id,
	}, bson.M{
		"$set": bson.M{
			"cids": wl.CoursesId,
		},
	})
}

func (wl *WatchList) RemoveCourseFromWatchList(cid primitive.ObjectID) error {
	index := wlIndexOfCid(cid, wl.CoursesId)
	wl.CoursesId = wlRemoveCidFromIndex(wl.CoursesId, index)

	return db.Collection(wl.collName()).UpdateOne(ctx, bson.M{
		"_id": wl.Id,
	}, bson.M{
		"$set": bson.M{
			"cids": wl.CoursesId,
		},
	})
}

// Hooks
func (wl *WatchList) BeforeInsert() error {
	wl.CoursesId = wlUnique(wl.CoursesId)
	return nil
}

func wlUnique(data []primitive.ObjectID) []primitive.ObjectID {
	keys := make(map[primitive.ObjectID]bool)
	var list []primitive.ObjectID
	for _, entry := range data {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func wlRemoveCidFromIndex(s []primitive.ObjectID, index int) []primitive.ObjectID {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func wlIndexOfCid(cid primitive.ObjectID, data []primitive.ObjectID) (index int) {
	for index, val := range data {
		if val == cid {
			return index
		}
	}
	return -1 //not found.
}

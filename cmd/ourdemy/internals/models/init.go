package models

import (
	"context"
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client *qmgo.Client
var db *qmgo.Database
var ctx context.Context

func InitDb(config *ultis.Config) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", config.DbUsername, config.DbPassword, config.DbUrl)
	ctx = context.Background()

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	cli, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	indexes := CreateCourseTextIndexModels()
	_, err = cli.Database(config.DbName).Collection(Course{}.collName()).Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}
	err = cli.Disconnect(ctx)
	if err != nil {
		return err
	}

	client, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: uri})
	if err != nil {
		return err
	}
	db = client.Database(config.DbName)
	if err != nil {
		return err
	}
	return nil
}

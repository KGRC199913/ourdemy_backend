package models

import (
	"context"
	"fmt"
	"github.com/KGRC199913/ourdemy_backend/cmd/ourdemy/internals/ultis"
	"github.com/qiniu/qmgo"
)

var client *qmgo.Client
var db *qmgo.Database
var ctx context.Context

func InitDb(config *ultis.Config) error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s", config.DbUsername, config.DbPassword, config.DbUrl)
	ctx = context.Background()

	var err error
	client, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: uri})
	if err != nil {
		return err
	}
	db = client.Database("ourdemy")
	return nil
}

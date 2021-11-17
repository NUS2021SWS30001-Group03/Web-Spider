package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sws3001spider/utils"
	"time"
)

var UserCollection *mongo.Collection

func init() {
	mongoURI:=utils.MongoURL
	fmt.Println("connection string is:", mongoURI)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*100)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	UserCollection = client.Database("admin").Collection("researcher")
}

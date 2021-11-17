package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type CoAuthors struct {
	CoAuthorName string `json:"name"bson:"name"`
	Weight       int    `json:"weight"bson:"weight"`
}

type Researcher struct {
	Name         string      `bson:"name"`
	Affiliation  string      `bson:"affiliation"`
	AuthorID     string      `bson:"author_id"`
	CoAuthorList []CoAuthors `bson:"coauthors"`
}

func CheckExist(author_id string) bool {
	filter:=bson.M{"author_id":author_id}
	var result Researcher
	err:=UserCollection.FindOne(context.TODO(),filter).Decode(&result)
	if err != nil{
		return false
	}
	return true
}

func InsertNode(request Researcher) error {
	_, err := UserCollection.InsertOne(context.TODO(), request)
	fmt.Println(request)
	return err
}

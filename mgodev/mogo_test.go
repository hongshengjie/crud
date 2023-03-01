package mgo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMgo(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("example").Collection("user")
	//DeleteUser(coll)
	//UpdateUser(coll)
	//FindUser(coll)
	Insert(coll)
}

func FindUser(coll *mongo.Collection) {
	//id, _ := primitive.ObjectIDFromHex("63ff2f14983bef62a8c881c0")
	u, err := Find(coll).Filter(NameEq("aa"), AgeLT(5)).All(context.Background())
	fmt.Println(u, err)
}

func UpdateUser(coll *mongo.Collection) {
	id, _ := primitive.ObjectIDFromHex("63ff2f14983bef62a8c881c0")
	Update(coll).SetName("woqu").SetAge(100).ByID(context.Background(), id)
}

func DeleteUser(coll *mongo.Collection) {
	id, _ := primitive.ObjectIDFromHex("63ff2f14983bef62a8c881c1")
	Delete(coll).ByID(context.Background(), id)
}

func Insert(coll *mongo.Collection) {

	var list []*User
	for i := 0; i < 10; i++ {
		u := &User{
			Name:  "aa",
			Age:   i,
			Sex:   false,
			Mtime: time.Now(),
		}
		list = append(list, u)
	}

	err := Create(coll).SetUsers(list...).Save(context.TODO())
	b, _ := json.Marshal(list)
	fmt.Println(err, string(b))
}

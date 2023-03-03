package mgo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hongshengjie/crud/mgo"
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
	FindUser(coll)
	//Insert(coll)
}

func FindUser(coll *mongo.Collection) {
	//id, _ := primitive.ObjectIDFromHex("63ff2f14983bef62a8c881c0")
	q := mgo.And(mgo.In(Name, "aa"))
	qq, _ := json.Marshal(q.Query())
	fmt.Println(string(qq))
	u, err := Find(coll).Filter(q.Query()...).All(context.Background())
	b, _ := json.Marshal(u)
	fmt.Println(string(b), err)
}

func UpdateUser(coll *mongo.Collection) {
	//id, _ := primitive.ObjectIDFromHex("63ff2f14983bef62a8c881c0")
	//Update(coll).SetName("woqu").SetAge(100).ByID(context.Background(), id)
}

func DeleteUser(coll *mongo.Collection) {
	//id, _ := primitive.ObjectIDFromHex("63ff2f14983bef62a8c881c1")
	//Delete(coll).ByID(context.Background(), id)
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

func TestPredicatt(t *testing.T) {
	//d := NEQ(Name, "dxxx")
	//x := In(Age, []int{1, 2, 3})
	//p := Nor(d, x).Query()
	//b, _ := json.Marshal(p)
	//fmt.Println(string(b))
}

func TestParse(t *testing.T) {
	//temp, _ := os.ReadFile("../templates/builder_mgo.tmpl")
	//r, _ := template.New("").Parse(string(temp))
	//x := mgo.ParseMongoStruct("./user.go", "User")

	// f, err := os.Create("user.go")
	// if err != nil {
	// 	panic(err)
	// }
	// r.Execute(f, x)
	// f.Close()

}

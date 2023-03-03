package mgo

import (
	"context"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string
	Age   int
	Sex   bool
	Mtime time.Time
}

const (
	tableName = "user"
	ID        = "_id"
	Name      = "name"
	Age       = "age"
	Sex       = "sex"
	Mtime     = "mtime"
)

func Collection(db *mongo.Database) *mongo.Collection {
	return db.Collection(tableName)
}

type FinderBuilder struct {
	col     *mongo.Collection
	filters primitive.D
	opts    *options.FindOptions
}

func Find(col *mongo.Collection) *FinderBuilder {
	return &FinderBuilder{col: col, opts: options.Find()}
}

func (f *FinderBuilder) Filter(filter ...primitive.E) *FinderBuilder {
	f.filters = append(f.filters, filter...)
	return f
}
func (f *FinderBuilder) Limit(l int64) *FinderBuilder {
	f.opts.SetLimit(l)
	return f
}

func (f *FinderBuilder) Sort(field string, desc bool) *FinderBuilder {
	i := 1
	if desc {
		i = -1
	}
	f.opts.SetSort(primitive.E{Key: field, Value: i})
	return f
}

func (f *FinderBuilder) Skip(s int64) *FinderBuilder {
	f.opts.SetSkip(s)
	return f
}
func (f *FinderBuilder) One(ctx context.Context) (*User, error) {
	f.opts = f.opts.SetLimit(1)
	ret, err := f.All(ctx)
	if err != nil {
		return nil, err
	}
	if len(ret) == 1 {
		return ret[0], nil
	}
	return nil, mongo.ErrNoDocuments
}

func (f *FinderBuilder) All(ctx context.Context) ([]*User, error) {
	cursor, err := f.col.Find(ctx, f.filters, f.opts)
	if err != nil {
		return nil, err
	}
	var results []*User
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

type InsertBuilder struct {
	col *mongo.Collection
	a   []interface{}
}

func Create(col *mongo.Collection) *InsertBuilder {
	return &InsertBuilder{col: col}
}
func (i *InsertBuilder) SetUsers(u ...*User) *InsertBuilder {
	for _, v := range u {
		i.a = append(i.a, v)
	}
	return i
}

func (i *InsertBuilder) Save(ctx context.Context) error {
	ret, err := i.col.InsertMany(ctx, i.a)
	if err != nil {
		return err
	}
	for idx, v := range ret.InsertedIDs {
		item := i.a[idx].(*User)
		id := v.(primitive.ObjectID)
		item.ID = id
	}
	return nil
}

type UpdateBuilder struct {
	col *mongo.Collection
	a   primitive.D
}

func Update(col *mongo.Collection) *UpdateBuilder {
	return &UpdateBuilder{col: col}
}
func (u *UpdateBuilder) SetID(a primitive.ObjectID) *UpdateBuilder {
	u.a = append(u.a, primitive.E{
		Key:   ID,
		Value: a,
	})
	return u
}
func (u *UpdateBuilder) SetName(a string) *UpdateBuilder {
	u.a = append(u.a, primitive.E{
		Key:   Name,
		Value: a,
	})
	return u
}
func (u *UpdateBuilder) SetAge(a int) *UpdateBuilder {
	u.a = append(u.a, primitive.E{
		Key:   Age,
		Value: a,
	})
	return u
}
func (u *UpdateBuilder) SetSex(a bool) *UpdateBuilder {
	u.a = append(u.a, primitive.E{
		Key:   Sex,
		Value: a,
	})
	return u
}
func (u *UpdateBuilder) SetMtime(a time.Time) *UpdateBuilder {
	u.a = append(u.a, primitive.E{
		Key:   Mtime,
		Value: a,
	})
	return u
}

func (u *UpdateBuilder) ByID(ctx context.Context, a primitive.ObjectID) (int64, error) {
	ret, err := u.col.UpdateByID(ctx, a, primitive.D{primitive.E{Key: "$set", Value: u.a}})
	if err != nil {
		return 0, err
	}
	return ret.ModifiedCount, nil
}

type DeleteBuilder struct {
	col *mongo.Collection
}

func Delete(col *mongo.Collection) *DeleteBuilder {
	return &DeleteBuilder{col: col}
}

func (d *DeleteBuilder) ByID(ctx context.Context, a primitive.ObjectID) (int64, error) {
	ret, err := d.col.DeleteOne(ctx, primitive.D{primitive.E{Key: "_id", Value: a}})
	if err != nil {
		return 0, err
	}
	return ret.DeletedCount, nil
}

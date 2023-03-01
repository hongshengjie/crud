package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Builder struct {
	filter bson.D
}
type Predicate struct {
	Builder
	fns []func(*Builder)
}

func (p *Predicate) EQ(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key:   field,
			Value: value,
		})
	})
	return p
}
func (p *Predicate) GT(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: "$gt", Value: value},
			},
		})
	})
	return p
}
func (p *Predicate) GTE(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: "$gte", Value: value},
			},
		})
	})
	return p
}
func (p *Predicate) LT(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: "$lt", Value: value},
			},
		})
	})
	return p
}
func (p *Predicate) LTE(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: "$lte", Value: value},
			},
		})
	})
	return p
}
func (p *Predicate) In(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: "$in", Value: value},
			},
		})
	})
	return p
}

func (p *Predicate) NE(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: "$ne", Value: value},
			},
		})
	})
	return p
}

func (p *Predicate) NIn(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: "$nin", Value: value},
			},
		})
	})
	return p
}

func P(fns ...func(*Builder)) *Predicate {
	return &Predicate{fns: fns}
}
func EQ(field string, value interface{}) *Predicate {
	return P().EQ(field, value)
}
func LT()
func LTE()

func _()

package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
)

// An Op represents a predicate operator.
type Op int

// Predicate operators
const (
	// Comparision Query Operators
	OpEQ Op = iota
	OpNEQ
	OpGT
	OpGTE
	OpLT
	OpLTE
	OpIn
	OpNotIn
	// Logical Query Operators
	OpAnd
	OpNot
	OpNor
	OpOr
	// Element Query Operators
	OpExists
	OpType

	// Field Update Operators
	OpSet
	OpInc

	// Array Update Operatos
	OpPop
	OpPull
	OpPush
)

var ops = [...]string{
	// Comparision Query Operators
	OpEQ:    "$eq",
	OpNEQ:   "$ne",
	OpGT:    "$gt",
	OpGTE:   "$gte",
	OpLT:    "$lt",
	OpLTE:   "$lte",
	OpIn:    "$in",
	OpNotIn: "$nin",
	// Logical Query Operators
	OpAnd: "$and",
	OpNot: "$not",
	OpNor: "$nor",
	OpOr:  "$or",
	// Element Query Operators
	OpExists: "$exists",
	OpType:   "$type",
	// Field Update Operators
	OpSet: "$set",
	OpInc: "$inc",

	// Array Update Operatos
	OpPop:  "$pop",
	OpPull: "$pull",
	OpPush: "$push",
}

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
			Key: field,
			Value: bson.D{
				{Key: ops[OpEQ], Value: value},
			},
		})
	})
	return p
}
func (p *Predicate) GT(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: ops[OpGT], Value: value},
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
				{Key: ops[OpGTE], Value: value},
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
				{Key: ops[OpLT], Value: value},
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
				{Key: ops[OpLTE], Value: value},
			},
		})
	})
	return p
}

func (p *Predicate) In(field string, value ...interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: ops[OpIn], Value: bson.A(value)},
			},
		})
	})
	return p
}

func (p *Predicate) NEQ(field string, value interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: ops[OpNEQ], Value: value},
			},
		})
	})
	return p
}

func (p *Predicate) NotIn(field string, value ...interface{}) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: ops[OpNotIn], Value: bson.A(value)},
			},
		})
	})
	return p
}

func (p *Predicate) Not(field string, p2 *Predicate) *Predicate {
	p.fns = append(p.fns, func(b *Builder) {
		b.filter = append(b.filter, bson.E{
			Key: field,
			Value: bson.D{
				{Key: ops[OpNot], Value: p2.Query()},
			},
		})
	})
	return p
}

func (p *Predicate) Query() bson.D {
	for _, f := range p.fns {
		f(&p.Builder)
	}
	return p.filter
}

func Or(pp ...*Predicate) *Predicate {
	return logic(OpOr, pp...)
}
func Nor(pp ...*Predicate) *Predicate {
	return logic(OpNor, pp...)
}
func And(pp ...*Predicate) *Predicate {
	return logic(OpAnd, pp...)
}

func logic(op Op, pp ...*Predicate) *Predicate {
	p := P()
	p.fns = append(p.fns, func(b *Builder) {
		var conds bson.A
		for _, v := range pp {
			c := v.Query()
			conds = append(conds, c)

		}
		item := bson.D{{
			Key:   ops[op],
			Value: conds,
		},
		}
		b.filter = append(b.filter, item...)
	})
	return p
}

func P(fns ...func(*Builder)) *Predicate {
	return &Predicate{fns: fns}
}
func EQ(field string, value interface{}) *Predicate {
	return P().EQ(field, value)
}
func NEQ(field string, value interface{}) *Predicate {
	return P().NEQ(field, value)
}

func LT(field string, value interface{}) *Predicate {
	return P().LT(field, value)
}
func LTE(field string, value interface{}) *Predicate {
	return P().LTE(field, value)
}
func GT(field string, value interface{}) *Predicate {
	return P().GT(field, value)
}
func GTE(field string, value interface{}) *Predicate {
	return P().GTE(field, value)
}
func In(field string, value ...interface{}) *Predicate {
	return P().In(field, value...)
}

func NotIn(field string, value ...interface{}) *Predicate {
	return P().NotIn(field, value...)
}

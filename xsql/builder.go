// Package xsql provides wrappers around the standard database/sql package
// to allow the generated code to interact with a statically-typed API.
package xsql

import (
	"context"
	"database/sql/driver"
	"fmt"

	"strconv"
	"strings"
)

// Querier wraps the basic Query method that is implemented
// by the different builders in this file.
type Querier interface {
	// Query returns the query representation of the element
	// and its arguments (if any).
	Query() (string, []interface{})
}

// InsertBuilder is a builder for `INSERT INTO` statement.
type InsertBuilder struct {
	Builder
	table    string
	schema   string
	columns  []string
	defaults string
	values   [][]interface{}
	// OnDuplicateKeyUpdate Expr
	updateExpr []string
	updateArgs []interface{}
}

// Insert creates a builder for the `INSERT INTO` statement.
//
//	Insert("users").
//		Columns("name", "age").
//		Values("a8m", 10).
//		Values("foo", 20)
//
// Note: Insert inserts all values in one batch.
func Insert(table string) *InsertBuilder { return &InsertBuilder{table: table} }

// Schema sets the database name for the insert table.
func (i *InsertBuilder) Schema(name string) *InsertBuilder {
	i.schema = name
	return i
}

// Table sets the table name for the insert table.
func (i *InsertBuilder) Table(name string) *InsertBuilder {
	i.table = name
	return i
}

// Set is a syntactic sugar API for inserting only one row.
func (i *InsertBuilder) Set(column string, v interface{}) *InsertBuilder {
	i.columns = append(i.columns, column)
	if len(i.values) == 0 {
		i.values = append(i.values, []interface{}{v})
	} else {
		i.values[0] = append(i.values[0], v)
	}
	return i
}

// Columns sets the columns of the insert statement.
func (i *InsertBuilder) Columns(columns ...string) *InsertBuilder {
	i.columns = append(i.columns, columns...)
	return i
}

// Values append a value tuple for the insert statement.
func (i *InsertBuilder) Values(values ...interface{}) *InsertBuilder {
	i.values = append(i.values, values)
	return i
}

// OnDuplicateKeyUpdate UpdateColumns generate like  "ON DUPLICATE KEY UPDATE `id` = VALUES (`id`) " sql statement
func (i *InsertBuilder) OnDuplicateKeyUpdate(columns ...string) *InsertBuilder {
	if len(columns) > 0 {
		b := &Builder{}
		for j, v := range columns {
			if j > 0 {
				b.Comma()
			}
			// `id` = VALUES (`id`)
			b.Ident(v).WriteOp(OpEQ).WriteString("VALUES ")
			b.Nested(func(bb *Builder) {
				bb.Ident(v)
			})
		}
		i.updateExpr = append(i.updateExpr, b.String())
	}
	return i
}

// OnDuplicateKeyUpdateExpr OnDuplicateKeyUpdateExpr generate like "ON DUPLICATE KEY UPDATE `c` = VALUES(`a`) + Values(`b`)"
func (i *InsertBuilder) OnDuplicateKeyUpdateExpr(q ...Querier) *InsertBuilder {
	for _, v := range q {
		expr, args := v.Query()
		i.updateExpr = append(i.updateExpr, expr)
		i.updateArgs = append(i.updateArgs, args...)
	}
	return i
}

// Query returns query representation of an `INSERT INTO` statement.
func (i *InsertBuilder) Query() (string, []interface{}) {
	i.WriteString("INSERT INTO ")
	i.writeSchema(i.schema)
	i.Ident(i.table).Pad()
	if i.defaults != "" && len(i.columns) == 0 {
		i.WriteString(i.defaults)
	} else {
		i.Nested(func(b *Builder) {
			b.IdentComma(i.columns...)
		})
		i.WriteString(" VALUES ")
		for j, v := range i.values {
			if j > 0 {
				i.Comma()
			}
			i.Nested(func(b *Builder) {
				b.Args(v...)
			})
		}
		if len(i.updateExpr) > 0 {
			i.WriteString(" ON DUPLICATE KEY UPDATE ")
			if len(i.updateExpr) > 0 {
				for j, v := range i.updateExpr {
					if j > 0 {
						i.Comma()
					}
					i.WriteString(v)
				}
				i.args = append(i.args, i.updateArgs...)
			}

		}

	}

	statement, args := i.String(), i.args
	return statement, args
}

// UpdateBuilder is a builder for `UPDATE` statement.
type UpdateBuilder struct {
	Builder
	table   string
	schema  string
	where   *Predicate
	nulls   []string
	columns []string
	values  []interface{}
}

// Update creates a builder for the `UPDATE` statement.
//
//	Update("users").Set("name", "foo").Set("age", 10)
func Update(table string) *UpdateBuilder { return &UpdateBuilder{table: table} }

// Schema sets the database name for the updated table.
func (u *UpdateBuilder) Schema(name string) *UpdateBuilder {
	u.schema = name
	return u
}

// Table sets the table name for the updated table.
func (u *UpdateBuilder) Table(name string) *UpdateBuilder {
	u.table = name
	return u
}

// Set sets a column and a its value.
func (u *UpdateBuilder) Set(column string, v interface{}) *UpdateBuilder {
	u.columns = append(u.columns, column)
	u.values = append(u.values, v)
	return u
}

// Add adds a numeric value to the given column.
func (u *UpdateBuilder) Add(column string, v interface{}) *UpdateBuilder {
	u.columns = append(u.columns, column)
	u.values = append(u.values, P().Append(func(b *Builder) {
		b.WriteString("COALESCE")
		b.Nested(func(b *Builder) {
			b.Ident(column).Comma().Arg(0)
		})
		b.WriteString(" + ")
		b.Arg(v)
	}))
	return u
}

// SetNull sets a column as null value.
func (u *UpdateBuilder) SetNull(column string) *UpdateBuilder {
	u.nulls = append(u.nulls, column)
	return u
}

// Where adds a where predicate for update statement.
func (u *UpdateBuilder) Where(p *Predicate) *UpdateBuilder {
	if u.where != nil {
		u.where = And(u.where, p)
	} else {
		u.where = p
	}
	return u
}

// FromSelect makes it possible to update entities that match the sub-query.
func (u *UpdateBuilder) FromSelect(s *Selector) *UpdateBuilder {
	u.Where(s.where)
	if table, _ := s.from.(*SelectTable); table != nil {
		u.table = table.name
	}
	return u
}

// Empty reports whether this builder does not contain update changes.
func (u *UpdateBuilder) Empty() bool {
	return len(u.columns) == 0 && len(u.nulls) == 0
}

// Query returns query representation of an `UPDATE` statement.
func (u *UpdateBuilder) Query() (string, []interface{}) {
	u.WriteString("UPDATE ")
	u.writeSchema(u.schema)
	u.Ident(u.table).WriteString(" SET ")
	for i, c := range u.nulls {
		if i > 0 {
			u.Comma()
		}
		u.Ident(c).WriteString(" = NULL")
	}
	if len(u.nulls) > 0 && len(u.columns) > 0 {
		u.Comma()
	}
	for i, c := range u.columns {
		if i > 0 {
			u.Comma()
		}
		u.Ident(c).WriteString(" = ")
		switch v := u.values[i].(type) {
		case Querier:
			u.Join(v)
		default:
			u.Arg(v)
		}
	}
	if u.where != nil {
		u.WriteString(" WHERE ")
		u.Join(u.where)
	}
	statement, args := u.String(), u.args

	return statement, args
}

// DeleteBuilder is a builder for `DELETE` statement.
type DeleteBuilder struct {
	Builder
	table  string
	schema string
	where  *Predicate
}

// Delete creates a builder for the `DELETE` statement.
//
//	Delete("users").
//		Where(
//			Or(
//				EQ("name", "foo").And().EQ("age", 10),
//				EQ("name", "bar").And().EQ("age", 20),
//				And(
//					EQ("name", "qux"),
//					EQ("age", 1).Or().EQ("age", 2),
//				),
//			),
//		)
func Delete(table string) *DeleteBuilder { return &DeleteBuilder{table: table} }

// Schema sets the database name for the table whose row will be deleted.
func (d *DeleteBuilder) Schema(name string) *DeleteBuilder {
	d.schema = name
	return d
}

// Table sets the table name for the table whose row will be deleted.
func (d *DeleteBuilder) Table(name string) *DeleteBuilder {
	d.table = name
	return d
}

// Where appends a where predicate to the `DELETE` statement.
func (d *DeleteBuilder) Where(p *Predicate) *DeleteBuilder {
	if d.where != nil {
		d.where = And(d.where, p)
	} else {
		d.where = p
	}
	return d
}

// FromSelect makes it possible to delete a sub query.
func (d *DeleteBuilder) FromSelect(s *Selector) *DeleteBuilder {
	d.Where(s.where)
	if table, _ := s.from.(*SelectTable); table != nil {
		d.table = table.name
	}
	return d
}

// Query returns query representation of a `DELETE` statement.
func (d *DeleteBuilder) Query() (string, []interface{}) {
	d.WriteString("DELETE FROM ")
	d.writeSchema(d.schema)
	d.Ident(d.table)
	if d.where != nil {
		d.WriteString(" WHERE ")
		d.Join(d.where)
	}
	statement, args := d.String(), d.args
	return statement, args
}

// Predicate is a where predicate.
type Predicate struct {
	Builder
	depth int
	fns   []func(*Builder)
}

// P creates a new predicate.
//
//	P().EQ("name", "a8m").And().EQ("age", 30)
func P(fns ...func(*Builder)) *Predicate {
	return &Predicate{fns: fns}
}

// ExprP creates a new predicate from the given expression.
//
//	ExprP("A = ? AND B > ?", args...)
func ExprP(exr string, args ...interface{}) *Predicate {
	return P(func(b *Builder) {
		b.Join(Expr(exr, args...))
	})
}

// Or combines all given predicates with OR between them.
//
//	Or(EQ("name", "foo"), EQ("name", "bar"))
func Or(preds ...*Predicate) *Predicate {
	p := P()
	return p.Append(func(b *Builder) {
		p.mayWrap(preds, b, "OR")
	})
}

// False appends the FALSE keyword to the predicate.
//
//	Delete().From("users").Where(False())
func False() *Predicate {
	return P().False()
}

// False appends FALSE to the predicate.
func (p *Predicate) False() *Predicate {
	return p.Append(func(b *Builder) {
		b.WriteString("FALSE")
	})
}

// Not wraps the given predicate with the not predicate.
//
//	Not(Or(EQ("name", "foo"), EQ("name", "bar")))
func Not(pred *Predicate) *Predicate {
	return P().Not().Append(func(b *Builder) {
		b.Nested(func(b *Builder) {
			b.Join(pred)
		})
	})
}

// Not appends NOT to the predicate.
func (p *Predicate) Not() *Predicate {
	return p.Append(func(b *Builder) {
		b.WriteString("NOT ")
	})
}

// And combines all given predicates with AND between them.
func And(preds ...*Predicate) *Predicate {
	p := P()
	return p.Append(func(b *Builder) {
		p.mayWrap(preds, b, "AND")
	})
}

// EQ returns a "=" predicate.
func EQ(col string, value interface{}) *Predicate {
	return P().EQ(col, value)
}

// EQ appends a "=" predicate.
func (p *Predicate) EQ(col string, arg interface{}) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col)
		b.WriteOp(OpEQ)
		b.Arg(arg)
	})
}

// NEQ returns a "<>" predicate.
func NEQ(col string, value interface{}) *Predicate {
	return P().NEQ(col, value)
}

// NEQ appends a "<>" predicate.
func (p *Predicate) NEQ(col string, arg interface{}) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col)
		b.WriteOp(OpNEQ)
		b.Arg(arg)
	})
}

// LT returns a "<" predicate.
func LT(col string, value interface{}) *Predicate {
	return P().LT(col, value)
}

// LT appends a "<" predicate.
func (p *Predicate) LT(col string, arg interface{}) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col)
		p.WriteOp(OpLT)
		b.Arg(arg)
	})
}

// LTE returns a "<=" predicate.
func LTE(col string, value interface{}) *Predicate {
	return P().LTE(col, value)
}

// LTE appends a "<=" predicate.
func (p *Predicate) LTE(col string, arg interface{}) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col)
		p.WriteOp(OpLTE)
		b.Arg(arg)
	})
}

// GT returns a ">" predicate.
func GT(col string, value interface{}) *Predicate {
	return P().GT(col, value)
}

// GT appends a ">" predicate.
func (p *Predicate) GT(col string, arg interface{}) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col)
		p.WriteOp(OpGT)
		b.Arg(arg)
	})
}

// GTE returns a ">=" predicate.
func GTE(col string, value interface{}) *Predicate {
	return P().GTE(col, value)
}

// GTE appends a ">=" predicate.
func (p *Predicate) GTE(col string, arg interface{}) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col)
		p.WriteOp(OpGTE)
		b.Arg(arg)
	})
}

// NotNull returns the `IS NOT NULL` predicate.
func NotNull(col string) *Predicate {
	return P().NotNull(col)
}

// NotNull appends the `IS NOT NULL` predicate.
func (p *Predicate) NotNull(col string) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col).WriteString(" IS NOT NULL")
	})
}

// IsNull returns the `IS NULL` predicate.
func IsNull(col string) *Predicate {
	return P().IsNull(col)
}

// IsNull appends the `IS NULL` predicate.
func (p *Predicate) IsNull(col string) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col).WriteString(" IS NULL")
	})
}

// In returns the `IN` predicate.
func In(col string, args ...interface{}) *Predicate {
	return P().In(col, args...)
}

// In appends the `IN` predicate.
func (p *Predicate) In(col string, args ...interface{}) *Predicate {
	if len(args) == 0 {
		return p
	}
	return p.Append(func(b *Builder) {
		b.Ident(col).WriteOp(OpIn)
		b.Nested(func(b *Builder) {
			if s, ok := args[0].(*Selector); ok {
				b.Join(s)
			} else {
				b.Args(args...)
			}
		})
	})
}

// InInts returns the `IN` predicate for ints.
func InInts(col string, args ...int) *Predicate {
	return P().InInts(col, args...)
}

// InValues adds the `IN` predicate for slice of driver.Value.
func InValues(col string, args ...driver.Value) *Predicate {
	return P().InValues(col, args...)
}

// InInts adds the `IN` predicate for ints.
func (p *Predicate) InInts(col string, args ...int) *Predicate {
	iface := make([]interface{}, len(args))
	for i := range args {
		iface[i] = args[i]
	}
	return p.In(col, iface...)
}

// InValues adds the `IN` predicate for slice of driver.Value.
func (p *Predicate) InValues(col string, args ...driver.Value) *Predicate {
	iface := make([]interface{}, len(args))
	for i := range args {
		iface[i] = args[i]
	}
	return p.In(col, iface...)
}

// NotIn returns the `Not IN` predicate.
func NotIn(col string, args ...interface{}) *Predicate {
	return P().NotIn(col, args...)
}

// NotIn appends the `Not IN` predicate.
func (p *Predicate) NotIn(col string, args ...interface{}) *Predicate {
	if len(args) == 0 {
		return p
	}
	return p.Append(func(b *Builder) {
		b.Ident(col).WriteOp(OpNotIn)
		b.Nested(func(b *Builder) {
			if s, ok := args[0].(*Selector); ok {
				b.Join(s)
			} else {
				b.Args(args...)
			}
		})
	})
}

// Like returns the `LIKE` predicate.
func Like(col, pattern string) *Predicate {
	return P().Like(col, pattern)
}

// Like appends the `LIKE` predicate.
func (p *Predicate) Like(col, pattern string) *Predicate {
	return p.Append(func(b *Builder) {
		b.Ident(col).WriteOp(OpLike)
		b.Arg(pattern)
	})
}

// HasPrefix is a helper predicate that checks prefix using the LIKE predicate.
func HasPrefix(col, prefix string) *Predicate {
	return P().HasPrefix(col, prefix)
}

// HasPrefix is a helper predicate that checks prefix using the LIKE predicate.
func (p *Predicate) HasPrefix(col, prefix string) *Predicate {
	return p.Like(col, prefix+"%")
}

// HasSuffix is a helper predicate that checks suffix using the LIKE predicate.
func HasSuffix(col, suffix string) *Predicate { return P().HasSuffix(col, suffix) }

// HasSuffix is a helper predicate that checks suffix using the LIKE predicate.
func (p *Predicate) HasSuffix(col, suffix string) *Predicate {
	return p.Like(col, "%"+suffix)
}

// EqualFold is a helper predicate that applies the "=" predicate with case-folding.
func EqualFold(col, sub string) *Predicate { return P().EqualFold(col, sub) }

// EqualFold is a helper predicate that applies the "=" predicate with case-folding.
func (p *Predicate) EqualFold(col, sub string) *Predicate {
	return p.Append(func(b *Builder) {
		f := &Func{}
		f.SetDialect(b.dialect)
		f.Lower(col)
		b.WriteString(f.String())
		b.WriteOp(OpEQ)
		b.Arg(strings.ToLower(sub))
	})
}

// Contains is a helper predicate that checks substring using the LIKE predicate.
func Contains(col, sub string) *Predicate { return P().Contains(col, sub) }

// Contains is a helper predicate that checks substring using the LIKE predicate.
func (p *Predicate) Contains(col, sub string) *Predicate {
	return p.Like(col, "%"+sub+"%")
}

// CompositeGT returns a comiposite ">" predicate
func CompositeGT(columns []string, args ...interface{}) *Predicate {
	return P().CompositeGT(columns, args...)
}

// CompositeLT returns a comiposite "<" predicate
func CompositeLT(columns []string, args ...interface{}) *Predicate {
	return P().CompositeLT(columns, args...)
}

func (p *Predicate) compositeP(operator string, columns []string, args ...interface{}) *Predicate {
	return p.Append(func(b *Builder) {
		b.Nested(func(nb *Builder) {
			nb.IdentComma(columns...)
		})
		b.WriteString(operator)
		b.WriteString("(")
		b.Args(args...)
		b.WriteString(")")
	})
}

// CompositeGT returns a composite ">" predicate.
func (p *Predicate) CompositeGT(columns []string, args ...interface{}) *Predicate {
	const operator = " > "
	return p.compositeP(operator, columns, args...)
}

// CompositeLT appends a composite "<" predicate.
func (p *Predicate) CompositeLT(columns []string, args ...interface{}) *Predicate {
	const operator = " < "
	return p.compositeP(operator, columns, args...)
}

// Append appends a new function to the predicate callbacks.
// The callback list are executed on call to 	ry.
func (p *Predicate) Append(f func(*Builder)) *Predicate {
	p.fns = append(p.fns, f)
	return p
}

// Query returns query representation of a predicate.
func (p *Predicate) Query() (string, []interface{}) {
	if p.Len() > 0 || len(p.args) > 0 {
		p.Reset()
		p.args = nil
	}
	for _, f := range p.fns {
		f(&p.Builder)
	}
	return p.String(), p.args
}

// clone returns a shallow clone of p.
func (p *Predicate) clone() *Predicate {
	if p == nil {
		return p
	}
	return &Predicate{fns: append([]func(*Builder){}, p.fns...)}
}

func (p *Predicate) mayWrap(preds []*Predicate, b *Builder, op string) {
	switch n := len(preds); {
	case n == 1:
		b.Join(preds[0])
		return
	case n > 1 && p.depth != 0:
		b.WriteByte('(')
		defer b.WriteByte(')')
	}
	for i := range preds {
		preds[i].depth = p.depth + 1
		if i > 0 {
			b.WriteByte(' ')
			b.WriteString(op)
			b.WriteByte(' ')
		}
		if len(preds[i].fns) > 1 {
			b.Nested(func(b *Builder) {
				b.Join(preds[i])
			})
		} else {
			b.Join(preds[i])
		}
	}
}

// Func represents an SQL function.
type Func struct {
	Builder
	fns []func(*Builder)
}

// Lower wraps the given column with the LOWER function.
//
//	P().EQ(sql.Lower("name"), "a8m")
func Lower(ident string) string {
	f := &Func{}
	f.Lower(ident)
	return f.String()
}

// Lower wraps the given ident with the LOWER function.
func (f *Func) Lower(ident string) {
	f.byName("LOWER", ident)
}

// Count wraps the ident with the COUNT aggregation function.
func Count(ident string) string {
	f := &Func{}
	f.Count(ident)
	return f.String()
}

// Count wraps the ident with the COUNT aggregation function.
func (f *Func) Count(ident string) {
	f.byName("COUNT", ident)
}

// Max wraps the ident with the MAX aggregation function.
func Max(ident string) string {
	f := &Func{}
	f.Max(ident)
	return f.String()
}

// Max wraps the ident with the MAX aggregation function.
func (f *Func) Max(ident string) {
	f.byName("MAX", ident)
}

// Min wraps the ident with the MIN aggregation function.
func Min(ident string) string {
	f := &Func{}
	f.Min(ident)
	return f.String()
}

// Min wraps the ident with the MIN aggregation function.
func (f *Func) Min(ident string) {
	f.byName("MIN", ident)
}

// Sum wraps the ident with the SUM aggregation function.
func Sum(ident string) string {
	f := &Func{}
	f.Sum(ident)
	return f.String()
}

// Sum wraps the ident with the SUM aggregation function.
func (f *Func) Sum(ident string) {
	f.byName("SUM", ident)
}

// Avg wraps the ident with the AVG aggregation function.
func Avg(ident string) string {
	f := &Func{}
	f.Avg(ident)
	return f.String()
}

// Avg wraps the ident with the AVG aggregation function.
func (f *Func) Avg(ident string) {
	f.byName("AVG", ident)
}

// byName wraps an identifier with a function name.
func (f *Func) byName(fn, ident string) {
	f.Append(func(b *Builder) {
		f.WriteString(fn)
		f.Nested(func(b *Builder) {
			b.Ident(ident)
		})
	})
}

// Append appends a new function to the function callbacks.
// The callback list are executed on call to String.
func (f *Func) Append(fn func(*Builder)) *Func {
	f.fns = append(f.fns, fn)
	return f
}

// String implements the fmt.Stringer.
func (f *Func) String() string {
	for _, fn := range f.fns {
		fn(&f.Builder)
	}
	return f.Builder.String()
}

// As suffixed the given column with an alias (`a` AS `b`).
func As(ident string, as string) string {
	b := &Builder{}
	b.Ident(ident).Pad().WriteString("AS")
	b.Pad().Ident(as)
	return b.String()
}

// Distinct prefixed the given columns with the `DISTINCT` keyword (DISTINCT `id`).
func Distinct(idents ...string) string {
	b := &Builder{}
	b.WriteString("DISTINCT")
	b.Pad().IdentComma(idents...)
	return b.String()
}

// TableView is a view that returns a table view. Can ne a Table, Selector or a View (WITH statement).
type TableView interface {
	view()
}

// SelectTable is a table selector.
type SelectTable struct {
	Builder
	as     string
	name   string
	schema string
	quote  bool
}

// Table returns a new table selector.
//
//	t1 := Table("users").As("u")
//	return Select(t1.C("name"))
func Table(name string) *SelectTable {
	return &SelectTable{quote: true, name: name}
}

// Schema sets the schema name of the table.
func (s *SelectTable) Schema(name string) *SelectTable {
	s.schema = name
	return s
}

// As adds the AS clause to the table selector.
func (s *SelectTable) As(alias string) *SelectTable {
	s.as = alias
	return s
}

// C returns a formatted string for the table column.
func (s *SelectTable) C(column string) string {
	name := s.name
	if s.as != "" {
		name = s.as
	}
	b := &Builder{dialect: s.dialect}
	if s.as == "" {
		b.writeSchema(s.schema)
	}
	b.Ident(name).WriteByte('.').Ident(column)
	return b.String()
}

// Columns returns a list of formatted strings for the table columns.
func (s *SelectTable) Columns(columns ...string) []string {
	names := make([]string, 0, len(columns))
	for _, c := range columns {
		names = append(names, s.C(c))
	}
	return names
}

// Unquote makes the table name to be formatted as raw string (unquoted).
// It is useful whe you don't want to query tables under the current database.
// For example: "INFORMATION_SCHEMA.TABLE_CONSTRAINTS" in MySQL.
func (s *SelectTable) Unquote() *SelectTable {
	s.quote = false
	return s
}

// ref returns the table reference.
func (s *SelectTable) ref() string {
	if !s.quote {
		return s.name
	}
	b := &Builder{dialect: s.dialect}
	b.writeSchema(s.schema)
	b.Ident(s.name)
	if s.as != "" {
		b.WriteString(" AS ")
		b.Ident(s.as)
	}
	return b.String()
}

// implement the table view.
func (*SelectTable) view() {}

// join table option.
type join struct {
	on    *Predicate
	kind  string
	table TableView
}

// clone a joiner.
func (j join) clone() join {
	if sel, ok := j.table.(*Selector); ok {
		j.table = sel.Clone()
	}
	j.on = j.on.clone()
	return j
}

// Selector is a builder for the `SELECT` statement.
type Selector struct {
	Builder
	// ctx stores contextual data typically from
	// generated code such as alternate table schemas.
	ctx      context.Context
	as       string
	columns  []string
	from     TableView
	joins    []join
	where    *Predicate
	or       bool
	not      bool
	order    []interface{}
	group    []string
	having   *Predicate
	limit    *int
	offset   *int
	distinct bool
	lock     *LockOptions
	index    *IndexOptions
}

// WithContext sets the context into the *Selector.
func (s *Selector) WithContext(ctx context.Context) *Selector {
	if ctx == nil {
		panic("nil context")
	}
	s.ctx = ctx
	return s
}

// Context returns the Selector context or Background
// if nil.
func (s *Selector) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}
	return context.Background()
}

// Select returns a new selector for the `SELECT` statement.
//
//	t1 := Table("users").As("u")
//	t2 := Select().From(Table("groups")).Where(EQ("user_id", 10)).As("g")
//	return Select(t1.C("id"), t2.C("name")).
//			From(t1).
//			Join(t2).
//			On(t1.C("id"), t2.C("user_id"))
func Select(columns ...string) *Selector {
	return (&Selector{}).Select(columns...)
}

// Select changes the columns selection of the SELECT statement.
// Empty selection means all columns *.
func (s *Selector) Select(columns ...string) *Selector {
	s.columns = columns
	return s
}

// From sets the source of `FROM` clause.
func (s *Selector) From(t TableView) *Selector {
	s.from = t
	if st, ok := t.(state); ok {
		st.SetDialect(s.dialect)
	}
	return s
}

// Distinct adds the DISTINCT keyword to the `SELECT` statement.
func (s *Selector) Distinct() *Selector {
	s.distinct = true
	return s
}

// SetDistinct sets explicitly if the returned rows are distinct or indistinct.
func (s *Selector) SetDistinct(v bool) *Selector {
	s.distinct = v
	return s
}

// Limit adds the `LIMIT` clause to the `SELECT` statement.
func (s *Selector) Limit(limit int) *Selector {
	s.limit = &limit
	return s
}

// Offset adds the `OFFSET` clause to the `SELECT` statement.
func (s *Selector) Offset(offset int) *Selector {
	s.offset = &offset
	return s
}

// Where sets or appends the given predicate to the statement.
func (s *Selector) Where(p *Predicate) *Selector {
	if s.not {
		p = Not(p)
		s.not = false
	}
	switch {
	case s.where == nil:
		s.where = p
	case s.where != nil && s.or:
		s.where = Or(s.where, p)
		s.or = false
	default:
		s.where = And(s.where, p)
	}
	return s
}

// P returns the predicate of a selector.
func (s *Selector) P() *Predicate {
	return s.where
}

// SetP sets explicitly the predicate function for the selector and clear its previous state.
func (s *Selector) SetP(p *Predicate) *Selector {
	s.where = p
	s.or = false
	s.not = false
	return s
}

// FromSelect copies the predicate from a selector.
func (s *Selector) FromSelect(s2 *Selector) *Selector {
	s.where = s2.where
	return s
}

// Not sets the next coming predicate with not.
func (s *Selector) Not() *Selector {
	s.not = true
	return s
}

// Or sets the next coming predicate with OR operator (disjunction).
func (s *Selector) Or() *Selector {
	s.or = true
	return s
}

// Table returns the selected table.
func (s *Selector) Table() *SelectTable {
	return s.from.(*SelectTable)
}

// Join appends a `JOIN` clause to the statement.
func (s *Selector) Join(t TableView) *Selector {
	return s.join("JOIN", t)
}

// LeftJoin appends a `LEFT JOIN` clause to the statement.
func (s *Selector) LeftJoin(t TableView) *Selector {
	return s.join("LEFT JOIN", t)
}

// RightJoin appends a `RIGHT JOIN` clause to the statement.
func (s *Selector) RightJoin(t TableView) *Selector {
	return s.join("RIGHT JOIN", t)
}

// join adds a join table to the selector with the given kind.
func (s *Selector) join(kind string, t TableView) *Selector {
	s.joins = append(s.joins, join{
		kind:  kind,
		table: t,
	})
	switch view := t.(type) {
	case *SelectTable:
		if view.as == "" {
			view.as = "t0"
		}
	case *Selector:
		if view.as == "" {
			view.as = "t" + strconv.Itoa(len(s.joins))
		}
	}
	if st, ok := t.(state); ok {
		st.SetDialect(s.dialect)
	}
	return s
}

// C returns a formatted string for a selected column from this statement.
func (s *Selector) C(column string) string {
	if s.as != "" {
		b := &Builder{dialect: s.dialect}
		b.Ident(s.as)
		b.WriteByte('.')
		b.Ident(column)
		return b.String()
	}
	return s.Table().C(column)
}

// Columns returns a list of formatted strings for a selected columns from this statement.
func (s *Selector) Columns(columns ...string) []string {
	names := make([]string, 0, len(columns))
	for _, c := range columns {
		names = append(names, s.C(c))
	}
	return names
}

func (s *Selector) SelectedColumns() []string {
	columns := make([]string, 0, len(s.columns))
	columns = append(columns, s.columns...)
	return columns
}

// SelectColumnsLen returns len of select columns
func (s *Selector) SelectColumnsLen() int {
	return len(s.columns)
}

// OnP sets or appends the given predicate for the `ON` clause of the statement.
func (s *Selector) OnP(p *Predicate) *Selector {
	if len(s.joins) > 0 {
		join := &s.joins[len(s.joins)-1]
		switch {
		case join.on == nil:
			join.on = p
		default:
			join.on = And(join.on, p)
		}
	}
	return s
}

// On sets the `ON` clause for the `JOIN` operation.
func (s *Selector) On(c1, c2 string) *Selector {
	s.OnP(P(func(builder *Builder) {
		builder.Ident(c1).WriteOp(OpEQ).Ident(c2)
	}))
	return s
}

// As give this selection an alias.
func (s *Selector) As(alias string) *Selector {
	s.as = alias
	return s
}

// Count sets the Select statement to be a `SELECT COUNT(*)`.
func (s *Selector) Count(columns ...string) *Selector {
	column := "*"
	if len(columns) > 0 {
		b := &Builder{}
		b.IdentComma(columns...)
		column = b.String()
	}
	s.columns = []string{Count(column)}
	return s
}

type IndexOptions struct {
	action    IndexAction
	forWhat   ForWhat
	indexName []string
}

type ForWhat string
type IndexAction string

const (
	ForGroupBy ForWhat = "FOR GROUP BY"
	ForOrderBy ForWhat = "FOR ORDER BY"
	ForJoin    ForWhat = "FOR JOIN"
	ForDefault ForWhat = ""

	USE    IndexAction = "USE INDEX"
	IGNORE IndexAction = "IGNORE INDEX"
	FORCE  IndexAction = "FORCE INDEX"
)

// UseIndex  generate  USE INDEX (`indexName`)
func (s *Selector) UseIndex(indexName ...string) *Selector {
	s.UseIndexFor(ForDefault, indexName...)
	return s
}

// ForceIndex generate FORCE INDEX (`indexName`)
func (s *Selector) ForceIndex(indexName ...string) *Selector {
	s.ForceIndexFor(ForDefault, indexName...)
	return s
}

// IgnoreIndex generate IGNORE INDEX (`indexName`)
func (s *Selector) IgnoreIndex(indexName ...string) *Selector {
	s.IgnoreIndexFor(ForDefault, indexName...)
	return s
}

// UseIndexFor generate USE INDEX FOR xxx (`indexName`)
func (s *Selector) UseIndexFor(forwhat ForWhat, indexName ...string) *Selector {
	s.index = &IndexOptions{
		action:    USE,
		forWhat:   forwhat,
		indexName: indexName,
	}
	return s
}

// ForceIndexFor generate FORCE INDEX FOR xxx (`indexName`)
func (s *Selector) ForceIndexFor(forwhat ForWhat, indexName ...string) *Selector {
	s.index = &IndexOptions{
		action:    FORCE,
		forWhat:   forwhat,
		indexName: indexName,
	}
	return s
}

// IgnoreIndexFor generate IGNORE INDEX FOR xxx (`indexName`)
func (s *Selector) IgnoreIndexFor(forwhat ForWhat, indexName ...string) *Selector {
	s.index = &IndexOptions{
		action:    IGNORE,
		forWhat:   forwhat,
		indexName: indexName,
	}
	return s
}

// LockAction tells the transaction what to do in case of
// requesting a row that is locked by other transaction.
type LockAction string

const (
	// NoWait means never wait and returns an error.
	NoWait LockAction = "NOWAIT"
	// SkipLocked means never wait and skip.
	SkipLocked LockAction = "SKIP LOCKED"
)

// LockStrength defines the strength of the lock (see the list below).
type LockStrength string

// A list of all locking clauses.
const (
	LockShare  LockStrength = "SHARE"
	LockUpdate LockStrength = "UPDATE"
)

type (
	// LockOptions defines a SELECT statement
	// lock for protecting concurrent updates.
	LockOptions struct {
		// Strength of the lock.
		Strength LockStrength
		// Action of the lock.
		Action LockAction
	}
	// LockOption allows configuring the LockConfig using functional options.
	LockOption func(*LockOptions)
)

// WithLockAction sets the Action of the lock.
func WithLockAction(action LockAction) LockOption {
	return func(c *LockOptions) {
		c.Action = action
	}
}

// For sets the lock configuration for suffixing the `SELECT`
// statement with the `FOR [SHARE | UPDATE] ...` clause.
func (s *Selector) For(l LockStrength, opts ...LockOption) *Selector {
	s.lock = &LockOptions{Strength: l}
	for _, opt := range opts {
		opt(s.lock)
	}
	return s
}

// ForShare sets the lock configuration for suffixing the
// `SELECT` statement with the `FOR SHARE` clause.
func (s *Selector) ForShare(opts ...LockOption) *Selector {
	return s.For(LockShare, opts...)
}

// ForUpdate sets the lock configuration for suffixing the
// `SELECT` statement with the `FOR UPDATE` clause.
func (s *Selector) ForUpdate(opts ...LockOption) *Selector {
	return s.For(LockUpdate, opts...)
}

// Clone returns a duplicate of the selector, including all associated steps. It can be
// used to prepare common SELECT statements and use them differently after the clone is made.
func (s *Selector) Clone() *Selector {
	if s == nil {
		return nil
	}
	joins := make([]join, len(s.joins))
	for i := range s.joins {
		joins[i] = s.joins[i].clone()
	}
	return &Selector{
		Builder:  s.Builder.clone(),
		ctx:      s.ctx,
		as:       s.as,
		or:       s.or,
		not:      s.not,
		from:     s.from,
		limit:    s.limit,
		offset:   s.offset,
		distinct: s.distinct,
		where:    s.where.clone(),
		having:   s.having.clone(),
		joins:    append([]join{}, joins...),
		group:    append([]string{}, s.group...),
		order:    append([]interface{}{}, s.order...),
		columns:  append([]string{}, s.columns...),
	}
}

// Asc adds the ASC suffix for the given column.
func Asc(column string) string {
	b := &Builder{}
	b.Ident(column).WriteString(" ASC")
	return b.String()
}

// Desc adds the DESC suffix for the given column.
func Desc(column string) string {
	b := &Builder{}
	b.Ident(column).WriteString(" DESC")
	return b.String()
}

// OrderBy appends the `ORDER BY` clause to the `SELECT` statement.
func (s *Selector) OrderBy(columns ...string) *Selector {
	for i := range columns {
		s.order = append(s.order, columns[i])
	}
	return s
}

// OrderExpr appends the `ORDER BY` clause to the `SELECT`
// statement with custom list of expressions.
func (s *Selector) OrderExpr(exprs ...Querier) *Selector {
	for i := range exprs {
		s.order = append(s.order, exprs[i])
	}
	return s
}

// GroupBy appends the `GROUP BY` clause to the `SELECT` statement.
func (s *Selector) GroupBy(columns ...string) *Selector {
	s.group = append(s.group, columns...)
	return s
}

// Having appends a predicate for the `HAVING` clause.
func (s *Selector) Having(p *Predicate) *Selector {
	s.having = p
	return s
}

// Query returns query representation of a `SELECT` statement.
func (s *Selector) Query() (string, []interface{}) {
	b := s.Builder.clone()
	b.WriteString("SELECT ")
	if s.distinct {
		b.WriteString("DISTINCT ")
	}
	if len(s.columns) > 0 {
		b.IdentComma(s.columns...)
	} else {
		b.WriteString("*")
	}
	b.WriteString(" FROM ")
	switch t := s.from.(type) {
	case *SelectTable:
		t.SetDialect(s.dialect)
		b.WriteString(t.ref())
	case *Selector:
		t.SetDialect(s.dialect)
		b.Nested(func(b *Builder) {
			b.Join(t)
		})
		b.WriteString(" AS ")
		b.Ident(t.as)
	}
	if s.index != nil {
		s.joinIndex(&b)
	}
	for _, join := range s.joins {
		b.WriteString(" " + join.kind + " ")
		switch view := join.table.(type) {
		case *SelectTable:
			view.SetDialect(s.dialect)
			b.WriteString(view.ref())
		case *Selector:
			view.SetDialect(s.dialect)
			b.Nested(func(b *Builder) {
				b.Join(view)
			})
			b.WriteString(" AS ")
			b.Ident(view.as)
		}
		if join.on != nil {
			b.WriteString(" ON ")
			b.Join(join.on)
		}
	}
	if s.where != nil {
		b.WriteString(" WHERE ")
		b.Join(s.where)
	}
	if len(s.group) > 0 {
		b.WriteString(" GROUP BY ")
		b.IdentComma(s.group...)
	}
	if s.having != nil {
		b.WriteString(" HAVING ")
		b.Join(s.having)
	}
	if len(s.order) > 0 {
		s.joinOrder(&b)
	}
	if s.limit != nil {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.Itoa(*s.limit))
	}
	if s.offset != nil {
		b.WriteString(" OFFSET ")
		b.WriteString(strconv.Itoa(*s.offset))
	}
	s.joinLock(&b)
	s.total = b.total
	statement, args := b.String(), b.args

	return statement, args
}

func (s *Selector) joinIndex(b *Builder) {
	if len(s.index.indexName) == 0 {
		return
	}
	b.Pad()
	b.WriteString(string(s.index.action))
	b.Pad()
	b.WriteString(string(s.index.forWhat))
	b.Nested(func(b *Builder) {
		for i := range s.index.indexName {
			if i > 0 {
				b.Comma()
			}
			b.Ident(s.index.indexName[i])
		}
	})
}

func (s *Selector) joinOrder(b *Builder) {
	b.WriteString(" ORDER BY ")
	for i := range s.order {
		if i > 0 {
			b.Comma()
		}
		switch order := s.order[i].(type) {
		case string:
			b.Ident(order)
		case Querier:
			b.Join(order)
		}
	}
}
func (s *Selector) joinLock(b *Builder) {
	if s.lock == nil {
		return
	}
	b.Pad()
	switch s.lock.Strength {
	case LockShare:
		b.WriteString("LOCK IN SHARE MODE")
	case LockUpdate:
		b.WriteString("FOR ").WriteString(string(s.lock.Strength))

	}
	if s.lock.Action != "" {
		b.Pad().WriteString(string(s.lock.Action))
	}
}

// implement the table view interface.
func (*Selector) view() {}

// Raw returns a raw SQL query that is placed as-is in the query.
func Raw(s string) Querier { return &raw{s} }

type raw struct{ s string }

func (r *raw) Query() (string, []interface{}) { return r.s, nil }

// Expr returns an SQL expression that implements the Querier interface.
func Expr(exr string, args ...interface{}) Querier { return &expr{s: exr, args: args} }

type expr struct {
	s    string
	args []interface{}
}

func (e *expr) Query() (string, []interface{}) { return e.s, e.args }

// Queries are list of queries join with space between them.
type Queries []Querier

// Query returns query representation of Queriers.
func (n Queries) Query() (string, []interface{}) {
	b := &Builder{}
	for i := range n {
		if i > 0 {
			b.Pad()
		}
		query, args := n[i].Query()
		b.WriteString(query)
		b.args = append(b.args, args...)
	}
	return b.String(), b.args
}

// Builder is the base query builder for the sql dsl.
type Builder struct {
	sb      *strings.Builder // underlying builder.
	dialect string           // configured dialect.
	args    []interface{}    // query parameters.
	total   int              // total number of parameters in query tree.
	errs    []error          // errors that added during the query construction.
}

// Quote quotes the given identifier with the characters based
// on the configured dialect. It defaults to "`".
func (b *Builder) Quote(ident string) string {
	switch {

	// An identifier for unknown dialect.
	case b.dialect == "" && strings.ContainsAny(ident, "`\""):
		return ident
	default:
		return fmt.Sprintf("`%s`", ident)
	}
}

// Ident appends the given string as an identifier.
func (b *Builder) Ident(s string) *Builder {
	switch {
	case len(s) == 0:
	case s != "*" && !b.isIdent(s) && !isFunc(s) && !isModifier(s):
		b.WriteString(b.Quote(s))
	default:
		b.WriteString(s)
	}
	return b
}

// IdentComma calls Ident on all arguments and adds a comma between them.
func (b *Builder) IdentComma(s ...string) *Builder {
	for i := range s {
		if i > 0 {
			b.Comma()
		}
		b.Ident(s[i])
	}
	return b
}

// String returns the accumulated string.
func (b *Builder) String() string {
	if b.sb == nil {
		return ""
	}
	return b.sb.String()
}

// WriteByte wraps the Buffer.WriteByte to make it chainable with other methods.
func (b *Builder) WriteByte(c byte) *Builder {
	if b.sb == nil {
		b.sb = &strings.Builder{}
	}
	b.sb.WriteByte(c)
	return b
}

// WriteString wraps the Buffer.WriteString to make it chainable with other methods.
func (b *Builder) WriteString(s string) *Builder {
	if b.sb == nil {
		b.sb = &strings.Builder{}
	}
	b.sb.WriteString(s)
	return b
}

// Len returns the number of accumulated bytes.
func (b *Builder) Len() int {
	if b.sb == nil {
		return 0
	}
	return b.sb.Len()
}

// Reset resets the Builder to be empty.
func (b *Builder) Reset() *Builder {
	if b.sb != nil {
		b.sb.Reset()
	}
	return b
}

// AddError appends an error to the builder errors.
func (b *Builder) AddError(err error) *Builder {
	b.errs = append(b.errs, err)
	return b
}

func (b *Builder) writeSchema(schema string) {
	if schema != "" {
		b.Ident(schema).WriteByte('.')
	}
}

// Err returns a concatenated error of all errors encountered during
// the query-building, or were added manually by calling AddError.
func (b *Builder) Err() error {
	if len(b.errs) == 0 {
		return nil
	}
	br := strings.Builder{}
	for i := range b.errs {
		if i > 0 {
			br.WriteString("; ")
		}
		br.WriteString(b.errs[i].Error())
	}
	return fmt.Errorf(br.String())
}

// An Op represents a predicate operator.
type Op int

// Predicate operators
const (
	OpEQ      Op = iota // logical and.
	OpNEQ               // <>
	OpGT                // >
	OpGTE               // >=
	OpLT                // <
	OpLTE               // <=
	OpIn                // IN
	OpNotIn             // NOT IN
	OpLike              // LIKE
	OpIsNull            // IS NULL
	OpNotNull           // IS NOT NULL
)

var ops = [...]string{
	OpEQ:      "=",
	OpNEQ:     "<>",
	OpGT:      ">",
	OpGTE:     ">=",
	OpLT:      "<",
	OpLTE:     "<=",
	OpIn:      "IN",
	OpNotIn:   "NOT IN",
	OpLike:    "LIKE",
	OpIsNull:  "IS NULL",
	OpNotNull: "IS NOT NULL",
}
var opsmap = map[string]Op{
	"=":      OpEQ,
	"<>":     OpNEQ,
	">":      OpGT,
	">=":     OpGTE,
	"<":      OpLT,
	"<=":     OpLTE,
	"IN":     OpIn,
	"NOT IN": OpNotIn,
	"LIKE":   OpLike,
}

func VailedOp(op string) (vailed bool, t Op) {
	o := strings.ToUpper(strings.TrimSpace(op))
	if p, ok := opsmap[o]; ok {
		return true, p
	}
	return false, Op(-1)
}

func GenP(field, op, value string) (*Predicate, error) {
	v, o := VailedOp(op)
	if !v {
		return nil, fmt.Errorf("op:%s is not support", op)
	}
	switch o {
	case OpEQ:
		return EQ(field, value), nil
	case OpNEQ:
		return NEQ(field, value), nil
	case OpGT:
		return GT(field, value), nil
	case OpGTE:
		return GTE(field, value), nil
	case OpLT:
		return LT(field, value), nil
	case OpLTE:
		return LTE(field, value), nil
	case OpIn:
		vs := strings.Split(value, ",")
		is := make([]interface{}, 0, len(vs))
		for _, i := range vs {
			is = append(is, i)
		}
		return In(field, is...), nil
	case OpNotIn:
		vs := strings.Split(value, ",")
		is := make([]interface{}, 0, len(vs))
		for _, i := range vs {
			is = append(is, i)
		}
		return NotIn(field, is...), nil
	case OpLike:
		return Like(field, value), nil
	default:
		return nil, fmt.Errorf("op:%s is not support", op)
	}

}

// WriteOp writes an operator to the builder.
func (b *Builder) WriteOp(op Op) *Builder {
	switch {
	case op >= OpEQ && op <= OpLike:
		b.Pad().WriteString(ops[op]).Pad()
	case op == OpIsNull || op == OpNotNull:
		b.Pad().WriteString(ops[op])
	default:
		panic(fmt.Sprintf("invalid op %d", op))
	}
	return b
}

type (
	// StmtInfo holds an information regarding
	// the statement
	StmtInfo struct {
		// The Dialect of the SQL driver.
		Dialect string
	}
	// ParamFormatter wraps the FormatPram function.
	ParamFormatter interface {
		// The FormatParam function lets users to define
		// custom placeholder formatting for their types.
		// For example, formatting the default placeholder
		// from '?' to 'ST_GeomFromWKB(?)' for MySQL dialect.
		FormatParam(placeholder string, info *StmtInfo) string
	}
)

// Arg appends an input argument to the builder.
func (b *Builder) Arg(a interface{}) *Builder {
	switch a := a.(type) {
	case *raw:
		b.WriteString(a.s)
		return b
	case Querier:
		b.Join(a)
		return b
	}
	b.total++
	b.args = append(b.args, a)
	// Default placeholder param (MySQL and SQLite).
	param := "?"

	if f, ok := a.(ParamFormatter); ok {
		param = f.FormatParam(param, &StmtInfo{
			Dialect: b.dialect,
		})
	}
	b.WriteString(param)
	return b
}

// Args appends a list of arguments to the builder.
func (b *Builder) Args(a ...interface{}) *Builder {
	for i := range a {
		if i > 0 {
			b.Comma()
		}
		b.Arg(a[i])
	}
	return b
}

// Comma adds a comma to the query.
func (b *Builder) Comma() *Builder {
	b.WriteString(", ")
	return b
}

// Pad adds a space to the query.
func (b *Builder) Pad() *Builder {
	b.WriteString(" ")
	return b
}

// Join joins a list of Queries to the builder.
func (b *Builder) Join(qs ...Querier) *Builder {
	return b.join(qs, "")
}

// JoinComma joins a list of Queries and adds comma between them.
func (b *Builder) JoinComma(qs ...Querier) *Builder {
	return b.join(qs, ", ")
}

// join joins a list of Queries to the builder with a given separator.
func (b *Builder) join(qs []Querier, sep string) *Builder {
	for i, q := range qs {
		if i > 0 {
			b.WriteString(sep)
		}
		st, ok := q.(state)
		if ok {
			st.SetDialect(b.dialect)
			st.SetTotal(b.total)
		}
		query, args := q.Query()
		b.WriteString(query)
		b.args = append(b.args, args...)
		b.total += len(args)
	}
	return b
}

// Nested gets a callback, and wraps its result with parentheses.
func (b *Builder) Nested(f func(*Builder)) *Builder {
	nb := &Builder{dialect: b.dialect, total: b.total, sb: &strings.Builder{}}
	nb.WriteByte('(')
	f(nb)
	nb.WriteByte(')')
	b.WriteString(nb.String())
	b.args = append(b.args, nb.args...)
	b.total = nb.total
	return b
}

// SetDialect sets the builder dialect. It's used for garnering dialect specific queries.
func (b *Builder) SetDialect(dialect string) {
	b.dialect = dialect
}

// Dialect returns the dialect of the builder.
func (b Builder) Dialect() string {
	return b.dialect
}

// Total returns the total number of arguments so far.
func (b Builder) Total() int {
	return b.total
}

// SetTotal sets the value of the total arguments.
// Used to pass this information between sub queries/expressions.
func (b *Builder) SetTotal(total int) {
	b.total = total
}

// Query implements the Querier interface.
func (b Builder) Query() (string, []interface{}) {
	return b.String(), b.args
}

// clone returns a shallow clone of a builder.
func (b Builder) clone() Builder {
	c := Builder{dialect: b.dialect, total: b.total, sb: &strings.Builder{}}
	if len(b.args) > 0 {
		c.args = append(c.args, b.args...)
	}
	if b.sb != nil {
		c.sb.WriteString(b.sb.String())
	}

	return c
}

// isIdent reports if the given string is a dialect identifier.
func (b *Builder) isIdent(s string) bool {
	return strings.Contains(s, "`")
}

// state wraps the all methods for setting and getting
// update state between all queries in the query tree.
type state interface {
	Dialect() string
	SetDialect(string)
	Total() int
	SetTotal(int)
}

func isFunc(s string) bool {
	return strings.Contains(s, "(") && strings.Contains(s, ")")
}

func isModifier(s string) bool {
	for _, m := range [...]string{"DISTINCT", "ALL", "WITH ROLLUP"} {
		if strings.HasPrefix(s, m) {
			return true
		}
	}
	return false
}

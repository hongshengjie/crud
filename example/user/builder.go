// Code generated by bcurd. DO NOT EDIT.

package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hongshengjie/crud/xsql"

	"time"
)

// InsertBuilder InsertBuilder
type InsertBuilder struct {
	eq      xsql.ExecQuerier
	builder *xsql.InsertBuilder
	a       []*User
	upsert  bool
}

// Create Create
func Create(eq xsql.ExecQuerier) *InsertBuilder {
	return &InsertBuilder{
		builder: xsql.Insert(table),
		eq:      eq,
	}
}

// SetUser SetUser
func (in *InsertBuilder) SetUser(a ...*User) *InsertBuilder {
	in.a = append(in.a, a...)
	return in
}

// Upsert update all field when insert conflict
func (in *InsertBuilder) Upsert(ctx context.Context) (int64, error) {
	in.upsert = true
	return in.Save(ctx)
}

// Save Save one or many records set by SetUser method
// if insert a record , the LastInsertId  will be setted on the struct's  PrimeKey field
// if insert many records , every struct's PrimeKey field will not be setted
// return number of RowsAffected or error
func (in *InsertBuilder) Save(ctx context.Context) (int64, error) {
	if len(in.a) == 0 {
		return 0, errors.New("please set a User")
	}
	in.builder.Columns(Id, Name, Age, Ctime, Mtime)
	if in.upsert {
		in.builder.OnDuplicateKeyUpdate(Id, Name, Age, Ctime, Mtime)
	}
	for _, a := range in.a {
		if a == nil {
			return 0, errors.New("can not insert a nil User")
		}
		in.builder.Values(a.Id, a.Name, a.Age, a.Ctime, a.Mtime)
	}

	ins, args := in.builder.Query()
	result, err := in.eq.ExecContext(ctx, ins, args...)
	if err != nil {
		return 0, err
	}

	if len(in.a) == 1 {
		id, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		if in.a[0].Id == 0 && id > 0 {
			in.a[0].Id = int64(id)
		}
	}

	return result.RowsAffected()
}

// DeleteBuilder DeleteBuilder
type DeleteBuilder struct {
	builder *xsql.DeleteBuilder
	eq      xsql.ExecQuerier
}

// Delete Delete
func Delete(eq xsql.ExecQuerier) *DeleteBuilder {
	return &DeleteBuilder{
		builder: xsql.Delete(table),
		eq:      eq,
	}
}

// Where  UserWhere
func (d *DeleteBuilder) Where(p ...UserWhere) *DeleteBuilder {
	s := &xsql.Selector{}
	for _, v := range p {
		v(s)
	}
	d.builder = d.builder.Where(s.P())
	return d
}

// Exec Exec
func (d *DeleteBuilder) Exec(ctx context.Context) (int64, error) {
	del, args := d.builder.Query()
	res, err := d.eq.ExecContext(ctx, del, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// SelectBuilder SelectBuilder
type SelectBuilder struct {
	builder   *xsql.Selector
	eq        xsql.ExecQuerier
	selectAll bool
}

// Find Find
func Find(eq xsql.ExecQuerier) *SelectBuilder {
	sel := &SelectBuilder{
		builder: xsql.Select(),
		eq:      eq,
	}
	sel.builder = sel.builder.From(xsql.Table(table))
	return sel
}

// Select Select
func (s *SelectBuilder) Select(columns ...string) *SelectBuilder {
	s.builder.Select(columns...)
	return s
}

// Count Count
func (s *SelectBuilder) Count(columns ...string) *SelectBuilder {
	s.builder.Count(columns...)
	return s
}

//Where where
func (s *SelectBuilder) Where(p ...UserWhere) *SelectBuilder {
	sel := &xsql.Selector{}
	for _, v := range p {
		v(sel)
	}
	s.builder = s.builder.Where(sel.P())
	return s
}

// Offset Offset
func (s *SelectBuilder) Offset(offset int64) *SelectBuilder {
	s.builder = s.builder.Offset(int(offset))
	return s
}

// Limit Limit
func (s *SelectBuilder) Limit(limit int64) *SelectBuilder {
	s.builder = s.builder.Limit(int(limit))
	return s
}

// OrderDesc OrderDesc
func (s *SelectBuilder) OrderDesc(field string) *SelectBuilder {
	s.builder = s.builder.OrderBy(xsql.Desc(field))
	return s
}

// OrderAsc OrderAsc
func (s *SelectBuilder) OrderAsc(field string) *SelectBuilder {
	s.builder = s.builder.OrderBy(xsql.Asc(field))
	return s
}

//  ForceIndex ForceIndex  FORCE INDEX (`index_name`)
func (s *SelectBuilder) ForceIndex(indexName ...string) *SelectBuilder {
	s.builder.ForceIndex(indexName...)
	return s
}

// GroupBy GroupBy
func (s *SelectBuilder) GroupBy(fields ...string) *SelectBuilder {
	s.builder.GroupBy(fields...)
	return s
}

// Having Having
func (s *SelectBuilder) Having(p *xsql.Predicate) *SelectBuilder {
	s.builder.Having(p)
	return s
}

// Slice Slice scan query result to slice
func (s *SelectBuilder) Slice(ctx context.Context, dstSlice interface{}) error {
	sqlstr, args := s.builder.Query()
	q, err := s.eq.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return err
	}
	defer q.Close()
	return xsql.ScanSlice(q, dstSlice)
}

// One One
func (s *SelectBuilder) One(ctx context.Context) (*User, error) {
	s.builder.Limit(1)
	results, err := s.All(ctx)
	if err != nil {
		return nil, err
	}
	if len(results) <= 0 {
		return nil, sql.ErrNoRows
	}
	return results[0], nil
}

// Int64 count or select only one int field
func (s *SelectBuilder) Int64(ctx context.Context) (int64, error) {
	sqlstr, args := s.builder.Query()
	q, err := s.eq.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return 0, err
	}
	defer q.Close()
	return xsql.ScanInt64(q)
}

// Int64s return int64 slice
func (s *SelectBuilder) Int64s(ctx context.Context) ([]int64, error) {
	sqlstr, args := s.builder.Query()
	q, err := s.eq.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var arr []int64
	for q.Next() {
		var a int64
		if err := q.Scan(&a); err != nil {
			return nil, err
		}
		arr = append(arr, a)
	}
	if q.Err() != nil {
		return nil, q.Err()
	}

	return arr, nil
}

// String  String
func (s *SelectBuilder) String(ctx context.Context) (string, error) {
	sqlstr, args := s.builder.Query()
	q, err := s.eq.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return "", err
	}
	defer q.Close()
	return xsql.ScanString(q)
}

// Strings return string slice
func (s *SelectBuilder) Strings(ctx context.Context) ([]string, error) {
	sqlstr, args := s.builder.Query()
	q, err := s.eq.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	var arr []string
	for q.Next() {
		var a string
		if err := q.Scan(&a); err != nil {
			return nil, err
		}
		arr = append(arr, a)
	}
	if q.Err() != nil {
		return nil, q.Err()
	}

	return arr, nil
}

// All  return all results
func (s *SelectBuilder) All(ctx context.Context) ([]*User, error) {
	if s.builder.SelectColumnsLen() <= 0 {
		s.builder.Select(columns...)
		s.selectAll = true
	}
	sqlstr, args := s.builder.Query()
	q, err := s.eq.QueryContext(ctx, sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	result := []*User{}
	if s.selectAll {
		for q.Next() {
			a := &User{}
			if err := q.Scan(&a.Id, &a.Name, &a.Age, &a.Ctime, &a.Mtime); err != nil {
				return nil, err
			}
			result = append(result, a)
		}
		if q.Err() != nil {
			return nil, q.Err()
		}
		return result, nil
	}

	if err := xsql.ScanSlice(q, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateBuilder UpdateBuilder
type UpdateBuilder struct {
	builder *xsql.UpdateBuilder
	eq      xsql.ExecQuerier
}

// Update return a UpdateBuilder
func Update(eq xsql.ExecQuerier) *UpdateBuilder {
	return &UpdateBuilder{
		eq:      eq,
		builder: xsql.Update(table),
	}
}

// Where Where
func (u *UpdateBuilder) Where(p ...UserWhere) *UpdateBuilder {
	s := &xsql.Selector{}
	for _, v := range p {
		v(s)
	}
	u.builder = u.builder.Where(s.P())
	return u
}

// SetId  set id
func (u *UpdateBuilder) SetId(arg int64) *UpdateBuilder {
	u.builder.Set(Id, arg)
	return u
}

// SetName  set name
func (u *UpdateBuilder) SetName(arg string) *UpdateBuilder {
	u.builder.Set(Name, arg)
	return u
}

// SetAge  set age
func (u *UpdateBuilder) SetAge(arg int64) *UpdateBuilder {
	u.builder.Set(Age, arg)
	return u
}

// AddAge  add  age set x = x + arg
func (u *UpdateBuilder) AddAge(arg interface{}) *UpdateBuilder {
	u.builder.Add(Age, arg)
	return u
}

// SetCtime  set ctime
func (u *UpdateBuilder) SetCtime(arg time.Time) *UpdateBuilder {
	u.builder.Set(Ctime, arg)
	return u
}

// SetMtime  set mtime
func (u *UpdateBuilder) SetMtime(arg time.Time) *UpdateBuilder {
	u.builder.Set(Mtime, arg)
	return u
}

// Save do a update statmxsql  if tx can without context
func (u *UpdateBuilder) Save(ctx context.Context) (int64, error) {
	up, args := u.builder.Query()
	result, err := u.eq.ExecContext(ctx, up, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

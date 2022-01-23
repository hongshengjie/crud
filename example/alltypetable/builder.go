// Code generated by bcurd. DO NOT EDIT.

package alltypetable

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
	a       []*AllTypeTable
	upsert  bool
}

// Create Create
func Create(eq xsql.ExecQuerier) *InsertBuilder {
	return &InsertBuilder{
		builder: xsql.Insert(table),
		eq:      eq,
	}
}

// SetAllTypeTable SetAllTypeTable
func (in *InsertBuilder) SetAllTypeTable(a ...*AllTypeTable) *InsertBuilder {
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
		return 0, errors.New("please set a AllTypeTable")
	}
	in.builder.Columns(Id, TInt, SInt, MInt, BInt, F32, F64, DecimalMysql, CharM, VarcharM, JsonM, NvarcharM, NcharM, TimeM, DateM, DataTimeM, TimestampM, TimestampUpdate, YearM, TText, MText, TextM, LText, BinaryM, BlobM, LBlob, MBlob, TBlob, BitM, EnumM, SetM, BoolM)
	if in.upsert {
		in.builder.OnDuplicateKeyUpdate(Id, TInt, SInt, MInt, BInt, F32, F64, DecimalMysql, CharM, VarcharM, JsonM, NvarcharM, NcharM, TimeM, DateM, DataTimeM, TimestampM, TimestampUpdate, YearM, TText, MText, TextM, LText, BinaryM, BlobM, LBlob, MBlob, TBlob, BitM, EnumM, SetM, BoolM)
	}
	for _, a := range in.a {
		if a == nil {
			return 0, errors.New("can not insert a nil AllTypeTable")
		}
		in.builder.Values(a.Id, a.TInt, a.SInt, a.MInt, a.BInt, a.F32, a.F64, a.DecimalMysql, a.CharM, a.VarcharM, a.JsonM, a.NvarcharM, a.NcharM, a.TimeM, a.DateM, a.DataTimeM, a.TimestampM, a.TimestampUpdate, a.YearM, a.TText, a.MText, a.TextM, a.LText, a.BinaryM, a.BlobM, a.LBlob, a.MBlob, a.TBlob, a.BitM, a.EnumM, a.SetM, a.BoolM)
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

// Where  AllTypeTableWhere
func (d *DeleteBuilder) Where(p ...AllTypeTableWhere) *DeleteBuilder {
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
func (s *SelectBuilder) Where(p ...AllTypeTableWhere) *SelectBuilder {
	sel := &xsql.Selector{}
	for _, v := range p {
		v(sel)
	}
	s.builder = s.builder.Where(sel.P())
	return s
}

func (s *SelectBuilder) WhereP(ps ...*xsql.Predicate) *SelectBuilder {
	for _, v := range ps {
		s.builder.Where(v)
	}
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
func (s *SelectBuilder) One(ctx context.Context) (*AllTypeTable, error) {
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

// Int64 count or select only one int64 field
func (s *SelectBuilder) Int64(ctx context.Context) (int64, error) {
	return xsql.Int64(ctx, s.builder, s.eq)
}

// Int64s return int64 slice
func (s *SelectBuilder) Int64s(ctx context.Context) ([]int64, error) {
	return xsql.Int64s(ctx, s.builder, s.eq)
}

// String  String
func (s *SelectBuilder) String(ctx context.Context) (string, error) {
	return xsql.String(ctx, s.builder, s.eq)
}

// Strings return string slice
func (s *SelectBuilder) Strings(ctx context.Context) ([]string, error) {
	return xsql.Strings(ctx, s.builder, s.eq)
}

// All  return all results
func (s *SelectBuilder) All(ctx context.Context) ([]*AllTypeTable, error) {
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
	result := []*AllTypeTable{}
	if s.selectAll {
		for q.Next() {
			a := &AllTypeTable{}
			if err := q.Scan(&a.Id, &a.TInt, &a.SInt, &a.MInt, &a.BInt, &a.F32, &a.F64, &a.DecimalMysql, &a.CharM, &a.VarcharM, &a.JsonM, &a.NvarcharM, &a.NcharM, &a.TimeM, &a.DateM, &a.DataTimeM, &a.TimestampM, &a.TimestampUpdate, &a.YearM, &a.TText, &a.MText, &a.TextM, &a.LText, &a.BinaryM, &a.BlobM, &a.LBlob, &a.MBlob, &a.TBlob, &a.BitM, &a.EnumM, &a.SetM, &a.BoolM); err != nil {
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
func (u *UpdateBuilder) Where(p ...AllTypeTableWhere) *UpdateBuilder {
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

// SetTInt  set t_int
func (u *UpdateBuilder) SetTInt(arg int64) *UpdateBuilder {
	u.builder.Set(TInt, arg)
	return u
}

// AddTInt  add  t_int set x = x + arg
func (u *UpdateBuilder) AddTInt(arg interface{}) *UpdateBuilder {
	u.builder.Add(TInt, arg)
	return u
}

// SetSInt  set s_int
func (u *UpdateBuilder) SetSInt(arg int64) *UpdateBuilder {
	u.builder.Set(SInt, arg)
	return u
}

// AddSInt  add  s_int set x = x + arg
func (u *UpdateBuilder) AddSInt(arg interface{}) *UpdateBuilder {
	u.builder.Add(SInt, arg)
	return u
}

// SetMInt  set m_int
func (u *UpdateBuilder) SetMInt(arg int64) *UpdateBuilder {
	u.builder.Set(MInt, arg)
	return u
}

// AddMInt  add  m_int set x = x + arg
func (u *UpdateBuilder) AddMInt(arg interface{}) *UpdateBuilder {
	u.builder.Add(MInt, arg)
	return u
}

// SetBInt  set b_int
func (u *UpdateBuilder) SetBInt(arg int64) *UpdateBuilder {
	u.builder.Set(BInt, arg)
	return u
}

// AddBInt  add  b_int set x = x + arg
func (u *UpdateBuilder) AddBInt(arg interface{}) *UpdateBuilder {
	u.builder.Add(BInt, arg)
	return u
}

// SetF32  set f32
func (u *UpdateBuilder) SetF32(arg float32) *UpdateBuilder {
	u.builder.Set(F32, arg)
	return u
}

// AddF32  add  f32 set x = x + arg
func (u *UpdateBuilder) AddF32(arg interface{}) *UpdateBuilder {
	u.builder.Add(F32, arg)
	return u
}

// SetF64  set f64
func (u *UpdateBuilder) SetF64(arg float64) *UpdateBuilder {
	u.builder.Set(F64, arg)
	return u
}

// AddF64  add  f64 set x = x + arg
func (u *UpdateBuilder) AddF64(arg interface{}) *UpdateBuilder {
	u.builder.Add(F64, arg)
	return u
}

// SetDecimalMysql  set decimal_mysql
func (u *UpdateBuilder) SetDecimalMysql(arg float64) *UpdateBuilder {
	u.builder.Set(DecimalMysql, arg)
	return u
}

// AddDecimalMysql  add  decimal_mysql set x = x + arg
func (u *UpdateBuilder) AddDecimalMysql(arg interface{}) *UpdateBuilder {
	u.builder.Add(DecimalMysql, arg)
	return u
}

// SetCharM  set char_m
func (u *UpdateBuilder) SetCharM(arg string) *UpdateBuilder {
	u.builder.Set(CharM, arg)
	return u
}

// SetVarcharM  set varchar_m
func (u *UpdateBuilder) SetVarcharM(arg string) *UpdateBuilder {
	u.builder.Set(VarcharM, arg)
	return u
}

// SetJsonM  set json_m
func (u *UpdateBuilder) SetJsonM(arg string) *UpdateBuilder {
	u.builder.Set(JsonM, arg)
	return u
}

// SetNvarcharM  set nvarchar_m
func (u *UpdateBuilder) SetNvarcharM(arg string) *UpdateBuilder {
	u.builder.Set(NvarcharM, arg)
	return u
}

// SetNcharM  set nchar_m
func (u *UpdateBuilder) SetNcharM(arg string) *UpdateBuilder {
	u.builder.Set(NcharM, arg)
	return u
}

// SetTimeM  set time_m
func (u *UpdateBuilder) SetTimeM(arg string) *UpdateBuilder {
	u.builder.Set(TimeM, arg)
	return u
}

// SetDateM  set date_m
func (u *UpdateBuilder) SetDateM(arg time.Time) *UpdateBuilder {
	u.builder.Set(DateM, arg)
	return u
}

// SetDataTimeM  set data_time_m
func (u *UpdateBuilder) SetDataTimeM(arg time.Time) *UpdateBuilder {
	u.builder.Set(DataTimeM, arg)
	return u
}

// SetTimestampM  set timestamp_m
func (u *UpdateBuilder) SetTimestampM(arg time.Time) *UpdateBuilder {
	u.builder.Set(TimestampM, arg)
	return u
}

// SetTimestampUpdate  set timestamp_update
func (u *UpdateBuilder) SetTimestampUpdate(arg time.Time) *UpdateBuilder {
	u.builder.Set(TimestampUpdate, arg)
	return u
}

// SetYearM  set year_m
func (u *UpdateBuilder) SetYearM(arg string) *UpdateBuilder {
	u.builder.Set(YearM, arg)
	return u
}

// SetTText  set t_text
func (u *UpdateBuilder) SetTText(arg string) *UpdateBuilder {
	u.builder.Set(TText, arg)
	return u
}

// SetMText  set m_text
func (u *UpdateBuilder) SetMText(arg string) *UpdateBuilder {
	u.builder.Set(MText, arg)
	return u
}

// SetTextM  set text_m
func (u *UpdateBuilder) SetTextM(arg string) *UpdateBuilder {
	u.builder.Set(TextM, arg)
	return u
}

// SetLText  set l_text
func (u *UpdateBuilder) SetLText(arg string) *UpdateBuilder {
	u.builder.Set(LText, arg)
	return u
}

// SetBinaryM  set binary_m
func (u *UpdateBuilder) SetBinaryM(arg []byte) *UpdateBuilder {
	u.builder.Set(BinaryM, arg)
	return u
}

// SetBlobM  set blob_m
func (u *UpdateBuilder) SetBlobM(arg []byte) *UpdateBuilder {
	u.builder.Set(BlobM, arg)
	return u
}

// SetLBlob  set l_blob
func (u *UpdateBuilder) SetLBlob(arg []byte) *UpdateBuilder {
	u.builder.Set(LBlob, arg)
	return u
}

// SetMBlob  set m_blob
func (u *UpdateBuilder) SetMBlob(arg []byte) *UpdateBuilder {
	u.builder.Set(MBlob, arg)
	return u
}

// SetTBlob  set t_blob
func (u *UpdateBuilder) SetTBlob(arg []byte) *UpdateBuilder {
	u.builder.Set(TBlob, arg)
	return u
}

// SetBitM  set bit_m
func (u *UpdateBuilder) SetBitM(arg []byte) *UpdateBuilder {
	u.builder.Set(BitM, arg)
	return u
}

// SetEnumM  set enum_m
func (u *UpdateBuilder) SetEnumM(arg string) *UpdateBuilder {
	u.builder.Set(EnumM, arg)
	return u
}

// SetSetM  set set_m
func (u *UpdateBuilder) SetSetM(arg string) *UpdateBuilder {
	u.builder.Set(SetM, arg)
	return u
}

// SetBoolM  set bool_m
func (u *UpdateBuilder) SetBoolM(arg int64) *UpdateBuilder {
	u.builder.Set(BoolM, arg)
	return u
}

// AddBoolM  add  bool_m set x = x + arg
func (u *UpdateBuilder) AddBoolM(arg interface{}) *UpdateBuilder {
	u.builder.Add(BoolM, arg)
	return u
}

// Save do a update statment  if tx can without context
func (u *UpdateBuilder) Save(ctx context.Context) (int64, error) {
	up, args := u.builder.Query()
	result, err := u.eq.ExecContext(ctx, up, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

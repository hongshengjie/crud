package service

import (
	"context"
	"errors"
	"github.com/hongshengjie/crud/xsql"
	"math"
	"strings"
	"time"

	"github.com/hongshengjie/crud/example/alltypetable"
	"github.com/hongshengjie/crud/example/alltypetable/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AllTypeTableServiceImpl AllTypeTableServiceImpl
type AllTypeTableServiceImpl struct {
	db xsql.ExecQuerier
}

func (s *AllTypeTableServiceImpl) SetDB(db xsql.ExecQuerier) {
	s.db = db
}

// CreateAllTypeTable CreateAllTypeTable
func (s *AllTypeTableServiceImpl) CreateUser(ctx context.Context, req *api.AllTypeTable) (*api.AllTypeTable, error) {

	// do some parameter check
	// if req.GetXXXX() != 0 {
	// 	return nil, errors.New(-1, "参数错误")
	// }
	a := &alltypetable.AllTypeTable{
		Id:              0,
		TInt:            req.GetTInt(),
		SInt:            req.GetSInt(),
		MInt:            req.GetMInt(),
		BInt:            req.GetBInt(),
		F32:             req.GetF32(),
		F64:             req.GetF64(),
		DecimalMysql:    req.GetDecimalMysql(),
		CharM:           req.GetCharM(),
		VarcharM:        req.GetVarcharM(),
		JsonM:           req.GetJsonM(),
		NvarcharM:       req.GetNvarcharM(),
		NcharM:          req.GetNcharM(),
		TimeM:           req.GetTimeM(),
		TimestampM:      time.Now(),
		TimestampUpdate: time.Now(),
		YearM:           req.GetYearM(),
		TText:           req.GetTText(),
		MText:           req.GetMText(),
		TextM:           req.GetTextM(),
		LText:           req.GetLText(),
		BinaryM:         req.GetBinaryM(),
		BlobM:           req.GetBlobM(),
		LBlob:           req.GetLBlob(),
		MBlob:           req.GetMBlob(),
		TBlob:           req.GetTBlob(),
		BitM:            req.GetBitM(),
		EnumM:           req.GetEnumM(),
		SetM:            req.GetSetM(),
		BoolM:           req.GetBoolM(),
	}
	var err error
	if a.DateM, err = time.ParseInLocation("2006-01-02", req.GetDateM(), time.Local); err != nil {
		return nil, err
	}
	if a.DataTimeM, err = time.ParseInLocation("2006-01-02 15:04:05", req.GetDataTimeM(), time.Local); err != nil {
		return nil, err
	}
	_, err = alltypetable.
		Create(s.db).
		SetAllTypeTable(a).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// query after create and return
	a2, err := alltypetable.
		Find(s.db).
		Where(
			alltypetable.IdEQ(a.Id),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertAllTypeTable(a2), nil
}

// DeleteAllTypeTable DeleteAllTypeTable
func (s *AllTypeTableServiceImpl) DeletesAllTypeTable(ctx context.Context, req *api.AllTypeTableId) (*emptypb.Empty, error) {
	_, err := alltypetable.
		Delete(s.db).
		Where(
			alltypetable.IdEQ(req.GetId()),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Updatealltypetable UpdateAllTypeTable
func (s *AllTypeTableServiceImpl) UpdateAllTypeTable(ctx context.Context, req *api.UpdateAllTypeTableReq) (*api.AllTypeTable, error) {

	if len(req.GetUpdateMask()) == 0 {
		return nil, errors.New("update_mask empty")
	}
	update := alltypetable.Update(s.db)
	for _, v := range req.GetUpdateMask() {
		switch v {
		case "alltypetable.t_int":
			update.SetTInt(req.GetAllTypeTable().GetTInt())
		case "alltypetable.s_int":
			update.SetSInt(req.GetAllTypeTable().GetSInt())
		case "alltypetable.m_int":
			update.SetMInt(req.GetAllTypeTable().GetMInt())
		case "alltypetable.b_int":
			update.SetBInt(req.GetAllTypeTable().GetBInt())
		case "alltypetable.f32":
			update.SetF32(req.GetAllTypeTable().GetF32())
		case "alltypetable.f64":
			update.SetF64(req.GetAllTypeTable().GetF64())
		case "alltypetable.decimal_mysql":
			update.SetDecimalMysql(req.GetAllTypeTable().GetDecimalMysql())
		case "alltypetable.char_m":
			update.SetCharM(req.GetAllTypeTable().GetCharM())
		case "alltypetable.varchar_m":
			update.SetVarcharM(req.GetAllTypeTable().GetVarcharM())
		case "alltypetable.json_m":
			update.SetJsonM(req.GetAllTypeTable().GetJsonM())
		case "alltypetable.nvarchar_m":
			update.SetNvarcharM(req.GetAllTypeTable().GetNvarcharM())
		case "alltypetable.nchar_m":
			update.SetNcharM(req.GetAllTypeTable().GetNcharM())
		case "alltypetable.time_m":
			update.SetTimeM(req.GetAllTypeTable().GetTimeM())
		case "alltypetable.date_m":
			t, err := time.ParseInLocation("2006-01-02", req.GetAllTypeTable().GetDateM(), time.Local)
			if err != nil {
				return nil, err
			}
			update.SetDateM(t)
		case "alltypetable.data_time_m":
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetAllTypeTable().GetDataTimeM(), time.Local)
			if err != nil {
				return nil, err
			}
			update.SetDataTimeM(t)
		case "alltypetable.timestamp_m":
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetAllTypeTable().GetTimestampM(), time.Local)
			if err != nil {
				return nil, err
			}
			update.SetTimestampM(t)
		case "alltypetable.timestamp_update":
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetAllTypeTable().GetTimestampUpdate(), time.Local)
			if err != nil {
				return nil, err
			}
			update.SetTimestampUpdate(t)
		case "alltypetable.year_m":
			update.SetYearM(req.GetAllTypeTable().GetYearM())
		case "alltypetable.t_text":
			update.SetTText(req.GetAllTypeTable().GetTText())
		case "alltypetable.m_text":
			update.SetMText(req.GetAllTypeTable().GetMText())
		case "alltypetable.text_m":
			update.SetTextM(req.GetAllTypeTable().GetTextM())
		case "alltypetable.l_text":
			update.SetLText(req.GetAllTypeTable().GetLText())
		case "alltypetable.binary_m":
			update.SetBinaryM(req.GetAllTypeTable().GetBinaryM())
		case "alltypetable.blob_m":
			update.SetBlobM(req.GetAllTypeTable().GetBlobM())
		case "alltypetable.l_blob":
			update.SetLBlob(req.GetAllTypeTable().GetLBlob())
		case "alltypetable.m_blob":
			update.SetMBlob(req.GetAllTypeTable().GetMBlob())
		case "alltypetable.t_blob":
			update.SetTBlob(req.GetAllTypeTable().GetTBlob())
		case "alltypetable.bit_m":
			update.SetBitM(req.GetAllTypeTable().GetBitM())
		case "alltypetable.enum_m":
			update.SetEnumM(req.GetAllTypeTable().GetEnumM())
		case "alltypetable.set_m":
			update.SetSetM(req.GetAllTypeTable().GetSetM())
		case "alltypetable.bool_m":
			update.SetBoolM(req.GetAllTypeTable().GetBoolM())
		}
	}
	_, err := update.
		Where(
			alltypetable.IdEQ(req.GetAllTypeTable().GetId()),
		).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// query after update and return
	a, err := alltypetable.
		Find(s.db).
		Where(
			alltypetable.IdEQ(req.GetAllTypeTable().GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertAllTypeTable(a), nil
}

// GetAllTypeTable GetAllTypeTable
func (s *AllTypeTableServiceImpl) GetAllTypeTable(ctx context.Context, req *api.AllTypeTableId) (*api.AllTypeTable, error) {
	a, err := alltypetable.
		Find(s.db).
		Where(
			alltypetable.IdEQ(req.GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertAllTypeTable(a), nil
}

// ListAllTypeTables ListAllTypeTables
func (s *AllTypeTableServiceImpl) ListAllTypeTables(ctx context.Context, req *api.ListAllTypeTablesReq) (*api.ListAllTypeTablesResp, error) {
	page := req.GetPage()
	size := req.GetPageSize()
	if size <= 0 {
		size = 20
	}
	offset := size * (page - 1)
	if offset < 0 {
		offset = 0
	}
	finder := alltypetable.
		Find(s.db).
		Offset(offset).
		Limit(size)

	if req.GetOrderby() != "" {
		odb := strings.TrimPrefix(req.GetOrderby(), "-")
		if odb == req.GetOrderby() {
			finder.OrderAsc(odb)
		} else {
			finder.OrderDesc(odb)
		}
	}
	counter := alltypetable.
		Find(s.db).
		Count()

	var ps []*xsql.Predicate
	for _, v := range req.GetFilter() {
		p, err := xsql.GenP(v.Field, v.Op, v.Value)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	list, err := finder.WhereP(ps...).All(ctx)
	if err != nil {
		return nil, err
	}

	count, err := counter.WhereP(ps...).Int64(ctx)
	if err != nil {
		return nil, err
	}
	pageCount := int64(math.Ceil(float64(count) / float64(size)))

	return &api.ListAllTypeTablesResp{AllTypeTables: convertAllTypeTableList(list), TotalCount: count, PageCount: pageCount}, nil
}

func convertAllTypeTable(a *alltypetable.AllTypeTable) *api.AllTypeTable {
	return &api.AllTypeTable{
		Id:              a.Id,
		TInt:            a.TInt,
		SInt:            a.SInt,
		MInt:            a.MInt,
		BInt:            a.BInt,
		F32:             a.F32,
		F64:             a.F64,
		DecimalMysql:    a.DecimalMysql,
		CharM:           a.CharM,
		VarcharM:        a.VarcharM,
		JsonM:           a.JsonM,
		NvarcharM:       a.NvarcharM,
		NcharM:          a.NcharM,
		TimeM:           a.TimeM,
		DateM:           a.DateM.Format("2006-01-02"),
		DataTimeM:       a.DataTimeM.Format("2006-01-02 15:04:05"),
		TimestampM:      a.TimestampM.Format("2006-01-02 15:04:05"),
		TimestampUpdate: a.TimestampUpdate.Format("2006-01-02 15:04:05"),
		YearM:           a.YearM,
		TText:           a.TText,
		MText:           a.MText,
		TextM:           a.TextM,
		LText:           a.LText,
		BinaryM:         a.BinaryM,
		BlobM:           a.BlobM,
		LBlob:           a.LBlob,
		MBlob:           a.MBlob,
		TBlob:           a.TBlob,
		BitM:            a.BitM,
		EnumM:           a.EnumM,
		SetM:            a.SetM,
		BoolM:           a.BoolM,
	}
}

func convertAllTypeTableList(list []*alltypetable.AllTypeTable) []*api.AllTypeTable {
	ret := make([]*api.AllTypeTable, 0, len(list))
	for _, v := range list {
		ret = append(ret, convertAllTypeTable(v))
	}
	return ret
}

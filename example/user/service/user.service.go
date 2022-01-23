package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/hongshengjie/crud/xsql"
	"math"
	"strings"
	"time"

	"github.com/hongshengjie/crud/example/user"
	"github.com/hongshengjie/crud/example/user/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserServiceImpl UserServiceImpl
type UserServiceImpl struct {
	db *sql.DB
}

// CreateUser CreateUser
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *api.User) (*api.User, error) {

	// do some parameter check
	// if req.GetXXXX() != 0 {
	// 	return nil, errors.New(-1, "参数错误")
	// }
	a := &user.User{
		Id:    0,
		Name:  req.GetName(),
		Age:   req.GetAge(),
		Ctime: time.Now(),
		Mtime: time.Now(),
	}
	var err error
	_, err = user.
		Create(s.db).
		SetUser(a).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// query after create and return
	a2, err := user.
		Find(s.db).
		Where(
			user.IdEQ(a.Id),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertUser(a2), nil
}

// DeleteUser DeleteUser
func (s *UserServiceImpl) DeletesUser(ctx context.Context, req *api.UserId) (*emptypb.Empty, error) {
	_, err := user.
		Delete(s.db).
		Where(
			user.IdEQ(req.GetId()),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Updateuser UpdateUser
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.User, error) {

	if len(req.GetUpdateMask()) == 0 {
		return nil, errors.New("update_mask empty")
	}
	update := user.Update(s.db)
	for _, v := range req.GetUpdateMask() {
		switch v {
		case "user.name":
			update.SetName(req.GetUser().GetName())
		case "user.age":
			update.SetAge(req.GetUser().GetAge())
		case "user.ctime":
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetUser().GetCtime(), time.Local)
			if err != nil {
				return nil, err
			}
			update.SetCtime(t)
		case "user.mtime":
			t, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetUser().GetMtime(), time.Local)
			if err != nil {
				return nil, err
			}
			update.SetMtime(t)
		}
	}
	_, err := update.
		Where(
			user.IdEQ(req.GetUser().GetId()),
		).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// query after update and return
	a, err := user.
		Find(s.db).
		Where(
			user.IdEQ(req.GetUser().GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertUser(a), nil
}

// GetUser GetUser
func (s *UserServiceImpl) GetUser(ctx context.Context, req *api.UserId) (*api.User, error) {
	a, err := user.
		Find(s.db).
		Where(
			user.IdEQ(req.GetId()),
		).
		One(ctx)
	if err != nil {
		return nil, err
	}
	return convertUser(a), nil
}

// ListUsers ListUsers
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *api.ListUsersReq) (*api.ListUsersResp, error) {
	page := req.GetPage()
	size := req.GetPageSize()
	if size <= 0 {
		size = 20
	}
	offset := size * (page - 1)
	if offset < 0 {
		offset = 0
	}
	finder := user.
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
	counter := user.
		Find(s.db).
		Count()

	var ps []*xsql.Predicate
	for _, v := range req.GetFilter() {
		if !xsql.VailedOp(v.Op) {
			return nil, fmt.Errorf("invalid Op %s", v.Op)
		}
		if strings.Contains(v.Field, " ") {
			return nil, fmt.Errorf("invalid field %s ", v.Field)
		}
		exp := fmt.Sprintf("`%s` %s ?", v.Field, v.Op)
		ap := xsql.ExprP(exp, v.GetValue())
		ps = append(ps, ap)
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

	return &api.ListUsersResp{Users: convertUserList(list), TotalCount: count, PageCount: pageCount}, nil
}

func convertUser(a *user.User) *api.User {
	return &api.User{
		Id:    a.Id,
		Name:  a.Name,
		Age:   a.Age,
		Ctime: a.Ctime.Format("2006-01-02 15:04:05"),
		Mtime: a.Mtime.Format("2006-01-02 15:04:05"),
	}
}

func convertUserList(list []*user.User) []*api.User {
	ret := make([]*api.User, 0, len(list))
	for _, v := range list {
		ret = append(ret, convertUser(v))
	}
	return ret
}

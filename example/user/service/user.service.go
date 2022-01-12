package service

import (
	"context"
	"database/sql"
	"errors"
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
	find := user.Find(s.db).Offset(offset).Limit(size)
	for _, v := range req.GetOrderby() {
		if strings.HasPrefix(v, "-") {
			find.OrderDesc(strings.TrimPrefix(v, "-"))
			continue
		}
		find.OrderAsc(v)
	}
	// costom filter
	// {
	// 	find.Where(user.NameContains(req.GetFilter()))
	// }
	list, err := find.All(ctx)
	if err != nil {
		return nil, err
	}
	count, err := user.
		Find(s.db).
		Count().
		Int64(ctx)
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

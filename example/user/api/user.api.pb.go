// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: user.api.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//id字段
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	//名称
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	//年龄
	Age int64 `protobuf:"varint,3,opt,name=age,proto3" json:"age,omitempty"`
	//创建时间
	Ctime string `protobuf:"bytes,4,opt,name=ctime,proto3" json:"ctime,omitempty"`
	//更新时间
	Mtime string `protobuf:"bytes,5,opt,name=mtime,proto3" json:"mtime,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_api_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetAge() int64 {
	if x != nil {
		return x.Age
	}
	return 0
}

func (x *User) GetCtime() string {
	if x != nil {
		return x.Ctime
	}
	return ""
}

func (x *User) GetMtime() string {
	if x != nil {
		return x.Mtime
	}
	return ""
}

type UserId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UserId) Reset() {
	*x = UserId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserId) ProtoMessage() {}

func (x *UserId) ProtoReflect() protoreflect.Message {
	mi := &file_user_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserId.ProtoReflect.Descriptor instead.
func (*UserId) Descriptor() ([]byte, []int) {
	return file_user_api_proto_rawDescGZIP(), []int{1}
}

func (x *UserId) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type UpdateUserReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User       *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	UpdateMask []string `protobuf:"bytes,2,rep,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
}

func (x *UpdateUserReq) Reset() {
	*x = UpdateUserReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserReq) ProtoMessage() {}

func (x *UpdateUserReq) ProtoReflect() protoreflect.Message {
	mi := &file_user_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserReq.ProtoReflect.Descriptor instead.
func (*UpdateUserReq) Descriptor() ([]byte, []int) {
	return file_user_api_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateUserReq) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *UpdateUserReq) GetUpdateMask() []string {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

type ListUsersReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//
	Page int64 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	// default 20
	PageSize int64 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// order by  for example :  [name] [-id]  -表示：倒序排序
	Orderby []string `protobuf:"bytes,4,rep,name=orderby,proto3" json:"orderby,omitempty"`
	// 过滤条件需要自定义 for example  query name has
	Filter string `protobuf:"bytes,3,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *ListUsersReq) Reset() {
	*x = ListUsersReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUsersReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUsersReq) ProtoMessage() {}

func (x *ListUsersReq) ProtoReflect() protoreflect.Message {
	mi := &file_user_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUsersReq.ProtoReflect.Descriptor instead.
func (*ListUsersReq) Descriptor() ([]byte, []int) {
	return file_user_api_proto_rawDescGZIP(), []int{3}
}

func (x *ListUsersReq) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListUsersReq) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListUsersReq) GetOrderby() []string {
	if x != nil {
		return x.Orderby
	}
	return nil
}

func (x *ListUsersReq) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

type ListUsersResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users      []*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	TotalCount int64   `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	PageCount  int64   `protobuf:"varint,3,opt,name=page_count,json=pageCount,proto3" json:"page_count,omitempty"`
}

func (x *ListUsersResp) Reset() {
	*x = ListUsersResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUsersResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUsersResp) ProtoMessage() {}

func (x *ListUsersResp) ProtoReflect() protoreflect.Message {
	mi := &file_user_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUsersResp.ProtoReflect.Descriptor instead.
func (*ListUsersResp) Descriptor() ([]byte, []int) {
	return file_user_api_proto_rawDescGZIP(), []int{4}
}

func (x *ListUsersResp) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *ListUsersResp) GetTotalCount() int64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

func (x *ListUsersResp) GetPageCount() int64 {
	if x != nil {
		return x.PageCount
	}
	return 0
}

var File_user_api_proto protoreflect.FileDescriptor

var file_user_api_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x68, 0x0a,
	0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x74, 0x69, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6d, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x18, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x4b, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x12, 0x19, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1f, 0x0a,
	0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x73, 0x6b, 0x22, 0x71,
	0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x62, 0x79, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x62, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x22, 0x6c, 0x0a, 0x0d, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x12, 0x1b, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x32,
	0xc4, 0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x1a, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x05, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x1a, 0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x0a, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x07, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x23, 0x0a, 0x0a, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x19, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x07, 0x2e, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x1a, 0x05, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x09, 0x4c, 0x69,
	0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x0d, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x42, 0x06, 0x5a, 0x04, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_api_proto_rawDescOnce sync.Once
	file_user_api_proto_rawDescData = file_user_api_proto_rawDesc
)

func file_user_api_proto_rawDescGZIP() []byte {
	file_user_api_proto_rawDescOnce.Do(func() {
		file_user_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_api_proto_rawDescData)
	})
	return file_user_api_proto_rawDescData
}

var file_user_api_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_user_api_proto_goTypes = []interface{}{
	(*User)(nil),          // 0: User
	(*UserId)(nil),        // 1: UserId
	(*UpdateUserReq)(nil), // 2: UpdateUserReq
	(*ListUsersReq)(nil),  // 3: ListUsersReq
	(*ListUsersResp)(nil), // 4: ListUsersResp
	(*emptypb.Empty)(nil), // 5: google.protobuf.Empty
}
var file_user_api_proto_depIdxs = []int32{
	0, // 0: UpdateUserReq.user:type_name -> User
	0, // 1: ListUsersResp.users:type_name -> User
	0, // 2: UserService.CreateUser:input_type -> User
	1, // 3: UserService.DeleteUser:input_type -> UserId
	2, // 4: UserService.UpdateUser:input_type -> UpdateUserReq
	1, // 5: UserService.GetUser:input_type -> UserId
	3, // 6: UserService.ListUsers:input_type -> ListUsersReq
	0, // 7: UserService.CreateUser:output_type -> User
	5, // 8: UserService.DeleteUser:output_type -> google.protobuf.Empty
	0, // 9: UserService.UpdateUser:output_type -> User
	0, // 10: UserService.GetUser:output_type -> User
	4, // 11: UserService.ListUsers:output_type -> ListUsersResp
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_user_api_proto_init() }
func file_user_api_proto_init() {
	if File_user_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_user_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_user_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_user_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUsersReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_user_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUsersResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_api_proto_goTypes,
		DependencyIndexes: file_user_api_proto_depIdxs,
		MessageInfos:      file_user_api_proto_msgTypes,
	}.Build()
	File_user_api_proto = out.File
	file_user_api_proto_rawDesc = nil
	file_user_api_proto_goTypes = nil
	file_user_api_proto_depIdxs = nil
}

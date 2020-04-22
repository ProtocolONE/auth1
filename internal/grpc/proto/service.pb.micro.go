// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: internal/grpc/proto/service.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Service service

func NewServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Service service

type Service interface {
	GetProfile(ctx context.Context, in *GetProfileRequest, opts ...client.CallOption) (*ProfileResponse, error)
	SetProfile(ctx context.Context, in *SetProfileRequest, opts ...client.CallOption) (*ProfileResponse, error)
	//
	SetPassword(ctx context.Context, in *SetPasswordRequest, opts ...client.CallOption) (*SetPasswordResponse, error)
	//
	GetUserIdentities(ctx context.Context, in *GetUserIdentitiesRequest, opts ...client.CallOption) (*UserIdentitiesResponse, error)
}

type service struct {
	c    client.Client
	name string
}

func NewService(name string, c client.Client) Service {
	return &service{
		c:    c,
		name: name,
	}
}

func (c *service) GetProfile(ctx context.Context, in *GetProfileRequest, opts ...client.CallOption) (*ProfileResponse, error) {
	req := c.c.NewRequest(c.name, "Service.GetProfile", in)
	out := new(ProfileResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) SetProfile(ctx context.Context, in *SetProfileRequest, opts ...client.CallOption) (*ProfileResponse, error) {
	req := c.c.NewRequest(c.name, "Service.SetProfile", in)
	out := new(ProfileResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) SetPassword(ctx context.Context, in *SetPasswordRequest, opts ...client.CallOption) (*SetPasswordResponse, error) {
	req := c.c.NewRequest(c.name, "Service.SetPassword", in)
	out := new(SetPasswordResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *service) GetUserIdentities(ctx context.Context, in *GetUserIdentitiesRequest, opts ...client.CallOption) (*UserIdentitiesResponse, error) {
	req := c.c.NewRequest(c.name, "Service.GetUserIdentities", in)
	out := new(UserIdentitiesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Service service

type ServiceHandler interface {
	GetProfile(context.Context, *GetProfileRequest, *ProfileResponse) error
	SetProfile(context.Context, *SetProfileRequest, *ProfileResponse) error
	//
	SetPassword(context.Context, *SetPasswordRequest, *SetPasswordResponse) error
	//
	GetUserIdentities(context.Context, *GetUserIdentitiesRequest, *UserIdentitiesResponse) error
}

func RegisterServiceHandler(s server.Server, hdlr ServiceHandler, opts ...server.HandlerOption) error {
	type service interface {
		GetProfile(ctx context.Context, in *GetProfileRequest, out *ProfileResponse) error
		SetProfile(ctx context.Context, in *SetProfileRequest, out *ProfileResponse) error
		SetPassword(ctx context.Context, in *SetPasswordRequest, out *SetPasswordResponse) error
		GetUserIdentities(ctx context.Context, in *GetUserIdentitiesRequest, out *UserIdentitiesResponse) error
	}
	type Service struct {
		service
	}
	h := &serviceHandler{hdlr}
	return s.Handle(s.NewHandler(&Service{h}, opts...))
}

type serviceHandler struct {
	ServiceHandler
}

func (h *serviceHandler) GetProfile(ctx context.Context, in *GetProfileRequest, out *ProfileResponse) error {
	return h.ServiceHandler.GetProfile(ctx, in, out)
}

func (h *serviceHandler) SetProfile(ctx context.Context, in *SetProfileRequest, out *ProfileResponse) error {
	return h.ServiceHandler.SetProfile(ctx, in, out)
}

func (h *serviceHandler) SetPassword(ctx context.Context, in *SetPasswordRequest, out *SetPasswordResponse) error {
	return h.ServiceHandler.SetPassword(ctx, in, out)
}

func (h *serviceHandler) GetUserIdentities(ctx context.Context, in *GetUserIdentitiesRequest, out *UserIdentitiesResponse) error {
	return h.ServiceHandler.GetUserIdentities(ctx, in, out)
}

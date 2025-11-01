package user

import (
	"context"
	"fmt"

	server "go-interview-practice/challenge14/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---- proto types ----

// User represents a user in the system
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

// Request/Response types (normally generated from .proto)
type GetUserRequest struct {
	UserID int64 `json:"user_id"`
}

type GetUserResponse struct {
	User *User `json:"user"`
}

// ---- server ----

const (
	userServiceName   = "UserService"
	getUserMethodName = "GetUser"
)

// UserService interface
type UserService interface {
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
}

// userServer implements the UserService
type userServer struct {
	users map[int64]*User
}

// newUserService creates a new UserService
func newUserService() UserService {
	users := map[int64]*User{
		1: {ID: 1, Username: "alice", Email: "alice@example.com", Active: true},
		2: {ID: 2, Username: "bob", Email: "bob@example.com", Active: true},
		3: {ID: 3, Username: "charlie", Email: "charlie@example.com", Active: false},
	}
	return &userServer{users: users}
}

// GetUser retrieves a user by ID
func (s *userServer) GetUser(
	_ context.Context,
	req *GetUserRequest,
) (*GetUserResponse, error) {
	user, exists := s.users[req.UserID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.UserID)
	}
	if !user.Active {
		return nil, status.Errorf(codes.PermissionDenied, "user %d is not active", req.UserID)
	}
	return &GetUserResponse{user}, nil
}

// gRPC service registration helpers
func RegisterUserService(srv *server.GRPCServer) error {
	srv.Server.RegisterService(&grpc.ServiceDesc{
		ServiceName: userServiceName,
		HandlerType: (*UserService)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: getUserMethodName,
				Handler:    userServiceGetUserHandler,
			},
		},
	}, newUserService())
	return nil
}

//nolint:revive // context-as-argument: This is the signature required for gRPC registration
func userServiceGetUserHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	service, ok := srv.(UserService)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"server type %T does not implement UserService",
			srv,
		)
	}
	req := new(GetUserRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return service.GetUser(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: fmt.Sprintf("/%s/%s", userServiceName, getUserMethodName),
	}
	handler := func(ctx context.Context, req any) (any, error) {
		r, ok := req.(*GetUserRequest)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid request type %T", req)
		}
		return service.GetUser(ctx, r)
	}
	return interceptor(ctx, req, info, handler)
}

// ---- client ----

type userServiceClient struct {
	conn *grpc.ClientConn
}

func NewUserServiceClient(conn *grpc.ClientConn) UserService {
	return &userServiceClient{conn: conn}
}

func (cl *userServiceClient) GetUser(
	ctx context.Context,
	req *GetUserRequest,
) (*GetUserResponse, error) {
	var resp GetUserResponse
	err := cl.conn.Invoke(
		ctx,
		fmt.Sprintf("/%s/%s", userServiceName, getUserMethodName),
		req,
		&resp,
	)
	return &resp, err
}

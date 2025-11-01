package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"strings"

	codec "go-interview-practice/challenge14/codec"
	product "go-interview-practice/challenge14/product"
	server "go-interview-practice/challenge14/server"
	user "go-interview-practice/challenge14/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// +---------------------+------------------------------+--------------------------------------+
// |        Step         |      Generated version       |            Manual version            |
// +---------------------+------------------------------+--------------------------------------+
// | .proto              | Defines service and rpc      | We define ProductService interface   |
// | protoc-gen-go       | Generates pb.go with handler | We handwrite RegisterProductService, |
// |					 | + client                     | handler, and Invoke helper.		   |
// | pb.RegisterServer() | Registers methods to gRPC    | We call s.RegisterService() manually |
// | pb.NewClient()      | Creates typed client         | We use conn.Invoke() manually        |
// +---------------------+---------------------------------------+-----------------------------+

// ---- proto types ----

// `json:` is a struct tag — a string literal attached to a struct field.
// It provides metadata used by reflection-based libraries like encoding/json.
// It tells Go's JSON encoder/decoder that when converting to or from JSON:
//   - The Go field ProductId
//   - Should appear in JSON as "product_id" (snake_case)
//   - Instead of the default "ProductId" (camel-case from the Go field name).

// The gRPC server and client serialize and deserialize the request/response using
// the gRPC codec (by default protobuf, unless we explicitly use a JSON codec).
// The json tags are completely ignored by gRPC.

// Order represents an order in the system
type Order struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Total     float64 `json:"total"`
}

type GetOrderRequest struct {
	OrderID int64 `json:"order_id"`
}

type GetOrderResponse struct {
	Order *Order `json:"order"`
}

type CreateOrderRequest struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type CreateOrderResponse struct {
	Order *Order `json:"order"`
}

// ---- server ----

const (
	orderServiceName      = "OrderService"
	getOrderMethodName    = "GetOrder"
	createOrderMethodName = "CreateOrder"
	bearerToken           = "Bearer token123"
)

// OrderService interface
type OrderService interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error)
}

// orderServer implements the OrderService
type orderServer struct {
	orders         map[int64]*Order
	userService    user.UserService
	productService product.ProductService
}

// newOrderService creates a new OrderService
func newOrderService(
	userService user.UserService,
	productService product.ProductService,
) OrderService {
	return &orderServer{
		orders:         make(map[int64]*Order),
		userService:    userService,
		productService: productService,
	}
}

// CreateOrder creates a new order
func (srv *orderServer) CreateOrder(
	ctx context.Context,
	req *CreateOrderRequest,
) (*CreateOrderResponse, error) {
	// Hint: 1. Validate user, 2. Get product and check inventory, 3. Create order
	_, err := srv.userService.GetUser(ctx, &user.GetUserRequest{UserID: req.UserID})
	if err != nil {
		return nil, err
	}
	product, err := srv.productService.GetProduct(
		ctx,
		&product.GetProductRequest{ProductID: req.ProductID},
	)
	if err != nil {
		return nil, err
	}
	if product.Product.Inventory < req.Quantity {
		return nil, errors.New("insufficient inventory")
	}
	order := Order{
		ID:        rand.Int64(),
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Total:     float64(req.Quantity) * product.Product.Price,
	}
	srv.orders[order.ID] = &order
	return &CreateOrderResponse{&order}, nil
}

// GetOrder retrieves an order by ID
func (srv *orderServer) GetOrder(
	_ context.Context,
	req *GetOrderRequest,
) (*GetOrderResponse, error) {
	order, exists := srv.orders[req.OrderID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "order %d not found", req.OrderID)
	}
	return &GetOrderResponse{order}, nil
}

func RegisterOrderService(srv *server.GRPCServer) error {
	conn, err := grpc.NewClient(
		srv.Listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.JSONCodec())),
	)
	if err != nil {
		return err
	}
	srv.AddCleanupFn(func() {
		err := conn.Close()
		if err != nil {
			log.Println(`failed to close connection to 
UserService and ProductService from OrderService`)
		}
	})
	srv.Server.RegisterService(&grpc.ServiceDesc{
		ServiceName: orderServiceName,
		HandlerType: (*OrderService)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: getOrderMethodName,
				Handler:    orderServiceGetOrderHandler,
			},
			{
				MethodName: createOrderMethodName,
				Handler:    orderServiceCreateOrderHandler,
			},
		},
	}, newOrderService(user.NewUserServiceClient(conn), product.NewProductServiceClient(conn)))
	return nil
}

// When we compile a .proto file using the official plugin:
//
//	protoc --go-grpc_out=. product.proto
//
// The generated code includes one handler function per RPC:
//
//	_<ServiceName>_<MethodName>_Handler

// authServerInterceptor verifies the "authorization" header for a Bearer token.
func authInterceptor(expectedToken string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Extract metadata from the incoming context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Check the Authorization header
		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		token := authHeaders[0]
		if !strings.HasPrefix(token, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		if token != expectedToken {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Token is valid, continue to the actual RPC handler
		return handler(ctx, req)
	}
}

//nolint:revive // context-as-argument: This is the signature required for gRPC registration
func orderServiceGetOrderHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	globalInterceptor grpc.UnaryServerInterceptor,
) (any, error) {
	service, ok := srv.(OrderService)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"server type %T does not implement OrderService",
			srv,
		)
	}
	req := new(GetOrderRequest)
	if err := dec(req); err != nil {
		return nil, err
	}

	handler := func(ctx context.Context, req any) (any, error) {
		r, ok := req.(*GetOrderRequest)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid request type %T", req)
		}
		return service.GetOrder(ctx, r)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: fmt.Sprintf("/%s/%s", orderServiceName, getOrderMethodName),
	}

	return chainUnaryInterceptors(ctx, req, info, handler, globalInterceptor)
}

func chainUnaryInterceptors(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	globalInterceptor grpc.UnaryServerInterceptor,
) (any, error) {
	// Wrap the handler with service-specific interceptor
	wrappedHandler := func(ctx context.Context, req any) (any, error) {
		res, err := authInterceptor(bearerToken)(ctx, req, info, handler)
		return res, err
	}

	// If global interceptor exists, pass the wrappedHandler to it
	if globalInterceptor != nil {
		return globalInterceptor(ctx, req, info, wrappedHandler)
	}

	// No global interceptor; just call wrappedHandler
	return wrappedHandler(ctx, req)
}

// Helper to chain multiple unary interceptors.
// func chainUnaryInterceptors(
// 	interceptors ...grpc.UnaryServerInterceptor,
// ) grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req any,
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (any, error) {
// 		// Start with the actual service handler
// 		chained := handler
// 		// Wrap interceptors in reverse order so the first in the slice runs first
// 		for i := len(interceptors) - 1; i >= 0; i-- {
// 			ic := interceptors[i]
// 			next := chained
// 			chained = func(ctx context.Context, req any) (any, error) {
// 				return ic(ctx, req, info, next)
// 			}
// 		}
// 		return chained(ctx, req)
// 	}
// }

//nolint:revive // context-as-argument: This is the signature required for gRPC registration
func orderServiceCreateOrderHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	globalInterceptor grpc.UnaryServerInterceptor,
) (any, error) {
	service, ok := srv.(OrderService)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"server type %T does not implement OrderService",
			srv,
		)
	}
	req := new(CreateOrderRequest)
	if err := dec(req); err != nil {
		return nil, err
	}

	handler := func(ctx context.Context, req any) (any, error) {
		r, ok := req.(*CreateOrderRequest)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid request type %T", req)
		}
		return service.CreateOrder(ctx, r)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: fmt.Sprintf("/%s/%s", orderServiceName, createOrderMethodName),
	}

	return chainUnaryInterceptors(ctx, req, info, handler, globalInterceptor)
}

// ---- client ----

type orderServiceClient struct {
	conn *grpc.ClientConn
}

func NewOrderServiceClient(conn *grpc.ClientConn) OrderService {
	return &orderServiceClient{conn}
}

func (cl *orderServiceClient) CreateOrder(
	ctx context.Context,
	req *CreateOrderRequest,
) (*CreateOrderResponse, error) {
	var resp CreateOrderResponse
	err := cl.conn.Invoke(
		addAuthToken(ctx),
		fmt.Sprintf("/%s/%s", orderServiceName, createOrderMethodName),
		req,
		&resp,
	)
	return &resp, err
}

// GetOrder retrieves an order by ID
func (cl *orderServiceClient) GetOrder(
	ctx context.Context,
	req *GetOrderRequest,
) (*GetOrderResponse, error) {
	var resp GetOrderResponse
	err := cl.conn.Invoke(
		addAuthToken(ctx),
		fmt.Sprintf("/%s/%s", orderServiceName, getOrderMethodName),
		req,
		&resp,
	)
	return &resp, err
}

func addAuthToken(ctx context.Context) context.Context {
	md := metadata.Pairs("authorization", bearerToken)
	return metadata.NewOutgoingContext(ctx, md)
}

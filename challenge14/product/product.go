package product

import (
	"context"
	"fmt"

	server "go-interview-practice/challenge14/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Product represents a product in the catalog
type Product struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Inventory int32   `json:"inventory"`
}

// Request/Response types

type GetProductRequest struct {
	ProductID int64 `json:"product_id"`
}

type GetProductResponse struct {
	Product *Product `json:"product"`
}

// ---- server ----

const (
	productServiceName   = "ProductService"
	getProductMethodName = "GetProduct"
)

// ProductService interface
type ProductService interface {
	GetProduct(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error)
}

// productServer implements the ProductService
type productServer struct {
	products map[int64]*Product
}

// newProductService creates a new ProductService
func newProductService() ProductService {
	products := map[int64]*Product{
		1: {ID: 1, Name: "Laptop", Price: 999.99, Inventory: 10},
		2: {ID: 2, Name: "Phone", Price: 499.99, Inventory: 20},
		3: {ID: 3, Name: "Headphones", Price: 99.99, Inventory: 1},
	}
	return &productServer{products: products}
}

// GetProduct retrieves a product by ID
func (srv *productServer) GetProduct(
	_ context.Context,
	req *GetProductRequest,
) (*GetProductResponse, error) {
	// Hint: check if product exists, return product or gRPC NotFound error
	p, exists := srv.products[req.ProductID]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "product %d not found", req.ProductID)
	}
	return &GetProductResponse{Product: p}, nil
}

func RegisterProductService(srv *server.GRPCServer) error {
	srv.Server.RegisterService(&grpc.ServiceDesc{
		ServiceName: productServiceName,
		HandlerType: (*ProductService)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: getProductMethodName,
				Handler:    productServiceGetProductHandler,
			},
		},
	}, newProductService())
	return nil
}

//nolint:revive // context-as-argument: This is the signature required for gRPC registration
func productServiceGetProductHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	service, ok := srv.(ProductService)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"server type %T does not implement ProductService",
			srv,
		)
	}
	req := new(GetProductRequest)
	if err := dec(req); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return service.GetProduct(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: fmt.Sprintf("/%s/%s", productServiceName, getProductMethodName),
	}
	handler := func(ctx context.Context, req any) (any, error) {
		r, ok := req.(*GetProductRequest)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid request type %T", req)
		}
		return service.GetProduct(ctx, r)
	}
	return interceptor(ctx, req, info, handler)
}

// ---- client ----

type productServiceClient struct {
	conn *grpc.ClientConn
}

func NewProductServiceClient(conn *grpc.ClientConn) ProductService {
	return &productServiceClient{conn: conn}
}

func (cl *productServiceClient) GetProduct(
	ctx context.Context,
	req *GetProductRequest,
) (*GetProductResponse, error) {
	var resp GetProductResponse
	err := cl.conn.Invoke(
		ctx,
		fmt.Sprintf("/%s/%s", productServiceName, getProductMethodName),
		req,
		&resp,
	)
	return &resp, err
}

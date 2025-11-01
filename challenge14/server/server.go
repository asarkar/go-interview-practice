package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	codec "go-interview-practice/challenge14/codec"
)

// LoggingInterceptor is a server interceptor for logging
func loggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	log.Printf("request received: %s", info.FullMethod)
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("request completed: %s in %v", info.FullMethod, time.Since(start))
	return resp, err
}

type GRPCServer struct {
	Server     *grpc.Server
	Listener   net.Listener
	cleanupFns []func()
}

// AuthInterceptor is a client interceptor for authentication
// func authInterceptor(
// 	ctx context.Context,
// 	method string,
// 	req any,
// 	reply any,
// 	cc *grpc.ClientConn,
// 	invoker grpc.UnaryInvoker,
// 	opts ...grpc.CallOption,
// ) error {
// 	// Add auth token to metadata
// 	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer token123")
// 	return invoker(ctx, method, req, reply, cc, opts...)
// }

// ServiceRegistrar defines a function that registers a service on a gRPC server.
type ServiceRegistrar func(*GRPCServer) error

// NewGRPCServer creates a gRPC server and registers the provided services.
// It returns a server ready to Serve, but does not block.
func NewGRPCServer(registrars ...ServiceRegistrar) (*GRPCServer, error) {
	// Hint: create listener, gRPC server with interceptor, register service, serve
	lc := net.ListenConfig{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lis, err := lc.Listen(ctx, "tcp", "127.0.0.1:0") // start on random port
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
		grpc.ForceServerCodec(codec.JSONCodec()),
	)

	srv := &GRPCServer{
		Server:   grpcServer,
		Listener: lis,
	}

	if len(registrars) == 0 {
		log.Println("expected at least one ServiceRegistrar, got none")
	}

	// Register all services
	for _, register := range registrars {
		if err := register(srv); err != nil {
			return nil, err
		}
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Panicf("failed to start gRPC Server: %v", err)
		}
		log.Printf("started gRPC Server on port %s", lis.Addr().String())
	}()

	return srv, nil
}

func (srv *GRPCServer) Shutdown() {
	// Gracefully stop gRPC
	srv.Server.GracefulStop()

	// Run cleanup hooks
	for _, fn := range srv.cleanupFns {
		fn()
	}

	log.Println("shutdown complete.")
}

func (srv *GRPCServer) AddCleanupFn(fn func()) {
	srv.cleanupFns = append(srv.cleanupFns, fn)
}

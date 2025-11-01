package product

import (
	"context"
	"log"
	"os"
	"testing"

	codec "go-interview-practice/challenge14/codec"
	server "go-interview-practice/challenge14/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestMain(m *testing.M) {
	// Setup
	srv, err := server.NewGRPCServer(
		RegisterProductService,
	)
	if err != nil {
		log.Panic(err)
	}

	conn, err := grpc.NewClient(
		srv.Listener.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.JSONCodec())),
	)
	if err != nil {
		log.Panic(err)
	}

	srv.AddCleanupFn(func() {
		err := conn.Close()
		if err != nil {
			log.Println("failed to close connection to ProductService from test")
		}
	})

	client = NewProductServiceClient(conn)

	// Run tests
	code := m.Run()

	// Teardown
	srv.Shutdown()

	os.Exit(code)
}

var client ProductService

func TestProductService(t *testing.T) {
	t.Run("GetProduct", func(t *testing.T) {
		t.Run("Existing", func(t *testing.T) {
			resp, err := client.GetProduct(context.Background(), &GetProductRequest{1})
			if err != nil {
				t.Errorf("GetProduct failed: %v", err)
			} else {
				if resp.Product.ID != 1 || resp.Product.Name != "Laptop" {
					t.Errorf(
						"Expected product with ID 1 and name 'Laptop', got ID %d and name '%s'",
						resp.Product.ID, resp.Product.Name)
				}
			}
		})
		t.Run("NonExistent", func(t *testing.T) {
			_, err := client.GetProduct(context.Background(), &GetProductRequest{999})
			if err == nil {
				t.Errorf("Expected error for non-existent product, got nil")
			} else if status.Code(err) != codes.NotFound {
				t.Errorf("Expected NotFound error, got %v", err)
			}
		})
	})
}

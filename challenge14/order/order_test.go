package order

import (
	"context"
	"log"
	"math/rand/v2"
	"os"
	"testing"

	codec "go-interview-practice/challenge14/codec"
	product "go-interview-practice/challenge14/product"
	server "go-interview-practice/challenge14/server"
	user "go-interview-practice/challenge14/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestMain(m *testing.M) {
	// Setup
	srv, err := server.NewGRPCServer(
		user.RegisterUserService,
		product.RegisterProductService,
		RegisterOrderService,
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
			log.Println("failed to close connection to OrderService from test")
		}
	})

	client = NewOrderServiceClient(conn)

	// Run tests
	code := m.Run()

	// Teardown
	srv.Shutdown()

	os.Exit(code)
}

var client OrderService

func TestOrderService(t *testing.T) {
	t.Run("CreateOrder", func(t *testing.T) {
		t.Run("ValidUserAndProduct", func(t *testing.T) {
			resp, err := client.CreateOrder(
				context.Background(),
				&CreateOrderRequest{UserID: 1, ProductID: 1, Quantity: 2})
			if err != nil {
				t.Errorf("CreateOrder failed: %v", err)
			} else {
				order := resp.Order
				if order.UserID != 1 || order.ProductID != 1 || order.Quantity != 2 {
					t.Errorf(`expected order with UserID 1, ProductID 1, and Quantity 2, 
got UserID %d, ProductID %d, and Quantity %d`,
						order.UserID, order.ProductID, order.Quantity)
				}
				if order.Total != 1999.98 {
					t.Errorf("expected total 1999.98, got %f", order.Total)
				}
			}
		})

		t.Run("InactiveUser", func(t *testing.T) {
			_, err := client.CreateOrder(
				context.Background(),
				&CreateOrderRequest{UserID: 3, ProductID: 1, Quantity: 2})
			if err == nil {
				t.Errorf("Expected error for inactive user, got nil")
			}
		})

		t.Run("NonExistentUser", func(t *testing.T) {
			_, err := client.CreateOrder(
				context.Background(),
				&CreateOrderRequest{UserID: 999, ProductID: 1, Quantity: 2})
			if err == nil {
				t.Errorf("Expected error for non-existent user, got nil")
			}
		})

		t.Run("InsufficientInventory", func(t *testing.T) {
			_, err := client.CreateOrder(
				context.Background(),
				&CreateOrderRequest{UserID: 1, ProductID: 1, Quantity: 15})
			if err == nil {
				t.Errorf("Expected error for insufficient inventory, got nil")
			}
		})

		t.Run("NonExistentProduct", func(t *testing.T) {
			_, err := client.CreateOrder(
				context.Background(),
				&CreateOrderRequest{UserID: 3, ProductID: 999, Quantity: 2})
			if err == nil {
				t.Errorf("Expected error for non-existent product, got nil")
			}
		})
	})

	t.Run("GetOrder", func(t *testing.T) {
		t.Run("New", func(t *testing.T) {
			coResp, err := client.CreateOrder(
				context.Background(),
				&CreateOrderRequest{UserID: 1, ProductID: 1, Quantity: 2})
			if err != nil {
				t.Errorf("CreateOrder failed: %v", err)
			}
			resp, err := client.GetOrder(context.Background(), &GetOrderRequest{coResp.Order.ID})
			if err != nil {
				t.Errorf("GetOrder failed: %v", err)
			} else {
				expected := Order{
					ID:        coResp.Order.ID,
					UserID:    1,
					ProductID: 1,
					Quantity:  2,
					Total:     999.99 * 2,
				}
				if *resp.Order != expected {
					t.Errorf("want %v, got %v", *resp.Order, expected)
				}
			}
		})
		t.Run("NonExistent", func(t *testing.T) {
			_, err := client.GetOrder(context.Background(), &GetOrderRequest{rand.Int64()})
			if err == nil {
				t.Errorf("Expected error for non-existent order, got nil")
			} else if status.Code(err) != codes.NotFound {
				t.Errorf("Expected NotFound error, got %v", err)
			}
		})
	})
}

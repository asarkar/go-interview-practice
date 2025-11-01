package user

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
		RegisterUserService,
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
			log.Println("failed to close connection to UserService from test")
		}
	})

	client = NewUserServiceClient(conn)

	// Run tests
	code := m.Run()

	// Teardown
	srv.Shutdown()

	os.Exit(code)
}

var client UserService

func TestUserService(t *testing.T) {
	t.Run("GetUser", func(t *testing.T) {
		t.Run("Existing", func(t *testing.T) {
			resp, err := client.GetUser(context.Background(), &GetUserRequest{1})
			if err != nil {
				t.Errorf("GetUser failed: %v", err)
			} else {
				if resp.User.ID != 1 || resp.User.Username != "alice" {
					t.Errorf(
						`Expected user with ID 1 and username 'alice', 
got ID %d and username '%s'`,
						resp.User.ID, resp.User.Username)
				}
			}
		})

		t.Run("Inactive", func(t *testing.T) {
			_, err := client.GetUser(context.Background(), &GetUserRequest{3})
			if err == nil {
				t.Errorf("Expected error for inactive user, got nil")
			} else if status.Code(err) != codes.PermissionDenied {
				t.Errorf("Expected PermissionDenied error, got %v", err)
			}
		})

		t.Run("NonExistent", func(t *testing.T) {
			_, err := client.GetUser(context.Background(), &GetUserRequest{999})
			if err == nil {
				t.Errorf("Expected error for non-existent user, got nil")
			} else if status.Code(err) != codes.NotFound {
				t.Errorf("Expected NotFound error, got %v", err)
			}
		})
	})
}

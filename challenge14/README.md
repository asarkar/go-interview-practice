[View the Scoreboard](SCOREBOARD.md)

# Challenge 14: Microservices with gRPC

In this challenge, you will implement a microservice architecture using gRPC concepts for service-to-service communication. You'll create a User service and a Product service that communicate to provide an order management system.

## Learning Objectives

- Understand microservices architecture principles
- Learn gRPC concepts and error handling
- Implement service-to-service communication
- Practice with gRPC interceptors for cross-cutting concerns
- Handle network-based service interactions

## Requirements

Implement the following TODO methods in `solution-template.go`:

### 1. Service Business Logic
- `UserServiceServer.GetUser()`: Retrieve user by ID with proper error handling
- `UserServiceServer.ValidateUser()`: Check if user exists and is active
- `ProductServiceServer.GetProduct()`: Retrieve product by ID with proper error handling  
- `ProductServiceServer.CheckInventory()`: Check product availability

### 2. Order Service
- `OrderService.CreateOrder()`: Orchestrate user validation, product checking, and order creation

### 3. gRPC Infrastructure
- `StartUserService()`: Set up and start the user service with interceptors
- `StartProductService()`: Set up and start the product service with interceptors
- `ConnectToServices()`: Create clients and connect to both services

### 4. gRPC Clients
- `UserServiceClient.GetUser()`: Make gRPC calls to user service
- `UserServiceClient.ValidateUser()`: Make gRPC calls for user validation
- `ProductServiceClient.GetProduct()`: Make gRPC calls to product service
- `ProductServiceClient.CheckInventory()`: Make gRPC calls for inventory checking

### 5. Interceptors
- `LoggingInterceptor`: Log method calls and execution time
- `AuthInterceptor`: Add authentication metadata to requests

## Key Concepts Covered

- **gRPC Status Codes**: Use appropriate status codes (NotFound, PermissionDenied, etc.)
- **Service Interfaces**: Clean separation between service logic and transport
- **Network Communication**: Real service-to-service calls
- **Error Propagation**: Proper error handling across service boundaries
- **Interceptors**: Cross-cutting concerns like logging and authentication
- **Microservices Patterns**: Service orchestration and coordination

## Hints

- Use `status.Errorf()` to return proper gRPC errors
- Remember to handle both business logic errors and network errors
- Services should validate inputs and return appropriate error codes
- The OrderService orchestrates calls to both User and Product services
- Interceptors wrap all service calls for logging/auth

## Running Tests

```bash
go test -v
```

The tests will:
1. Start individual services and test their functionality
2. Test service-to-service communication through OrderService
3. Verify proper error handling for various scenarios
4. Check that interceptors work correctly

## Example gRPC Status Codes

- `codes.OK`: Successful operation
- `codes.NotFound`: Resource doesn't exist
- `codes.PermissionDenied`: User not authorized
- `codes.ResourceExhausted`: Insufficient inventory
- `codes.InvalidArgument`: Invalid input parameters

Complete all TODO methods to make the tests pass! 
# Learning Materials for Bank Account with Error Handling

## Error Handling in Go

Error handling is a critical aspect of writing robust Go programs. This challenge focuses on implementing a banking system with proper error handling techniques.

### Basic Error Handling

Go uses explicit error handling with return values instead of exceptions:

```go
// Function that may return an error
func divideNumbers(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Error handling with if checks
result, err := divideNumbers(10, 0)
if err != nil {
    fmt.Println("Error:", err)
    return // or handle the error
}
fmt.Println("Result:", result)
```

### Creating Custom Errors

Standard ways to create errors:

```go
// Using errors.New for simple error messages
err := errors.New("insufficient funds")

// Using fmt.Errorf for formatted error messages
amount := 100
err := fmt.Errorf("insufficient funds: need $%d more", amount)
```

### Custom Error Types

Creating custom error types allows for more detailed error handling:

```go
// Define a custom error type
type InsufficientFundsError struct {
    Balance float64
    Amount  float64
}

// Implement the error interface
func (e *InsufficientFundsError) Error() string {
    return fmt.Sprintf("insufficient funds: balance $%.2f, attempted to withdraw $%.2f", 
        e.Balance, e.Amount)
}
```

**Key concepts for custom errors:**
- Implement the `Error() string` method
- Include relevant context in error fields
- Use pointer receivers when checking error types
- Type assertion with `ok` pattern for error type checking

### Error Wrapping (Go 1.13+)

Go 1.13 introduced error wrapping for better error context:

**Key concepts:**
- Use `fmt.Errorf` with `%w` verb to wrap errors
- Use `errors.As()` to check for specific error types in a chain
- Use `errors.Is()` to check for specific error values in a chain
- Preserve original error context while adding meaningful information

### Sentinel Errors

Predefined errors that can be compared directly:

```go
// Define sentinel errors as package-level variables
var (
    ErrAccountNotFound   = errors.New("account not found")
    ErrInsufficientFunds = errors.New("insufficient funds")
    ErrInvalidAmount     = errors.New("invalid amount")
)
```

**Key concepts:**
- Use package-level variables for reusable errors
- Compare errors using `==` or `errors.Is()`
- Provide clear, descriptive error messages

### Error Handling Patterns

#### 1. Return Early Pattern
Validate inputs first and return errors immediately to avoid deep nesting.

#### 2. Error Handler Functions
Create functions that can handle multiple error-prone operations in sequence.

#### 3. Error Context
Always provide meaningful context when returning or wrapping errors.

### Banking Application Specific Considerations

#### Account Operations
- **Balance validation**: Check sufficient funds before withdrawal
- **Amount validation**: Ensure positive amounts for deposits/withdrawals
- **Account existence**: Verify account exists before operations
- **Input sanitization**: Validate all user inputs

#### Error Types for Banking
- **InsufficientFundsError**: Specific error for balance issues
- **InvalidAmountError**: For negative or zero amounts
- **AccountNotFoundError**: When account lookup fails
- **ValidationError**: For input validation failures

### Thread Safety in Banking Applications

Banking applications need to handle concurrent access:

**Key concepts:**
- Use `sync.Mutex` to protect account operations
- Lock before checking balance and modifying it
- Use `defer` to ensure mutex is always unlocked
- Consider read/write locks for read-heavy operations

### Testing Error Scenarios

Testing error handling is critical:

**Testing strategies:**
- Test each error condition separately
- Verify error types and messages
- Test successful operations after handling errors
- Use table-driven tests for multiple error scenarios
- Mock dependencies to simulate error conditions

### Error Logging and Reporting

Proper error logging is essential:

**Logging best practices:**
- Log errors with sufficient context
- Include relevant IDs (account, transaction, user)
- Log at appropriate levels (error, warning, info)
- Don't log the same error multiple times in the call stack
- Structure logs for easy parsing and monitoring

### Panic and Recover

While Go prefers explicit error handling, `panic` and `recover` are available for exceptional cases:

**When to use panic:**
- Unrecoverable errors that indicate programmer mistakes
- Initialization failures that prevent the program from working
- Internal consistency violations

**Recovery patterns:**
- Use `defer` with `recover()` to catch panics
- Convert panics to errors when appropriate
- Log panics for debugging
- Only recover at appropriate boundaries

## Further Reading

- [Error Handling in Go](https://blog.golang.org/error-handling-and-go)
- [Working with Errors in Go 1.13+](https://blog.golang.org/go1.13-errors)
- [Effective Error Handling in Go](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- [Sync Package Documentation](https://pkg.go.dev/sync) 
[View the Scoreboard](SCOREBOARD.md)

# Challenge 7: Bank Account with Error Handling

## Problem Statement

Implement a simple banking system with proper error handling. You'll create a `BankAccount` struct that manages balance operations and implements appropriate error handling.

## Requirements

1. Implement a `BankAccount` struct that has the following fields:
   - `ID` (string): Unique identifier for the account
   - `Owner` (string): Name of the account owner
   - `Balance` (float64): Current balance of the account
   - `MinBalance` (float64): Minimum balance that must be maintained

2. Implement the following methods:
   - `NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error)`: Constructor that validates input parameters
   - `Deposit(amount float64) error`: Adds money to the account
   - `Withdraw(amount float64) error`: Removes money from the account
   - `Transfer(amount float64, target *BankAccount) error`: Transfers money from one account to another

3. You must implement custom error types:
   - `InsufficientFundsError`: When withdrawal/transfer would bring balance below minimum
   - `NegativeAmountError`: When deposit/withdraw/transfer amount is negative
   - `ExceedsLimitError`: When deposit/withdrawal amount exceeds your defined limits
   - `AccountError`: A general bank account error with appropriate subtypes

## Function Signatures

```go
// Constructor
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error)

// Methods
func (a *BankAccount) Deposit(amount float64) error
func (a *BankAccount) Withdraw(amount float64) error
func (a *BankAccount) Transfer(amount float64, target *BankAccount) error

// Error types
type AccountError struct {
    // Implement custom error type with appropriate fields
}

type InsufficientFundsError struct {
    // Implement custom error type with appropriate fields
}

type NegativeAmountError struct {
    // Implement custom error type with appropriate fields
}

type ExceedsLimitError struct {
    // Implement custom error type with appropriate fields
}

// Each error type should implement the Error() string method
func (e *AccountError) Error() string
func (e *InsufficientFundsError) Error() string
func (e *NegativeAmountError) Error() string
func (e *ExceedsLimitError) Error() string
```

## Constraints

- All amounts must be valid values (non-negative).
- Withdrawals/transfers cannot bring account balance below the minimum balance.
- Define a reasonable limit for deposits and withdrawals (e.g., $10,000).
- Error messages should be descriptive and include relevant information.
- All operations should be thread-safe (use proper synchronization mechanisms).

## Sample Usage

```go
// Create new bank accounts
account1, err := NewBankAccount("ACC001", "Alice", 1000.0, 100.0)
if err != nil {
    // Handle error
}

account2, err := NewBankAccount("ACC002", "Bob", 500.0, 50.0) 
if err != nil {
    // Handle error
}

// Deposit money
if err := account1.Deposit(200.0); err != nil {
    // Handle error
}

// Withdraw money
if err := account1.Withdraw(50.0); err != nil {
    // Handle error
}

// Transfer money
if err := account1.Transfer(300.0, account2); err != nil {
    // Handle error
}
```

## Instructions

- **Fork** the repository.
- **Clone** your fork to your local machine.
- **Create** a directory named after your GitHub username inside `challenge-7/submissions/`.
- **Copy** the `solution-template.go` file into your submission directory.
- **Implement** the required structs and methods.
- **Test** your solution locally by running the test file.
- **Commit** and **push** your code to your fork.
- **Create** a pull request to submit your solution.

## Testing Your Solution Locally

Run the following command in the `challenge-7/` directory:

```bash
go test -v
``` 
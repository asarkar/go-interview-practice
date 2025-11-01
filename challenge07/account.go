// Package challenge7 contains the solution for Challenge 7: Bank Account with Error Handling.
package account

import (
	"fmt"
	"sync"
	// Add any other necessary imports
)

// BankAccount represents a bank account with balance management and minimum balance requirements.
type BankAccount struct {
	ID         string
	Owner      string
	Balance    float64
	MinBalance float64
	mu         sync.Mutex
}

// Constants for account operations
const (
	MaxTransactionAmount = 10000.0 // Example limit for deposits/withdrawals
)

// Custom error types

// AccountError is a general error type for bank account operations.
type AccountError struct {
	accountID string
	message   string
}

func (e *AccountError) Error() string {
	return fmt.Sprintf(
		"account ID: %s, %s",
		e.accountID,
		e.message,
	)
}

// InsufficientFundsError occurs when a withdrawal or transfer would bring the balance below minimum.
type InsufficientFundsError struct {
	accountID  string
	amount     float64
	minBalance float64
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf(
		"account ID: %s, transaction of amount %f would bring the balance below the minimum %f",
		e.accountID,
		e.amount,
		e.minBalance,
	)
}

// NegativeAmountError occurs when an amount for deposit, withdrawal, or transfer is negative.
type NegativeAmountError struct {
	accountID string
}

func (e *NegativeAmountError) Error() string {
	return fmt.Sprintf(
		"account ID: %s, amount is negative",
		e.accountID,
	)
}

// ExceedsLimitError occurs when a deposit or withdrawal amount exceeds the defined limit.
type ExceedsLimitError struct {
	accountID string
	amount    float64
}

func (e *ExceedsLimitError) Error() string {
	return fmt.Sprintf(
		"account ID: %s, transaction of amount %f exceeds the limit %f",
		e.accountID,
		e.amount,
		MaxTransactionAmount,
	)
}

// NewBankAccount creates a new bank account with the given parameters.
// It returns an error if any of the parameters are invalid.
func NewBankAccount(id, owner string, initialBalance, minBalance float64) (*BankAccount, error) {
	if id == "" {
		return nil, &AccountError{accountID: id, message: "account ID is blank"}
	}
	if owner == "" {
		return nil, &AccountError{accountID: id, message: "account owner is blank"}
	}
	if initialBalance < 0 || minBalance < 0 {
		return nil, &NegativeAmountError{accountID: id}
	}
	if initialBalance < minBalance {
		return nil, &InsufficientFundsError{
			accountID:  id,
			amount:     initialBalance,
			minBalance: minBalance,
		}
	}

	return &BankAccount{
		ID:         id,
		Owner:      owner,
		Balance:    initialBalance,
		MinBalance: minBalance,
	}, nil
}

// Deposit adds the specified amount to the account balance.
// It returns an error if the amount is invalid or exceeds the transaction limit.
func (a *BankAccount) Deposit(amount float64) error {
	if err := validateAmount(a.ID, amount); err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.transact(amount)
}

// Withdraw removes the specified amount from the account balance.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.
func (a *BankAccount) Withdraw(amount float64) error {
	// Implement withdrawal functionality with proper error handling
	if err := validateAmount(a.ID, amount); err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.transact(-amount)
}

func validateAmount(id string, amount float64) error {
	if amount < 0 {
		return &NegativeAmountError{accountID: id}
	}
	if amount > MaxTransactionAmount {
		return &ExceedsLimitError{accountID: id, amount: amount}
	}
	return nil
}

// Transfer moves the specified amount from this account to the target account.
// It returns an error if the amount is invalid, exceeds the transaction limit,
// or would bring the balance below the minimum required balance.

//nolint:revive // receiver-naming: `from` is clearer than `a`.
func (from *BankAccount) Transfer(amount float64, to *BankAccount) error {
	if err := validateAmount(from.ID, amount); err != nil {
		return err
	}
	// Lock accounts in a consistent order to prevent deadlocks.
	// For example, always lock the account with the smaller ID first.
	if from.ID < to.ID {
		from.mu.Lock()
		to.mu.Lock()
	} else {
		to.mu.Lock()
		from.mu.Lock()
	}
	defer from.mu.Unlock()
	defer to.mu.Unlock()

	// The exported methods try to acquire the locks again, and block forever
	// since `sync.Mutex` is not reentrant. Call the unexported internal method.
	if err := from.transact(-amount); err != nil {
		return err
	}
	return to.transact(amount)
}

func (a *BankAccount) transact(amount float64) error {
	if amount < 0 && a.Balance+amount < a.MinBalance {
		return &InsufficientFundsError{accountID: a.ID, amount: amount, minBalance: a.MinBalance}
	}
	a.Balance += amount
	return nil
}

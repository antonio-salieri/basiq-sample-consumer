package entity

// Transaction type describes transaction entity
type Transaction struct {
	ID            string
	Status        TransactionStatus
	Description   string
	Amount        float64
	Balance       string
	InstitutionID string
	ConnectionID  string
	Class         string
	Direction     TransactionDirection
	Type          TransactionType
}

// TransactionType transaction type structure
type TransactionType struct {
	Code  string
	Title string
}

// TransactionCollection type
type TransactionCollection []Transaction

// TransactionStatus status type
type TransactionStatus string

// TransactionDirection direction type
type TransactionDirection string

const (
	// Pending transaction status
	Pending TransactionStatus = "pending"
	// Posted transaction status
	Posted TransactionStatus = "posted"

	// Debit transaction direction
	Debit TransactionDirection = "debit"
	// Credit transaction direction
	Credit TransactionDirection = "credit"
)

package service

import (
	"fmt"

	"github.com/antonio-salieri/basiq-sample-consumer/client"
	"github.com/antonio-salieri/basiq-sample-consumer/entity"
)

// TransactionService transaction service structure
type TransactionService struct {
	apiClient client.Client
}

// TransactionAmountPerType amount and transaction count per transaction type
type TransactionAmountPerType struct {
	transactionCount int32
	totalSum         float64
	transactionType  entity.TransactionType
	average          float64
}

// AggregatedTransactionsPerType map of transactions aggregated per transaction type
type AggregatedTransactionsPerType map[string]TransactionAmountPerType

// NewTransactionService creates new transaction service
func NewTransactionService(apiClient client.Client) *TransactionService {
	return &TransactionService{
		apiClient: apiClient,
	}
}

// AggregateTransactionPerDebitCategory calculate average transaction amount per transaction type
func (ts *TransactionService) AggregateTransactionPerDebitCategory(userID string, institutionData entity.ConnectionData) (AggregatedTransactionsPerType, error) {
	transactions, err := ts.fetchUserTransactions(userID, institutionData)
	if err != nil {
		return nil, err
	}

	var aggregatedTransactions = make(AggregatedTransactionsPerType)

	for _, t := range transactions {
		if t.Direction == entity.Debit {
			if existing, ok := aggregatedTransactions[t.Type.Code]; ok {
				existing.totalSum += t.Amount
				existing.transactionCount++

				aggregatedTransactions[t.Type.Code] = existing
			} else {
				aggregatedTransactions[t.Type.Code] = TransactionAmountPerType{
					totalSum:         t.Amount,
					transactionCount: 1,
					transactionType:  t.Type,
				}
			}
		}
	}

	return aggregatedTransactions, nil
}

func (ts *TransactionService) fetchUserTransactions(userID string, institutionData entity.ConnectionData) (entity.TransactionCollection, error) {
	var connection *entity.Connection

	// Try to find existing connection
	connections, err := ts.apiClient.GetConnectionsToInstitution(userID, institutionData.InstitutionID)
	if err != nil {
		return nil, err
	}
	if len(connections) > 0 {
		connection, err = connections.GetFirstActiveConnection()
	} else {
		// Create connection if it does not exists
		connection, err = ts.apiClient.CreateConnection(userID, institutionData)
	}
	if err != nil {
		return nil, err
	}

	// Fetch all user transactions in given insitution
	transactions, err := ts.apiClient.GetTransactions(userID, connection.ID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetAverageAmounts prints average transaction amount per category
func (a AggregatedTransactionsPerType) GetAverageAmounts(transactionTypeCode *string) AggregatedTransactionsPerType {
	var result = make(AggregatedTransactionsPerType)
	if transactionTypeCode == nil {
		for k, v := range a {
			result[k] = TransactionAmountPerType{
				totalSum:         v.totalSum,
				transactionCount: v.transactionCount,
				transactionType:  v.transactionType,
				average:          v.totalSum / float64(v.transactionCount),
			}
		}
	} else {
		v := a[*transactionTypeCode]
		result[*transactionTypeCode] = TransactionAmountPerType{
			totalSum:         v.totalSum,
			transactionCount: v.transactionCount,
			transactionType:  v.transactionType,
			average:          v.totalSum / float64(v.transactionCount),
		}
	}

	return result
}

// Print prints contents of AggregatedTransactionsPerType
func (a AggregatedTransactionsPerType) Print() {
	fmt.Printf("Code\t\t| Average\t\t| Total\t\t\t| Count\t\t\t |Title \n")
	for _, v := range a {
		fmt.Printf("%s\t\t| %0.2f\t\t| %0.2f\t\t| %d\t\t\t |%s\n",
			v.transactionType.Code, v.average, v.totalSum, v.transactionCount, v.transactionType.Title)
	}
}

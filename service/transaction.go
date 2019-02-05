package service

import (
	"fmt"
	"sort"

	"github.com/antonio-salieri/basiq-sample-consumer/client"
	"github.com/antonio-salieri/basiq-sample-consumer/entity"
)

// TransactionService transaction service structure
type TransactionService struct {
	apiClient client.TransactionsAPI
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
func NewTransactionService(apiClient client.TransactionsAPI) *TransactionService {
	return &TransactionService{
		apiClient: apiClient,
	}
}

// AggregateTransactionPerDebitCategory calculate average transaction amount per transaction type
func (ts *TransactionService) AggregateTransactionPerDebitCategory(userID string, institutionData entity.ConnectionData) (AggregatedTransactionsPerType, error) {
	transactions, err := ts.apiClient.GetUserTransactionsInInstitution(userID, institutionData)
	if err != nil {
		return nil, err
	}

	var aggregatedTransactions = make(AggregatedTransactionsPerType)

	for _, t := range transactions {
		// Include only debit transactions with some category
		if t.Direction == entity.Debit && len(t.Type.Code) > 0 {
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
	fmt.Printf("\nCode\t\t| Average\t\t| Total\t\t\t| Count\t\t\t | Title\n")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")

	// Sort output by transaction types
	keys := make([]string, 0)
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		value := a[k]
		fmt.Printf("%s\t\t| %0.2f\t\t| %0.2f\t\t| %d\t\t\t |%s\n",
			value.transactionType.Code, value.GetAverageAmount(), value.totalSum, value.transactionCount, value.transactionType.Title)
	}
}

// GetAverageAmount returns aggregated transaction average amount
func (a TransactionAmountPerType) GetAverageAmount() float64 {
	return a.average
}

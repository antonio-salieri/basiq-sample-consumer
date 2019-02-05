package service_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/antonio-salieri/basiq-sample-consumer/entity"
	"github.com/antonio-salieri/basiq-sample-consumer/service"
)

const expectedAggregatedDebitTransactios = 3
const averageAmountForTransaction12 = -60.33
const averageAmountForTransaction10 = -10.5

var userIDMock string
var connectionDataMock entity.ConnectionData
var mockedTransactionCollection = entity.TransactionCollection{
	entity.Transaction{Amount: -10, Direction: entity.Debit, Type: entity.TransactionType{Code: "10", Title: "Out 10"}},
	entity.Transaction{Amount: -11, Direction: entity.Debit, Type: entity.TransactionType{Code: "10", Title: "Out 10"}},

	entity.Transaction{Amount: 10, Direction: entity.Credit, Type: entity.TransactionType{Code: "1", Title: "In 1"}},

	entity.Transaction{Amount: -122, Direction: entity.Debit, Type: entity.TransactionType{Code: "11", Title: "Out 11"}},

	entity.Transaction{Amount: 10, Direction: entity.Credit, Type: entity.TransactionType{Code: "2", Title: "In 2"}},

	entity.Transaction{Amount: -112, Direction: entity.Debit, Type: entity.TransactionType{Code: "12", Title: "Out 12"}},
	entity.Transaction{Amount: -67, Direction: entity.Debit, Type: entity.TransactionType{Code: "12", Title: "Out 12"}},
	entity.Transaction{Amount: -2, Direction: entity.Debit, Type: entity.TransactionType{Code: "12", Title: "Out 12"}},

	entity.Transaction{Amount: 10, Direction: entity.Credit, Type: entity.TransactionType{Code: "3", Title: "In 3"}},
}

type transactionApiStub struct {
}

type faultyTransactionApiStub struct {
}

func (t transactionApiStub) GetTransactions(userID string, connectionID string) (entity.TransactionCollection, error) {
	return mockedTransactionCollection, nil
}

func (t transactionApiStub) GetUserTransactionsInInstitution(userID string, institutionData entity.ConnectionData) (entity.TransactionCollection, error) {
	return mockedTransactionCollection, nil
}

func (t faultyTransactionApiStub) GetTransactions(userID string, connectionID string) (entity.TransactionCollection, error) {
	return nil, fmt.Errorf("Error fetching transactions")
}

func (t faultyTransactionApiStub) GetUserTransactionsInInstitution(userID string, institutionData entity.ConnectionData) (entity.TransactionCollection, error) {
	return nil, fmt.Errorf("Error fetching institution transactions")
}

func TestFaultyTransactionApi(t *testing.T) {
	service := service.NewTransactionService(faultyTransactionApiStub{})
	_, err := service.AggregateTransactionPerDebitCategory(userIDMock, connectionDataMock)
	if err == nil {
		t.Fatal("Expected error not received")
	}
	if err.Error() != "Error fetching institution transactions" {
		t.Fatalf("Unexpected error '%s' received", err)
	}
}
func TestAverageTransactionCalculation(t *testing.T) {

	service := service.NewTransactionService(transactionApiStub{})
	transactions, err := service.AggregateTransactionPerDebitCategory(userIDMock, connectionDataMock)
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if len(transactions) != expectedAggregatedDebitTransactios {
		t.Fatalf("Expected transactions count %d, got %d", expectedAggregatedDebitTransactios, len(transactions))
	}

	// Verify average for multiple transaction types
	result := transactions.GetAverageAmounts(nil)
	verifyAverageAggregation(t, result, expectedAggregatedDebitTransactios, map[string]float64{
		"10": averageAmountForTransaction10,
		"11": -122,
		"12": averageAmountForTransaction12,
	})

	// Verify average for one transaction type
	specificTransactionType := "12"
	result = transactions.GetAverageAmounts(&specificTransactionType)
	verifyAverageAggregation(t, result, 1, map[string]float64{"12": averageAmountForTransaction12})
}

func ExampleAggregatedTransactionsPerType_Print() {
	service := service.NewTransactionService(transactionApiStub{})
	transactions, _ := service.AggregateTransactionPerDebitCategory(userIDMock, connectionDataMock)

	avgs := transactions.GetAverageAmounts(nil)
	avgs.Print()
	// Output:
	// Code		| Average		| Total			| Count			 | Title
	// -------------------------------------------------------------------------------------------------------------------------------
	// 10		| -10.50		| -21.00		| 2			 |Out 10
	// 11		| -122.00		| -122.00		| 1			 |Out 11
	// 12		| -60.33		| -181.00		| 3			 |Out 12

}

func verifyAverageAggregation(t *testing.T, transactions service.AggregatedTransactionsPerType, expectedCount int, truthTable map[string]float64) {
	if len(transactions) != expectedCount {
		t.Fatalf("Expected aggregated transactions count %d, got %d", expectedAggregatedDebitTransactios, len(transactions))
	}

	for code, expectedAvg := range truthTable {
		calculatedAvg := math.Round((transactions[code].GetAverageAmount() * 100)) / 100
		if calculatedAvg != expectedAvg {
			t.Fatalf("Expected average transactions amount for code %s: %02f, got %02f", code, expectedAvg, calculatedAvg)
		}
	}
}

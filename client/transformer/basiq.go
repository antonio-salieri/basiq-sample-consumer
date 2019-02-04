package transformer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/antonio-salieri/basiq-sample-consumer/entity"
	"github.com/basiqio/basiq-sdk-golang/errors"
	v2 "github.com/basiqio/basiq-sdk-golang/v2"
)

// UserToUserEntity converts API user format to user entity
func UserToUserEntity(user v2.User) *entity.User {
	return &entity.User{
		ID:     user.Id,
		Email:  user.Email,
		Mobile: user.Mobile,
	}
}

// ConnectionsToConnectionsEntityCollection converts API connection collection to entity collection
func ConnectionsToConnectionsEntityCollection(connections v2.ConnectionList) entity.ConnectionCollection {
	var collection = make(entity.ConnectionCollection, len(connections.Data))

	for i, c := range connections.Data {
		collection[i] = ConnectionToEntityConnection(c)
	}

	return collection
}

// ConnectionToEntityConnection converts API connection to entity connection
func ConnectionToEntityConnection(connection v2.Connection) entity.Connection {
	return entity.Connection{
		ID:            connection.Id,
		InstitutionID: connection.Institution.Id,
		Status:        connection.Status,
	}
}

// TransactionToEntityTransaction converts transaction to entity transaction
func TransactionToEntityTransaction(transaction v2.Transaction) entity.Transaction {
	amount, err := strconv.ParseFloat(transaction.Amount, 32)
	if err != nil {
		log.Printf("Error converting transaction amount %s. Using: %f. %s", transaction.Amount, amount, err)
	}

	return entity.Transaction{
		ID:            transaction.Id,
		Amount:        amount,
		Class:         transaction.Class,
		ConnectionID:  transaction.Connection,
		InstitutionID: transaction.Institution,
		Balance:       transaction.Balance,
		Description:   transaction.Description,
		Direction:     entity.TransactionDirection(transaction.Direction),
		Status:        entity.TransactionStatus(transaction.Status),
		Type: entity.TransactionType{
			Code:  transaction.SubClass.Code,
			Title: transaction.SubClass.Title,
		},
	}
}

// APIErrorToError converts API descriptive error object to error
func APIErrorToError(label string, apiError *errors.APIError) error {
	if apiError == nil {
		return nil
	}
	return fmt.Errorf("%s => API ERROR: [%d] %s", label, apiError.StatusCode, apiError.Message)
}

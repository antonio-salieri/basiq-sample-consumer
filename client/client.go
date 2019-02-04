package client

import "github.com/antonio-salieri/basiq-sample-consumer/entity"

// Client exposes api client interface
type Client interface {
	GetUser(userID string) (*entity.User, error)
	CreateUser(email string, mobile string) (*entity.User, error)
	DeleteUser(userID string) error

	GetConnectionsToInstitution(userID string, institution string) (entity.ConnectionCollection, error)
	CreateConnection(userID string, connectionData entity.ConnectionData) (*entity.Connection, error)
	DeleteConnection(userID string, connectionID string) error

	GetTransactions(userID string, connectionID string) (entity.TransactionCollection, error)
}

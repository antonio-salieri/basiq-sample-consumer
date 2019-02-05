package client

import "github.com/antonio-salieri/basiq-sample-consumer/entity"

// Client exposes api client interface
type Client interface {
	UsersAPI
	ConnectionsAPI
	TransactionsAPI
}

// UsersAPI contains methods for manipulating with users
type UsersAPI interface {
	GetUser(userID string) (*entity.User, error)
	CreateUser(email string, mobile string) (*entity.User, error)
	DeleteUser(userID string) error
}

// TransactionsAPI contains transactions related API
type TransactionsAPI interface {
	GetTransactions(userID string, connectionID string) (entity.TransactionCollection, error)
	GetUserTransactionsInInstitution(userID string, institutionData entity.ConnectionData) (entity.TransactionCollection, error)
}

// ConnectionsAPI contains methods for manipulating connections
type ConnectionsAPI interface {
	GetConnectionsToInstitution(userID string, institution string) (entity.ConnectionCollection, error)
	CreateConnection(userID string, connectionData entity.ConnectionData) (*entity.Connection, error)
	DeleteConnection(userID string, connectionID string) error
}

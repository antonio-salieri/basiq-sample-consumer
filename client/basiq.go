package client

import (
	"log"

	t "github.com/antonio-salieri/basiq-sample-consumer/client/transformer"
	"github.com/antonio-salieri/basiq-sample-consumer/entity"
	"github.com/basiqio/basiq-sdk-golang/basiq"
	"github.com/basiqio/basiq-sdk-golang/utilities"
	api "github.com/basiqio/basiq-sdk-golang/v2"
)

const (
	jobPollingIntervalMs = 2000
	jobPollingTimeoutMs  = 15000
)

// BasiqClient structure that implements Client interface
type BasiqClient struct {
	session            api.Session
	userService        api.UserService
	transactionService api.TransactionService
}

// NewBasiqClient creates new Basiq cilent instance
func NewBasiqClient(apiKey string) (Client, error) {
	session, err := basiq.NewSessionV2(apiKey)
	if err != nil {
		return nil, t.APIErrorToError("Create new session", err)
	}

	return &BasiqClient{
		session:            *session,
		userService:        *api.NewUserService(session),
		transactionService: *api.NewTransactionService(session),
	}, nil
}

// GetUser retrieves user entitiy using Basiq API
func (c *BasiqClient) GetUser(userID string) (*entity.User, error) {
	user, err := c.userService.GetUser(userID)
	if err != nil {
		return nil, t.APIErrorToError("Fetching user", err)
	}

	return t.UserToUserEntity(user), nil
}

// CreateUser creates new user using Basiq API
func (c *BasiqClient) CreateUser(email string, mobile string) (*entity.User, error) {
	user, err := c.userService.CreateUser(&api.UserData{
		Email:  email,
		Mobile: mobile,
	})
	if err != nil {
		return nil, t.APIErrorToError("Creating user", err)
	}

	return t.UserToUserEntity(user), nil
}

// DeleteUser deletes user given user id
func (c *BasiqClient) DeleteUser(userID string) error {
	return t.APIErrorToError("Deleting user", c.userService.DeleteUser(userID))
}

// GetConnectionsToInstitution retrieves users connections for specified institution
func (c *BasiqClient) GetConnectionsToInstitution(userID string, institution string) (entity.ConnectionCollection, error) {
	user := c.userService.ForUser(userID)

	filter := utilities.FilterBuilder{}
	filter.Eq("institution.id", institution)

	connections, err := user.ListAllConnections(&filter)
	if err != nil {
		return nil, t.APIErrorToError("Fetching connections to institution for user", err)
	}

	return t.ConnectionsToConnectionsEntityCollection(connections), nil
}

// CreateConnection creates new connection to institution
func (c *BasiqClient) CreateConnection(userID string, connectionData entity.ConnectionData) (*entity.Connection, error) {
	user := c.userService.ForUser(userID)
	connectionService := api.NewConnectionService(&c.session, &user)
	connJob, err := connectionService.NewConnection(&api.ConnectionData{
		Institution: &api.InstitutionData{
			Id: connectionData.InstitutionID,
		},
		LoginId:  connectionData.LoginID,
		Password: connectionData.LoginPassword,
	})
	if err != nil {
		return nil, t.APIErrorToError("Creating connection job", err)
	}

	conn, err := connJob.WaitForCredentials(jobPollingIntervalMs, jobPollingTimeoutMs)
	if err != nil {
		return nil, t.APIErrorToError("Waiting create connection job to finish", err)
	}

	connection := t.ConnectionToEntityConnection(conn)
	return &connection, nil
}

// DeleteConnection deletes user connection given connection id and user id
func (c *BasiqClient) DeleteConnection(userID string, connectionID string) error {
	user := c.userService.ForUser(userID)
	connectionService := api.NewConnectionService(&c.session, &user)

	return t.APIErrorToError("Deleting connection", connectionService.DeleteConnection(connectionID))
}

// GetTransactions retireves all transactions given connection id
func (c *BasiqClient) GetTransactions(userID string, connectionID string) (entity.TransactionCollection, error) {
	user := c.userService.ForUser(userID)

	fb := utilities.FilterBuilder{}
	fb.Eq("connection.id", connectionID)

	transactions, err := user.GetTransactions(&fb)
	if err != nil {
		return nil, t.APIErrorToError("Initial user transactions fetch", err)
	}

	var collection entity.TransactionCollection

	for {
		next, err := transactions.Next()
		if err != nil {
			return nil, t.APIErrorToError("Fetching transactions page", err)
		}
		if next == false {
			break
		}

		for _, data := range transactions.Data {
			collection = append(collection, t.TransactionToEntityTransaction(data))
		}
	}
	log.Printf("Fetched %d transactions", len(collection))

	return collection, nil
}

// GetUserTransactionsInInstitution retrieves all user transactions for
func (c *BasiqClient) GetUserTransactionsInInstitution(userID string, institutionData entity.ConnectionData) (entity.TransactionCollection, error) {
	var connection *entity.Connection

	// Try to find existing connection
	connections, err := c.GetConnectionsToInstitution(userID, institutionData.InstitutionID)
	if err != nil {
		return nil, err
	}
	if len(connections) > 0 {
		connection, err = connections.GetFirstActiveConnection()
	} else {
		// Create connection if it does not exists
		connection, err = c.CreateConnection(userID, institutionData)
	}
	if err != nil {
		return nil, err
	}

	// Fetch all user transactions in given insitution
	transactions, err := c.GetTransactions(userID, connection.ID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

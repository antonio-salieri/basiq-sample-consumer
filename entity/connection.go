package entity

import "fmt"

// Connection describes connection entity
type Connection struct {
	ID            string
	Status        string
	InstitutionID string
}

// ConnectionCollection collection of connections
type ConnectionCollection []Connection

// ConnectionData used for creating connection
type ConnectionData struct {
	InstitutionID string
	LoginID       string
	LoginPassword string
}

// IsActive returns true when connection is active
func (c Connection) IsActive() bool {
	return c.Status == "active"
}

// GetFirstActiveConnection traverse connection collection and returns first active connection.
// If no connection is found error is returned
func (collection ConnectionCollection) GetFirstActiveConnection() (*Connection, error) {
	var connection *Connection
	err := fmt.Errorf("None of %d found connection is active. Delete existing and create new connection with valid credentials", len(collection))

	for _, c := range collection {
		if c.Status == "active" {
			connection = &c
			err = nil
			break
		}
	}

	return connection, err
}

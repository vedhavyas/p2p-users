package db

import (
	"fmt"
	"os"

	"github.com/globalsign/mgo"
)

type Service struct {
	db *mgo.Database
}

// GetDBService returns the Database service
func GetDBService() (*Service, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	session, err := mgo.Dial(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	err = session.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	db := session.DB("")
	err = addUserIndexes(db)
	if err != nil {
		return nil, fmt.Errorf("failed to add user indexes: %v", err)
	}

	return &Service{db: db}, nil
}

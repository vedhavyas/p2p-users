package db

import (
	"fmt"

	"github.com/globalsign/mgo"
)

type Service struct {
	db *mgo.Session
}

// GetDBService returns the Database service
func GetDBService(dbURL string) (*Service, error) {
	session, err := mgo.Dial(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	err = session.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %v", err)
	}

	return &Service{db: session}, nil
}

func (s *Service) EnsureIndexes() error {
	err := addUserIndexes(s.db)
	if err != nil {
		return fmt.Errorf("failed to add user indexes: %v", err)
	}

	return nil
}

package db

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const userCollection = "users"

// User represents a user in DB
type User struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Phone     string        `json:"phone" bson:"phone"`
	FirstName string        `json:"first_name" bson:"first_name"`
	LastName  string        `json:"last_name" bson:"last_name"`
	Picture   string        `json:"picture" bson:"picture"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

// UserOperator is an abstraction for CRUD operations on User
type UserOperator interface {
	CreateUser(user User) (User, error)
	GetUser(id string) (User, error)
	GetUserByPhone(phone string) (User, error)
}

// addUserIndexes adds indexes(if not) for user collection
func addUserIndexes(db *mgo.Database) error {
	c := db.C(userCollection)
	i := mgo.Index{
		Key:    []string{"phone"},
		Unique: true,
	}
	err := c.EnsureIndex(i)
	if err != nil {
		return fmt.Errorf("failed to add index: %v", err)
	}

	return nil
}

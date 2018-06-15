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
func addUserIndexes(db *mgo.Session) error {
	session := db.Copy()
	defer session.Close()
	c := session.DB("").C(userCollection)
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

// CreateUser creates a new user
func (s *Service) CreateUser(user User) (User, error) {
	session := s.db.Copy()
	defer session.Close()

	// TODO rpc server should take care of cleaning up phone number
	c := session.DB("").C(userCollection)
	user.ID = bson.NewObjectId()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := c.Insert(&user)
	if err == nil {
		return user, nil
	}

	if mgo.IsDup(err) {
		return user, fmt.Errorf("user already exists: %v", err)
	}

	return user, fmt.Errorf("unknown error during insert user: %v", err)
}

// GetUser returns the user by ID
func (s *Service) GetUser(id string) (User, error) {
	session := s.db.Copy()
	defer session.Close()

	c := session.DB("").C(userCollection)
	var user User
	err := c.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return user, fmt.Errorf("user not found: %v", err)
	}

	return user, nil
}

// GetUserByPhone returns user by phone
func (s *Service) GetUserByPhone(phone string) (User, error) {
	session := s.db.Copy()
	defer session.Close()

	c := session.DB("").C(userCollection)
	var user User
	err := c.Find(bson.M{"phone": phone}).One(&user)
	if err != nil {
		return user, fmt.Errorf("user not found: %v", err)
	}

	return user, nil
}

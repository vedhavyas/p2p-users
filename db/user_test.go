package db

import (
	"fmt"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	service *Service
}

func (u *UserTestSuite) SetupSuite() {
	var err error
	u.service, err = GetDBService("localhost")
	if err != nil {
		panic(fmt.Errorf("unable to create service: %v", err))
	}

	u.service.db.DB("").C(userCollection).DropCollection()

	err = u.service.EnsureIndexes()
	if err != nil {
		panic(fmt.Errorf("failed to add indexes: %v", err))
	}
}

func (u *UserTestSuite) TearDownSuite() {
	u.service.db.DB("").C(userCollection).DropCollection()
	u.service.db.Close()
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (u *UserTestSuite) TestUser() {
	user := u.UserCreateOK()
	u.UserCreateError()

	u1 := u.UserGetByPhoneOK(user.Phone)
	checkUser(u, u1)
	err := u.UserGetByPhoneError("+919556445567")
	require.Errorf(u.T(), err, "user not found")

	u2 := u.UserGetByID(user.ID.Hex())
	checkUser(u, u2)
	err = u.UserGetByIDError(bson.NewObjectId().Hex())
	require.Errorf(u.T(), err, "user not found")
}

func checkUser(u *UserTestSuite, user User) {
	r := require.New(u.T())
	r.NotEmpty(user.ID)
	r.Equal(user.Phone, "+919663556657")
	r.Equal(user.FirstName, "Vedhavyas")
	r.Equal(user.LastName, "Singareddi")
	r.Equal(user.Picture, "http://some.image/vedhavyas")
	r.NotEmpty(user.CreatedAt)
	r.NotEmpty(user.UpdatedAt)
}

func (u *UserTestSuite) UserCreateOK() User {

	user := User{
		Phone:     "+919663556657",
		FirstName: "Vedhavyas",
		LastName:  "Singareddi",
		Picture:   "http://some.image/vedhavyas",
	}

	var err error
	user, err = u.service.CreateUser(user)
	require.NoError(u.T(), err)
	checkUser(u, user)
	return user
}

func (u *UserTestSuite) UserCreateError() {
	r := require.New(u.T())
	user := User{
		Phone: "+919663556657",
	}

	user, err := u.service.CreateUser(user)
	r.Errorf(err, "user already exists")
}

func (u *UserTestSuite) UserGetByPhoneOK(phone string) User {
	user, err := u.service.GetUserByPhone(phone)
	require.NoError(u.T(), err)
	return user
}

func (u *UserTestSuite) UserGetByPhoneError(phone string) error {
	_, err := u.service.GetUserByPhone(phone)
	return err
}

func (u *UserTestSuite) UserGetByID(id string) User {
	user, err := u.service.GetUser(id)
	require.NoError(u.T(), err)
	return user
}

func (u *UserTestSuite) UserGetByIDError(id string) error {
	_, err := u.service.GetUser(id)
	return err
}

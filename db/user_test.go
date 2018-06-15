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
	u.UserGetByPhoneOK(user.Phone)
	u.UserGetByPhoneError("+919556445567")
	u.UserGetByID(user.ID)
	u.UserGetByIDError(bson.NewObjectId())
	u.UpdateUserOK(user.ID)
	u.UpdateUserError(bson.NewObjectId())
}

func checkUser(user User, r *require.Assertions) {
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

	r := require.New(u.T())
	var err error
	user, err = u.service.CreateUser(user)
	require.NoError(u.T(), err)
	checkUser(user, r)
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

func (u *UserTestSuite) UserGetByPhoneOK(phone string) {
	user, err := u.service.GetUserByPhone(phone)
	r := require.New(u.T())
	r.NoError(err)
	checkUser(user, r)
}

func (u *UserTestSuite) UserGetByPhoneError(phone string) {
	_, err := u.service.GetUserByPhone(phone)
	require.Errorf(u.T(), err, "user not found")
}

func (u *UserTestSuite) UserGetByID(id bson.ObjectId) {
	user, err := u.service.GetUser(id)
	r := require.New(u.T())
	r.NoError(err)
	checkUser(user, r)
}

func (u *UserTestSuite) UserGetByIDError(id bson.ObjectId) {
	_, err := u.service.GetUser(id)
	require.Errorf(u.T(), err, "user not found")
}

func (u *UserTestSuite) UpdateUserOK(id bson.ObjectId) {
	updates := map[string]interface{}{
		"first_name": "Ved",
		"last_name":  "singareddi",
		"picture":    "",
	}
	err := u.service.UpdateUserByID(id, updates)
	r := require.New(u.T())
	r.NoError(err)

	user, err := u.service.GetUser(id)
	r.NoError(err)
	r.NotEmpty(user.ID)
	r.Equal(user.Phone, "+919663556657")
	r.Equal(user.FirstName, "Ved")
	r.Equal(user.LastName, "singareddi")
	r.Equal(user.Picture, "")
	r.NotEmpty(user.CreatedAt)
	r.NotEmpty(user.UpdatedAt)
}

func (u *UserTestSuite) UpdateUserError(id bson.ObjectId) {
	r := require.New(u.T())
	err := u.service.UpdateUserByID(id, nil)
	r.Errorf(err, "nothing to update")

	err = u.service.UpdateUserByID(id, map[string]interface{}{
		"_id": id,
	})
	r.Errorf(err, "nothing to update")

	err = u.service.UpdateUserByID(id, map[string]interface{}{
		"first_name": "Vedhavyas",
	})
	r.Errorf(err, "user not found")
}

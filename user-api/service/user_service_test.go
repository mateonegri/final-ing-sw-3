package service

import (
	"testing"
	"user-api/dto"
	"user-api/model"
	e "user-api/utils/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the userClient to simulate client responses
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) GetUserById(id int) model.User {
	args := m.Called(id)
	return args.Get(0).(model.User)
}

func (m *MockUserClient) GetUsers() model.Users {
	args := m.Called()
	return args.Get(0).(model.Users)
}

func (m *MockUserClient) GetUserByEmail(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func (m *MockUserClient) InsertUser(user model.User) model.User {
	args := m.Called(user)
	return args.Get(0).(model.User)
}

func (m *MockUserClient) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserClient) UpdateUser(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestGetUserById_Success(t *testing.T) {

	mockUserClient := new(MockUserClient)
	UserClient = mockUserClient

	mockUser := model.User{
		Id:       1,
		Name:     "John",
		LastName: "Doe",
		UserName: "jdoe",
		Email:    "jdoe@example.com",
	}

	mockUserClient.On("GetUserById", 1).Return(mockUser)

	userDto, err := UserService.GetUserById(1)

	assert.Nil(t, err)
	assert.Equal(t, "John", userDto.Name)
	assert.Equal(t, "Doe", userDto.LastName)
	mockUserClient.AssertExpectations(t)
}

func TestGetUserById_NotFound(t *testing.T) {

	mockUserClient := new(MockUserClient)
	UserClient = mockUserClient

	mockUserClient.On("GetUserById", 2).Return(model.User{Id: 0})

	userDto, err := UserService.GetUserById(2)

	message := string(err.Message())

	assert.Nil(t, userDto)
	assert.Equal(t, "Usuario no encontrado", message)
	mockUserClient.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {

	mockUserClient := new(MockUserClient)
	UserClient = mockUserClient

	mockUsers := model.Users{
		{Id: 1, Name: "John", LastName: "Doe", UserName: "jdoe"},
		{Id: 2, Name: "Jane", LastName: "Smith", UserName: "jsmith"},
	}

	mockUserClient.On("GetUsers").Return(mockUsers)

	usersDto, err := UserService.GetUsers()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(usersDto))
	assert.Equal(t, "John", usersDto[0].Name)
	assert.Equal(t, "Jane", usersDto[1].Name)
	mockUserClient.AssertExpectations(t)
}

func TestInsertUser_EmailExists(t *testing.T) {

	mockUserClient := new(MockUserClient)
	UserClient = mockUserClient

	mockUserDto := &dto.UserDto{Email: "jdoe@example.com"}
	mockUserClient.On("GetUserByEmail", "jdoe@example.com").Return(true)

	user, err := UserService.InsertUser(mockUserDto)

	message := string(err.Message())

	assert.Nil(t, user)
	assert.Equal(t, "El email ya está registrado", message)
	mockUserClient.AssertExpectations(t)
}

func TestInsertUser_Success(t *testing.T) {

	mockUserClient := new(MockUserClient)
	UserClient = mockUserClient

	mockUserDto := &dto.UserDto{
		Name:     "John",
		LastName: "Doe",
		UserName: "jdoe",
		Email:    "jdoe@example.com",
		Password: "password123",
	}

	mockUserClient.On("GetUserByEmail", "jdoe@example.com").Return(false)
	mockUserClient.On("InsertUser", mock.Anything).Return(model.User{Id: 1})

	user, err := UserService.InsertUser(mockUserDto)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)
	mockUserClient.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {

	mockUserClient := new(MockUserClient)
	UserClient = mockUserClient

	mockUserClient.On("DeleteUser", 1).Return(nil)

	err := UserService.DeleteUser(1)

	assert.Nil(t, err)
	mockUserClient.AssertExpectations(t)
}

func TestDeleteUser_Failure(t *testing.T) {

	mockUserClient := new(MockUserClient)
	UserClient = mockUserClient

	mockUserClient.On("DeleteUser", 2).Return(e.NewBadRequestApiError("Error deleting user"))

	err := UserService.DeleteUser(2)

	message := string(err.Error())

	assert.NotNil(t, err)
	assert.Equal(t, "Message: Error deleting user;Error Code: bad_request;Status: 400;Cause: []", message)
	mockUserClient.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {

	mockUserClient := new(MockUserClient)
	UserClient = mockUserClient

	mockUser := model.User{Id: 1, Name: "John", LastName: "Doe", UserName: "jdoe"}
	mockUserDto := &dto.UserDto{Name: "John Updated", LastName: "Doe Updated", UserName: "jdoeupdated"}

	mockUserClient.On("GetUserById", 1).Return(mockUser)
	mockUserClient.On("UpdateUser", mock.Anything).Return(nil)

	updatedUser, err := UserService.UpdateUser(1, mockUserDto)

	assert.Nil(t, err)
	assert.Equal(t, "John Updated", updatedUser.Name)
	assert.Equal(t, "Doe Updated", updatedUser.LastName)
	mockUserClient.AssertExpectations(t)
}

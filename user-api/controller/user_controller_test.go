package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-api/dto"
	"user-api/model"
	"user-api/service"
	"user-api/utils/errors"
	e "user-api/utils/errors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock the service layer
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(id int, userDto *dto.UserDto) (model.User, error) {
	args := m.Called(id, userDto)

	user := args.Get(0).(model.User)
	var apiErr e.ApiError
	if args.Get(1) != nil {
		apiErr = args.Get(1).(e.ApiError)
	}

	return user, apiErr
}

func (m *MockUserService) GetUserById(id int) (*dto.UserDto, e.ApiError) {
	args := m.Called(id)
	userDto := args.Get(0).(*dto.UserDto)
	var apiErr e.ApiError
	if args.Get(1) != nil {
		apiErr = args.Get(1).(e.ApiError)
	}
	return userDto, apiErr
}

func (m *MockUserService) GetUsers() (dto.UsersDto, e.ApiError) {
	args := m.Called()
	usersDto := args.Get(0).(dto.UsersDto)
	var apiErr e.ApiError
	if args.Get(1) != nil {
		apiErr = args.Get(1).(e.ApiError)
	}
	return usersDto, apiErr
}

func (m *MockUserService) InsertUser(userDto *dto.UserDto) (*dto.UserDto, errors.ApiError) {
	args := m.Called(userDto)
	newUserDto := args.Get(0).(*dto.UserDto)
	var apiErr errors.ApiError
	if args.Get(1) != nil {
		apiErr = args.Get(1).(errors.ApiError)
	}
	return newUserDto, apiErr
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestDeleteUser(t *testing.T) {
	mockService := new(MockUserService)
	service.UserService = mockService // Replace service with mock

	// Test case: Successful deletion
	mockService.On("DeleteUser", 1).Return(nil)
	router := setupRouter()
	router.DELETE("/users/:id", DeleteUser)

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertCalled(t, "DeleteUser", 1)

	// Test case: User not found
	mockService.On("DeleteUser", 999).Return(e.NewBadRequestApiError("User not found"))
	req, _ = http.NewRequest("DELETE", "/users/999", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	mockService.AssertCalled(t, "DeleteUser", 999)
}

func TestGetUserById(t *testing.T) {

	expectedUser := &dto.UserDto{
		Id:       1,
		Name:     "test",
		LastName: "test",
		UserName: "testuser",
	}

	mockService := new(MockUserService)
	service.UserService = mockService
	mockService.On("GetUserById", 1).Return(expectedUser, nil)

	router := setupRouter()
	router.GET("/user/:id", GetUserById)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response dto.UserDto
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Id, response.Id)
	assert.Equal(t, expectedUser.UserName, response.UserName)

	mockService.AssertCalled(t, "GetUserById", 1)

	/*
		 	// Test case: User not found
			mockService.On("GetUserById", 999).Return(nil, e.NewBadRequestApiError("User not found"))
			req, _ = http.NewRequest("GET", "/user/999", nil)
			resp = httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusNotFound, resp.Code)
	*/
}

func TestUserInsert(t *testing.T) {
	mockService := new(MockUserService)
	service.UserService = mockService

	userDto := &dto.UserDto{UserName: "newuser"}
	mockService.On("InsertUser", userDto).Return(userDto, nil)

	router := setupRouter()
	router.POST("/users", UserInsert)

	userJSON, _ := json.Marshal(userDto)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response dto.UserDto
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, userDto.UserName, response.UserName)
}

func TestUpdateUser(t *testing.T) {
	mockService := new(MockUserService)
	service.UserService = mockService

	userModel := model.User{UserName: "updateduser"}
	userDto := &dto.UserDto{UserName: "updateduser"}
	mockService.On("UpdateUser", 1, userDto).Return(userModel, nil)

	router := setupRouter()
	router.PUT("/users/:id", UpdateUser)

	userJSON, _ := json.Marshal(userDto)
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response dto.UserDto
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, userModel.UserName, response.UserName)

	// Test case: invalid user ID format
	req, _ = http.NewRequest("PUT", "/users/invalid", bytes.NewBuffer(userJSON))
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetUsers(t *testing.T) {
	mockService := new(MockUserService)
	service.UserService = mockService

	usersDto := dto.UsersDto{{Id: 1, UserName: "testuser1"}, {Id: 2, UserName: "testuser2"}}
	mockService.On("GetUsers").Return(usersDto, nil)

	router := setupRouter()
	router.GET("/users", GetUsers)

	req, _ := http.NewRequest("GET", "/users", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response dto.UsersDto
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, usersDto[0].UserName, response[0].UserName)
}

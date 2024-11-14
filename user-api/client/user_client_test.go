package user

import (
	"testing"
	"user-api/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *gorm.DB {
	// Create a new in-memory SQLite database for testing
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.User{}) // Assuming model.User exists and has correct structure
	Db = db                       // Assign the test DB to the package variable
	return db
}

func TestGetUserByUsername(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Seed the database
	testUser := model.User{UserName: "testuser", Email: "testuser@example.com"}
	db.Create(&testUser)

	// Test case: user exists
	retrievedUser, err := GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, testUser.UserName, retrievedUser.UserName)

	// Test case: user does not exist
	_, err = GetUserByUsername("nonexistentuser")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGetUserByEmail(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Seed the database
	testUser := model.User{UserName: "testuser", Email: "testuser@example.com"}
	db.Create(&testUser)

	// Test case: email exists
	assert.True(t, GetUserByEmail("testuser@example.com"))

	// Test case: email does not exist
	assert.False(t, GetUserByEmail("nonexistent@example.com"))
}

func TestGetUserById(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Seed the database
	testUser := model.User{UserName: "testuser", Email: "testuser@example.com"}
	db.Create(&testUser)

	// Test case: user exists
	retrievedUser := GetUserById(int(testUser.Id))
	assert.Equal(t, testUser.Id, retrievedUser.Id)

	// Test case: user does not exist
	retrievedUser = GetUserById(999)          // Assuming 999 does not exist
	assert.Equal(t, int(0), retrievedUser.Id) // Default Id should be zero if not found
}

func TestCheckUserById(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Seed the database
	testUser := model.User{UserName: "testuser", Email: "testuser@example.com"}
	db.Create(&testUser)

	// Test case: user exists
	assert.True(t, CheckUserById(int(testUser.Id)))

	// Test case: user does not exist
	assert.False(t, CheckUserById(999))
}

func TestInsertUser(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Create a new user
	testUser := model.User{UserName: "newuser", Email: "newuser@example.com"}
	insertedUser := InsertUser(testUser)

	// Verify user was inserted correctly
	var foundUser model.User
	db.Where("user_name = ?", "newuser").First(&foundUser)
	assert.Equal(t, insertedUser.UserName, foundUser.UserName)
	assert.Equal(t, insertedUser.Email, foundUser.Email)
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Seed the database
	testUser := model.User{UserName: "testuser", Email: "testuser@example.com"}
	db.Create(&testUser)

	// Test case: delete existing user
	err := DeleteUser(int(testUser.Id))
	assert.NoError(t, err)

	// Verify user was deleted
	var foundUser model.User
	err = db.Where("Id = ?", testUser.Id).First(&foundUser).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// Test case: delete non-existing user
	err = DeleteUser(999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	// Seed the database
	testUser := model.User{UserName: "testuser", Email: "testuser@example.com"}
	db.Create(&testUser)

	// Update user information
	testUser.Email = "updated@example.com"
	err := UpdateUser(testUser)
	assert.NoError(t, err)

	// Verify user was updated
	var foundUser model.User
	db.Where("Id = ?", testUser.Id).First(&foundUser)
	assert.Equal(t, "updated@example.com", foundUser.Email)
}

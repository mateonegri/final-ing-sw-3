package user

import (
	"fmt"
	"user-api/model"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var Db *gorm.DB

// UserClientInterface defines the interface for user operations.
type UserClientInterface interface {
	GetUserById(id int) model.User
	GetUsers() model.Users
	GetUserByEmail(email string) bool
	InsertUser(user model.User) model.User
	DeleteUser(id int) error
	UpdateUser(user model.User) error
}

type UserClient struct{}

func (UserClient) GetUserById(id int) model.User {
	return GetUserById(id)
}

func (UserClient) GetUsers() model.Users {
	return GetUsers()
}

func (UserClient) GetUserByEmail(email string) bool {
	return GetUserByEmail(email)
}

func (UserClient) InsertUser(user model.User) model.User {
	return InsertUser(user)
}

func (UserClient) DeleteUser(id int) error {
	return DeleteUser(id)
}

func (UserClient) UpdateUser(user model.User) error {
	return UpdateUser(user)
}

func GetUserByUsername(username string) (model.User, error) {
	var user model.User
	result := Db.Where("user_name = ?", username).First(&user)

	log.Debug("User: ", user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func GetUserByEmail(email string) bool {
	var user model.User
	result := Db.Where("email = ?", email).First(&user)

	log.Debug("User: ", user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false // No se encontró el usuario, el email no está registrado
		}
		// Manejo de otros errores, podría ser útil añadir un log aquí
		log.Error("Error buscando usuario por email: ", result.Error)
		return false // Asumimos que el email no está registrado si hay un error distinto
	}

	return true // El usuario existe, el email está registrado
}

func GetUserById(id int) model.User {
	var user model.User

	Db.Where("id = ?", id).First(&user)
	log.Debug("User: ", user)

	return user
}

//Checkear si existe un usuario en el sistema

func CheckUserById(id int) bool {
	var user model.User

	// realza consulta a la base de datos: (con el id proporcionado como parametro)
	result := Db.Where("id = ?", id).First(&user)

	if result.Error != nil {
		return false
	}

	return true
}

func GetUsers() model.Users {
	var users model.Users
	Db.Find(&users)

	log.Debug("Users: ", users)

	return users
}

func InsertUser(user model.User) model.User {
	result := Db.Create(&user)

	if result.Error != nil {
		//TODO Manage Errors
		log.Error("")
	}
	log.Debug("User Created: ", user.Id)
	return user
}

func DeleteUser(id int) error {

	var user model.User

	// Find the user by ID first
	result := Db.Where("id = ?", id).First(&user)

	// Check if the user exists
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Warn("User not found for deletion, ID: ", id)
			return result.Error // User not found
		}
		log.Error("Error finding user for deletion: ", result.Error)
		return result.Error // Other error occurred
	}

	// Proceed to delete the user
	deleteResult := Db.Delete(&user)

	if deleteResult.Error != nil {
		log.Error("Error deleting user: ", deleteResult.Error)
		return deleteResult.Error // Deletion failed
	}

	log.Info("User deleted successfully, ID: ", id)
	return nil // Deletion successful
}

func UpdateUser(user model.User) error {
	if err := Db.Save(user).Error; err != nil {
		fmt.Println(err)
		return err // Return error if the update fails
	}
	return nil
}

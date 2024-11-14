package service

import (
	"fmt"
	userClient "user-api/client"

	"golang.org/x/crypto/bcrypt"

	"user-api/dto"
	"user-api/model"
	e "user-api/utils/errors"
)

type userService struct{}

type userServiceInterface interface {
	GetUsers() (dto.UsersDto, e.ApiError)
	InsertUser(userDto *dto.UserDto) (*dto.UserDto, e.ApiError)
	GetUserById(id int) (*dto.UserDto, e.ApiError)
	DeleteUser(id int) error
	UpdateUser(id int, userDto *dto.UserDto) (model.User, error)
}

var (
	UserService userServiceInterface
	UserClient  userClient.UserClientInterface
)

func init() {
	UserService = &userService{}
	UserClient = &userClient.UserClient{}

}

func (s *userService) GetUserById(id int) (*dto.UserDto, e.ApiError) {
	user := UserClient.GetUserById(id)
	if user.Id == 0 {
		return nil, e.NewBadRequestApiError("Usuario no encontrado")
	}

	userDto := &dto.UserDto{
		Name:     user.Name,
		LastName: user.LastName,
		UserName: user.UserName,
		Phone:    user.Phone,
		Address:  user.Address,
		Email:    user.Email,
		Id:       user.Id,
		Type:     user.Type, // Suponiendo que user.Type es string
	}
	return userDto, nil
}

func (s *userService) GetUsers() (dto.UsersDto, e.ApiError) {
	users := UserClient.GetUsers()
	var usersDto dto.UsersDto

	for _, user := range users {
		userDto := dto.UserDto{
			Name:     user.Name,
			LastName: user.LastName,
			UserName: user.UserName,
			Phone:    user.Phone,
			Address:  user.Address,
			Email:    user.Email,
			Id:       user.Id,
			Type:     user.Type, // Suponiendo que user.Type es string
		}
		usersDto = append(usersDto, userDto)
	}
	return usersDto, nil
}

func (s *userService) InsertUser(userDto *dto.UserDto) (*dto.UserDto, e.ApiError) {
	if UserClient.GetUserByEmail(userDto.Email) {
		return nil, e.NewBadRequestApiError("El email ya está registrado")
	}

	hashedPassword, err := s.HashPassword(userDto.Password)
	if err != nil {
		return nil, e.NewBadRequestApiError("No se puede utilizar esa contraseña")
	}

	user := model.User{
		Name:     userDto.Name,
		LastName: userDto.LastName,
		UserName: userDto.UserName,
		Password: hashedPassword,
		Phone:    userDto.Phone,
		Address:  userDto.Address,
		Email:    userDto.Email,
		Type:     userDto.Type, // Suponiendo que userDto.Type es string
	}

	user = UserClient.InsertUser(user)
	if user.Id == 0 {
		return nil, e.NewBadRequestApiError("Nombre de usuario repetido")
	}

	userDto.Id = user.Id
	return userDto, nil
}

func (s *userService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("no se pudo hashear la contraseña: %w", err)
	}
	return string(hashedPassword), nil
}

func (s *userService) VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

func (s *userService) DeleteUser(id int) error {

	result := UserClient.DeleteUser(id)

	if result != nil {
		return result
	}

	return nil

}

func (s *userService) UpdateUser(id int, userDto *dto.UserDto) (model.User, error) {
	// Check if the user exists
	user := UserClient.GetUserById(id)

	// Update the user's fields with the new data from userDto
	user.Name = userDto.Name
	user.LastName = userDto.LastName
	user.UserName = userDto.UserName
	user.Email = userDto.Email
	user.Phone = userDto.Phone
	user.Address = userDto.Address

	// Save the updated user to the database
	if err := UserClient.UpdateUser(user); err != nil {
		return user, err
	}

	return user, nil
}

package user

import (
	"fmt"
	"net/http"
	"strconv"
	"user-api/dto"
	"user-api/service"
	e "user-api/utils/errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.UserService.DeleteUser(id)

	if err != nil {
		apiErr, ok := err.(e.ApiError)
		if !ok {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, true)

}

func GetUserById(c *gin.Context) {
	log.Debug("User id to load: " + c.Param("id"))

	id, _ := strconv.Atoi(c.Param("id"))
	var userDto *dto.UserDto

	userDto, err := service.UserService.GetUserById(id)

	if err != nil {
		/* if !ok {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		} */
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, userDto)
}

func GetUsers(c *gin.Context) {
	var usersDto dto.UsersDto
	usersDto, err := service.UserService.GetUsers()

	if err != nil {
		apiErr, ok := err.(e.ApiError)
		if !ok {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, usersDto)
}

func UserInsert(c *gin.Context) {
	var userDto dto.UserDto
	err := c.BindJSON(&userDto)

	// Error Parsing json param
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userDtoPtr, er := service.UserService.InsertUser(&userDto)
	// Error del Insert
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": er.Message()})
		log.Debug(er.Message())
		return
	}

	c.JSON(http.StatusCreated, userDtoPtr)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	var userDto dto.UserDto
	if err := c.BindJSON(&userDto); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Datos invalidos"})
		return
	}

	// Call the service layer to update the user
	updatedUser, updateErr := service.UserService.UpdateUser(id, &userDto)
	if updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": updateErr.Error()})
		fmt.Println(updateErr.Error())
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

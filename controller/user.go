package controller

import (
	"fmt"
	"net/http"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(c *gin.Context) {
	fmt.Println("===========")
	val := c.MustGet("user")
	user, ok := val.(models.User)
	if !ok {
		c.JSON(200, "failed di ok")
	}
	fmt.Println(user.Email)

	fmt.Println(val)
	fmt.Println("==========")

	c.JSON(http.StatusOK, gin.H{
		"user": &user,
	})

}
func GetUsers(c *gin.Context) {
	user := []models.User{}
	config.DB.Find(&user)
	c.JSON(200, &user)
}
func InsertUser(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	config.DB.Create(&user)
	c.JSON(200, &user)
}
func DeleteUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id = ?", c.Param("id")).Delete(&user)
	c.JSON(200, &user)
}
func UpdateUser(c *gin.Context) {
	var user models.User
	config.DB.Where("id =?", c.Param("id")).First(&user)
	c.BindJSON(&user)
	config.DB.Save(&user)
	c.JSON(200, &user)
}
func createnewshop(c *gin.Context) {
	var body struct {
		Banner []byte
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": body,
	})
}
func CreateShop(c *gin.Context) {
	var body struct {
		FirstName string
		IsBan     bool
		Role      string
		Email     string
		Password  string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	// gorm.Model
	// Email       string  `json:"email"`
	// Name        string  `json:"name"`
	// Description string  `json:"description"`
	// Status      string  `json:"status"`
	// Banner      string  `json:"banner"`
	// Sales       int     `json:"sales"`
	// Service     float64 `json:"service"`
	//hash the pass
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body2",
		})
		return
	}
	shop := models.Shop{
		Email:       body.Email,
		Name:        body.FirstName,
		Description: "",
		IsBan:       false,
		Banner:      "",
		Sales:       0,
		Service:     0.0,
	}
	//create the user
	user := models.User{
		FirstName: body.FirstName,
		IsBan:     body.IsBan,
		Role:      body.Role,
		Email:     body.Email,
		Password:  string(hash)}
	config.DB.Create(&user)
	config.DB.Create(&shop)
	c.JSON(http.StatusOK, gin.H{
		"message": "Seller Account Successfuly created",
		"user":    user,
	})

}

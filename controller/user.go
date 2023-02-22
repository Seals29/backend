package controller

import (
	"fmt"
	"net/http"
	"net/mail"

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
func EmailValidation(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
func GetUsers(c *gin.Context) {
	user := []models.User{}
	config.DB.Find(&user)
	c.JSON(200, &user)
}
func GetShops(c *gin.Context) {
	shop := []models.Shop{}
	config.DB.Find(&shop)
	c.JSON(200, &shop)
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
	if len(body.FirstName) == 0 {
		c.JSON(200, gin.H{
			"error": "Name cannot be empty",
		})
		return
	}
	if len(body.Email) == 0 || len(body.Password) == 0 {
		c.JSON(200, gin.H{
			"error": "Email or Password cannot be empty",
		})
		return
	}
	email := body.Email
	if EmailValidation(email) {

	} else {
		c.JSON(200, gin.H{
			"error": "Email is not in an email format!",
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
func InsertProduct(c *gin.Context) {
	var body struct {
		Name        string   `json:"name"`
		Type        string   `json:"type"`
		ID          int      `json:"id"`
		Stock       int      `json:"stock"`
		Image       []string `json:"image"`
		Price       int      `json:"price"`
		Description string   `json:"description"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
}
func SetBanUser(c *gin.Context) {
	var body struct {
		IsBan bool `json:"isban"`
		ID    int  `json:"id"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println(body)
	if body.IsBan == false {
		ban := true
		var user models.User
		config.DB.Where("id = ?", body.ID).First(&user)
		user.IsBan = ban
		config.DB.Save(&user)
		c.JSON(200, gin.H{
			"message": "You have banned " + user.FirstName + user.LastName,
		})
		return
	} else {
		ban := false
		var user models.User
		config.DB.Where("id = ?", body.ID).First(&user)
		user.IsBan = ban
		config.DB.Save(&user)
		c.JSON(200, gin.H{
			"message": "You have unbanned " + user.FirstName + user.LastName,
		})
		return
	}
}
func SetBan(c *gin.Context) {
	var body struct {
		IsBan bool   `json:"isban"`
		Email string `json:"email"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	if body.IsBan == false {
		ban := true
		var user models.User
		var shop models.Shop
		config.DB.Where("email =?", body.Email).First(&user)
		config.DB.Where("email =?", body.Email).First(&shop)
		user.IsBan = ban
		shop.IsBan = ban
		config.DB.Save(&user)
		config.DB.Save(&shop)
		c.JSON(200, gin.H{
			"message": "You have banned " + user.FirstName + user.LastName,
		})
		return
	} else {
		ban := false
		var user models.User
		var shop models.Shop
		config.DB.Where("email =?", body.Email).First(&user)
		config.DB.Where("email =?", body.Email).First(&shop)
		user.IsBan = ban
		shop.IsBan = ban
		config.DB.Save(&user)
		config.DB.Save(&shop)
		c.JSON(200, gin.H{
			"message": "You have banned " + user.FirstName + user.LastName,
		})
		return
	}

	// var user models.User
	// config.DB.Where("id =?", c.Param("id")).First(&user)
	// c.BindJSON(&user)
	// config.DB.Save(&user)
	// c.JSON(200, &user)
	// c.JSON(200, gin.H{
	// 	"message": body,
	// })
}

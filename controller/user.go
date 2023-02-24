package controller

import (
	"fmt"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		var checkuser models.User
		checkUniqueEmail := config.DB.Where("email = ?", body.Email).First(&checkuser)
		if checkUniqueEmail.Error == gorm.ErrRecordNotFound {
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
		} else {
			c.JSON(200, gin.H{
				"error": "Email is not Unique",
			})
			return
		}
	} else {
		c.JSON(200, gin.H{
			"error": "Email is not in an email format!",
		})
		return
	}
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
func UpdateShopProfile(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	if len(body.Name) <= 0 {
		c.JSON(200, gin.H{
			"error": "Name cannot be empty",
		})
		return
	}
	if EmailValidation(body.Email) {
		var shop models.Shop
		config.DB.Where("email = ?", body.Email).First(&shop)
		shop.Name = body.Name
		shop.Banner = body.Image
		config.DB.Save(&shop)
		c.JSON(200, gin.H{
			"message": shop,
		})
	} else {
		c.JSON(200, gin.H{
			"error": "Invalid Email",
		})
	}
	return
}
func UpdateShopPassword(c *gin.Context) {
	var body struct {
		Email   string `json:"email"`
		OldPass string `json:"oldpass"`
		NewPass string `json:"newpass"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println("===")
	fmt.Println(body)
	if len(body.OldPass) <= 0 || len(body.NewPass) == 0 {
		c.JSON(200, gin.H{
			"error": "Password cannot be empty",
		})
		return
	}
	var user models.User
	config.DB.Where("email = ?", body.Email).First(&user)
	//dapet akunnya
	// c.JSON(200, &user)
	// if()
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPass)) == nil {
		newHashed, err := bcrypt.GenerateFromPassword([]byte(body.NewPass), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read body2",
			})
			return
		}
		stringHased := string(newHashed)
		user.Password = stringHased
		config.DB.Save(&user)
		c.JSON(200, gin.H{
			"success": "New password successfully changed!",
		})

	} else {
		c.JSON(200, gin.H{
			"error": "Old password is not match!",
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
}
func CreateProduct(c *gin.Context) {
	fmt.Println("ga rusak")

	var body struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		Price       string `json:"price"`
		Email       string `json:"email"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Stock       string `json:"stock"`
		Rating      string `json:"rating"`
		Detail      string `json:"detail"`
		ShopID      string `json:"shopid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println("=====")
	fmt.Println(body)
	fmt.Println("=====")
	price, errprice := strconv.Atoi(body.Price)
	if errprice != nil {
		c.JSON(200, gin.H{
			"error": "Failed convert to int",
		})
	}
	stock, errstock := strconv.Atoi(body.Stock)
	if errstock != nil {
		c.JSON(200, gin.H{
			"error": "Failed convert",
		})
	}
	if len(body.Name) == 0 {
		c.JSON(200, gin.H{
			"error": "Name must not be empty",
		})
		return
	}
	if len(body.Category) == 0 {
		c.JSON(200, gin.H{
			"error": "Name must not be empty",
		})
		return
	}
	if price <= 0 {
		c.JSON(200, gin.H{
			"error": "Price cannot be zero",
		})
		return
	}
	if len(body.Description) < 5 {
		c.JSON(200, gin.H{
			"error": "Description must be at least 5 characters",
		})
		return
	}
	if len(body.Detail) == 0 {
		c.JSON(200, gin.H{
			"error": "Detail cannot be empty!",
		})
		return
	}
	shopid, errshopid := strconv.Atoi(body.ShopID)
	if errshopid != nil {
		c.JSON(200, gin.H{
			"error": "Failed convert to int",
		})
	}
	product := models.Product{
		Name:        body.Name,
		ShopEmail:   body.Email,
		Category:    body.Category,
		Price:       price,
		Description: body.Description,
		Image:       body.Image,
		Rating:      0,
		Stock:       stock,
		Detail:      body.Detail,
		ShopID:      shopid,
	}
	fmt.Println(product)
	config.DB.Create(&product)
	c.JSON(200, gin.H{
		"message": "New Product Successfuly Created!",
	})
	return
}
func getProduct(c *gin.Context) {

}

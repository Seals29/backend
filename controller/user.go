package controller

import (
	"fmt"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
	"github.com/nyaruka/phonenumbers"
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
	// fmt.Println(user.Email)
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
func GetCategoryByShopId(c *gin.Context) {
	shopid := c.Param("shopid")
	var product []string
	config.DB.Model(models.Product{}).Joins(`join shops on products.shop_id = shops.ID`).
		Select("category").Where("shop_id = ?", shopid).Find(&product)
	c.JSON(200, &product)
}
func GetShopIDByUserID(c *gin.Context) {
	userid := c.Query("userid")
	userID, err := strconv.Atoi(userid)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid conversion!",
		})
		return
	}
	var user models.User
	config.DB.Where("id = ?", userID).First(&user)
	var shop models.Shop
	config.DB.Where("email = ?", user.Email).First(&shop)
	c.JSON(200, &shop)
}
func GetUserShopId(c *gin.Context) {
	id := c.Param("id")
	var shop models.Shop
	config.DB.Where("id = ?", id).First(&shop)
	var user models.User
	config.DB.Where("email = ?", shop.Email).First(&user)
	// fmt.Println(user)
	// fmt.Println("id : " + id + " shop : " + shop.Name)
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
		SubCategory string `json:"subcategory"`
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
		SubCategory: body.SubCategory,
	}
	// fmt.Println(product)
	config.DB.Create(&product)
	c.JSON(200, gin.H{
		"message": "New Product Successfuly Created!",
	})
	return
}
func getProduct(c *gin.Context) {

}
func UpdateAccountEmail(c *gin.Context) {
	var body struct {
		UserID string `json:"userid"`
		Email  string `json:"email"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	userid, err := strconv.Atoi(body.UserID)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}
	if EmailValidation(body.Email) {
		var checkuser models.User

		checkUniqueEmail := config.DB.Where("email = ?", body.Email).First(&checkuser)
		if checkUniqueEmail.Error == gorm.ErrRecordNotFound {
			checkuser.Email = body.Email
			var userEmail models.User
			config.DB.Where("id = ?", userid).First(&userEmail)
			userEmail.Email = body.Email
			// fmt.Println(userEmail)
			// fmt.Println(checkuser)
			config.DB.Save(&userEmail)
			c.JSON(200, gin.H{
				"message": "Email has been changed successfully!",
			})
			return

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
func UpdateAccountPhoneNumber(c *gin.Context) {
	var body struct {
		UserID      string `json:"userid"`
		PhoneNumber string `json:"phonenumber"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	parsed, err := phonenumbers.Parse(body.PhoneNumber, "ID")

	if err != nil || !phonenumbers.IsValidNumber(parsed) {
		c.JSON(200, gin.H{
			"error": "Invalid Phone Number",
		})
		return
	}
	last := phonenumbers.Format(parsed, phonenumbers.INTERNATIONAL)
	var user models.User
	userID, err := strconv.Atoi(body.UserID)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "Failed convert userid",
		})
		return
	}
	config.DB.Where("ID = ?", userID).First(&user)

	if user.PhoneNumber == last {
		c.JSON(200, gin.H{
			"message": "You did not changed anything!",
		})
		return

	}
	user.PhoneNumber = last
	config.DB.Save(&user)
	c.JSON(200, gin.H{
		"message": "Phone Number successfully changed!",
	})
	return

}
func SubscribeFromHome(c *gin.Context) {
	var body struct {
		UserEmail string `json:"useremail"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	var user models.User
	config.DB.Where("email = ?", body.UserEmail).First(&user)
	var subscribeUser models.UserSubscribe
	config.DB.Where("user_email = ?", body.UserEmail).First(&subscribeUser)
	// config.DB.Create(&subscribeUser)
	c.JSON(200, &subscribeUser)
	if subscribeUser.ID == 0 {
		var newSubscriber models.UserSubscribe
		newSubscriber.UserEmail = body.UserEmail
		config.DB.Create(&newSubscriber)
		c.JSON(200, &newSubscriber)
		return
	} else {
		c.JSON(200, gin.H{
			"error": "Invalid email!",
		})
		return
	}
}
func UpdateAccountPassword(c *gin.Context) {
	var body struct {
		UserID      string `json:"userid"`
		OldPassword string `json:"oldpassword"`
		NewPassword string `json:"newpassword"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	if body.NewPassword == body.OldPassword {
		c.JSON(200, gin.H{
			"error": "cannot be the same with old password",
		})
		return
	}
	if len(body.OldPassword) <= 0 || len(body.NewPassword) <= 0 {
		c.JSON(200, gin.H{
			"error": "Password cannot be empty",
		})
		return
	}
	if len(body.NewPassword) <= 5 {
		c.JSON(200, gin.H{
			"error": "New Password must be above 5 characters!",
		})
		return
	}
	var user models.User
	config.DB.Where("ID = ?", body.UserID).First(&user)
	if user.ID == 0 {
		c.JSON(200, gin.H{
			"error": "User not Found!",
		})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)) == nil {
		newHashed, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
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
	c.JSON(200, &user)

}
func GetSubscribeStatus(c *gin.Context) {
	var body struct {
		UserID string `json:"userid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	var currUser models.User
	var subscribeCheck models.UserSubscribe
	config.DB.Where("id = ?", body.UserID).First(&currUser)
	config.DB.Where("user_email = ?", currUser.Email).First(&subscribeCheck)
	if subscribeCheck.ID == 0 {
		//ga ada subscribe
		var newSubscribe models.UserSubscribe
		newSubscribe.UserEmail = currUser.Email
		config.DB.Create(&newSubscribe)
		c.JSON(200, gin.H{
			"message": "Not Subscribed",
		})
		return
	} else {
		//subscribed
		c.JSON(200, gin.H{
			"error": "Subscribed",
		})
		return
	}
}

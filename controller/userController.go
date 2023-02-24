package controller

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nyaruka/phonenumbers"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(c *gin.Context) {
	//get email or pass4

	var body struct {
		FirstName   string
		LastName    string
		IsBan       bool
		Role        string
		Email       string
		Password    string
		PhoneNumber string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}

	parsed, err := phonenumbers.Parse(body.PhoneNumber, "ID")
	fmt.Println("masuk")
	last := ""

	if len(body.PhoneNumber) > 0 {
		fmt.Println(parsed)
		if err != nil || !phonenumbers.IsValidNumber(parsed) {
			fmt.Println(parsed)
			c.JSON(200, gin.H{
				"error": "Invalid Phone Number",
			})
			return
		}
		last = phonenumbers.Format(parsed, phonenumbers.INTERNATIONAL)

	}
	fmt.Println(last)
	fmt.Println(body.PhoneNumber)
	if len(body.FirstName) == 0 {
		c.JSON(200, gin.H{
			"error": "First Name must not be empty",
		})
		return
	}
	if len(body.LastName) == 0 {
		c.JSON(200, gin.H{
			"error": "First Name must not be empty",
		})
		return
	}
	if len(body.Email) == 0 || len(body.Password) == 0 {
		c.JSON(200, gin.H{
			"error": "Email Or Password must not be empty",
		})
		return
	}
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

		//create the user
		user := models.User{
			FirstName:   body.FirstName,
			LastName:    body.LastName,
			IsBan:       body.IsBan,
			Role:        body.Role,
			Email:       body.Email,
			PhoneNumber: last,
			Password:    string(hash)}
		config.DB.Create(&user)
		log.Println(&body)

		log.Println(body.Password)
		// if res != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{
		// 		"error": "Failed to read body3",
		// 	})
		// 	return
		// }
		//respon

		//generate JWT TOKEN
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
		TokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid create token",
			})
			return
		}

		//sign in and get the complete encoded token as as tring using the secret
		//sent it back
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorized", TokenString, 3600*24, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"token": TokenString,
		})
	} else {
		c.JSON(200, gin.H{
			"error": "Email is not Unique",
		})
		return

	}
	//hash the pass

}
func Getproduct(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	products := []models.Product{}
	config.DB.Where("shop_email = ? ", body.Email).Find(&products)
	c.JSON(200, &products)
}
func validates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
	})
}
func LoginAssistance(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println("====body===")
	fmt.Println(body)
	if len(body.Email) == 0 {
		c.JSON(200, gin.H{
			"error": "Email cannot be empty!",
		})
		return
	}
	var user models.User
	config.DB.Where("email = ?", body.Email).First(&user)
	if user.ID == 0 {
		c.JSON(200, gin.H{
			"error": "Email not found!",
		})
		return
	} else {
		auth := smtp.PlainAuth(
			"",
			"lionelriyadi13@gmail.com",
			"tkfuhgsqnhhidrnb",
			"smtp.gmail.com",
		)
		num, err := rand.Int(rand.Reader, big.NewInt(900000))
		if err != nil {
			c.JSON(200, gin.H{
				"error": err,
			})
			return
		}
		intcode := int(num.Int64() + 100000)
		msg := "Subject: Login Assistance\n\nHere is your Verification code that last 5 minutes. Code : " + strconv.Itoa(intcode)
		smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			"lionelriyadi13@gmail.com",
			[]string{body.Email},
			[]byte(msg),
		)
		fmt.Println(intcode)
		fmt.Println("=======")
		forgotuser := models.ForgotUser{
			ExpiredDate: time.Now().Add(time.Minute * 5),
			ResetCode:   int(intcode),
			UserID:      int(user.ID),
			Used:        false,
		}
		fmt.Println(forgotuser)
		config.DB.Create(&forgotuser)
		c.JSON(http.StatusOK, gin.H{
			"message": "New Verification code has been sent!",
		})
		return
	}
}
func ResetPassword(c *gin.Context) {
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
	fmt.Println("==")
	fmt.Println(body)
	c.JSON(200, &body)
}
func ResendForgotPassword(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	var user models.User
	config.DB.Where("email = ?", body.Email).First(&user)
	config.DB.Model(&models.ForgotUser{}).Where("user_id = ?", user.ID).Update("used", true)
	fmt.Println("===")
	// if res.Error != nil {
	// 	c.JSON(200, gin.H{
	// 		"error": res.Error,
	// 	})
	// 	return
	// }
	//==========
	auth := smtp.PlainAuth(
		"",
		"lionelriyadi13@gmail.com",
		"tkfuhgsqnhhidrnb",
		"smtp.gmail.com",
	)
	num, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})
		return
	}
	intcode := int(num.Int64() + 100000)
	msg := "Subject: Login Assistance\n\nHere is your Verification code that last 5 minutes. Code : " + strconv.Itoa(intcode)
	smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"lionelriyadi13@gmail.com",
		[]string{body.Email},
		[]byte(msg),
	)

	//=
	forgotuser := models.ForgotUser{
		ExpiredDate: time.Now().Add(time.Minute * 5),
		ResetCode:   int(intcode),
		UserID:      int(user.ID),
		Used:        false,
	}
	config.DB.Create(&forgotuser)
	c.JSON(200, gin.H{
		"message": "New Verification code has been sent!",
	})
	return

}
func UpdateUserPassword(c *gin.Context) {
	var body struct {
		Id       string `json:"id"`
		Password string `json:"password"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println(body)
	var user models.User
	config.DB.Where("id = ?", body.Id).First(&user)
	fmt.Println(user)
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body2",
		})
		return
	}
	fmt.Println(string(hash))
	user.Password = string(hash)
	config.DB.Save(&user)
	c.JSON(200, gin.H{
		"message": "Password has been changed!",
	})
}
func VerifForgotPassword(c *gin.Context) {
	var body struct {
		ResetCode string `json:"resetcode"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	var forgotuser models.ForgotUser
	config.DB.Where("reset_code = ?", body.ResetCode).First(&forgotuser)
	fmt.Println(forgotuser)
	fmt.Println(body)
	if forgotuser.Used == false {

		if time.Now().Before(forgotuser.ExpiredDate) {
			//belom expired
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub": forgotuser.UserID,
				"exp": time.Now().Add(time.Hour * 24).Unix(),
			})
			tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
			if err != nil {
				c.JSON(200, gin.H{
					"error": "Invalid Create Token!",
				})
				return
			}
			c.SetSameSite(http.SameSiteLaxMode)
			c.JSON(200, gin.H{
				"token": tokenString,
			})
			return
		} else {
			//expired
			c.JSON(200, gin.H{
				"error": "Code is Expired!",
			})
			return
		}
	}

}
func UserSubscribe(c *gin.Context) {
	//get email or pass
	var body struct {
		UserEmail string `json:"useremail"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	Subscribe := models.UserSubscribe{
		// Model:     gorm.Model{},
		UserEmail: body.UserEmail,
	}
	fmt.Println(Subscribe)
	config.DB.Create(&Subscribe)
	c.JSON(http.StatusOK, gin.H{
		"Message ": "User has been subscribed",
	})
	// //create the user
	// user := models.User{
	// 	FirstName: body.FirstName,
	// 	LastName:  body.LastName,
	// 	IsBan:     body.IsBan,
	// 	Role:      body.Role,
	// 	Email:     body.Email,
	// 	Password:  string(hash)}
	// config.DB.Create(&user)
	// log.Println(&body)
}
func Login(c *gin.Context) {
	//get email and pass off req body
	var body struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	pass := body.Password
	if len(pass) == 0 || len(body.Email) == 0 {
		c.JSON(200, gin.H{
			"Error": "Email or Password cannot be empty",
		})
		return
	}

	//lookup requested user
	var user models.User
	config.DB.Model(models.User{}).Where("email = ?", body.Email).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or password1",
		})
		return
	}
	//compare sent in pass with saved user pass hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BCrypt Failed",
		})
		return
	}
	//generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	TokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid create token",
		})
		return
	}

	//sign in and get the complete encoded token as as tring using the secret
	//sent it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorized", TokenString, 3600*24, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": TokenString,
	})
}
func Announce(c *gin.Context) {
	var body struct {
		Subject string `json:"subject"`
		Message string `json:"message"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	if len(body.Subject) == 0 || len(body.Message) == 0 {
		c.JSON(200, gin.H{
			"error": "Subject/ Message cannot be empty",
		})
		return
	}
	user := []models.UserSubscribe{}
	config.DB.Find(&user)
	allEmails := []string{}
	for _, u := range user {
		email := u.UserEmail
		fmt.Println("Email:", email)
		allEmails = append(allEmails, email)
	}

	//email dapet di body
	auth := smtp.PlainAuth(
		"",
		"lionelriyadi13@gmail.com",
		"tkfuhgsqnhhidrnb",
		"smtp.gmail.com",
	)

	msg := "Subject: " + body.Subject + "\n" + body.Message
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"lionelriyadi13@gmail.com",
		allEmails,
		[]byte(msg),
	)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "",
	})
}
func NotifyShop(c *gin.Context) {
	var body struct {
		Email string `json:"email"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}

	auth := smtp.PlainAuth(
		"",
		"lionelriyadi13@gmail.com",
		"tkfuhgsqnhhidrnb",
		"smtp.gmail.com",
	)

	msg := "Subject: Your Shop Account has been verified" + "\n" + "You've been verified the shop account congrats"
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"lionelriyadi13@gmail.com",
		[]string{body.Email},
		[]byte(msg),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Shop account has been verified",
	})
}
func GetDetailProduct(c *gin.Context) {
	var body struct {
		ProductID string `json:"productid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println("===")
	fmt.Println(body)
	var product models.Product
	config.DB.Where("id = ?", body.ProductID).First(&product)
	c.JSON(200, &product)
	return
}
func UpdateProduct(c *gin.Context) {
	var body struct {
		ProductID   string `json:"productid"`
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Detail      string `json:"detail"`
		Price       string `json:"price"`
		Stock       string `json:"stock"`
		Image       string `json:"image"`
		Email       string `json:"email"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	// c.JSON(200, gin.H{
	// "message": body,
	// })
	var product models.Product
	config.DB.Where("id = ?", body.ProductID).First(&product)

	product.Name = body.Name
	product.Category = body.Category
	newprice, errprice := strconv.Atoi(body.Price)
	if errprice != nil {
		c.JSON(200, gin.H{
			"error": "Invalid conversion string to int1",
		})
		return
	}
	product.Price = newprice
	newstock, errstock := strconv.Atoi(body.Stock)
	if errstock != nil {
		c.JSON(200, gin.H{
			"error": "invalid conversion ",
		})
		return
	}
	product.Stock = newstock
	fmt.Println(body.Email + "++")
	product.Detail = body.Detail
	product.Description = body.Description
	if len(body.Image) == 0 {
		fmt.Println("kosong")
		config.DB.Save(&product)
	} else {
		fmt.Println("ga kosong")
		product.Image = body.Image
		config.DB.Save(&product)
	}
	c.JSON(200, gin.H{
		"message": "New Product Successfully edited",
	})
}

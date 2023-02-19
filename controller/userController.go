package controller

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	//get email or pass
	var body struct {
		FirstName string
		LastName  string
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

	//hash the pass
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body2",
		})
		return
	}

	//create the user
	user := models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		IsBan:     body.IsBan,
		Role:      body.Role,
		Email:     body.Email,
		Password:  string(hash)}
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
}
func validates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in",
	})
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
	user := []models.UserSubscribe{}
	config.DB.Find(&user)
	allEmails := []string{}
	for _, u := range user {
		email := u.UserEmail
		// Do something with the email, such as send an email to this address
		fmt.Println("Email:", email)
		allEmails = append(allEmails, email)
	}
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

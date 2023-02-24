package controller

import (
	"crypto/rand"
	"fmt"
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
)

func ResendOneTime(c *gin.Context) {
	var body struct {
		Email string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	if len(body.Email) == 0 {
		c.JSON(200, gin.H{
			"error": "Email cannot be empty",
		})
		return
	}
	auth := smtp.PlainAuth(
		"",
		"lionelriyadi13@gmail.com",
		"tkfuhgsqnhhidrnb",
		"smtp.gmail.com",
	)
	num, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Message": "Failed Error",
		})
		return
	}
	intcode := int(num.Int64() + 100000)
	code := strconv.Itoa(intcode)
	msg := "Subject: Newegg ONE TIME CODES! \n\nYou've request ONE TIME CODE, here is your code : " + code + "\n"
	var user models.User
	config.DB.Where("email = ?", body.Email).First(&user)
	fmt.Println(user.ID)
	if user.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "user not found",
		})
		return
	}
	// var oldresetuser models.ResetUser
	fmt.Println(user.ID)
	res := config.DB.Model(&models.ResetUser{}).Where("user_id = ?", user.ID).Update("used", true)
	if res.Error != nil {
		c.JSON(200, gin.H{
			"error": res.Error,
		})
		return
	}
	resetuser := models.ResetUser{
		ExpiredDate: time.Now().Add(time.Minute * 5),
		ResetCode:   int(intcode),
		UserID:      int(user.ID),
		Used:        false,
	}
	config.DB.Create(&resetuser)
	err2 := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"lionelriyadi13@gmail.com",
		[]string{body.Email},
		[]byte(msg),
	)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"resetuser": resetuser,
	})
}
func SendOneTime(c *gin.Context) {
	fmt.Println("===========")
	if c.Request.Method == "POST" {
		num, err := rand.Int(rand.Reader, big.NewInt(900000))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"Message": "Failed Error",
			})
			return
		}
		fmt.Println(num.Int64() + 100000) //6 digit code//save into database

		//send email
		var body struct {
			Email string `json:"email"`
		}

		if c.Bind(&body) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read body1",
			})
			return
		}
		if len(body.Email) == 0 {
			c.JSON(200, gin.H{
				"error": "Email cannot be empty",
			})
			return
		}
		fmt.Println(body.Email)
		//send email
		auth := smtp.PlainAuth(
			"",
			"lionelriyadi13@gmail.com",
			"tkfuhgsqnhhidrnb",
			"smtp.gmail.com",
		)
		intcode := int(num.Int64() + 100000)
		code := strconv.Itoa(intcode)
		msg := "Subject: Newegg ONE TIME CODES! \n\nYou've request ONE TIME CODE, here is your code : " + code + "\n"
		var user models.User
		config.DB.Where("email = ?", body.Email).First(&user)
		fmt.Println(user.ID)
		if user.ID == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "user not found",
			})
			return
		}
		resetuser := models.ResetUser{
			ExpiredDate: time.Now().Add(time.Minute * 5),
			ResetCode:   int(intcode),
			UserID:      int(user.ID),
			Used:        false,
		}
		config.DB.Create(&resetuser)
		err2 := smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			"lionelriyadi13@gmail.com",
			[]string{body.Email},
			[]byte(msg),
		)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"resetuser": resetuser,
		})
	} else {
		c.JSON(200, gin.H{
			"error": "Request method should be POST!",
		})
		return
	}
}

func OneTimeCode(c *gin.Context) {
	if c.Request.Method == "POST" {
		var body struct {
			ResetCode string
		}
		if c.Bind(&body) != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read body1",
			})
			return
		}

		//dh dapet resetcode
		var resetuser models.ResetUser
		fmt.Println("resetcode : " + body.ResetCode)
		config.DB.Where("reset_code = ?", body.ResetCode).First(&resetuser)
		fmt.Println(resetuser)
		if resetuser.Used == false {
			//user dapet
			fmt.Println(resetuser.ExpiredDate)
			fmt.Println(time.Now())
			if time.Now().Before(resetuser.ExpiredDate) {
				fmt.Println("Expired code belom expired")
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": resetuser.UserID,
					"exp": time.Now().Add(time.Hour * 24).Unix(),
				})
				TokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": "Invalid create token",
					})
					return
				}
				c.SetSameSite(http.SameSiteLaxMode)
				c.SetCookie("Authorized", TokenString, 3600*24, "", "", false, true)
				c.JSON(http.StatusOK, gin.H{
					"token": TokenString,
				})
				fmt.Println(TokenString)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": "Code is Expired",
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Code is Expired2",
			})
		}

		return

	}
}
func UpdatePassword(c gin.Context) {
	var body struct {
		Password string `json:"password"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to request body",
		})
		return
	}
	fmt.Println(body)
}

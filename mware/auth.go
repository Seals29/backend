package mware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func UserAuth(c *gin.Context) {
	body, errs := ioutil.ReadAll(c.Request.Body)
	if errs != nil {
		c.AbortWithError(http.StatusBadRequest, errs)
		return
	}
	fmt.Println(string(body))
	var data map[string]string
	errors := json.Unmarshal([]byte(string(body)), &data)
	if errors != nil {
		panic(errs)
	}
	tokenString := data["cookies"]
	fmt.Println(tokenString)
	fmt.Println("=========")
	// //decode
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	fmt.Println(err)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println(err)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		//cari di db
		var user models.User
		config.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		fmt.Println("users")
		fmt.Println(user)
		c.Set("user", user)
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}

}

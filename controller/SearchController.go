package controller

import (
	"fmt"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func SearchProduct(c *gin.Context) {
	fmt.Println("===searchingg...====")
	searchedName := c.Query("val")
	products := []models.Product{}
	fmt.Println(searchedName)
	config.DB.Where("name LIKE ?", "%"+searchedName+"%").Find(&products)
	c.JSON(200, &products)
}
func SaveSearchQuery(c *gin.Context) {
	jwttoken := c.Query("token")
	query := c.Query("query")

	currUser, err := extractUserIDFromToken(jwttoken)
	curruserID, errcur := strconv.Atoi(currUser)
	if err != nil || errcur != nil {
		c.JSON(400, gin.H{
			"error": "Invalid conversion!",
		})
		return
	}
	var user models.User
	config.DB.Where("id = ?", curruserID).First(&user)
	var countCheck int64
	config.DB.Table("save_queries").Where("user_id = ?", curruserID).Count(&countCheck)
	fmt.Println(countCheck)
	if countCheck >= 10 {
		c.JSON(200, gin.H{
			"error": "Your query has reach maxed out !",
		})
		return
	} else {
		var userQuery models.UserQuery
		userQuery.UserID = int(user.ID)
		userQuery.SearchQuery = query
		config.DB.Create(&userQuery)
		c.JSON(200, gin.H{
			"message": "Your Search Query has been saved successfully!",
		})
		return
	}

}

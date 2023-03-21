package controller

import (
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func GetStoreNotification(c *gin.Context) {
	jwttoken := c.Query("token")
	currUser, err := extractUserIDFromToken(jwttoken)
	curruserID, errcur := strconv.Atoi(currUser)
	if err != nil || errcur != nil {
		c.JSON(400, gin.H{
			"error": "Invalid conversion!",
		})
		return

	}
	allNotif := []models.StoreNotification{}
	var user models.User
	var shopuser models.Shop
	config.DB.Where("id = ?", curruserID).First(&user)
	config.DB.Where("email = ?", user.Email).First(&shopuser)
	config.DB.Where("shop_id= ?", shopuser.ID).Find(&allNotif)
	c.JSON(200, &allNotif)

}
func GetAnnounceNotification(c *gin.Context) {
	allnotif := []models.AnnouncementNotification{}
	config.DB.Find(&allnotif)
	c.JSON(200, &allnotif)
}

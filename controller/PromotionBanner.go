package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func TestData(c *gin.Context) {

	fmt.Println("haloghasdas")

	banner := []models.PromotionBanner{}
	config.DB.Find(&banner)
	fmt.Println(banner)
	c.JSON(200, &banner)
	// var banner models.PromotionBanner
	// config.DB.First(banner)
	// // config.DB.Find(banner)
	// fmt.Println(banner)
	// c.JSON(200, &banner)
}
func GetAllPromotionBanner(c *gin.Context) {
	// var banner models.PromotionBanner
	// banner := []models.PromotionBanner{}

}
func DeletePromotionBanner(c *gin.Context) {
	var body struct {
		BannerID string `json:"bannerid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	bannerid, errb := strconv.Atoi(body.BannerID)
	if errb != nil {
		c.JSON(400, gin.H{
			"error": errb,
		})
		return
	}
	var deleteBanner models.PromotionBanner
	config.DB.Where("id =?", bannerid).First(&deleteBanner)
	config.DB.Where("id = ?", bannerid).Delete(&deleteBanner)
	c.JSON(200, gin.H{
		"message": "Banner has been deleted successfully!",
	})
}
func AddPromotionBanner(c *gin.Context) {
	var body struct {
		BannerImage string `json:"bannerimage"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	newBanner := models.PromotionBanner{
		PromotionImage: body.BannerImage,
	}
	config.DB.Create(&newBanner)
	c.JSON(200, gin.H{
		"message": "New Banner has been added successfully!",
	})

}

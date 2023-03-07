package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func GetAllCart(c *gin.Context) {
	carts := []models.Cart{}
	config.DB.Find(&carts)
	c.JSON(200, &carts)

}
func DeleteProductInCart(c *gin.Context) {
	var body struct {
		UserID    string `json:"userid"`
		ProductID string `json:"productid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to request body",
		})
		return
	}
	productid, errpr := strconv.Atoi(body.ProductID)
	userid, errus := strconv.Atoi(body.UserID)
	if errus != nil || errpr != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Parsing",
		})
		return

	}
	var carts models.Cart
	config.DB.Where("user_id = ?", userid).Where("product_id= ?", productid).First(&carts)
	fmt.Println(carts)
	config.DB.Delete(&carts)
	c.JSON(200, gin.H{
		"message": "This Item has been deleted successfully!",
	})

}
func MoveCartToSave(c *gin.Context) {
	var body struct {
		UserID    string `json:"userid"`
		ProductID string `json:"productid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to request body",
		})
		return
	}
	productid, errpr := strconv.Atoi(body.ProductID)
	userid, errus := strconv.Atoi(body.UserID)
	if errus != nil || errpr != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Parsing",
		})
		return

	}
	var checkCart models.Cart
	config.DB.Where("product_id = ?", productid).Where("user_id = ?", userid).First(&checkCart)
	if checkCart.ID == 0 {
		c.JSON(200, gin.H{
			"error": "Cart Not Found!",
		})
		return
	}
	newSave := models.SaveLater{
		ProductID: checkCart.ProductID,
		UserID:    checkCart.UserID,
		Quantity:  checkCart.Quantity,
	}
	fmt.Println(newSave)
	config.DB.Create(&newSave)
	config.DB.Delete(&checkCart)
	c.JSON(200, gin.H{
		"message": "Cart has been moved to saved later!",
	})
}
func GetAllSavelater(c *gin.Context) {
	savelater := []models.SaveLater{}
	config.DB.Find(&savelater)
	c.JSON(200, &savelater)
}

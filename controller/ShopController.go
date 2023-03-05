package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func ViewCart(c *gin.Context) {
	carts := []models.Cart{}
	config.DB.Find(&carts)
	c.JSON(200, &carts)
}
func GetAllSubCategory(c *gin.Context) {
	subcategory := []models.ProductSubCategory{}
	config.DB.Find(&subcategory)
	c.JSON(200, &subcategory)
}
func GetAllCategory(c *gin.Context) {
	category := []models.ProductCategory{}
	config.DB.Find(&category)
	c.JSON(200, &category)
}
func LoadProductByPage(c *gin.Context){
	page := c.Query("page")
	pageSize := c.Query("pagesize")
	Page ,errp:= strconv.Atoi(page)
	fmt.Println(page)
	pagesize,errps := strconv.Atoi(pageSize)
	if errps !=nil||errp!=nil {
		c.JSON(200,gin.H{
			"error":"invalid Parsing!",
		})
		return
	}
	products:= []models.Product{}
	config.DB.Offset(pagesize * Page).Limit(pagesize).Find(&products)
	fmt.Println(products)

	// fmt.Println(products)
	c.JSON(200,&products)
}
func GetProductByCategory(c *gin.Context) {
	shopid := c.Param("id")
	products := []models.Product{}
	config.DB.Where("shop_id = ?", shopid).Find(&products)

	// shop := []models.Shop{}
	c.JSON(200, &products)
}
func InsertCart(c *gin.Context) {
	var body struct {
		ProductID string `json:"productid"`
		Quantity  string `json:"quantity"`
		UserID    string `json:"userid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	var cart models.Cart
	// fmt.Println(body)
	productId, errid := strconv.Atoi(body.ProductID)
	if errid != nil {
		c.JSON(200, gin.H{
			"error": "error parsing productId",
		})
		return
	}
	cart.ProductID = productId
	Quantity, errstock := strconv.Atoi(body.Quantity)
	if errstock != nil {
		c.JSON(200, gin.H{
			"error": "error parsing quantity",
		})
		return
	}
	cart.Quantity = Quantity

	UserID, erruserId := strconv.Atoi(body.UserID)
	if erruserId != nil {
		c.JSON(200, gin.H{
			"error": "error parsing UserID",
		})
		return
	}
	cart.UserID = UserID
	var checkcart models.Cart
	config.DB.Where("product_id = ?", productId).First(&checkcart)
	if checkcart.ID == 0 {
		config.DB.Create(&cart)
	} else {
		checkcart.Quantity = checkcart.Quantity + Quantity
		config.DB.Save(&checkcart)
	}

	c.JSON(200, &checkcart)

}

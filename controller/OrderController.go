package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {
	allOrders := []models.Order{}
	config.DB.Find(&allOrders)
	c.JSON(200, &allOrders)
}
func GetOrderByOrderID(c *gin.Context) {
	orderID := c.Query("orderid")
	var order models.Order
	orderid, errc := strconv.Atoi(orderID)
	if errc != nil {
		c.JSON(404, gin.H{
			"error": "Invalid Parsing!",
		})
		return
	}
	config.DB.Where("id = ?", orderid).First(&order)
	c.JSON(200, &order)
}
func GetAllOrdersByUserID(c *gin.Context) {
	var body struct {
		JwtToken string `json:"jwttoken"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println("jwt : " + body.JwtToken)
	currUser, err := extractUserIDFromToken(body.JwtToken)
	fmt.Println(body.JwtToken)
	var shop models.Shop
	var user models.User
	curruserID, errcur := strconv.Atoi(currUser)
	config.DB.Where("id =?", curruserID).First(&user)
	config.DB.Where("email = ?", user.Email).First(&shop)
	if errcur != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}
	fmt.Println(err)
	allOrders := []models.Order{}
	config.DB.Where("user_id = ?", curruserID).Find(&allOrders)
	c.JSON(200, &allOrders)
	// config.DB.

}
func GetShopNameByShopID(c *gin.Context) {
	var body struct {
		ShopID string `json:"shopid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	var shop models.Shop
	shopid, errshop := strconv.Atoi(body.ShopID)
	if errshop != nil {
		c.JSON(200, gin.H{
			"error": "Failed parsing!",
		})
		return
	}

	config.DB.Where("id = ?", shopid).First(&shop)
	fmt.Println(shop)
	c.JSON(200, &shop)
}
func GetUseDetailByOrderID(c *gin.Context) {
	var body struct {
		OrderID string `json:"orderid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	orderid, errorder := strconv.Atoi(body.OrderID)
	if errorder != nil {
		c.JSON(400, gin.H{
			"error": errorder,
		})
		return
	}
	type UserWithIDAndName struct {
		ID        uint
		FirstName string
	}
	var users []UserWithIDAndName
	config.DB.Table("orders").Joins("JOIN users ON users.id = orders.user_id").Where("orders.id = ?", orderid).
		Select("users.id, users.first_name").Scan(&users)
	fmt.Println(users)
	c.JSON(200, &users)
}
func GetAllOrdersByShopID(c *gin.Context) {
	var body struct {
		JwtToken string `json:"jwttoken"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println(body.JwtToken)
	currUser, err := extractUserIDFromToken(body.JwtToken)
	fmt.Println(currUser)
	fmt.Println(body.JwtToken)
	var user models.User

	var shop models.Shop
	fmt.Println(body.JwtToken)
	curruserID, errcur := strconv.Atoi(currUser)
	config.DB.Where("id = ?", curruserID).First(&user)
	fmt.Println(curruserID)
	config.DB.Where("email =?", user.Email).First(&shop)
	if errcur != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}

	fmt.Println(err)
	fmt.Println(shop.ID)
	allOrders := []models.Order{}
	config.DB.Where("shop_id = ?", shop.ID).Find(&allOrders)
	c.JSON(200, &allOrders)
}
func AddAllOrderDetailToCart(c *gin.Context) {
	var body struct {
		OrderID  string `json:"orderid"`
		JwtToken string `json:"jwttoken"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}

	currUser, err := extractUserIDFromToken(body.JwtToken)
	curruserID, errcur := strconv.Atoi(currUser)
	if err != nil || errcur != nil {
		c.JSON(400, gin.H{
			"error": "invalid parsing",
		})
		return
	}
	// fmt.Println(currUser)
	fmt.Println(body.JwtToken)
	// var user models.User

	// var shop models.Shop
	// fmt.Println(body.JwtToken)

	fmt.Println("===")
	fmt.Println(body)
	orderid, errorder := strconv.Atoi(body.OrderID)
	if errorder != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Parsing!",
		})
		return
	}
	allOrders := []models.OrderDetail{}
	config.DB.Where("order_id = ?", orderid).Find(&allOrders)
	fmt.Println(allOrders)
	// c.JSON(200, &allOrders)
	for _, order := range allOrders {
		var newCart models.Cart
		newCart.ProductID = order.ProductID
		newCart.Quantity = order.Quantity
		newCart.UserID = curruserID
		config.DB.Create(&newCart)
		fmt.Println(newCart)
	}
	c.JSON(200, gin.H{
		"message": "You have add all items from the order into the cart again!",
	})
	return
}

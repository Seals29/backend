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
func LoadProductByPage(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("pagesize")
	Page, errp := strconv.Atoi(page)
	fmt.Println(page)
	pagesize, errps := strconv.Atoi(pageSize)
	if errps != nil || errp != nil {
		c.JSON(200, gin.H{
			"error": "invalid Parsing!",
		})
		return
	}
	products := []models.Product{}
	config.DB.Offset(pagesize * Page).Limit(pagesize).Find(&products)
	fmt.Println(products)

	// fmt.Println(products)
	c.JSON(200, &products)
}
func GetProductByCategory(c *gin.Context) {
	shopid := c.Param("category")
	products := []models.Product{}
	config.DB.Where("shop_id = ?", shopid).Find(&products)
	fmt.Println(shopid)
	// shop := []models.Shop{}
	var categories []string
	if err := config.DB.Table("product_categories").
		Joins("JOIN products ON product_categories.name = products.category").
		Where("products.shop_id = ?", shopid).
		Distinct("product_categories.name").
		Pluck("product_categories.name", &categories).
		Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get product categories"})
		return
	}
	// stmt:=
	// fmt.Println("====getproductcategorybyshopid")
	fmt.Println(categories)

	c.JSON(200, &categories)
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
func GetReviewsByShop(c *gin.Context) {
	shopId := c.Query("shopid")
	shopid, errshop := strconv.Atoi(shopId)
	fmt.Println(shopId)
	if errshop != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Conversion!",
		})
		return
	}
	reviews := []models.ShopReview{}
	config.DB.Where("shop_id = ?", shopid).Find(&reviews)
	c.JSON(200, &reviews)
}
func AddNewReviewShop(c *gin.Context) {
	var body struct {
		UserID string `json:"userid"`
		ShopID string `json:"shopid"`
		Review string `json:"review"`
		Rating string `json:"rating"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	userid, errus := strconv.Atoi(body.UserID)
	shopid, errshop := strconv.Atoi(body.ShopID)
	rating, errrate := strconv.Atoi(body.Rating)
	if errus != nil || errshop != nil || errrate != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Parsing",
		})
		return
	}
	var user models.User
	config.DB.Where("id = ?", userid).First(&user)
	newReview := models.ShopReview{
		UserID:        userid,
		Rating:        rating,
		ReviewComment: body.Review,
		ShopID:        shopid,
		IsHelpFull:    false,
	}
	config.DB.Create(&newReview)
	allReview := []models.ShopReview{}
	config.DB.Where("shop_id = ?", shopid).Find(&allReview)
	// length := len(allReview)
	avg := config.DB.Table("shop_reviews").Select("AVG(rating)").Row()
	var total float64
	avg.Scan(&total)
	fmt.Println(total)
	var myshop models.Shop
	config.DB.Where("id = ?", shopid).First(&myshop)
	myshop.Rating = total
	config.DB.Save(&myshop)
	c.JSON(200, &newReview)
}
func GetSimilarProductCategory(c *gin.Context) {
	fmt.Println("asdas")
	cat := c.Query("category")
	fmt.Println(cat)
	products := []models.Product{}

	config.DB.Where("category= ?", cat).Find(&products)
	fmt.Println(products)
	c.JSON(200, &products)
}

// updatenotes
func UpdateNotes(c *gin.Context) {
	newNotes := c.Query("note")
	wishid := c.Query("wishid")
	wishIDint, errwish := strconv.Atoi(wishid)
	if errwish != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Parsing!",
		})
		return
	}
	fmt.Println(newNotes)
	fmt.Println(wishid)
	var currwish models.WishList
	config.DB.Where("id= ?", wishIDint).First(&currwish)
	currwish.Note = newNotes
	config.DB.Save(&currwish)
	fmt.Println(currwish)
	c.JSON(200, &currwish)
}

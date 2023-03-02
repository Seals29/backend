package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func GetAllWishList(c *gin.Context) {
	wishlist := []models.WishList{}
	config.DB.Find(&wishlist)
	c.JSON(200, &wishlist)
}
func GetWishListDetail(c *gin.Context) {
	wishlistID := c.Param("id")
	fmt.Println(wishlistID)
	wishlist := []models.WishList{}
	config.DB.Where("user_id = ?", wishlistID).Find(&wishlist)
	
	c.JSON(200, &wishlist)
}
func UpdateWishListStatus(c *gin.Context){
	var body struct {
		UserID string `json:"userid"`
		Status string `json:"status"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
}
func CreateNewWishlist(c *gin.Context){
	var body struct{
		Name string `json:"name"`
		
		UserID string `json:"userid"`
		Status string `json:"status"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println("==")
	fmt.Println(body.Name)
	if len(body.Name)==0{
		c.JSON(200,gin.H{
			"error":"Wishlist name cannot be empty!",
		})
		return
	}
	fmt.Println(body)
	
	var wishlist models.WishList
	wishlist.Name= body.Name
	wishlist.Note=""
	wishlist.Status=body.Status
	userID, err:= strconv.Atoi(body.UserID)
	if err!=nil{
		c.JSON(200, gin.H{
			"error" :"Invalid convert ID",
		})
		return
	} 
	wishlist.UserID = userID
	config.DB.Create(&wishlist)
	c.JSON(200,&wishlist)
	return
}

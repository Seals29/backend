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
func GetWishListDetailByWishListID(c *gin.Context) {
	//dapet id wishlist
	wishlistid := c.Param("id")
	wishlistDetail := []models.WishListDetail{}
	config.DB.Where("wish_list_id= ?", wishlistid).Find(&wishlistDetail)
	c.JSON(200, &wishlistDetail)
	return
}
func GetPublicWishList(c *gin.Context) {
	wishlist := []models.WishList{}
	config.DB.Where("status = ?", "Public").Find(&wishlist)
	c.JSON(200, &wishlist)
}
func NewFollowWishList(c *gin.Context) {
	var body struct {
		UserID     string `json:"userid"`
		WishListID string `json:"wishlistid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println(body)
	wishID, errwish := strconv.Atoi(body.WishListID)
	userID, erruser := strconv.Atoi(body.UserID)
	var checkwishlist models.FollowingWishList
	config.DB.Where("wishlist_id=?", wishID).Where("user_id= ?", userID).First(&checkwishlist)
	fmt.Println(checkwishlist)
	if checkwishlist.ID == 0 {
		//empty unqiue
		if errwish != nil || erruser != nil {
			c.JSON(200, gin.H{
				"error": "Error parsing string to int",
			})
			return
		}
		var following models.FollowingWishList
		following.UserID = userID
		following.WishlistID = wishID
		config.DB.Create(&following)

		c.JSON(200, &following)
		return
	} else {
		// gaunique
		//delete
		config.DB.Where("id = ?", checkwishlist.ID).Delete(&checkwishlist)
		// var user models.User
		// config.DB.Where("id = ?", c.Param("id")).Delete(&user)
		c.JSON(200, gin.H{
			"message": "You've Unfollowed this wishlist",
		})
		return
	}

}
func DuplicatePublicWishlistToMyWishList(c *gin.Context) {
	var body struct {
		UserID     string `json:"userid"`
		WishListID string `json:"wishlistid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	// allDetial := []models.WishListDetail{}
	userid, erruser := strconv.Atoi(body.UserID)
	wishlistid, errwish := strconv.Atoi(body.WishListID)
	var joinedDetail []struct{
		models.WishList
		// models.WishListDetail
		// models.Product
	}
	fmt.Println(joinedDetail)
	var wishName models.WishList
	if erruser != nil || errwish != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}
	fmt.Println(userid)
	// config.DB.Table("wish_list_details").Joins("JOIN products ON products.id = wish_list_details.product_id JOIN wish_lists ON wish_lists.id = wish_list_details.wish_list_id").
	// Where("wish_list_details.wish_list_id=?",wishlistid).
	// Select("wish_list_details.*,wish_lists.name").Find(&joinedDetail)

	config.DB.Table("wish_list_details").Joins("JOIN wish_lists ON wish_lists.id = wish_list_details.wish_list_id").
	Where("wish_list_details.wish_list_id=?",wishlistid).
	Select("wish_lists.*").Find(&wishName)
	//dapet wishnnya lalu  bisa di duplikat aja lgsung
	
	//dh dpet semua wishlistdetial
	//create wishlist baru
	duplicatedWishList := models.WishList{
		Name: wishName.Name,
		Note: wishName.Note,
		Status: "Private",
		UserID: userid,
	}
	//create duplicated wishlistnya 
	config.DB.Create(&duplicatedWishList)
	var duplicatedDetail []struct{
		models.WishListDetail
	}
	fmt.Println(duplicatedWishList)
	config.DB.Table("wish_list_details").Joins("JOIN wish_lists ON wish_lists.id = wish_list_details.wish_list_id").
	Where("wish_list_details.wish_list_id=?",wishlistid).
	Select("wish_list_details.*").Find(&duplicatedDetail)
	fmt.Println(wishName)
	for _, wishlist:= range duplicatedDetail{
		var inputallDetails models.WishListDetail
		inputallDetails.ProductID=wishlist.ProductID
		inputallDetails.Quantity=wishlist.Quantity
		inputallDetails.WishListID=int(duplicatedWishList.ID)
		config.DB.Create(&inputallDetails)
		fmt.Println(inputallDetails)
		// fmt.Println(wishlist.ProductID)
		fmt.Println("=====")
	}
	
	allDuplicatedDetails := []models.WishListDetail{}
	config.DB.Where("wish_list_id = ?",duplicatedWishList.ID).Find(&allDuplicatedDetails)


	// config.DB.Where("wish_list_id = ?", wishlistid).Find(&allDetial)
	// config.DB.Joins("JOIN products ON products.ID = wish_list_detail.product_id").Where("wish_list_details.wish_list_id = ?",wishlistid).Find(&allDetial)
	// fmt.Println(allDetial)

	c.JSON(200, &allDuplicatedDetails)
}
func GetFollowWishListByUserID(c *gin.Context) {
	userID := c.Param("id")
	fmt.Println(userID)
	allfollowwishlist := []models.FollowingWishList{}
	config.DB.Where("user_id =?", userID).Find(&allfollowwishlist)

	// config.DB.Model(models.Product{}).Joins(`join shops on products.shop_id = shops.ID`).
	// 	Select("category").Where("shop_id = ?", shopid).Find(&product)
	c.JSON(200, &allfollowwishlist)
	return
}
func GetWishListByFollowedID(c *gin.Context) {
	wishlistID := c.Param("id")
	fmt.Println(wishlistID)
	// allfollowed := []models.WishList{}
	var allFollowedwishList models.WishList
	config.DB.Where("id = ?", wishlistID).First(&allFollowedwishList)
	c.JSON(200, &allFollowedwishList)
}
func GetPrivateWishList(c *gin.Context) {
	wishlist := []models.WishList{}
	config.DB.Where("status = ?", "Private").Find(&wishlist)
	c.JSON(200, &wishlist)
	return
}
func GetWishListDetail(c *gin.Context) {
	wishlistID := c.Param("id")
	fmt.Println(wishlistID)
	wishlist := []models.WishList{}
	config.DB.Where("user_id = ?", wishlistID).Find(&wishlist)

	c.JSON(200, &wishlist)
}
func GetDetailWishListByWishListID(c *gin.Context) {
	// wishlistID := c.Param("id")
	var body struct {
		Id string `json:"id"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println(body)
	wishId, err := strconv.Atoi(body.Id)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}
	var wishlist models.WishList
	// fmt.Println(wishlistID)
	config.DB.Where("id = ?", wishId).First(&wishlist)
	fmt.Println(wishlist)
	c.JSON(200, &wishlist)
}
func UpdateWishListStatus(c *gin.Context) {
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
func AddNewProductIntoWishList(c *gin.Context) {
	var body struct {
		WishlistID string `json:"wishlistid"`
		ProductID  string `json:"productid"`
		Quantity string `json:"quantity"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	wishlistid, errwl := strconv.Atoi(body.WishlistID)
	productid, errpr := strconv.Atoi(body.ProductID)
	quantity , errq := strconv.Atoi(body.Quantity)
	if errwl != nil || errpr != nil ||errq!=nil{
		c.JSON(200, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}
	detailwishlist := models.WishListDetail{
		WishListID: wishlistid,
		ProductID:  productid,
		Quantity:  quantity,
	}
	config.DB.Create(&detailwishlist)
	c.JSON(200, gin.H{
		"message": "Product has been successfully added!",
	})
	return

}
func CreateNewWishlist(c *gin.Context) {
	var body struct {
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
	if len(body.Name) == 0 {
		c.JSON(200, gin.H{
			"error": "Wishlist name cannot be empty!",
		})
		return
	}
	fmt.Println(body)

	var wishlist models.WishList
	wishlist.Name = body.Name
	wishlist.Note = ""
	wishlist.Status = body.Status
	userID, err := strconv.Atoi(body.UserID)
	if err != nil {
		c.JSON(200, gin.H{
			"error": "Invalid convert ID",
		})
		return
	}
	wishlist.UserID = userID
	config.DB.Create(&wishlist)
	c.JSON(200, &wishlist)
	return
}

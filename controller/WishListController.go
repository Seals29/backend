package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GetAllCommentWishlist(c *gin.Context) {
	wishlistid := c.Query("id")
	allcomment := []models.CommentWishList{}
	config.DB.Where("wish_list_id = ?", wishlistid).Find(&allcomment)
	c.JSON(200, &allcomment)
	fmt.Println(allcomment)
	fmt.Println(wishlistid)
	fmt.Println("==ggetallwishlistcomment")
}
func CommentWishList(c *gin.Context) {
	var body struct {
		IsAnon         bool   `json:"isanon"`
		UserID         string `json:"userid"`
		WishListID     string `json:"wishlistid"`
		CommentMessage string `json:"commentmessage"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	userid, erruser := strconv.Atoi(body.UserID)
	wishlistid, errwish := strconv.Atoi(body.WishListID)
	fmt.Println(body)
	if erruser != nil || errwish != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Conversion string to int",
		})
		return
	}
	var checkusername models.User
	config.DB.Where("id = ?", userid).First(&checkusername)
	var username string

	if body.IsAnon == false {
		username = checkusername.FirstName
	}
	if body.IsAnon == true {
		username = "Anonymous!"
	}

	newcomment := models.CommentWishList{
		WishListID:     wishlistid,
		Username:       username,
		CommentMessage: body.CommentMessage,
	}
	fmt.Print(newcomment)
	config.DB.Create(&newcomment)
	c.JSON(200, &newcomment)
	// return

}
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
	fmt.Println(body)
	// allDetial := []models.WishListDetail{}
	userid, erruser := strconv.Atoi(body.UserID)
	wishlistid, errwish := strconv.Atoi(body.WishListID)
	var joinedDetail []struct {
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
		Where("wish_list_details.wish_list_id=?", wishlistid).
		Select("wish_lists.*").Find(&wishName)
	//dapet wishnnya lalu  bisa di duplikat aja lgsung
	fmt.Println(wishName)
	if wishName.ID == 0 {
		c.JSON(200, gin.H{
			"error": "No WishListDetial Found!",
		})
		return
	}
	//dh dpet semua wishlistdetial
	//create wishlist baru
	duplicatedWishList := models.WishList{
		Name:   wishName.Name,
		Note:   wishName.Note,
		Status: "Private",
		UserID: userid,
	}
	fmt.Println(duplicatedWishList)
	//create duplicated wishlistnya
	config.DB.Create(&duplicatedWishList)
	var duplicatedDetail []struct {
		models.WishListDetail
	}
	fmt.Println(duplicatedWishList)
	config.DB.Table("wish_list_details").Joins("JOIN wish_lists ON wish_lists.id = wish_list_details.wish_list_id").
		Where("wish_list_details.wish_list_id=?", wishlistid).
		Select("wish_list_details.*").Find(&duplicatedDetail)
	fmt.Println(wishName)
	for _, wishlist := range duplicatedDetail {
		var inputallDetails models.WishListDetail
		inputallDetails.ProductID = wishlist.ProductID
		inputallDetails.Quantity = wishlist.Quantity
		inputallDetails.WishListID = int(duplicatedWishList.ID)
		config.DB.Create(&inputallDetails)
		fmt.Println(inputallDetails)
		// fmt.Println(wishlist.ProductID)
		fmt.Println("=====")
	}

	allDuplicatedDetails := []models.WishListDetail{}
	config.DB.Where("wish_list_id = ?", duplicatedWishList.ID).Find(&allDuplicatedDetails)

	// config.DB.Where("wish_list_id = ?", wishlistid).Find(&allDetial)
	// config.DB.Joins("JOIN products ON products.ID = wish_list_detail.product_id").Where("wish_list_details.wish_list_id = ?",wishlistid).Find(&allDetial)
	// fmt.Println(allDetial)

	c.JSON(200, gin.H{
		"message": "Duplicated Succesfull to your wishlist!",
	})
	return
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
func extractUserIDFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify that the signing method is HMAC SHA256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key used to sign the token
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Extract the user ID from the token's claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("unexpected claims type: %T", token.Claims)
	}

	// userID, ok := claims["sub"].(string)
	// if !ok {
	//     return "", fmt.Errorf("unexpected user ID type: %T", claims["sub"])
	// }
	userID := fmt.Sprintf("%.0f", claims["sub"])
	return userID, nil
}
func DeleteProductFromWishListID(c *gin.Context) {
	var body struct {
		WishListID string `json:"wishlistid"`
		ProductID  string `json:"productid"`
		UserID     string `json:"userid"`
		JwtToken   string `json:"jwttoken"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	fmt.Println(body.JwtToken)
	curruser, err := extractUserIDFromToken(body.JwtToken)
	fmt.Println(err)
	fmt.Println(curruser)
	curruserID, errcurid := strconv.Atoi(curruser)
	var checkWishDetail models.WishListDetail
	userid, erruser := strconv.Atoi(body.UserID)
	productid, errpr := strconv.Atoi(body.ProductID)
	wishid, errwish := strconv.Atoi(body.WishListID)
	fmt.Println(body)
	if errpr != nil || errwish != nil || erruser != nil || errcurid != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Parsing!",
		})
		return
	}
	var validateUser models.WishList
	config.DB.Where("id = ?", wishid).Where("user_id=?", userid).First(&validateUser)
	fmt.Println(validateUser)
	config.DB.Where("product_id = ?", productid).Where("wish_list_id = ?", wishid).First(&checkWishDetail)
	fmt.Println(checkWishDetail)
	if checkWishDetail.ID == 0 {
		//ga nemu
		c.JSON(200, gin.H{
			"error": "Product Not Found in that wishlist!",
		})
		return
	} else {

		if curruserID == userid {
			fmt.Println("userid sama maka boleh delete")
			config.DB.Delete(&checkWishDetail)
			c.JSON(200, gin.H{
				"message": "Product Successfully deleted from your wishlist!",
			})
			return
		} else {
			c.JSON(404, gin.H{
				"error": "You Are not authorized to update the items!",
			})
			return
		}
	}
}
func UpdateWishListUser(c *gin.Context) {
	fmt.Println("====")
	var body struct {
		WishListId     string `json:"wishlistid"`
		WishListName   string `json:"wishlistname"`
		WishListStatus string `json:"wishliststatus"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	var wishlist models.WishList
	wishlisdid, errid := strconv.Atoi(body.WishListId)
	if errid != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Parsing",
		})
		return
	}
	config.DB.Where("id = ?", wishlisdid).First(&wishlist)
	wishlist.Name = body.WishListName
	wishlist.Status = body.WishListStatus
	config.DB.Save(&wishlist)
	c.JSON(200, &wishlist)
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
		Quantity   string `json:"quantity"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	wishlistid, errwl := strconv.Atoi(body.WishlistID)
	productid, errpr := strconv.Atoi(body.ProductID)
	quantity, errq := strconv.Atoi(body.Quantity)
	if errwl != nil || errpr != nil || errq != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}
	detailwishlist := models.WishListDetail{
		WishListID: wishlistid,
		ProductID:  productid,
		Quantity:   quantity,
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
func UpdateAddQuantityProduct(c *gin.Context){
	var body struct{
		JwtToken string `json:"jwttoken"`
		WishListID string `json:"wishlistid"`
		ProductID string `json:"productid"`
		Quantity string `json:"quantity"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	curruser, err := extractUserIDFromToken(body.JwtToken)
	curruserID, errcurid := strconv.Atoi(curruser)
	fmt.Println(err)
	wishlistid, errwl := strconv.Atoi(body.WishListID)
	productid, errpr := strconv.Atoi(body.ProductID)
	quantity, errq := strconv.Atoi(body.Quantity)
	var wishdetail models.WishListDetail
	config.DB.Where("wish_list_id = ?", wishlistid).Where("product_id=?", productid).First(&wishdetail)
	fmt.Println(wishdetail)
	if errwl != nil || errpr != nil || errq != nil || errcurid != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}
	var validateUser models.WishList
	config.DB.Where("id = ?",wishlistid).Where("user_id = ?",curruserID).First(&validateUser)
	if validateUser.ID==0{
		c.JSON(400,gin.H{
			"error":"You are unauthroized!",
		})
		return
	}
	if wishdetail.ID==0{
		c.JSON(400,gin.H{
			"error":"Product Not Found in wishlist",
		})
		return
	}
	wishdetail.Quantity=quantity;
	config.DB.Save(&wishdetail)
	c.JSON(200,gin.H{
		"message":"Successfully updated!",
	})
}
func AddToCartFromWishList(c *gin.Context) {
	var body struct {
		UserID     string `json:"userid"`
		WishlistID string `json:"wishlistid"`
		ProductID  string `json:"productid"`
		Quantity   string `json:"quantity"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	
	wishlistid, errwl := strconv.Atoi(body.WishlistID)
	productid, errpr := strconv.Atoi(body.ProductID)
	quantity, errq := strconv.Atoi(body.Quantity)
	var wishdetail models.WishListDetail
	config.DB.Where("wish_list_id = ?", wishlistid).Where("product_id=?", productid).First(&wishdetail)
	fmt.Println(wishdetail)
	userid, erruser := strconv.Atoi(body.UserID)
	if errwl != nil || errpr != nil || errq != nil || erruser != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Conversion",
		})
		return
	}
	cart := models.Cart{
		ProductID: productid,
		UserID:    userid,
		Quantity:  quantity,
	}
	// config.DB.Create(&cart)
	fmt.Println(cart)
	c.JSON(200, gin.H{
		"message": "Product has been succesfully added to your cart!",
	})
	return

}


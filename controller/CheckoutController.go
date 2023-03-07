package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func NewAddress(c *gin.Context) {
	var body struct {
		UserID       string `json:"userid"`
		AddressField string `json:"addressfield"`
		ReceiverName string `json:"receivername"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	// c.JSON(200, &body)
	userid, errid := strconv.Atoi(body.UserID)
	if errid != nil {
		c.JSON(200, gin.H{
			"error": "Invalid Parsing!",
		})
		return
	}
	if len(body.ReceiverName) == 0 {
		c.JSON(200, gin.H{
			"error": "Adress cannot be empty!",
		})
		return
	}
	if len(body.AddressField) == 0 {
		c.JSON(200, gin.H{
			"error": "Receiver Name cannot be empty!",
		})
		return
	}
	var checkaddr models.CustomerAddress
	config.DB.Where("user_id = ?", userid).Where("is_active = ?", true).First(&checkaddr)
	if checkaddr.ID == 0 {
		newAddr := models.CustomerAddress{
			UserID:       userid,
			ReceiverName: body.ReceiverName,
			AddressField: body.AddressField,
			IsActive:     true,
		}
		config.DB.Create(&newAddr)
	} else {
		newAddr := models.CustomerAddress{
			UserID:       userid,
			ReceiverName: body.ReceiverName,
			AddressField: body.AddressField,
			IsActive:     false,
		}
		config.DB.Create(&newAddr)
	}

	c.JSON(200, gin.H{
		"message": "New Address has been added!",
	})

}
func GetAllAddress(c *gin.Context) {
	addresses := []models.CustomerAddress{}
	config.DB.Find(&addresses)
	c.JSON(200, &addresses)
}
func CheckoutToOrderPage(c *gin.Context) {
	var body struct {
		UserID        string `json:"userid"`
		Address       string `json:"address"`
		Receiver      string `json:"receiver"`
		PaymentMethod string `json:"paymentmethod"`
		Delivery      string `json:"delivery"`
		ProductID     string `json:"productid"`
		// GrandTotal    float64 `json:"grandtotal"`
		ProductTotal string `json:"producttotal"`
		// AdminFee      float64 `json:"adminfee"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	// c.JSON(200, &body)
	productid, errpr := strconv.Atoi(body.ProductID)
	userid, erruser := strconv.Atoi(body.UserID)
	totalproduct, errtotal := strconv.Atoi(body.ProductTotal)
	if errtotal != nil || erruser != nil || errpr != nil {
		c.JSON(200, gin.H{
			"error": "Invalid conversion!",
		})
		return
	}
	fmt.Println(userid)
	var user models.User
	config.DB.Where("id = ?", userid).First(&user)
	if user.Balance >= totalproduct {
		user.Balance = user.Balance - totalproduct
		config.DB.Save(&user)
		var productshop models.Product
		config.DB.Where("id=?", productid).First(&productshop)

		newOrder := models.Order{
			UserID:        userid,
			ProductID:     productid,
			ShopID:        productshop.ShopID,
			PaymentMethod: body.PaymentMethod,
			Receiver:      body.Receiver,
			Address:       body.Address,
			Delivery:      body.Delivery,
			ProductTotal:  totalproduct,
		}
		config.DB.Create(&newOrder)
		var cart models.Cart
		config.DB.Where("user_id = ?",userid).Where("product_id =?",productid).First(&cart)
		config.DB.Delete(&cart)
		
		c.JSON(200, true)
		// c.JSON(200,)
		//create order

	} else {
		c.JSON(200,false)
		return
	}

}

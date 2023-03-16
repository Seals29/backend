package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLength  = 15
)

func generateRandomInvoiceCode(name string) string {
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	invoiceCode := fmt.Sprintf("INV-%s-%s", name, string(code))
	return invoiceCode
}
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
func CreateNewOrders(c *gin.Context) {
	var body struct {
		UserID        string `json:"userid"`
		Address       string `json:"address"`
		Receiver      string `json:"receiver"`
		PaymentMethod string `json:"paymentmethod"`
		Delivery      string `json:"delivery"`
		ProductTotal  string `json:"producttotal"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	userid, erruser := strconv.Atoi(body.UserID)
	totalPrice, errtotal := strconv.Atoi(body.ProductTotal)
	fmt.Println(totalPrice)
	if errtotal != nil || erruser != nil {
		c.JSON(200, gin.H{
			"error": "Invalid conversion!",
		})
		return
	}

	carts := []models.Cart{}
	var user models.User
	config.DB.Where("id = ?", userid).First(&user)

	float := float64(totalPrice) + float64(totalPrice)*0.05

	if user.Balance < float {
		c.JSON(400, gin.H{
			"error": "Your balance is insuficient!",
		})
		return
	}
	config.DB.Where("user_id = ?", userid).Find(&carts)

	fmt.Println(carts)
	var shopIds []int
	for _, cart := range carts {
		// fmt.Println(cart)
		var shopid int
		//dapet shopidnya
		config.DB.Table("products").Select("DISTINCT products.shop_id").Joins("JOIN carts ON carts.product_id=products.id").Where("carts.id = ?", cart.ID).Pluck("products.shop_id", &shopid)
		fmt.Println(shopid)
		fmt.Println("======")
		shopIds = append(shopIds, shopid)
	}
	// var shopIDs []int
	fmt.Println(shopIds)
	// config.DB.Table("products").Select("DISTINCT shop_id").Joins("JOIN carts ON products.id = carts.product_id").Where("user_id = ?", userid).Pluck("shop_id", &shopIDs)
	// config.DB.Table("carts").Select("DISTINCT shop_id").Joins("JOIN products ON carts.product_id = products.id").Joins("JOIN shops ON shops.id = products.shop_id").Where("carts.user_id = ?", userid).Pluck("shops.id", &shopIDs)
	fmt.Println(shopIds)
	fmt.Println("-----shopids")
	for i := 0; i < len(shopIds); i++ {
		shopID := shopIds[i]
		invoiceCode := generateRandomInvoiceCode(user.FirstName)

		var newOrder models.Order
		newOrder.ShopID = shopID

		newOrder.UserID = userid
		newOrder.Address = body.Address
		newOrder.Delivery = body.Delivery
		newOrder.Invoice = invoiceCode
		newOrder.PaymentMethod = body.PaymentMethod
		newOrder.Receiver = body.Receiver
		newOrder.Status = "Open"

		config.DB.Create(&newOrder)
		fmt.Println(newOrder)
		for _, cart := range carts {
			fmt.Println()
			fmt.Println("====looping dlm cart")
			// var newOrderDetail models.OrderDetail
			var product models.Product
			config.DB.Where("id = ?", cart.ProductID).First(&product)
			fmt.Println(product)
			if product.ID == uint(cart.ProductID) {
				var order models.Order
				config.DB.Where("shop_id = ?", product.ShopID).First(&order)
				newOrderDetail := models.OrderDetail{
					OrderID:   int(newOrder.ID),
					ProductID: cart.ProductID,
					Quantity:  cart.Quantity,
				}
				fmt.Println(newOrderDetail)
				// var product models.Product
				// config.DB.Where("id = ?",cart.ProductID).First(&product)

				config.DB.Create(&newOrderDetail)
				fmt.Println(newOrderDetail.OrderID)
				var deletedCart models.Cart
				config.DB.Where("id = ?", cart.ID).Delete(&deletedCart)
			}

		}
	}
	fmt.Println("carts")
	fmt.Println(carts)
	c.JSON(200, gin.H{
		"message": "success!",
	})

	// var user models.User
	// config.DB.Where("id = ?", userid).First(&user)
	// var float float64
	// float = float64(totalproduct) + float64(totalproduct)*0.05
	// if user.Balance >= float {
	// 	//cukup
	// 	//create order

	// 	var newOrder models.UserOrder
	// 	newOrder.Address = body.Address
	// 	newOrder.Delivery = body.Delivery
	// 	newOrder.TotalPrice = totalproduct
	// 	newOrder.UserID = userid
	// 	newOrder.ShopID = shopid
	// 	newOrder.Invoice = invoiceCode
	// 	fmt.Println(newOrder)
	// } else {
	// 	c.JSON(400, gin.H{
	// 		"error": "Balance insuficient!",
	// 	})
	// 	return
	// }
}
func CheckoutToOrderPage(c *gin.Context) {

}
func CalculateTotalPriceByEachUser(c *gin.Context) {
	var body struct {
		UserID string `json:"userid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	userid, errcart := strconv.Atoi(body.UserID)
	if errcart != nil {
		c.JSON(400, gin.H{
			"error": "Failed parsing",
		})
		return
	}
	fmt.Println("calculatinggg===")
	var totalPrice float64
	// var totalPrice float64
	err := config.DB.Raw("SELECT SUM(p.price * c.quantity) as Total FROM carts c JOIN products p ON c.product_id = p.id WHERE c.user_id = ? AND c.deleted_at IS NULL", userid).Scan(&totalPrice).Error
	if err != nil {
		// Handle error
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	// err := config.DB.Raw("SELECT SUM(p.price * c.quantity) as Total  FROM carts c JOIN products p ON c.product_id = p.id WHERE c.user_id = ?", userid).Scan(&totalPrice).Error
	// if err != nil {
	// 	// Handle error
	// 	c.JSON(400, gin.H{
	// 		"error": err,
	// 	})
	// 	return
	// }
	c.JSON(200, &totalPrice)
}
func CalculateTotalPriceByEachCart(c *gin.Context) {
	var body struct {
		CartID string `json:"cartid"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body1",
		})
		return
	}
	cartid, errcart := strconv.Atoi(body.CartID)
	if errcart != nil {
		c.JSON(400, gin.H{
			"error": "Failed parsing",
		})
		return
	}

	var totalPrice float64
	err := config.DB.Raw("SELECT SUM(p.price * c.quantity) AS total FROM carts c JOIN products p ON c.product_id = p.id WHERE c.id = ?", cartid).Scan(&totalPrice).Error
	if err != nil {
		// Handle error
		c.JSON(400, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(200, &totalPrice)
}

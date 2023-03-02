package routes

import (
	"github.com/Seals29/controller"
	"github.com/Seals29/mware"
	"github.com/gin-gonic/gin"
)

func UserRoute(route *gin.Engine) {
	route.GET("/getuser", controller.GetUsers)
	route.POST("/signup", controller.SignUp)
	route.POST("/login", controller.Login)
	route.DELETE("/:id", controller.DeleteUser)
	route.PUT("/:id", controller.UpdateUser)
	route.POST("/validate", mware.UserAuth, controller.ValidateUser)

	route.POST("/sendonetime", controller.SendOneTime)
	route.POST("/usersubscribe", controller.UserSubscribe)
	route.POST("/announce", controller.Announce)

	route.POST("/onetimecode", controller.OneTimeCode)

	route.POST("/setban", controller.SetBan)
	route.POST("/setbanuser", controller.SetBanUser)
	route.POST("/createproduct", controller.CreateProduct)
	route.POST("/getproduct", controller.Getproduct)
	route.POST("/updateshopprofile", controller.UpdateShopProfile)
	route.POST("/resendonetime", controller.ResendOneTime)
	route.POST("/updateshoppass", controller.UpdateShopPassword)
	route.POST("/resetpassword", controller.ResetPassword)
	route.POST("/getdetailproduct", controller.GetDetailProduct)
	route.POST("/updateproduct", controller.UpdateProduct)
	route.POST("/loginassistance", controller.LoginAssistance)
	route.POST("/resendassistance", controller.ResendForgotPassword)
	route.POST("/verifassistance", controller.VerifForgotPassword)
	route.POST("/updateuserpass", controller.UpdateUserPassword)
	route.GET("/getallproduct", controller.GetAllProduct)
	route.GET("/message", controller.SendingMessage)
	route.GET("/message/:id", controller.SendingMessage)

	route.POST("/insertcart", controller.InsertCart)
	route.GET("/getproduct/:id", controller.GetOneProduct)

	route.POST("/viewallmessage", controller.GetAllShopMsg)

	//Shop
	route.GET("/getoneshop/:id", controller.GetOneShopID)
	route.GET("/getallproductbyshop/:id", controller.GetAllProductByShop)
	route.POST("/getoneshop", controller.GetOneShop)
	route.GET("/getsingleshop/:id", controller.GetSingleShop)
	route.POST("/createshop", controller.CreateShop)
	route.POST("/notifyshop", controller.NotifyShop)
	route.GET("/getshop", controller.GetShops)
	route.GET("/getusershopid/:id", controller.GetUserShopId)
	route.GET("/getcategorybyshop/:shopid", controller.GetCategoryByShopId)
	//product
	route.GET("/getproductcategory/:category", controller.GetProductByCategory)
	route.GET("/getallsubcategory", controller.GetAllSubCategory)
	route.GET("/getallcategory", controller.GetAllCategory)

	//cart
	route.GET("/getallcarts")
	//user
	route.POST("/updateuserpassword", controller.UpdateAccountPassword)

	//wishlist
	route.GET("/getpublicwishlist")
	route.GET("/getallwishlist", controller.GetAllWishList)
	route.GET("/getwishlistbyid/:id", controller.GetWishListDetail)
	route.POST("/createnewwishlist", controller.CreateNewWishlist)
	route.POST("/updatewishliststatus",controller.UpdateWishListStatus)
}

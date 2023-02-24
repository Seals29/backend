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
	route.POST("/createshop", controller.CreateShop)
	route.POST("/sendonetime", controller.SendOneTime)
	route.GET("/sendmessage", controller.SendMessage)
	route.POST("/usersubscribe", controller.UserSubscribe)
	route.POST("/announce", controller.Announce)
	route.POST("/notifyshop", controller.NotifyShop)
	route.POST("/onetimecode", controller.OneTimeCode)
	route.GET("/getshop", controller.GetShops)
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
}

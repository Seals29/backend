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
	route.POST("/forgotpassword", controller.ForgotPassword)
	route.GET("/sendmessage", controller.SendMessage)
	route.POST("/usersubscribe", controller.UserSubscribe)
	route.POST("/announce", controller.Announce)
	route.POST("/notifyshop", controller.NotifyShop)
	route.POST("/onetimecode", controller.OneTimeCode)
}

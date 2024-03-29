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

	route.POST("/insertcart", controller.InsertCart)
	route.GET("/getproduct/:id", controller.GetOneProduct)

	route.POST("/viewallmessage", controller.GetAllShopMsg)
	//chat
	route.GET("/message", controller.SendingMessage)
	// route.GET("/message/:id", controller.SendingMessage)
	route.GET("/getallmessage", controller.GetAllMsg)
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
	route.GET("/getshopidbyuserid", controller.GetShopIDByUserID)
	//product
	route.GET("/getproductcategory/:category", controller.GetProductByCategory)
	route.GET("/getallsubcategory", controller.GetAllSubCategory)
	route.GET("/getallcategory", controller.GetAllCategory)
	route.GET("/loadProducts", controller.LoadProductByPage)
	route.GET("/getsimilarproductbycategory", controller.GetSimilarProductCategory)
	//cart
	// route.GET("/getallcarts")
	//user
	route.POST("/subcribenewsfromhome", controller.SubscribeFromHome)
	route.POST("/updateuserpassword", controller.UpdateAccountPassword)
	route.POST("/updateuseremail", controller.UpdateAccountEmail)
	route.POST("/updateuserphone", controller.UpdateAccountPhoneNumber)
	route.POST("/getsubscribestatus", controller.GetSubscribeStatus)

	route.GET("/sendverifuseremail", controller.SendVerifUserEmail)
	route.GET("/verifuseremail", controller.VerifUserEmail)
	//wishlist
	route.GET("/updatenotes", controller.UpdateNotes)

	route.GET("/getallcommentwishlist", controller.GetAllCommentWishlist)
	route.POST("/newcommentwishlist", controller.CommentWishList)
	route.GET("/getpublicwishlist", controller.GetPublicWishList)
	route.GET("/getallwishlist", controller.GetAllWishList)
	route.GET("/getwishlistbyid/:id", controller.GetWishListDetail)
	route.POST("/createnewwishlist", controller.CreateNewWishlist)
	route.POST("/AddNewProductIntoWishList", controller.AddNewProductIntoWishList)
	route.POST("/updatewishliststatus", controller.UpdateWishListStatus)
	route.GET("/getprivatewishlist", controller.GetPrivateWishList)
	route.POST("/NewFollowWishList", controller.NewFollowWishList)
	route.GET("/GetFollowWishListByUserId/:id", controller.GetFollowWishListByUserID)
	route.GET("/GetFollowedWishListByWishListID/:id", controller.GetWishListByFollowedID)
	route.GET("/getwishlistdetailbywishlistid/:id", controller.GetWishListDetailByWishListID)
	route.POST("/getDWbyID", controller.GetDetailWishListByWishListID)
	route.POST("/duplicatepublicwishlisttomywishlist", controller.DuplicatePublicWishlistToMyWishList)
	route.POST("/updatewishlistuser", controller.UpdateWishListUser)
	route.POST("/deleteproductfromwishlistid", controller.DeleteProductFromWishListID)
	route.POST("/addtocartfromwishlist", controller.AddToCartFromWishList)
	route.POST("/updatequantityproduct", controller.UpdateAddQuantityProduct)
	//voucher
	route.POST("/newvoucher", controller.NewVoucher)
	route.POST("/checkvoucher", controller.CheckVoucher)

	//cart
	route.GET("/getallcarts", controller.GetAllCart)
	route.GET("/getallsavelaters", controller.GetAllSavelater)
	route.POST("/deleteitemincart", controller.DeleteProductInCart)
	route.POST("/movecarttosavelater", controller.MoveCartToSave)
	route.POST("/calculatetotalprice", controller.CalculateTotalPriceByEachCart)
	route.POST("/calculatetotalpricebyuser", controller.CalculateTotalPriceByEachUser)
	//checkout
	route.POST("/newaddress", controller.NewAddress)
	route.GET("/getalladdress", controller.GetAllAddress)
	route.POST("/checkout", controller.CheckoutToOrderPage)
	route.POST("/createneworder", controller.CreateNewOrders)

	//review shop
	route.POST("/newreviewshop", controller.AddNewReviewShop)
	route.GET("/getreviewbyshopid", controller.GetReviewsByShop)

	//orders
	route.GET("/getorders", controller.GetOrders)
	route.POST("/getallordersbyuserid", controller.GetAllOrdersByUserID)
	route.POST("/getallordersbyshopid", controller.GetAllOrdersByShopID)
	route.POST("/getuserdetailbyorderid", controller.GetUseDetailByOrderID)
	route.POST("/getshopnamebyshopid", controller.GetShopNameByShopID)
	route.GET("/getorderbyorderid", controller.GetOrderByOrderID)
	route.POST("/getorderdetailbyorderid", controller.AddAllOrderDetailToCart)
	//banner
	route.GET("/getbanner", controller.TestData)
	route.POST("/addpromotionbanner", controller.AddPromotionBanner)
	route.POST("/rmvpromotionbanner", controller.DeletePromotionBanner)

	//reviews
	route.GET("/getreviewbyuserid", controller.GetReviewByUserID)
	route.GET("/deletereviewbyID", controller.DeleteReviewByRevID)
	route.GET("/updatereviewbyid",controller.UpdateReviewByID)
	//search
	route.GET("/searchproduct", controller.SearchProduct)
	route.GET("/savequery", controller.SaveSearchQuery)

	//notification
	route.GET("/getnotif", controller.GetStoreNotification)
	route.GET("/getannouncenotif", controller.GetAnnounceNotification)



	//visuali
	route.GET("/getcatcount",controller.GetCategoryCount)
}

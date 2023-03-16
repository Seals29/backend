package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string  `json:"firstname"`
	LastName    string  `json:"lastname"`
	IsBan       bool    `json:"isban"`
	Role        string  `json:"role"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	PhoneNumber string  `json:"phonenumber"`
	Balance     float64 `json:"balance"`
	IsVerif     bool    `json:"isverif"`
}
type ResetUser struct {
	gorm.Model
	ExpiredDate time.Time `json:"expireddate"`
	ResetCode   int       `json:"resetcode"`
	UserID      int       `json:"userid"`
	Used        bool      `json:"used"`
}
type UserSubscribe struct {
	gorm.Model
	UserEmail string `json:"useremail"`
}
type Shop struct {
	gorm.Model
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsBan       bool    `json:"isban"`
	Banner      string  `json:"banner"`
	Sales       int     `json:"sales"`
	Service     float64 `json:"service"`
	Followers   int     `json:"followers"`
	Rating      float64 `json:"rating"`
}
type Product struct {
	gorm.Model
	Name        string `json:"name"`
	Category    string `json:"category"`
	ShopEmail   string `json:"email"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Stock       int    `json:"stock"`
	Rating      int    `json:"rating"`
	Detail      string `json:"detail"`
	ShopID      int    `json:"shopid"`
	TotalSales  int    `json:"totalsales"`
	TotalStar   int    `json:"totalstar"`
	SubCategory string `json:"subcategory"`
}
type ProductDetailImage struct {
	gorm.Model
	ProductID int    `json:"productid"`
	Image     string `json:"image"`
}
type ProductCategory struct {
	gorm.Model
	Name string `json:"name"`
}
type ProductSubCategory struct {
	gorm.Model
	ProductCategoryID int    `json:"productcategoryid"`
	Name              string `json:"name"`
}
type ForgotUser struct {
	gorm.Model
	ExpiredDate time.Time `json:"expireddate"`
	ResetCode   int       `json:"resetcode"`
	UserID      int       `json:"userid"`
	Used        bool      `json:"used"`
}
type Message struct {
	gorm.Model
	SenderID    string `json:"from"`
	RecipientID string `json:"to"`
	Content     string `json:"content"`
	Type        string `json:"type"`
}
type Cart struct {
	gorm.Model
	ProductID int `json:"productid"`
	UserID    int `json:"userid"`
	Quantity  int `json:"quantity"`
}
type CustomerAddress struct {
	gorm.Model
	IsActive     bool   `json:"isactive"`
	UserID       int    `json:"userid"`
	ReceiverName string `json:"receivername"`
	AddressField string `json:"addressfield"`
}
type SaveLater struct {
	gorm.Model
	ProductID int `json:"productid"`
	UserID    int `json:"userid"`
	Quantity  int `json:"quantity"`
}
type Follow struct {
	gorm.Model
	FollowTo   int `json:"followto"`
	FollowedBy int `json:"followedby"`
	IsFollow   int `json:"isfollow"`
}
type Order struct {
	gorm.Model
	ShopID        int    `json:"shopid"`
	UserID        int    `json:"userid"`
	Status        string `json:"status"`
	Invoice       string `json:"invoice"`
	PaymentMethod string `json:"paymentmethod"`
	Receiver      string `json:"receiver"`
	Address       string `json:"address"`
	Delivery      string `json:"delivery"`
}

type OrderDetail struct {
	gorm.Model
	OrderID   int `json:"orderid"`
	ProductID int `json:"productid"`
	Quantity  int `json:"quantity"`
}
type WishList struct {
	gorm.Model
	Name   string `json:"name"`
	Status string `json:"status"`
	Note   string `json:"note"`
	UserID int    `json:"userid"`
}
type WishListDetail struct {
	gorm.Model
	WishListID int `json:"wishlistid"`
	ProductID  int `json:"productid"`
	Quantity   int `json:"quantity"`
}
type FollowingWishList struct {
	gorm.Model
	ID         int `json:"id"`
	WishlistID int `json:"wishlistid"`
	UserID     int `json:"userid"`
}
type Voucher struct {
	gorm.Model
	VoucherCode     string `json:"vouchercode"`
	VoucherCurrency int    `json:"vouchercurrency"`
}
type CommentWishList struct {
	gorm.Model
	Username       string `json:"username"`
	CommentMessage string `json:"commentmessage"`
	WishListID     int    `json:"wishlistid"`
}
type ShopReview struct {
	gorm.Model
	UserID        int    `json:"userid"`
	ShopID        int    `json:"shopid"`
	Rating        int    `json:"starreview"`
	ReviewComment string `json:"reviewcomment"`
	IsHelpFull    bool   `json:"ishelpfull"`
}
type Notification struct {
	gorm.Model
	UserID   int    `json:"userid"`
	Message  string `json:"message"`
	IsMarked bool   `json:"ismarked"`
}
type PromotionBanner struct {
	gorm.Model
	PromotionImage string `json:"promotionimage"`
}

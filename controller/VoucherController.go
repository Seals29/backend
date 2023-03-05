package controller

import (
	"strconv"

	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)
func CheckVoucher(c *gin.Context){
	var body struct{
		UserID string `json:"userid"`
		VoucherCode string `json:"vouchercode"`
	}
	if c.Bind(&body)!=nil{
		c.JSON(200,gin.H{
			"error":"Invalid Ready Body!",
		})
		return
	}
	userid ,erruser:= strconv.Atoi(body.UserID)
	if erruser!=nil{
		c.JSON(200,gin.H{
			"error":"Invalid Parsing data",
		})
		return
	}
	var checkCode models.Voucher
	config.DB.Where("voucher_code =?",body.VoucherCode).First(&checkCode)
	if checkCode.ID==0{
		//code ga ada
		c.JSON(200,gin.H{
			"error":"Invalid Voucher Code!",
		})
		return
	}else{
		//ada code
		var user models.User
		config.DB.Where("id = ?",userid).First(&user)
		
		user.Balance = user.Balance + checkCode.VoucherCurrency
		config.DB.Save(&user)
		c.JSON(200,gin.H{
			"message":"Code has been redeemed succesfully!",
		})
		return
	}
}
func NewVoucher(c *gin.Context){
	var body struct{
		VoucherCode string `json:"vouchercode"`
		VoucherCurrency string `json:"vouchercurrency"`
	}
	if c.Bind(&body)!=nil{
		c.JSON(200,gin.H{
			"error":"Invalid Ready Body!",
		})
		return
	}
	currency,errcurr := strconv.Atoi(body.VoucherCurrency)
	if errcurr!=nil{
		c.JSON(200,gin.H{
			"error":"Invalid Parsing",
		})
		return
	}
	var newVoucher models.Voucher
	newVoucher.VoucherCode = body.VoucherCode
	newVoucher.VoucherCurrency=currency
	config.DB.Create(&newVoucher)
	c.JSON(200,gin.H{
		"message":"New Voucher has been added successfully",
	})
}
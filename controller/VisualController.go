package controller

import (
	"github.com/Seals29/config"
	"github.com/Seals29/models"
	"github.com/gin-gonic/gin"
)

func GetCategoryCount(c *gin.Context) {
	allCateogry:= []models.ProductCategory{}
	config.DB.Find(&allCateogry)
	// c.JSON(200,&allCateogry)
	var counts []struct{
		Name string
		Count int
	}
	res := config.DB.Model(&models.Product{}).
	Joins("JOIN product_categories ON product_categories.name=products.category").Select("product_categories.name, COUNT(*) as count").Group("product_categories.name").Scan(&counts)
	if res.Error!=nil{
		c.JSON(400,&res)
	}
	c.JSON(200,&counts)

}
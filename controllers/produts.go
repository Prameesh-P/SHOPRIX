package controllers

import (
	D "github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path/filepath"
	"strconv"
)

var Products []struct {
	Product_ID  uint
	ProductName string
	ActualPrice uint
	Price       string
	Image       string
	SubPic1     string
	SubPic2     string
	Description string
	Color       string
	Brands      string
	Stock       uint
	Category    string
	Size        uint
}

func ListAllCategory(c *gin.Context) {
	var brandss models.Brand
	var categorys models.Category
	var shoesizes models.ShoeSize
	brand := D.Db.Find(&brandss)
	if brandSearch := c.Query("brandsearch"); brandSearch != "" {
		if brand.Error != nil {
			c.JSON(400, gin.H{
				"error": brand.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	category := D.Db.Find(&categorys)
	if categorySearch := c.Query("categorysearch"); categorySearch != "" {
		category = D.Db.Where("category LIKE=?", "%"+categorySearch+"%").Find(&categorys)

		if category.Error != nil {
			c.JSON(404, gin.H{
				"err": category.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	size := D.Db.Find(&shoesizes)
	if sizeSearch := c.Query("sizesearch"); sizeSearch != "" {
		sizes, _ := strconv.Atoi(sizeSearch)
		size = D.Db.Where("size = ?", sizes).Find(&shoesizes)
		if size.Error != nil {
			c.JSON(404, gin.H{
				"err": size.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	c.JSON(200, gin.H{
		"available brands":     brandss,
		"available categories": categorys,
		"available sizes":      shoesizes,
	})
}
func ProductAdding(c *gin.Context) {
	prodname := c.PostForm("productname")
	price := c.PostForm("price")
	Price, _ := strconv.Atoi(price)
	description := c.PostForm("description")
	color := c.PostForm("color")
	brand := c.PostForm("brandID")
	brands, _ := strconv.Atoi(brand)
	stock := c.PostForm("stock")
	Stock, _ := strconv.Atoi(stock)
	catogory := c.PostForm("catogoryID")
	catogoryy, _ := strconv.Atoi(catogory)
	size := c.PostForm("sizeID")

	Size, _ := strconv.Atoi(size)
	// images adding
	imagepath, _ := c.FormFile("image")
	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)

	subpic1, _ := c.FormFile("subpic1")
	extension = filepath.Ext(subpic1.Filename)
	subpic11 := uuid.New().String() + extension
	c.SaveUploadedFile(subpic1, "./public/images"+image)

	subpic2path, _ := c.FormFile("subpic2")
	extension = filepath.Ext(subpic2path.Filename)
	subpic2 := uuid.New().String() + extension
	c.SaveUploadedFile(subpic2path, "./public/images"+image)
	discont := c.PostForm("discount")
	discount, _ := strconv.Atoi(discont)
	BrandDiscount := c.PostForm("BrandDiscount")
	brandDiscount, _ := strconv.Atoi(BrandDiscount)
}

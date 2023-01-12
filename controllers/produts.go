package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Prameesh-P/SHOPRIX/database"
	D "github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Products []struct {
	ProductID   uint
	ProductName string
	ActualPrice uint
	Price       string
	Image       string
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
func ApplyDiscount(c *gin.Context) {
	var brand struct {
		Brand_id uint
		Discount uint
	}
	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	record := database.Db.Model(&models.Brand{}).Where("id=?", brand.Brand_id).Update("discount", brand.Discount)
	if record.Error == nil {
		c.JSON(200, gin.H{
			"discount": brand.Discount,
			"msg":      "Brand discount added succesfully",
		})

	}
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

	discont := c.PostForm("discount")
	discount, _ := strconv.Atoi(discont)
	BrandDiscount := c.PostForm("BrandDiscount")
	brandDiscount, _ := strconv.Atoi(BrandDiscount)
	var Discount int
	//inserting brand discount on to the products
	insert := database.Db.Raw("update brands set discount=? where id=?", brandDiscount, brands).Scan(&models.Brand{})
	if insert.Error != nil {
		c.JSON(404, gin.H{
			"err": insert.Error.Error(),
		})
		c.Abort()
		return
	}
	//comparing whcih type of discount is greater
	if brandDiscount > discount {
		Discount = (Price * brandDiscount) / 100

	} else {
		Discount = (Price * discount) / 100
	}

	// Discount = (Price * discount) / 100
	var count uint
	database.Db.Raw("select count(*) from products where product_name=?", prodname).Scan(&count)
	fmt.Println(count)
	if count > 0 {
		c.JSON(404, gin.H{
			"msg": "A product with same name already exists",
		})
		c.Abort()
		return
	}
	product := models.Product{

		ProductName: prodname,

		Price:       uint(Price) - uint(Discount),
		Color:       color,
		Description: description,
		ActualPrice: uint(Price),
		Discount:    uint(discount),

		BrandId:    uint(brands),
		CategoryID: uint(catogoryy),
		ShoeSizeID: uint(Size),
		Image:      image,

		Stock: uint(Stock),
	}

	record := database.Db.Create(&product)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"msg": "product already exists",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"msg": "added succesfully",
	})

}

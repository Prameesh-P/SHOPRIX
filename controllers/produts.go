package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Products []struct {
	Product_ID   uint
	Product_Name string
	Actual_price uint
	Price        string
	Image        string
	Description  string
	Color        string
	Brands       string
	Stock        uint
	Catogory     string
	Size         uint
}

func ListAllCategory(c *gin.Context) {
	var brandss models.Brand
	var categorys models.Category
	var shoesizes models.ShoeSize
	brand := database.Db.Find(&brandss)
	if brandSearch := c.Query("brandsearch"); brandSearch != "" {
		if brand.Error != nil {
			c.JSON(400, gin.H{
				"error": brand.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	category := database.Db.Find(&categorys)
	if categorySearch := c.Query("categorysearch"); categorySearch != "" {
		category = database.Db.Where("category LIKE=?", "%"+categorySearch+"%").Find(&categorys)

		if category.Error != nil {
			c.JSON(404, gin.H{
				"err": category.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	size := database.Db.Find(&shoesizes)
	if sizeSearch := c.Query("sizesearch"); sizeSearch != "" {
		sizes, _ := strconv.Atoi(sizeSearch)
		size = database.Db.Where("size = ?", sizes).Find(&shoesizes)
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
	
	prodname := c.Request.FormValue("productname")
	price := c.Request.FormValue("price")
	Price, _ := strconv.Atoi(price)
	description := c.Request.FormValue("description")
	color := c.Request.FormValue("color")
	brand := c.Request.FormValue("brandID")
	brands, _ := strconv.Atoi(brand)
	stock := c.Request.FormValue("stock")
	Stock, _ := strconv.Atoi(stock)
	catogory := c.Request.FormValue("catogoryID")
	catogoryy, _ := strconv.Atoi(catogory)
	size := c.Request.FormValue("sizeID")
	
	Size, _ := strconv.Atoi(size)
	// images adding
	imagepath, _ := c.FormFile("image")
	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)

	discont := c.Request.FormValue("discount")
	discount, _ := strconv.Atoi(discont)
	BrandDiscount := c.Request.FormValue("BrandDiscount")
	brandDiscount, _ := strconv.Atoi(BrandDiscount)
	var Discount int
	// inserting brand discount on to the products
	insert := database.Db.Raw("update brands set discount=? where id=?", brandDiscount, brands).Scan(&models.Brand{})
	if insert.Error != nil {
		c.JSON(404, gin.H{
			"err": insert.Error.Error(),
		})
		c.Abort()
		return
	}
	// comparing whcih type of discount is greater
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
	// fmt.Println(Discount)
	product := models.Product{
		ProductName: prodname,
		Price:       uint(Price),
		Color:       color,
		Description: description,
		ActualPrice: uint(Price),
		Discount:    uint(Discount),
		BrandId:    uint(brands),
		CategoryID: uint(catogoryy),
		ShoeSizeID: uint(Size),
		Image: image,
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

type EditProductsData struct {
	ProductName string `json:"productName"`
	Price       uint   `json:"price"`
	Image       string `json:"image"`
	Color       string `json:"color"`
}

func EditProducts(c *gin.Context) {
	params := c.Param("id")
	var edit EditProductsData
	// var editProduct models.Product
	// if err := c.ShouldBindJSON(&editProduct); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	c.Abort()
	// 	return
	// }
	edit.ProductName = c.Request.FormValue("productname")
	editprice := c.Request.FormValue("price")
	edit_price,_:=strconv.Atoi(editprice)
	edit.Price=uint(edit_price)
	edit.Color = c.Request.FormValue("color")
	imagePath, _ := c.FormFile("image")
	extension := filepath.Ext(imagePath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagePath, "./public/images"+image)
	edit.Image = image
	record := database.Db.Model(Products).Where("id=?", params).Updates(models.Product{ProductName: edit.ProductName, Price: edit.Price,
		Image: edit.Image, Color: edit.Color})
	if record.Error != nil {
		c.JSON(404, gin.H{
			"error": record.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Edit product successfully..!!!",
	})

}

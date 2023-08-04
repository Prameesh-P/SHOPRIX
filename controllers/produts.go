package controllers

import (
	"fmt"
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
	"strconv"
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
//show product list
// @Summary Admin productlist viewer
// @ID admin-productlist-view
// @Description admin can view productlist 
// @Tags Admin Product
// @Produce json
// @Param brandsearch query string true "brand of the  product"
// @Param categorysearch query string true "catogery of the  product"
// @Param sizesearch query string true "size of the  product"
// @Success 200 
// @Failure 400 
// @Router /admin/getcategory [get]
func ListAllCategory(c *gin.Context) {
	var brandss models.Brand
	var categorys models.Category
	var shoesizes models.ShoeSize
	if brandSearch := c.Query("brandsearch"); brandSearch != "" {
		brand := database.Db.Raw("SElECT * FROM brands WHERE brands=?", brandSearch).Scan(&brandss)
		if brand.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": brand.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	if categorySearch := c.Query("categorysearch"); categorySearch != "" {
		category := database.Db.Raw("SElECT * FROM categories WHERE Category=?", categorySearch).Scan(&categorys)
		if category.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": category.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	if sizeSearch := c.Query("sizesearch"); sizeSearch != "" {
		Sizes, _ := strconv.Atoi(sizeSearch)
		sizes := uint(Sizes)
		size := database.Db.Raw("SElECT * FROM shoe_sizes WHERE size=?", sizes).Scan(&shoesizes)
		if size.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": size.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"available brands":     brandss,
		"available categories": categorys,
		"available sizes":      shoesizes,
	})
}
type Brand struct {
	Brand_id uint `json:"brand_id"`
	Discount uint `json:"discount"`
}

//show product discount
// @Summary Admin product discount
// @ID admin-product-discount
// @Description admin can discount product 
// @Tags Admin Product
// @Produce json
// @Param Brand body  Brand true "brand discount"
// @Success 200 
// @Failure 400 
// @Router /admin/applydiscount [put]
func ApplyDiscount(c *gin.Context) {
	var brand Brand
	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	record := database.Db.Model(&models.Brand{}).Where("id=?", brand.Brand_id).Update("discount", brand.Discount)
	if record.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"discount": brand.Discount,
			"msg":      "Brand discount added succesfully",
		})
	}
}


//show product list
// @Summary Admin productadd 
// @ID admin-product add
// @Description admin can add product
// @Tags Admin Product
// @Produce json
// @Param productname formData string true "name of the  product"
// @Param price formData string true "price of the  product"
// @Param description formData string true "discription of the  product"
// @Param color formData string true "color of the  product"
// @Param brandID formData string true "brandID of the  product"
// @Param stock formData string true "stock of the  product"
// @Param catogoryID formData string true "cotogeryID of the  product"
// @Param sizeID formData string true "sizeID of the  product"
// @Param image formData file false "Upload a product image"
// @Param discount formData string true "discount of the  product"
// @Param BrandDiscount formData string true "brandDiscount of the  product"
// @Success 200 
// @Failure 400 
// @Router /admin/addproducts [post]
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
	fmt.Printf("jgsdigj %s", extension)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"err": insert.Error.Error(),
		})
		c.Abort()
		return
	}
	// comparing which type of discount is greater
	if brandDiscount > discount {
		Discount = (Price * brandDiscount) / 100
	} else {
		Discount = (Price * discount) / 100
	}
	var count uint
	database.Db.Raw("select count(*) from products where product_name=?", prodname).Scan(&count)
	fmt.Println(count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "A product with same name already exists",
		})
		c.Abort()
		return
	}
	product := models.Product{
		ProductName: prodname,
		Price:       uint(Price),
		Color:       color,
		Description: description,
		ActualPrice: uint(Price) - uint(Discount),
		Discount:    uint(Discount),
		BrandId:     uint(brands),
		CategoryID:  uint(catogoryy),
		ShoeSizeID:  uint(Size),
		Image:       image,
		Stock:       uint(Stock),
	}
	record := database.Db.Create(&product)
	if record.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "product already exists",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "added_successfully",
	})
}

type EditProductsData struct {
	ProductName string `json:"productName"`
	Price       uint   `json:"price"`
	Brand       string `json:"brand"`
	Color       string `json:"color"`
}

// @Summary Admin product edit
// @ID admin-product edit
// @Description admin can edit product
// @Tags Admin Product
// @Produce json
// @Param id query string true "id of the  product"
// @Param EditProductData body  EditProductsData true "edit product data"
// @Success 200 
// @Failure 400 
// @Router /admin/editproducts [put]
func EditProducts(c *gin.Context) { //admin
	params := c.Param("id")
	var editProducts EditProductsData
	if err := c.ShouldBindJSON(&editProducts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		c.Abort()
		return
	}
	var products models.Product
	record := database.Db.Model(products).Where("product_id=?", params).Updates(models.Product{ProductName: editProducts.ProductName,
		Price: editProducts.Price, Color: editProducts.Color})
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "updated_successfully"})

}

// @Summary Admin product delete
// @ID admin-product delete
// @Description admin can delete product
// @Tags Admin Product
// @Produce json
// @Param id query string true "id of the  product"
// @Success 200 
// @Failure 400 
// @Router /admin/deleteproducts/ [delete]
func DeleteProductById(c *gin.Context) { //admin
	params := c.Param("id")
	var products models.Product
	var count uint
	database.Db.Raw("select count(product_id) from products where product_id=?", params).Scan(&count)
	if count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "product does not exist",
		})
		c.Abort()
		return
	}

	record := database.Db.Raw("delete from products where product_id=?", params).Scan(&products)
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "deleted successfully"})
}


//show product 
// @Summary User productid viewer
// @ID user-productid-view
// @Description user can view productid with with name
// @Tags Users Product
// @Produce json
// @Param product-name formData string true "name of the  product"
// @Success 200 
// @Failure 400 
// @Router /user/show-product-id [get]
func ShowProductsID(c *gin.Context) {
	var product models.Product
	params := c.PostForm("product-name")
	fmt.Println(params)
	record := database.Db.Raw("select product_id from products where product_name=?", params).Scan(&product)
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": record.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "true",
		"Product id": product.ProductId,
	})
}

//show product 
// @Summary User product viewer
// @ID user-product-view with id
// @Description user can view product with with id
// @Tags Users Product
// @Produce json
// @Param id formData string true "id of the  product"
// @Success 200 
// @Failure 400 
// @Router /user/get-productbyid [get]
func GetProductByID(c *gin.Context) { //user
	params := c.Request.FormValue("id")
	fmt.Println(params)
	// var product models.Product
	record := database.Db.Raw("SELECT product_id,product_name,price,image,color,stock,brands.brands FROM products join brands on products.brand_id = brands.id where product_id=?", params).Scan(&Products)
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": Products})
}

//show product 
// @Summary User product viewer with search
// @ID user-product-view with search
// @Description user can view product with with id
// @Tags Users Product
// @Produce json
// @Param search formData string true "searched data of the  product"
// @Param sort formData string true "sort products"
// @Success 200 
// @Failure 400 
// @Router /user/view-products [get]
func ProductView(c *gin.Context) {
	record := database.Db.Raw("SELECT product_id,product_name,actual_price,price,image,color,description,stock,brands.brands,categories.category,shoe_sizes.size FROM products join brands on products.brand_id = brands.id join categories on products.category_id=categories.id join shoe_sizes on products.shoe_size_id=shoe_sizes.id").Scan(&Products)
	fmt.Println(record)
	if s := c.Request.FormValue("search"); s != "" { //search
		database.Db.Raw("SELECT product_id,product_name,actual_price,price,image,color,description,stock,brands.brands,categories.category,shoe_sizes.size FROM products join brands on products.brand_id = brands.id join categories on products.category_id=categories.id join shoe_sizes on products.shoe_size_id=shoe_sizes.id where product_name like ?", "%"+s+"%").Scan(&Products)
	}
	if sort := c.Request.FormValue("sort"); sort != "" { //sort
		database.Db.Raw("SELECT product_id,product_name,actual_price,price,image,color,description,stock,brands.brands,categories.category,shoe_sizes.size FROM products join brands on products.brand_id = brands.id join categories on products.category_id=categories.id join shoe_sizes on products.shoe_size_id=shoe_sizes.id  order by price ?", sort).Scan(&Products)
	}
	c.JSON(http.StatusOK, gin.H{
		"products": Products,
	})
}

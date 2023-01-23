package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	var user models.User
	var product models.Product
	userEmail := c.Request.FormValue("user")
	fmt.Printf("huuu%s", userEmail)
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	var ProductDetails struct {
		ProductID uint
		Quantity  uint
	}
	c.ShouldBindJSON(&ProductDetails)
	database.Db.Raw("select price,stock from products where product_id=?", ProductDetails.ProductID).Scan(&product)
	total := product.Price * ProductDetails.Quantity
	productId := ProductDetails.ProductID
	productQuantity := ProductDetails.Quantity
	cart := models.Cart{
		ProductID:  productId,
		Quantity:   productQuantity,
		UserId:     user.ID,
		TotalPrice: total,
	}
	record1 := database.Db.Create(&cart)
	if record1.Error != nil {
		c.JSON(404, gin.H{
			"err": record1.Error.Error(),
		})
		c.Abort()
		return
	}
	var Cart []models.Cart
	database.Db.Raw("select cart_id,product_id from carts where user_id=?", user.ID).Scan(&Cart)
	for _, v := range Cart {
		fmt.Println("loop started..!")
		if v.ProductID == productId {
			fmt.Println("in the condition..")
			database.Db.Raw("select quantity from carts where product_id=? and user_id=?", ProductDetails.ProductID, user.ID).Scan(&Cart)
			totlV := (productQuantity + cart.Quantity) * product.Price
			database.Db.Raw("update carts set quantity=?,total_price where  product_id=? and user_id=?", productQuantity+cart.Quantity, totlV, productId, user.ID)
			c.JSON(http.StatusOK, gin.H{
				"msg": "quantiity updated successfully",
			})
			c.Abort()
			return
		}
	}
	record := database.Db.Create(&cart)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	database.Db.Raw("select product_name,brands.brand from products join carts.product_id=products.product_id join brands on brands.id=products.brand_id where products.product_id=?", ProductDetails.ProductID).Scan(&cart)
	c.JSON(200, gin.H{
		"userId": user.ID,
		"msg":    "added to cart",
	})
}

type CartsInfo []struct {
	UserId      string
	ProductId   string
	ProductName string
	Price       string
	Email       string
	Quantity    string
	TotalAmount uint
	TotalPrice  string
}

func ViewCart(c *gin.Context) {
	var cart CartsInfo
	var products models.Product
	var user models.User
	userEmail := c.Request.FormValue("user")
	product := c.Request.FormValue("product")
	newProduct, _ := strconv.Atoi(product)
	quantity := c.Request.FormValue("quantity")
	newQuantity, _ := strconv.Atoi(quantity)
	database.Db.Raw("select price from prooducts where product_id=?", uint(newProduct)).Scan(&products)
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	total := products.Price * uint(newQuantity)
	if newQuantity >= 1 {
		database.Db.Raw("update carts set quantity=?,total_price=? where user_id=? and product_id=?", newQuantity, total, user.ID, newProduct).Scan(&cart)
	} else if newQuantity <= 0 {
		database.Db.Raw("dalete from carts where product_id=? and user_id=?", newProduct, user.ID).Scan(&cart)
	}
	record := database.Db.Raw("select  products.product_id, products.product_name,products.price,carts.user_id,users.email ,carts.quantity,total_price from carts join products on products.product_id=carts.product_id join users on carts.user_id=users.id where users.email=? ", userEmail).Scan(&cart)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}

	var totalcartvalue uint

	database.Db.Raw("select sum(total_price) as total from carts where user_id=?", user.ID).Scan(&totalcartvalue)
	c.JSON(200, gin.H{
		"cart":  cart,
		"total": totalcartvalue,
	})
}
func CheckOutAddress(c *gin.Context) {
	useremail := c.GetString("user")
	var user models.User
	database.Db.Raw("select id from users where email=?", useremail).Scan(&user)

	Name := c.PostForm("name")
	Phonenum := c.PostForm("phone_number")
	phonenum, _ := strconv.Atoi(Phonenum)
	pincod := c.PostForm("pincode")
	pincode, _ := strconv.Atoi(pincod)
	area := c.PostForm("area")
	houseadd := c.PostForm("house")
	landmark := c.PostForm("landmark")
	city := c.PostForm("city")
	address := models.Address{
		UserID:   user.ID,
		Name:     Name,
		PhoneNum: uint(phonenum),
		Pincode:  uint(pincode),
		Area:     area,
		House:    houseadd,
		LandMark: landmark,
		City:     city,
	}
	record := database.Db.Create(&address)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Address added",
	})
}

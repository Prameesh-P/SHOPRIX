package controllers

import (
	"fmt"
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	var cartss CartsInfo
	var products models.Product
	var user models.User
	userEmail := c.Request.FormValue("user")
	product := c.Request.FormValue("product")
	newProduct, _ := strconv.Atoi(product)
	quantity := c.Request.FormValue("quantity")
	newQuantity, _ := strconv.Atoi(quantity)
	database.Db.Raw("select price from products where product_id=?", uint(newProduct)).Scan(&products)
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	total := products.Price * uint(newQuantity)
	if newQuantity >= 1 {
		database.Db.Raw("update carts set quantity=?,total_price=? where user_id=? and product_id=?", newQuantity, total, user.ID, newProduct).Scan(&cartss)
	} else if newQuantity <= 0 {
		database.Db.Raw("dalete from carts where product_id=? and user_id=?", newProduct, user.ID).Scan(&cartss)
	}
	record := database.Db.Raw("select  products.product_id, products.product_name,products.price,carts.user_id,users.email ,carts.quantity,total_price from carts join products on products.product_id=carts.product_id join users on carts.user_id=users.id where users.email=? ", userEmail).Scan(&cartss)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	for _, v := range cartss {
		productID := v.ProductId
		//ProdctId,_:=strconv.Atoi(productID)
		userID := v.UserId
		//UserId,_:=strconv.Atoi(userID)
		ProName := v.ProductName
		quantity := v.Quantity
		//Quantity,_:=strconv.Atoi(quantity)
		PEmail := v.Email
		Tprice := v.TotalPrice
		//TotalPriceL,_:=strconv.Atoi(Tprice)
		carts := models.Cartsinfo{
			User_id:      userID,
			Product_id:   productID,
			Product_Name: ProName,
			Quantity:     quantity,
			Email:        PEmail,
			Total_Price:  Tprice,
		}
		database.Db.Create(&carts)
	}
	var totalcartvalue uint

	fmt.Println(cartss)
	database.Db.Raw("select sum(total_price) as total from carts where user_id=?", user.ID).Scan(&totalcartvalue)
	c.JSON(200, gin.H{
		"cartss": cartss,
		"total":  totalcartvalue,
	})
}
func CheckOutAddress(c *gin.Context) {
	useremail := c.PostForm("user")
	var user models.User
	database.Db.Raw("select id from users where email=?", useremail).Scan(&user)

	Name := c.PostForm("name")
	Phonenum1 := c.PostForm("phone_number")
	phonenum, _ := strconv.Atoi(Phonenum1)
	pincod := c.PostForm("pincode")
	pincode, _ := strconv.Atoi(pincod)
	area := c.PostForm("area")
	houseadd := c.PostForm("house")
	landmark := c.PostForm("landmark")
	city := c.PostForm("city")
	address := models.Address{
		UserId:    user.ID,
		Name:      Name,
		PhoneNum:  phonenum,
		Pincode:   pincode,
		Area:      area,
		House:     houseadd,
		Land_mark: landmark,
		City:      city,
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

var Address []struct {
	UserId       uint
	Address_id   uint
	Name         string
	Phone_number uint
	Pincode      uint
	Area         string
	House        string
	Landmark     string
	City         string
}
var Flag = 0

func CheckOut(c *gin.Context) {
	var user models.User
	var cart models.Cart
	fmt.Println(cart)
	var carts CartsInfo
	userEmail := c.PostForm("user")
	wallet := c.Query("wallet")
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	precord := database.Db.Raw("select  products.product_id, products.product_name,products.price,carts.user_id,users.email ,carts.quantity,total_price from carts join products on products.product_id=carts.product_id join users on carts.user_id=users.id where users.email=? ", userEmail).Scan(&carts)
	if precord.Error != nil {
		c.JSON(404, gin.H{"err": precord.Error.Error()})
		c.Abort()
		return
	}
	var totalCartValue uint
	var address models.Address
	addres := c.PostForm("addressID")
	addressID, _ := strconv.Atoi(addres)
	PaymentMethod := c.PostForm("PaymentMethod")
	//Coupons := c.PostForm("coupon")
	//cod := "COD"
	//fmt.Printf("codd..%s", cod)
	//razorpay := "RazorPay"
	database.Db.Raw("select sum(total_price) as total from carts where user_id=?", user.ID).Scan(&totalCartValue)
	if PaymentMethod == "COD" {
		fmt.Println(carts)
		fmt.Println("hyyyyyyyy..............")
		for _, v := range carts {
			pud := v.UserId
			puid, _ := strconv.Atoi(pud)
			fmt.Println(puid)
			pid := v.ProductId
			proID, _ := strconv.Atoi(pid)
			pname := v.ProductName
			fmt.Println(pname)
			pprice := v.Price
			Pprice, _ := strconv.Atoi(pprice)
			pQuantity := v.Quantity
			PQuantity, _ := strconv.Atoi(pQuantity)
			fmt.Println(CreateOrderId())
			//appliedCoupen := CoupenDisc.Coupen_code
			totalAmount := uint(PQuantity) * uint(Pprice)
			//discount := (totalAmount * CoupenDisc.Discount) / 100
			TotalAmount := totalAmount
			orderedItems := models.OrderedItems{UserId: uint(puid), Product_id: uint(proID),
				Product_Name: pname, Price: pprice, OrdersID: CreateOrderId(),
				Order_Status: "confirmed", PaymentMethod: PaymentMethod, Payment_Status: "pending", Total_amount: TotalAmount}
			database.Db.Create(&orderedItems)

		}

		c.JSON(200, gin.H{
			"msg": "Cash on delivery placed..!!",
		})
	}

	record := database.Db.Raw("select address_id, user_id,name,phone_num,pincode,house,area,land_mark,city from addresses where user_id=?", user.ID).Scan(&Address)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	database.Db.Raw("select address_id,user_id,name from addresses where address_id=?", addressID).Scan(&address)

	fmt.Printf("users %d", user.ID)
	fmt.Printf("address %d", address.UserId)
	if address.UserId != user.ID {
		c.JSON(200, gin.H{
			"msg": "enter valid address id",
		})
	}
	//addressID == int(address.AddressId) &&
	var wallets models.Wallet
	database.Db.Raw("select wallet_balance from wallets where user_id=?", user.ID).Scan(&user)
	if wallet == "use-wallet" && address.UserId == user.ID && PaymentMethod == "wallet" {
		//fmt.Println("adsfasfas")
		if wallets.WalletBalance > totalCartValue {
			c.JSON(400, gin.H{
				"msg": "can't apply wallet money on this transaction..Try another method..!!\n wallet money is low..",
			})
			c.Abort()
			return
		} else if wallets.WalletBalance > totalCartValue {
			wallets.WalletBalance = wallets.WalletBalance - totalCartValue
			database.Db.Model(&wallets).Where("user_id=?", user.ID).Update("wallet-balance", wallets.WalletBalance)
			walletOrder := models.Orders{
				UserId:         user.ID,
				Order_id:       CreateOrderId(),
				Total_Amount:   totalCartValue,
				PaymentMethod:  "wallet",
				Payment_Status: "payment completed",
				Order_Status:   "order placed",
				Address_id:     uint(addressID),
			}

			query := database.Db.Create(&walletOrder)
			if query.Error != nil {
				c.JSON(400, gin.H{
					"err": query.Error.Error(),
				})
				c.Abort()
			}
			var orderedItems models.OrderedItems
			database.Db.Raw("update orderd_items set  order_status=?,payment_status=?,payment_method=? where user_id=?", "orderplaced", "payment completed", "wallet", user.ID).Scan(&orderedItems)

		}

	}
	if Flag > 0 {
		c.JSON(400, gin.H{
			"msg": "item already placed",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "order Placed",
	})

	c.JSON(300, gin.H{
		"address":          Address,
		"total cart value": totalCartValue,
	})
	Flag++
	totalCartValue = 0
	carts = nil
}

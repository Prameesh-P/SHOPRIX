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
	var cart models.Cart
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

func CheckOut(c *gin.Context) {
	var user models.User
	var cart models.Cart
	var carts CartsInfo
	userEmail := c.Request.FormValue("user")
	wallet := c.Query("ApplyWallet")
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	precord := database.Db.Raw("select  products.product_id, products.product_name,products.price,carts.user_id,users.email ,carts.quantity,total_price from carts join products on products.product_id=carts.product_id join users on carts.user_id=users.id where users.email=? ", userEmail).Scan(&cart)
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
	Coupons := c.PostForm("coupon")
	cod := "COD"
	razorpay := "RazorPay"
	database.Db.Raw("select sum(total_price) as total from carts where user_id=?", user.ID).Scan(&totalCartValue)

	var CoupenDisc struct {
		Coupen_code string
		Discount    uint
		Count       uint
		Validity    uint
	}
	var AppliedCoupon struct {
		user_id     uint
		Coupen_code string
		count       uint
	}
	var flag = 0
	if Coupons == "" {
		c.JSON(300, gin.H{
			"msg": "Enter a coupon if you have any coupon",
		})
	} else if Coupons != "" {
		flag = 1
		database.Db.Raw("select coupon_code,discount,validity,count(*) as count from coupons where coupon_code=? group by discount,validity,coupon_code", Coupons).Scan(&CoupenDisc)
		database.Db.Raw("select user_id,coupon_code,count(*) as count from applied_coupons where coupon_code=? and user_id=? group by user_id,coupon_code", Coupons, user.ID).Scan(&AppliedCoupon)
		if AppliedCoupon.count > 0 {
			c.JSON(300, gin.H{
				"msg": "Coupon already applied",
			})
			flag = 2
		}
		if CoupenDisc.Count <= 0 {
			c.JSON(300, gin.H{
				"msg": "not a valid coupon",
			})
			flag = 2
		}
		//if CoupenDisc.Validity < time.Now().Local().Unix() && CoupenDisc.Validity > 1 {
		//
		//	c.JSON(300, gin.H{
		//		"msg": "coupon expired",
		//	})
		//	flag = 2
		//
		//}
		if flag == 1 {

			Discount := (totalCartValue * CoupenDisc.Discount) / 100
			totalCartValue = totalCartValue - Discount

		}
		if PaymentMethod == cod && PaymentMethod == razorpay {
			for _, v := range carts {
				pud := v.UserId
				puid, _ := strconv.Atoi(pud)
				pid := v.ProductId
				proID, _ := strconv.Atoi(pid)
				pname := v.ProductName
				pprice := v.Price
				Pprice, _ := strconv.Atoi(pprice)
				pQuantity := v.Quantity
				PQuantity, _ := strconv.Atoi(pQuantity)
				appliedCoupen := CoupenDisc.Coupen_code
				totalAmount := uint(PQuantity) * uint(Pprice)
				discount := (totalAmount * CoupenDisc.Discount) / 100
				TotalAmount := totalAmount - discount
				orderedItems := models.OrderedItems{UserId: uint(puid), Product_id: uint(proID),
					Product_Name: pname, Price: pprice, OrdersID: CreateOrderId(), Applied_Coupons: appliedCoupen,
					Order_Status: "confirmed", Payment_Status: "pending", Total_amount: TotalAmount}
				database.Db.Create(&orderedItems)
			}
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

		c.JSON(300, gin.H{
			"address":          Address,
			"total cart value": totalCartValue,
		})
		//rand.Seed(time.Now().UnixNano())
		//value := rand.Intn(9999999999-1000000000) + 1000000000
		//id := strconv.Itoa(value)
		//orderID := "OID" + id
		if address.UserId != user.ID {
			c.JSON(200, gin.H{
				"msg": "enter valid address id",
			})
		}
		database.Db.Raw("select wallet_balance from users where id=?", user.ID).Scan(&user)
		if wallet == "use-wallet" && addressID == int(address.AddressId) && address.UserId == user.ID {
			if user.WalletBalance > totalCartValue {
				c.JSON(400, gin.H{
					"msg": "can't apply wallet money on this transaction..Try another method..!!",
				})
				c.Abort()
				return
			} else if user.WalletBalance > totalCartValue {
				user.WalletBalance = user.WalletBalance - totalCartValue
				database.Db.Model(&user).Where("id=?", user.ID).Update("wallet-balance", user.WalletBalance)

				walletOrder := models.Orders{
					UserId:         user.ID,
					Order_id:       CreateOrderId(),
					Total_Amount:   totalCartValue,
					PaymentMethod:  "wallet",
					Payment_Status: "payment completed",
					Order_Status:   "order placed",
					Address_id:     uint(addressID),
				}
				database.Db.Create(&walletOrder)
				totalCartValue = 0
				var orderedItems models.OrderedItems
				database.Db.Raw("update orderd_items set  order_status=?,payment_status=?,payment_method=? where user_id=?", "orderplaced", "payment completed", "wallet", user.ID).Scan(&orderedItems)
			}
		}
	}
	c.JSON(200, gin.H{
		"msg": "order Placed",
	})
}

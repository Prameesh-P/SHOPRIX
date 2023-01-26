package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID            uint   `json:"id" gorm:"primaryKey;unique"  `
	FirstName     string `json:"first_name" validate:"required,min=2,max=100"`
	LastName      string `json:"last_name"  validate:"required,min=2,max=100"`
	Email         string `json:"email" gorm:"unique" validate:"email,required" `
	Password      string `json:"password" validate:"required,min=6"`
	Phone         string `json:"phone"  validate:"required"`
	BlockStatus   bool   `json:"block_Status" `
	Country       string `json:"country"`
	City          string `json:"city"`
	Pincode       string `json:"pincode"`
	LandMark      string `json:"landmark"`
	WalletBalance uint   `json:"wallet_balance"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Address struct {
	AddressId uint   `json:"address_id" gorm:"primaryKey"  `
	UserId    uint   `json:"user_id"  gorm:"not null" `
	Name      string `json:"name"  gorm:"not null" `
	PhoneNum  int    `json:"phone_number"  gorm:"not null" `
	Pincode   int    `json:"pincode"  gorm:"not null" `
	House     string `json:"house"   `
	Area      string `json:"area"   `
	Email     string `json:"email"`
	Land_mark string `json:"landmark"  gorm:"not null" `
	City      string `json:"city"  gorm:"not null" `
}
type Admin struct {
	gorm.Model
	Email    string `json:"email" validate:"required,min=5,max=100"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}
func (a *Admin) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	a.Password = string(bytes)
	return a.Password, nil
}
func (a *Admin) CheckPassword(incomingPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(incomingPass))
	if err != nil {
		return err
	}
	return nil
}

type Product struct {
	gorm.Model
	ProductId   uint   `json:"product_id" gorm:"autoIncrement" `
	ProductName string `json:"product_name" gorm:"not null"  `
	Price       uint   `json:"price" gorm:"not null"  `
	ActualPrice uint   `json:"actual_price" gorm:"not null"`
	Image       string `json:"image" gorm:"not null"  `
	Stock       uint   `json:"stock"  `
	Color       string `json:"color" gorm:"not null"  `
	Description string `json:"description"   `
	Brand       Brand
	BrandId     uint `json:"brand_id" `
	Cart        Cart
	CartId      uint `json:"cart_id" `
	Category    Category
	CategoryID  uint
	ShoeSize    ShoeSize
	ShoeSizeID  uint
	Discount    uint
}
type Brand struct {
	ID       uint   `json:"id" gorm:"primaryKey"  `
	Brands   string `json:"brands" gorm:"not null"  `
	Discount uint   `json:"discount"`
}
type Category struct {
	ID       uint `json:"id" gorm:"primaryKey"  `
	Category string
}
type ShoeSize struct {
	ID   uint `json:"id" gorm:"primaryKey"  `
	Size uint `json:"size"`
}
type Cart struct {
	CartId     uint `json:"cart_id" gorm:"primaryKey"`
	UserId     uint `json:"user_id"   `
	ProductID  uint `json:"product_id"  `
	Quantity   uint `json:"quantity" `
	TotalPrice uint `json:"total_price"   `
}
type Cartsinfo struct {
	gorm.Model
	User_id      string
	Product_id   string
	Product_Name string
	Price        string
	Email        string
	Quantity     string
	Total_Price  string
}

type Orders struct {
	gorm.Model
	UserId          uint   `json:"user_id"  gorm:"not null" `
	Order_id        string `json:"order_id"  gorm:"not null" `
	Total_Amount    uint   `json:"total_amount"  gorm:"not null" `
	Applied_Coupons string `json:"applied_coupons"  `
	Discount        uint   `json:"discount"   `
	PaymentMethod   string `json:"paymentmethod"  gorm:"not null" `
	Payment_Status  string `json:"payment_status"   `
	Order_Status    string `json:"order_status"   `
	Address         Address
	Address_id      uint `json:"address_id"`
}
type Coupon struct {
	gorm.Model

	Coupon_Code string `json:"coupon_code"`
	Discount    uint   `json:"discount"`
	Quantity    uint   `json:"quantity"`
	Validity    int64  `json:"validity"`
}
type OrderedItems struct {
	ID              uint `json:"id" gorm:"autoIncrement:true" `
	UserId          uint `json:"user_id" `
	Product_id      uint
	OrdersID        string
	Product_Name    string
	Price           string
	Order_Status    string
	Payment_Status  string
	PaymentMethod   string
	Applied_Coupons string
	Total_amount    uint
}
type Applied_Coupons struct {
	gorm.Model
	UserID      uint
	Coupon_Code string `json:"coupon_code" `
}

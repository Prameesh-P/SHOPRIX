package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string `json:"first_name" validate:"required,min=2,max=100"`
	LastName     string `json:"last_name"  validate:"required,min=2,max=100"`
	Email        string `json:"email" gorm:"unique" validate:"email,required" `
	Password     string `json:"password" validate:"required,min=6"`
	Phone        string `json:"phone"  validate:"required"`
	BlockStatus  bool   `json:"block_Status" `
	Token        string `json:"token"`
	RefreshToken string `json:"referesh_token" `
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

type Otp struct {
	gorm.Model
	Mobile string
	Otp    string
}
type Product struct {
	Product_id   uint   `json:"product_id" gorm:"primaryKey" `
	Product_name string `json:"product_name" gorm:"not null"  `
	Price        uint   `json:"price" gorm:"not null"  `
	Actual_Price uint   `json:"actual_price" gorm:"not null"`
	Image        string `json:"image" gorm:"not null"  `
	Cover        string `json:"cover"   `
	SubPic1      string `json:"subpic1"  `
	SubPic2      string `json:"subpic2"  `
	Stock        uint   `json:"stock"  `
	Color        string `json:"color" gorm:"not null"  `
	Description  string `json:"description"   `
	Discount     uint   `json:"discount"`

	Brand      Brand
	Brand_id   uint `json:"brand_id" `
	Cart       Cart
	Cart_id    uint `json:"cart_id" `
	Catogory   Category
	CatogoryID uint
	ShoeSize   ShoeSize
	ShoeSizeID uint
	WishList   WishList
	WishListID uint
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
	Cart_id     uint `json:"cart_id" gorm:"primaryKey"  `
	UserId      uint `json:"user_id"   `
	ProductID   uint `json:"product_id"  `
	Quantity    uint `json:"quantity" `
	Total_Price uint `json:"total_price"   `
}
type WishList struct {
	gorm.Model
	UserID     uint
	Product_id uint
}

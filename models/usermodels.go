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

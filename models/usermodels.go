package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	First_name    string `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     string `json:"last_name"  validate:"required,min=2,max=100"`
	Email         string `json:"email" gorm:"unique" validate:"email,required" `
	Password      string `json:"password" validate:"required,min=6"`
	Phone         string `json:"phone"  validate:"required"`
	Token         string `json:"token"`
	Refresh_token string `json:"referesh_token" `
}

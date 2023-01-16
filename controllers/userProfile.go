package controllers

import (
	"net/http"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
)

type Profile struct {
	FirstName string
	LastName  string
	UserEmail string
	Phone     string
	Country   string
	City      string
	Pincode   string
}

func UserProfileGet(c *gin.Context) {
	userEmail:=c.GetString("user")
	var profile Profile
	database.Db.Raw("select first_name,last_name,email,phone,country,city,pincode from users where email=?",userEmail).Scan(&profile)
	c.JSON(http.StatusOK,gin.H{
		"profile":profile,
	})
} 	
func AddUserAddress(c *gin.Context)  {
	var user models.User
	FirstName:=c.PostForm("first_name")
	LastName:=c.PostForm("last_name")
	Phone:=c.PostForm("phone")
	email:=c.PostForm("email")
	country:=c.PostForm("country")
	City:=c.PostForm("city")
	pincode:=c.PostForm("pincode")
	query:=database.Db.Raw("update users set fisrt_name=?,last_name,email=?,phone=?,country=?,city=?,piccode=?",FirstName,LastName,email,Phone,country,City,pincode).Scan(&user)
	if query.Error != nil {
		c.JSON(404,gin.H{
			"err":query.Error.Error(),
		})
		c.Abort()
		return	
	}
	c.JSON(http.StatusOK,gin.H{
		"success":"update successfully",
	})
}
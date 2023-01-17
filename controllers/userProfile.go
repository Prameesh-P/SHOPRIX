package controllers

import (
	"net/http"
	"strconv"

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
func AddUserAddressOnUserMedel(c *gin.Context)  {
	var user models.User
	country:=c.PostForm("country")
	City:=c.PostForm("city")
	pincode:=c.PostForm("pincode")
	landmark:=c.PostForm("landmark")
	query:=database.Db.Raw("update users set country=?,city=?,piccode=?,landmark=?",country,City,pincode,landmark).Scan(&user)
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
func AddAddress(c *gin.Context){
	var user models.User
	userEmail:=c.GetString("email")
	database.Db.Raw("select id from users where email=?",userEmail).Scan(&user)
	
	Name := c.GetString("name")
	Phonenum := c.GetString("phonenumber")
	phonenums, _ := strconv.Atoi(Phonenum)
	pincod := c.GetString("pincode")
	pincode, _ := strconv.Atoi(pincod)
	area := c.GetString("area")
	houseadd := c.GetString("house")
	landmark := c.GetString("landmark")
	city := c.GetString("city")

	address:=models.Address{
		UserID: user.ID,
		Name: Name,
		PhoneNum: uint(phonenums),
		Pincode: uint(pincode),
		Area: area,
		House: houseadd,
		LandMark: landmark,
		City: city,
	}
	record:=database.Db.Create(&address)
	if record.Error !=nil {
		c.JSON(404,gin.H{
			"error":record.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200,gin.H{
		"msg":"Address Added",
	})
}
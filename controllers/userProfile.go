package controllers

import (
	"net/http"
	"strconv"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
)

type Profile struct {
	Name     string
	PhoneNum string
	Area     string
	House    string
	LandMark string
	City     string
	Pincode  string
	Email    string
}

func UserProfileGet(c *gin.Context) {
	userEmail := c.Request.FormValue("user")
	var profile Profile
	database.Db.Raw("select name,phone_num,pincode,area,house,land_mark,city,email from addresses where email=?", userEmail).Scan(&profile)
	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
	})
}
func EditUserAddress(c *gin.Context) {
	var user models.User
	Name := c.Request.FormValue("name")
	Phonenum := c.Request.FormValue("phonenumber")
	phonenums, _ := strconv.Atoi(Phonenum)
	pincod := c.PostForm("pincode")
	pincode, _ := strconv.Atoi(pincod)
	area := c.PostForm("area")
	houseadd := c.PostForm("house")
	landmark := c.PostForm("landmark")
	city := c.Request.FormValue("city")
	query := database.Db.Raw("update addresses set name=?,phone_num=?,pincode=?,area=?,house=?,city=?,land_mark=?", Name, phonenums, pincode, area, houseadd, city, landmark).Scan(&user)
	if query.Error != nil {
		c.JSON(404, gin.H{
			"err": query.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "update successfully",
	})
}
func AddAddress(c *gin.Context) {
	var user models.User
	userEmail := c.PostForm("email")
	database.Db.Raw("select id from users where email=?", userEmail).Scan(&user)
	Name := c.Request.FormValue("name")
	Phonenum := c.Request.FormValue("phonenumber")
	phonenums, _ := strconv.Atoi(Phonenum)
	pincod := c.PostForm("pincode")
	pincode, _ := strconv.Atoi(pincod)
	area := c.PostForm("area")
	houseadd := c.PostForm("house")
	landmark := c.PostForm("landmark")
	city := c.Request.FormValue("city")
	address := models.Address{
		UserId:    user.ID,
		Name:      Name,
		PhoneNum:  phonenums,
		Pincode:   pincode,
		Area:      area,
		House:     houseadd,
		Land_mark: landmark,
		City:      city,
	}
	record := database.Db.Create(&address)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"error": record.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"msg": "Address Added",
	})
}

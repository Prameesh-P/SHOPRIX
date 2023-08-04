package controllers

import (
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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


// User profile get
// @Summary user details willget 
// @ID user profile get
// @Description user profile viewer
// @Tags Users Profile
// @Produce json
//@Param user formData string true "email of the user"
// @Success 200 
// @Failure 400 
// @Router /user/profile [get]
func UserProfileGet(c *gin.Context) {
	userEmail := c.Request.FormValue("user")
	var profile Profile
	database.Db.Raw("select name,phone_num,pincode,area,house,land_mark,city,email from addresses where email=?", userEmail).Scan(&profile)
	c.JSON(http.StatusOK, gin.H{
		"profile": profile,
	})
}

// User profile get
// @Summary user details edit 
// @ID user profile edit
// @Description user profile editor
// @Tags Users Profile
// @Produce json
// @Param name formData string true "name of the user"
// @Param phonenumber formData string true "phone number of the user"
// @Param pincode formData string true "pincode of the user"
// @Param area formData string true "area of the user"
// @Param house formData string true "house of the user"
// @Param landmark formData string true "landmark of the user"
// @Param city formData string true "city of the user"
// @Success 200 
// @Failure 400 
// @Router /user/profile/edit [put]
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

// User profile add
// @Summary user details add 
// @ID user profile add
// @Description user profile add
// @Tags Users Profile
// @Produce json
// @Param email formData string true "email of the user"
// @Param name formData string true "name of the user"
// @Param phonenumber formData string true "phone number of the user"
// @Param pincode formData string true "pincode of the user"
// @Param area formData string true "area of the user"
// @Param house formData string true "house of the user"
// @Param landmark formData string true "landmark of the user"
// @Param city formData string true "city of the user"
// @Success 200 
// @Failure 400 
// @Router /user/profile/add [post]
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
		Email:     userEmail,
	}
	record := database.Db.Create(&address)
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": record.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Address Added",
	})
}

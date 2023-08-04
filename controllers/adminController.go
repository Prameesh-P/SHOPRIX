package controllers

import (
	"fmt"
	"net/http"

	"github.com/Prameesh-P/SHOPRIX/authentification"
	i "github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
)
type Admin struct {
	Email    string
	Password string
}


// @Summary adminSignUp
// @ID admin-signup
// @Description Create a new admin with the specified details.
// @Tags Admin
// @Accept json
// @Produce json
// @Param user_details body  Admin true "User details"
// @Success 200 
// @Failure 400 
// @Router /admin/signup [post]
func AdminSignup(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	i.Db.First(&admin, "email=?", admin.Email)
	if admin.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This email already registered..",
		})
		return
	}
	bytes, err := admin.HashPassword(admin.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to hash",
		})
		return
	}

	admins := models.Admin{Email: admin.Email, Password: bytes}
	record := i.Db.Create(&admins)
	if record.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create Admin",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
// @Summary adminLogin
// @ID admin-login
// @Description admin login
// @Tags Admin
// @Accept json
// @Produce json
// @Param admin_details body Admin true "admin details"
// @Success 200 
// @Failure 400 
// @Router /admin/login [post]
func AdminLogin(c *gin.Context) {
	var admin Admin
	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	var admins models.Admin
	i.Db.First(&admins, "email=?", admin.Email)
	if admins.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}
	err := admins.CheckPassword(admin.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Password was wrong..!!!",
		})
		c.Abort()
		return
	}
	tokenString, err := authentification.GenerateJWT(admin.Email)
	token := tokenString["access_token"]
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminJWT", token, 3600*24*30, "", "", false, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      "true",
		"msg":         "OK",
		"tokenString": tokenString,
	})
}

// @Summary adminhome
// @ID admin-home
// @Description admin home
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 
// @Failure 400 
// @Router /admin/ [get]
func AdminHome(c *gin.Context) {
	
	
	c.JSON(http.StatusAccepted, gin.H{
		"status": "Welcome to admin home page ",
	})
}

func LogoutUser(c *gin.Context) {

	token := c.GetHeader("access_token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authorization Token is required"})
		c.Abort()
		return
	}
	c.SetCookie("AdminJWT", token, 0, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": "logout successfully",
	})
}

type UserDataStruct struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// @Summary admin-userdata
// @ID admin-userdata
// @Description admin userdata
// @Tags Admin User Handler
// @Accept json
// @Produce json
// @Param search query string true "user name"
// @Success 200 
// @Failure 400 
// @Router /admin/userdata/ [get]
func UserData(c *gin.Context) {
	var user UserDataStruct
	search := c.Param("search"); 
	if  search != "" {
		i.Db.Raw("SELECT id,first_name,last_name,email,phone FROM users where first_name like ? ORDER BY id ASC ", search).Scan(&user)
	}
	fmt.Println(search)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// @Summary admin-userblock
// @ID admin-userblock
// @Description admin user block
// @Tags Admin User Handler
// @Accept json
// @Produce json
// @Param id query string true "user id"
// @Success 200 
// @Failure 400 
// @Router /admin/userdata/block/ [put]
func BlockUser(c *gin.Context) {
	params := c.Param("id")
	var user models.User
	i.Db.Raw("UPDATE users SET block_status=true where id=?", params).Scan(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "Blocked successfully"})
}

// @Summary admin-userunblock
// @ID admin-userunblock
// @Description admin user unblock
// @Tags Admin User Handler
// @Accept json
// @Produce json
// @Param id query string true "user id"
// @Success 200 
// @Failure 400 
// @Router /admin/userdata/unblock/ [put]
func UnBlockUser(c *gin.Context) {
	params := c.Param("id")
	var user models.User
	i.Db.Raw("UPDATE users SET block_status=false where id=?", params).Scan(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "Unblocked successfully"})
}

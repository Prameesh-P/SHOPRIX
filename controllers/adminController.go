package controllers

import (
	"github.com/Prameesh-P/SHOPRIX/authentification"
	i "github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminSignup(c *gin.Context) {
	var admin models.Admin
	var count uint
	var Admin struct {
		Email    string
		Password string
	}
	if err := c.ShouldBind(&Admin); err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
	}
	i.Db.Raw("SELECT count(*) FROM admins WHERE email=?", admin.Email).Scan(&count)
	if count > 0 {
		c.JSON(404, gin.H{
			"err": "false",
			"msg": "Admin with same Email already exists",
		})
		c.Abort()
		return
	}
	bytes, err := admin.HashPassword(Admin.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to hash",
		})
		return
	}

	admins := models.Admin{Email: Admin.Email, Password: bytes}
	record := i.Db.Create(&admins)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"error": "failed to create Admin",
		})
	}
	c.JSON(200, gin.H{
		"status": "OK",
	})
}
func AdminLogin(c *gin.Context) {
	var admin struct {
		Email    string
		Password string
	}
	if err := c.ShouldBind(&admin); err != nil {
		c.JSON(404, gin.H{
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
		c.JSON(404, gin.H{
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
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"status":      "true",
		"msg":         "OK",
		"tokenstring": tokenString,
	})
}
func AdminHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "Welcome to admin home page ",
	})
}

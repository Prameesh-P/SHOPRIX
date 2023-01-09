package controllers

import (
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminSignup(c *gin.Context) {
	var admin models.Admin
	var Admin struct {
		Email    string
		Password string
	}
	if err := c.ShouldBind(&Admin); err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
	}
	bytes, err := admin.HashPassword(Admin.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to hash",
		})
		return
	}

	admins := models.Admin{Email: Admin.Email, Password: bytes}
	record := database.Db.Create(&admins)
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
	database.Db.First(&admins, "email=?", admin.Email)
	if admins.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

}

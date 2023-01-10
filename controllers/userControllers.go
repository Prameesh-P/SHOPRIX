package controllers

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		First_name string
		Last_name  string
		Email      string
		Password   string
		Phone      string
	}
	if c.ShouldBind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to hash",
		})
		return
	}
	user := models.User{First_name: body.First_name, Last_name: body.Last_name, Email: body.Email, Password: string(hash), Phone: body.Phone}
	result := database.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"success": "OK",
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if c.ShouldBind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get request",
		})
		return
	}
	var user models.User
	//result := database.Db.Where("email=?", body.Email).First(&user)
	database.Db.First(&user, "email=?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong password",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("Secret")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserJWT", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": "OK",
		"token":   tokenString,
	})
}
func UserHome(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": "Welcome to user home page..!!",
	})
}
func ForgetPassword(c *gin.Context) {
	var user models.User
	useremail := c.GetString("user")
	newpassword := c.PostForm("password")
	database.Db.Raw("select password,id from users where email=?", useremail).Scan(&user)
	// user.Password=newpassword
	if err := user.HashPassword(newpassword); err != nil {
		c.JSON(404, gin.H{"err": err.Error()})
		c.Abort()
		return
	}
	fmt.Println(user.HashPassword(newpassword))

	database.Db.Raw("update users set password=? where id=?", user.Password, user.ID).Scan(&user)
	fmt.Println(useremail)

	fmt.Println(newpassword)
	c.JSON(200, gin.H{
		"pass":  newpassword,
		"email": useremail,
	})

}

func Validate(c *gin.Context) {
	User, err := c.Get("user")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": User,
	})
}

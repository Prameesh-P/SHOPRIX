package controllers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		FirstName   string
		LastName    string
		Email       string
		Password    string
		Phone       string
		BlockStatus bool
	}
	if c.ShouldBind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		c.Abort()
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to hash",
		})
		c.Abort()
		return
	}
	user := models.User{FirstName: body.FirstName, LastName: body.LastName, Email: body.Email, Password: string(hash), Phone: body.Phone, BlockStatus: body.BlockStatus}
	result := database.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
		c.Abort()
		return
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
	if user.BlockStatus {
		c.JSON(404, gin.H{
			"msg": "user has been Blocked by admin",
		})
		c.Abort()
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
func ForgetPasswordEmail(c *gin.Context) {
	// var user models.User
	params:=c.Param("email")
	// var user models.User
	from := "prameepramee0@gmail.com"
	to := []string{params}
	msg := []byte("To:\r"+params+"\nFrom:"+"prameepramee0@gmail.com"+"\nSubject:Forget password request..!!!!\r\n"+"This is Email verification sent from SHOPRIX ECOMMERCE PLATFROM.\r\n"+"Please click on the link\r\n"+``)
	status := SentToEmail(from, to, msg)
	if status {
		c.JSON(http.StatusAccepted, gin.H{
			"Success": "true",
			"msg":     "Verification sent on email successfully",
		})
	} else {
		c.JSON(404, gin.H{
			"error": "something went wrong",
		})
		c.Abort()
		return
	}

}
func SentToEmail(from string, to []string, msg []byte) bool {
	auth := smtp.PlainAuth("", from, os.Getenv("SMT_PASSWORD"), os.Getenv("SMT_HOST"))
	smtpAddress := fmt.Sprintf("%s:%v", os.Getenv("SMT_HOST"), os.Getenv("SMT_PORT"))
	err := smtp.SendMail(smtpAddress, auth, from, to, msg)
	if err != nil {
		return false
	}else{
	return true
	}

}
func ForgetPassword(c *gin.Context) {
	var user models.User
	UserEmail := c.PostForm("user")
	NewPassword := c.PostForm("password")
	database.Db.Raw("select password,id from users where email=?", UserEmail).Scan(&user)
	// user.Password=NewPassword
	if err := user.HashPassword(NewPassword); err != nil {
		c.JSON(404, gin.H{"err": err.Error()})
		c.Abort()
		return
	}
	fmt.Println(user.HashPassword(NewPassword))

	database.Db.Raw("update users set password=? where email=?", user.Password, UserEmail).Scan(&user)
	fmt.Println(UserEmail)
	fmt.Println(NewPassword)
	c.JSON(200, gin.H{
		"pass":  NewPassword,
		"email": UserEmail,
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

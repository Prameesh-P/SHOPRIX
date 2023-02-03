package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var Body struct {
		FirstName   string
		LastName    string
		Email       string
		Password    string
		Phone       string
		BlockStatus bool
	}
	if c.ShouldBindJSON(&Body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read Body",
		})
		c.Abort()
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(Body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "failed to hash",
		})
		c.Abort()
		return
	}
	user := models.User{FirstName: Body.FirstName, LastName: Body.LastName, Email: Body.Email, Password: string(hash), Phone: Body.Phone, BlockStatus: Body.BlockStatus}
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
	var Body struct {
		Email    string
		Password string
	}
	if c.ShouldBind(&Body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get request",
		})
		return
	}
	var user models.User
	database.Db.First(&user, "email=?", Body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}
	if user.BlockStatus {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "user has been Blocked by admin",
		})
		c.Abort()
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Body.Password))
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
	c.JSON(http.StatusOK, gin.H{
		"success": "Welcome to user home page..!!",
	})
}

var OTP struct {
	Number string
}

func OtpGenerator() string {
	min := 1257
	max := 9871
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(max-min) + min
	otps := strconv.Itoa(otp)
	OTP.Number = otps
	return otps
}

func ForgetPasswordEmail(c *gin.Context) {
	otps := OtpGenerator()
	params := c.Param("email")
	from := "prameepramee0@gmail.com"
	to := []string{params}
	msg := []byte("To:" + params + "\r\n" +
		"From:prameepramee0@gmail.com\r\n" +
		"Subject: SHOPRIX verification!\r\n" +
		"\r\n" +
		"<html>This is the email is sent using golang and sendinblue.</html>\r\n" + "<html><h1 style=" + "color:red>" + otps + "</h1></html>")

	status := SentToEmail(from, to, msg)
	if status {
		c.JSON(http.StatusAccepted, gin.H{
			"Success": "true",
			"msg":     "Verification sent on email successfully",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please check your Email..!!",
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
	} else {

		return true

	}

}
func ForgetPassword(c *gin.Context) {
	UserEmail := c.Request.FormValue("useremail")
	var user models.User
	var count uint
	Userotp := c.Request.FormValue("otp")
	UserOtps, _ := strconv.Atoi(Userotp)
	Otp := OTP.Number
	Otps, _ := strconv.Atoi(Otp)
	fmt.Println(Otps)
	fmt.Println(UserOtps)
	NewPassword := c.Request.FormValue("password")
	database.Db.Raw("select count(*) from users where email=?", UserEmail).Scan(&count)
	if count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sorry we cant find no user with this email..!!",
		})
		c.Abort()
		return
	}
	if user.BlockStatus {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Sorry you are blocked by admin",
		})
		c.Abort()
		return
	}
	if err := user.HashPassword(NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		c.Abort()
		return
	}
	fmt.Println(user.HashPassword(NewPassword))
	hash, err := bcrypt.GenerateFromPassword([]byte(NewPassword), 10)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": "failed to hash",
		})
		c.Abort()
		return
	}
	database.Db.Raw("update users set password=? where email=?", hash, UserEmail).Scan(&user)
	fmt.Println(UserEmail)
	fmt.Println(NewPassword)
	if Otps == UserOtps {
		c.JSON(http.StatusOK, gin.H{
			"success":     "true",
			"UserEmail":   UserEmail,
			"NewPassword": NewPassword,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "error",
		})
	}

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
		"status": "true",
		"user":   User,
	})

}

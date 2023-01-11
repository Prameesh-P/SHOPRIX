package controllers

import (
	"fmt"
	"github.com/Prameesh-P/SHOPRIX/authentification"
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/initalizers"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"net/http"
	"os"
)

func init() {
	database.ConnectToDb()
	initalizers.LoadEnvVariables()
}

var (
	accountSid string
	authToken  string
	fromPhone  string
	client     *twilio.RestClient
)

func OtpLog(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.Query("number")
	result := CheckNumber(Mob)
	fmt.Println(result)
	if !result {
		c.JSON(400, gin.H{
			"status": "false",
			"msg":    "Mobile number does not exists..!! Please signup",
		})
		return
	}
	Mobile := "+91" + Mob
	params := &verify.CreateVerificationParams{}
	params.SetTo(Mobile)
	params.SetChannel("sms")
	response, err := client.VerifyV2.CreateVerification(fromPhone, params)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"error": "Error sending otp",
		})
	} else {
		fmt.Printf("Sent Verification '%s'\n", *response.Sid)
		c.JSON(200, gin.H{
			"status": "true",
			"msg":    "OTP sent successfully..!!",
		})
	}
}
func CheckNumber(str string) bool {
	mobileNumber := str
	var checkOtp models.User
	database.Db.Raw("SELECT phone FROM user WHERE phone=?", mobileNumber).Scan(&checkOtp)
	return checkOtp.Phone == mobileNumber
}
func CheckOTP(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.Query("number")
	code := c.Query("otps")
	CheckNumber(Mob)
	var user models.User
	database.Db.First(&user, "phone=?", Mob)
	mobile := "+91" + Mob
	fromPhone = os.Getenv("FROM_PHONE")
	fmt.Println(mobile)
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(mobile)
	params.SetCode(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(fromPhone, params)

	if err != nil {
		fmt.Println(err.Error())
	} else if *resp.Status == "approved" {
		tokenstring, err := authentification.GenerateJWT(user.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create token",
			})

			return
		}
		// Sent it back
		fmt.Println(tokenstring)
		token := tokenstring["access_token"]
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("UserAuth", token, 3600*24*30, "", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "ok",
			"data":    tokenstring,
		})
	} else {

		c.JSON(404, gin.H{
			"msg": "otp is invalid",
		})
	}
}

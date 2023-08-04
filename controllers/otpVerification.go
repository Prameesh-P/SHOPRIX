package controllers

import (
	"fmt"
	// "net/http"
	"os"

	// "net/http"
	// "os"

	// "github.com/Prameesh-P/SHOPRIX/authentification"
	"github.com/Prameesh-P/SHOPRIX/database"
	"github.com/Prameesh-P/SHOPRIX/initalizers"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/gin-gonic/gin"
	twilio "github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	// "github.com/nexmo-community/nexmo-go"
	// "github.com/kyokomi/emoji"
)

func init() {
	initalizers.LoadEnvVariables()
	database.ConnectToDb()
}

var (
	accountSid string
	authToken  string
	serviceSid string
	client     *twilio.RestClient
)
// func NexmoOtpVErification(c *gin.Context){
// 		// Initialize the Nexmo client with your API credentials
// 		apiKey :="264f4a65"
// 	apiSecret := "lHCmXeW8gIq00AeU"

// 	auth := nexmo.NewAuthSet()
// 	auth.SetAPISecret(apiKey, apiSecret)

// 	client := nexmo.NewClient(http.DefaultClient, auth)

// 	phoneNumber := "RECIPIENT_PHONE_NUMBER"

// 	// Generate an OTP
// 	otpRequest, _, err := client.Verify.Request(otpRequest{
// 		Number: phoneNumber,
// 		Brand:  "YOUR_BRAND_NAME",
// 	})
// 	if err != nil {
// 		log.Fatal("Failed to generate OTP:", err)
// 	}

// 	fmt.Println("OTP generated successfully!")

// 	// Verify the OTP
// 	otpCode := "123456" // Replace with the actual OTP code entered by the user

// 	_, err = client.Verify.Check(otpRequest{
// 		RequestId: otpRequest.RequestId,
// 		Code:      otpCode,
// 	})
// 	if err != nil {
// 		log.Fatal("Failed to verify OTP:", err)
// 	}

// 	fmt.Println("OTP verification successful!")
// }
func sendOTP()  {
	// Set your Twilio Account SID and Auth Token
      // Set your Twilio Account SID and Auth Token
	//   TWILIO_ACCOUNT_SID := "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	//   TWILIO_AUTH_TOKEN := "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"
  
	//   // Create a new Twilio REST client
	//   client := twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	  Username: TWILIO_ACCOUNT_SID,
	// 	  Password: TWILIO_AUTH_TOKEN,
	//   })
  
	//   // Create a new Verify service
	//   service := verify.New(client)
  
	//   // Send an SMS OTP
	//   phoneNumber := "+15555555555"
	//   channel := "sms"
	//   code := "123456"
  
	//   response, err := service.CreateSession(phoneNumber, channel, code)
	//   if err != nil {
	// 	  fmt.Println(err)
	// 	  return
	//   }
  
	//   // Check the OTP
	//   status := response.Status
	//   if status == "approved" {
	// 	  fmt.Println("OTP is correct")
	//   } else {
	// 	  fmt.Println("OTP is incorrect")
	//   }
}


func OtpLog(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone := os.Getenv("FROM_PHONE")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.Query("number")

	result := ChekNumber(Mob)
	fmt.Println(result)

	if !result {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Mobile number doesnt exist! Please SignUp",
		})
		return
	}

	// Get Twillio credentials from .env file

	mobile := "+91" + Mob
	// Creatin 4 digit OTP

	params := &verify.CreateVerificationParams{}
	params.SetTo(mobile)
	params.SetChannel("sms")
	resp, err := client.VerifyV2.CreateVerification(fromPhone, params)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"status":  false,
			"message": "error sending OTP",
		})
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
		c.JSON(200, gin.H{
			"status":  true,
			"message": "OTP Sent Succesfully",
		})
	}
}
func ChekNumber(str string) bool {
	mobileNumber := str
	var checkOtp models.User
	database.Db.Raw("SELECT phone FROM users WHERE phone=?", mobileNumber).Scan(&checkOtp)
	return checkOtp.Phone == mobileNumber
}
func CheckOTP(c *gin.Context) {
	// accountSid = os.Getenv("ACCOUNT_SID")
	// authToken = os.Getenv("AUTH_TOKEN")
	// serviceSid = os.Getenv("FROM_PHONE")
	// client = twilio.NewRestClientWithParams(twilio.ClientParams{
	// 	Username: accountSid,
	// 	Password: authToken,
	// })
	// Mob := c.Query("number")
	// code := c.Query("otp")
	// CheckNumber(Mob)
	// var user models.User
	// database.Db.First(&user, "phone=?", Mob)
	// mobile := "+91" + Mob
	// serviceSid = os.Getenv("FROM_PHONE")
	// fmt.Println(mobile)
	// params := &verify.CreateVerificationCheckParams{}
	// params.SetTo(mobile)
	// params.SetCode(code)
	// resp, err := client.VerifyV2.CreateVerificationCheck(serviceSid, params)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else if *resp.Status == "approved" {
	// 	tokenstring, err := authentification.GenerateJWT(user.Email)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"error": "Failed to create token",
	// 		})

	// 		return
	// 	}
	// 	// Sent it back
	// 	fmt.Println(tokenstring)
	// 	token := tokenstring["access_token"]
	// 	c.SetSameSite(http.SameSiteLaxMode)
	// 	c.SetCookie("UserAuth", token, 3600*24*30, "", "", false, true)

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"status":  true,
	// 		"message": "ok",
	// 		"data":    tokenstring,
	// 	})
	// } else {

	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"msg": "otp is invalid",
	// 	})
	// }
}

package main

import (
	"bytes"
	"encoding/json"
	"github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/Prameesh-P/SHOPRIX/models"
	"github.com/Prameesh-P/SHOPRIX/random"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func TestUserHome(t *testing.T) {
	mockResponse := `{"success":"Welcome to user home page..!!"}`
	r := SetUpRouter()
	r.GET("/", controllers.UserHome)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

//func TestForgetPassword(t *testing.T) {
//
//}

type Usesr struct {
	FirstName string
	LastName  string
	Password  string
	Email     string
	Phone     string
}

func TestSignup(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := SetUpRouter()
	r.POST("/signup", controllers.Signup)
	user := models.User{
		FirstName: random.RandomString(6),
		LastName:  random.RandomString(3),
		Password:  "pramee",
		Email:     random.RandomGmailGenerator(9),
		Phone:     strconv.Itoa(int(random.RandomInteger(3000000000, 60000000000))),
	}
	jsonValue, _ := json.Marshal(user)
	reqFound, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)
	assert.Equal(t, 200, w.Code)
}

type User struct {
	Email    string
	Password string
}

//	func TestLogin(t *testing.T) {
//		var users models.User
//		r := SetUpRouter()
//		r.POST("/login", controllers.Login)
//		user := models.User{
//			Email:    "prameesh@gmail.com",
//			Password: "pramee",
//		}
//		var count int
//		database.Db.Raw("select id from users where email=?", user.Email).Scan(&users)
//		database.Db.Raw("select count(*) from users where email=?", user.Email).Scan(&count)
//		database.Db.Raw("select * from users where email=?", user.Email)
//		if users.ID == 0 {
//			jsonValue, _ := json.Marshal(user)
//			reqFound, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
//			w := httptest.NewRecorder()
//			r.ServeHTTP(w, reqFound)
//			assert.Equal(t, 400, w.Code)
//		}
//		jsonValue, _ := json.Marshal(user)
//		reqFound, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
//		w := httptest.NewRecorder()
//		r.ServeHTTP(w, reqFound)
//		assert.Equal(t, 200, w.Code)
//	}
//func TestProductAdding(t *testing.T) {
//	gin.SetMode(gin.ReleaseMode)
//	r := SetUpRouter()
//	r.POST("/admin/addproducts", controllers.ProductAdding)
//
//}

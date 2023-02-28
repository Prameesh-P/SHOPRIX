package main

import (
	"github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

//func TestSignup(t *testing.T) {
//	gin.SetMode(gin.ReleaseMode)
//	r := SetUpRouter()
//	r.POST("/signup", controllers.Signup)
//	user := models.User{
//		FirstName: random.RandomString(6),
//		LastName:  random.RandomString(3),
//		Password:  "pramee",
//		Email:     random.RandomGmailGenerator(9),
//		Phone:     strconv.Itoa(int(random.RandomInteger(3000000000, 60000000000))),
//	}
//	jsonValue, _ := json.Marshal(user)
//	reqFound, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
//	w := httptest.NewRecorder()
//	r.ServeHTTP(w, reqFound)
//	assert.Equal(t, 200, w.Code)
//}

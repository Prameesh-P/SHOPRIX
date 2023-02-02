package controllers

import (
	"github.com/Prameesh-P/SHOPRIX/initalizers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	initalizers.LoadEnvVariables()
}
func TestUserHome(t *testing.T) {

	router := gin.Default()
	//routes.UserRoutes(router)
	//router.Run()
	mockResponse := `{"success":"Welcome to user home page..!!"}`
	router.GET("/", UserHome)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)

}

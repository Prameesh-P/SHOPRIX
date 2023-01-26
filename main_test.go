package main

import (
	"bytes"
	"encoding/json"
	"github.com/Prameesh-P/SHOPRIX/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHome(t *testing.T) {
	routess := gin.Default()
	routess.GET("/", controllers.UserHome)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	routess.ServeHTTP(resp, req)
	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}

func TestUserSignup(t *testing.T) {
	routes := gin.Default()
	routes.POST("/signup", controllers.Signup)
	input := struct {
		FirstName string
		LastName  string
		Email     string
		Password  string
		Phone     string
	}{
		FirstName: "Prameesh",
		LastName:  "P",
		Email:     "pramee@gmail.com",
		Password:  "pramee123",
		Phone:     "6767883734",
	}
	bodyReq, _ := json.Marshal(input)
	routes.POST("/signup", controllers.Signup)
	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(bodyReq))
	resp := httptest.NewRecorder()
	routes.ServeHTTP(resp, req)
	t.Log(resp.Result().StatusCode)
	t.Log(resp.Body)
	assert.Equal(t, 200, resp.Result().StatusCode)
}

package main

import (
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

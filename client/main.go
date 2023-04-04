package main

import(
	"github.com/gin-gonic/gin"
)


func init(){
	LoadEnv()
	ConnectToDb()
	SyncDB()
}


func main(){
	router :=gin.Default()
	UserRoutes(router)
	router.Run(":9000")
}
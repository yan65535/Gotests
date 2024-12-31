package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	engine := gin.Default() //创建gin引擎
	engine.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	err := engine.Run()
	if err != nil {
		return
	} //开启服务器，默认监听localhost:8080
}

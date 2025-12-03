package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zyaaco/wowdle_backend/controllers"
	"github.com/zyaaco/wowdle_backend/models"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/ping", Ping)

	word := router.Group("/word")
	word.GET("", controllers.GetWord)
	word.PUT("", controllers.UpsertWord)
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	SetupRoutes(r)

	r.Run(":6666") // listen and serve on 0.0.0.0:6666
}

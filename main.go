package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lkcsi/goauth/controller"
	"github.com/lkcsi/goauth/service"
)

func main() {
	godotenv.Load()

	server := gin.Default()

	s := service.NewInMemoryUserService()
	c := controller.NewUserController(&s)

	api := server.Group("/users")
	api.GET("/:username", c.FindByUsername)
	api.POST("", c.Save)
	api.POST("/login", c.Login)

	server.Run("localhost:8081")
}

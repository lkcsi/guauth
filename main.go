package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lkcsi/goauth/controller"
	"github.com/lkcsi/goauth/repository"
	"github.com/lkcsi/goauth/service"
)

func main() {
	godotenv.Load()

	server := gin.Default()

	r := repository.SqlUserRepository()
	s := service.NewUserService(&r)
	c := controller.NewUserController(&s)

	api := server.Group("/users")
	api.GET("/:username", c.FindByUsername)
	api.POST("", c.Save)
	api.POST("/login", c.Login)

	server.Run("0.0.0.0:8081")
}

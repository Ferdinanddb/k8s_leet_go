package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"k8s_leet_code/api"
	"k8s_leet_code/controller"
	"k8s_leet_code/database"
	"k8s_leet_code/env"
	"k8s_leet_code/helper"
	"k8s_leet_code/middleware"
	"k8s_leet_code/model"
	"k8s_leet_code/redis"
)

func init() {
	env.LoadEnv()

	// Connect to DB
	database.Connect()
	database.Database.AutoMigrate(&model.User{})
	database.Database.AutoMigrate(&model.Entry{})

	redis.Connect()

}

func main() {
	router := gin.Default()

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)
	publicRoutes.GET("/health", helper.HealthCheck)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())

	protectedRoutes.POST("/run_code", api.PostK8sJob)

	log.Fatal(router.Run(":8080"))
}

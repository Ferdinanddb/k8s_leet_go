package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"k8s_leet_code_backend/api"
	"k8s_leet_code_backend/asynq_client"
	"k8s_leet_code_backend/controller"
	"k8s_leet_code_backend/database"
	"k8s_leet_code_backend/env"
	"k8s_leet_code_backend/helper"
	"k8s_leet_code_backend/middleware"
	"k8s_leet_code_backend/model"
)

func init() {
	env.LoadEnv()

	// Connect to DB
	database.Connect()
	database.Database.AutoMigrate(&model.User{})
	// database.Database.AutoMigrate(&model.Entry{})

	database.Database.AutoMigrate(&model.UserCodeRequest{})

	asynq_client.Connect()

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
	protectedRoutes.GET("/get_history", api.GetUserCodeReqHistory)

	log.Fatal(router.Run(":8080"))
}

package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	jwtConfig()
	redisInit()
}

func main() {
	var r = gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	r.Use(cors.New(config))

	r.POST("/login", Login)
	r.POST("/logout", Logout)
	r.POST("/token/refresh", Refresh)
	r.POST("/update_password", TokenAuthMiddleware(), HandlerUpdatePassword)

	r.GET("/employees", TokenAuthMiddleware(), handleGetEmployees)
	r.POST("/employee", TokenAuthMiddleware(), handleGetEmployee)
	r.POST("/update_employee", TokenAuthMiddleware(), handleUpdateEmployee)
	r.POST("/add_employee", TokenAuthMiddleware(), handleAddEmployee)
	r.POST("/change_employee_status", TokenAuthMiddleware(), handleChangeEmployeeStatus)

	r.POST("/calculate", TokenAuthMiddleware(), handleCalculate)
	r.POST("/report", TokenAuthMiddleware(), handleGenerateReport)
	r.GET("/download_report", TokenAuthMiddleware(), handleDownloadReport)

	r.POST("/mail_to_employee", TokenAuthMiddleware(), handleSendMailToEmployee)
	r.POST("/mail_to_admin", TokenAuthMiddleware(), handleSendMailToAdmin)

	err := r.Run(":3000")
	if err != nil {
		log.Fatal("Unable run the application @ port 3000")
	}
}

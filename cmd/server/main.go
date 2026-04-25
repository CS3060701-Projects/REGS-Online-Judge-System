package main

import (
	"log"
	"net/http"
	"regs-backend/internal/api/handlers"
	"regs-backend/internal/database"
	jwtPkg "regs-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	if err := jwtPkg.InitKeys(); err != nil {
		log.Fatal("JWT 初始化失敗:", err)
	}

	r := gin.Default()

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	handlers.InitJudger(3)

	r.POST("/api/users/register", handlers.Register)
	r.POST("/api/users/login", handlers.Login)
	r.POST("/api/submissions", handlers.SubmitAssignment)
	r.GET("/api/submissions/:operatorId/status", handlers.GetSubmissionStatus)
	r.Run(":8081")
}

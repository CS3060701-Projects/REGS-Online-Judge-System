package main

import (
	"log"
	"net/http"
	"regs-backend/internal/api/handlers"
	"regs-backend/internal/api/middleware"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	jwtPkg "regs-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	err := database.DB.AutoMigrate(
		&models.Submission{},
		&models.User{},
		&models.Problem{},
		&models.JwtBlacklist{},
	)

	if err != nil {
		log.Fatalf("資料庫遷移失敗: %v", err)
	}

	if err := jwtPkg.InitKeys(); err != nil {
		log.Fatal("JWT 初始化失敗:", err)
	}

	r := gin.Default()

	handlers.InitJudger(3) // initialize 3 judge workers

	api := r.Group("/api")
	{
		// Guest
		guest := api.Group("/")
		guest.Use(middleware.AuthMiddleware("Guest"))
		{
			guest.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })

			guest.POST("/users/register", handlers.Register)
			guest.POST("/users/login", handlers.Login)
			guest.GET("/problems", handlers.GetProblems)
			guest.GET("/problems/:id", handlers.GetProblem)
			guest.GET("/users/:user_id/submissions", handlers.GetUserSubmissions)
			guest.GET("/stats/problems/:problem_id", handlers.GetProblemStats)
			guest.GET("/stats/users/:user_id", handlers.GetUserStats)
		}

		// User
		user := api.Group("/")
		user.Use(middleware.AuthMiddleware("User"))
		{
			user.POST("/users/logout", handlers.Logout)
			user.POST("/submissions", handlers.SubmitAssignment)
			user.GET("/submissions", handlers.GetSubmissions)
			user.GET("/submissions/:operatorId", handlers.GetSubmissionStatus)
			user.GET("/submissions/:operatorId/source", handlers.GetSubmissionSource)
			user.GET("/submissions/:operatorId/logs/:type", handlers.GetSubmissionLog)
			user.GET("/users/me", handlers.GetMe)
		}

		// Admin
		admin := api.Group("/")
		admin.Use(middleware.AuthMiddleware("Admin"))
		{
			admin.POST("/problems", handlers.CreateProblem)
			admin.POST("/problems/:id/testdata", handlers.UploadTestData)
			admin.DELETE("/problems/:id", handlers.DeleteProblem)
		}
	}

	r.Run(":8081")
}

package main

import (
	"log"
	"net/http"
	"regs-backend/internal/api/handlers"
	"regs-backend/internal/api/middleware"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	jwtPkg "regs-backend/pkg/jwt"

	_ "regs-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title REGS Online Judge API
// @version 1.0
// @openapi 3.0.0
// @description This is the API server for the REGS Online Judge system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/your-repo

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	database.Connect()

	if err := jwtPkg.InitKeys(); err != nil {
		log.Fatal("JWT 初始化失敗:", err)
	}

	handlers.InitJudger(3) // initialize 3 judge workers

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		// Public routes (no auth required)
		api.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })
		api.POST("/users/register", handlers.Register)
		api.POST("/users/login", handlers.Login)
		api.GET("/problems", handlers.GetProblems)
		api.GET("/problems/:id", handlers.GetProblem)
		api.GET("/problems/:id/examples", handlers.GetProblemExamples)
		api.GET("/users/:user_id/submissions", handlers.GetUserSubmissions)
		api.GET("/stats/problems/:problem_id", handlers.GetProblemStats)
		api.GET("/stats/users/:user_id", handlers.GetUserStats)

		// Authenticated routes
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware("User"))
		{
			// Routes for any authenticated user (User, Admin)
			auth.POST("/users/logout", handlers.Logout)
			auth.POST("/submissions", handlers.SubmitAssignment)
			auth.GET("/submissions", handlers.GetSubmissions)
			auth.GET("/submissions/:operatorId", handlers.GetSubmissionStatus)
			auth.GET("/submissions/:operatorId/source", handlers.GetSubmissionSource)
			auth.GET("/submissions/:operatorId/logs/:type", handlers.GetSubmissionLog)
			auth.GET("/users/me", handlers.GetMe)

			// Admin-only routes
			admin := auth.Group("/")
			admin.Use(middleware.AuthMiddleware("Admin"))
			{
				admin.PUT("/problems", handlers.CreateProblem)
				admin.GET("/problems/:id/testcases", handlers.DownloadTestCases)
				admin.POST("/problems/:id/testdata", handlers.UploadTestData)
				admin.DELETE("/problems/:id", handlers.DeleteProblem)
			}
		}
	}

	r.Run(":8081")
}

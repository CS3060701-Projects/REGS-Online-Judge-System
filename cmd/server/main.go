package main

import (
	"log"
	"net/http"
	"regs-backend/internal/api/handlers"
	"regs-backend/internal/api/middleware"
	"regs-backend/internal/database"
	"regs-backend/internal/judge"
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
// @host localhost:8081
// @BasePath /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	database.Connect()

	if err := jwtPkg.InitKeys(); err != nil {
		log.Fatal("JWT 初始化失敗:", err)
	}

	if err := judge.EnsureBuildNetwork(); err != nil {
		log.Printf("警告: %v", err)
	}

	handlers.InitJudger(3)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.StaticFile("/openapi.yaml", "./docs/openapi.yaml")
	r.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `<!DOCTYPE html>
<html><head><title>REGS API Docs</title>
<meta charset="utf-8"/>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/redoc@2.4.0/bundles/redoc.standalone.css"/>
</head><body>
<redoc spec-url="/openapi.yaml"></redoc>
<script src="https://cdn.jsdelivr.net/npm/redoc@2.4.0/bundles/redoc.standalone.js"></script>
</body></html>`)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		public := api.Group("/")
		public.Use(middleware.OptionalAuthMiddleware())
		{
			public.POST("/users/register", handlers.Register)
			public.POST("/users/login", handlers.Login)
			public.GET("/problems", handlers.GetProblems)
			public.GET("/problems/:id", handlers.GetProblem)
			public.GET("/users/:user_id/submissions", handlers.GetUserSubmissions)
			public.GET("/stats/problems/:problem_id", handlers.GetProblemStats)
			public.GET("/stats/users/:user_id", handlers.GetUserStats)
		}

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware("User"))
		{
			auth.POST("/users/logout", handlers.Logout)
			auth.POST("/submissions", handlers.SubmitAssignment)
			auth.GET("/submissions", handlers.GetSubmissions)
			auth.GET("/submissions/:operatorId", handlers.GetSubmissionStatus)
			auth.GET("/submissions/:operatorId/source", handlers.GetSubmissionSource)
			auth.GET("/users/me", handlers.GetMe)

			admin := auth.Group("/")
			admin.Use(middleware.AuthMiddleware("Admin"))
			{
				admin.PUT("/problems", handlers.CreateProblem)
				admin.GET("/problems/:id/testcases", handlers.DownloadTestCases)
				admin.DELETE("/problems/:id", handlers.DeleteProblem)
			}
		}
	}

	r.Run(":8081")
}

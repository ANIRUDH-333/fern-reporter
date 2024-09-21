package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/guidewire/fern-reporter/config"
	"github.com/guidewire/fern-reporter/pkg/api/handlers"
	"github.com/guidewire/fern-reporter/pkg/auth"
	"github.com/guidewire/fern-reporter/pkg/db"
)

var (
	testRun *gin.RouterGroup
)

func RegisterRouters(router *gin.Engine) {
	handler := handlers.NewHandler(db.GetDb())

	authEnabled := config.GetAuth().Enabled

	var api *gin.RouterGroup
	if authEnabled {
		api = router.Group("/api", auth.ScopeMiddleware())
	} else {
		api = router.Group("/api")
	}

	api.Use()
	{
		testRun = api.Group("/testrun/")
		testRun.GET("/", handler.GetTestRunAll)
		testRun.GET("/:id", handler.GetTestRunByID)
		testRun.POST("/", handler.CreateTestRun)
		testRun.PUT("/:id", handler.UpdateTestRun)
		testRun.DELETE("/:id", handler.DeleteTestRun)

		testReport := api.Group("/reports/testruns")
		testReport.GET("/", handler.ReportTestRunAll)
		testReport.GET("/:id", handler.ReportTestRunById)
	}

	var reports *gin.RouterGroup
	if authEnabled {
		reports = router.Group("/reports/testruns", auth.ScopeMiddleware())
	} else {
		reports = router.Group("/reports/testruns")
	}

	reports.Use()
	{
		reports.GET("/", handler.ReportTestRunAllHTML)
		reports.GET("/:id", handler.ReportTestRunByIdHTML)
	}

	var ping *gin.RouterGroup
	if authEnabled {
		ping = router.Group("/ping", auth.ScopeMiddleware())
	} else {
		ping = router.Group("/ping")
	}

	ping.Use()
	{
		ping.GET("/", handler.Ping)
	}
	insights := router.Group("/insights")
	{
		insights.GET("/:name", handler.ReportTestInsights)
		insights.POST("/gemini", handler.GetGeminiInsights)
	}
}

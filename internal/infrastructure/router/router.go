package router

import (
	"go-server/internal/api/handler"
	"go-server/internal/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func Initialize(h handler.AppHandler) {
	router := gin.Default()
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	router.Use(configSwagger).GET("/swagger/*any", swaggerHandler)

	appVersion1Group := router.Group("/v1")
	{
		interactionGroup := appVersion1Group.Group("interactions")
		{
			interactionGroup.POST("/:video_id", h.InteractionHandler.CreateNewInteraction)
		}
		rankingGroup := appVersion1Group.Group("rankings")
		{
			rankingGroup.GET("", h.ScoreHandler.GetGlobalRanking)
			rankingGroup.GET("/:user_id", h.ScoreHandler.GetPersonalRanking)
		}
	}

	router.Run(":8080")
}

func configSwagger(c *gin.Context) {
	docs.SwaggerInfo.Host = c.Request.Host
}

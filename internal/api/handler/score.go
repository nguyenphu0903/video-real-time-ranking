package handler

import (
	"context"
	"go-server/internal/usecase/score"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScoreHandler interface {
	GetGlobalRanking(c *gin.Context)
	GetPersonalRanking(c *gin.Context)
}

type scoreHandler struct {
	ScoreUseCase score.UseCase
}

func NewScoreHandler(scoreUseCase score.UseCase) ScoreHandler {
	go scoreUseCase.StartEventConsumer(context.Background())
	return &scoreHandler{
		ScoreUseCase: scoreUseCase,
	}
}

// GetGlobalRanking godoc
// @Summary Get global ranking
// @Description Get global ranking
// @Tags rankings
// @Accept json
// @Produce json
// @Router /v1/rankings [get]
// @Param limit query int false "Limit"
// @Success 200 {object} []string
// @Failure 500
// @Failure 400
func (h *scoreHandler) GetGlobalRanking(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.AbortWithStatusJSON(400, "Invalid limit")
		return
	}
	ranking, err := h.ScoreUseCase.ListTopRankedVideos(c, limitInt)
	if err != nil {
		c.AbortWithStatusJSON(500, err.Error())
		return
	}
	c.JSON(200, ranking)
}

// GetPersonalRanking godoc
// @Summary Get personal ranking
// @Description Get personal ranking
// @Tags rankings
// @Accept json
// @Produce json
// @Router /v1/rankings/{user_id} [get]
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit"
// @Success 200 {object} []string
// @Failure 500
// @Failure 400
func (h *scoreHandler) GetPersonalRanking(c *gin.Context) {
	userID := c.Param("user_id")
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.AbortWithStatusJSON(400, "Invalid limit")
		return
	}
	ranking, err := h.ScoreUseCase.ListPersonalTopRankedVideos(c, userID, limitInt)
	if err != nil {
		c.AbortWithStatusJSON(500, err.Error())
		return
	}
	c.JSON(200, ranking)
}

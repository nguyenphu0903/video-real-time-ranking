// interaction handler
package handler

import (
	"go-server/internal/entity"
	"go-server/internal/usecase/interaction"

	"github.com/gin-gonic/gin"
)

// InteractionHandler interface
type InteractionHandler interface {
	CreateNewInteraction(c *gin.Context)
}

type interactionHandler struct {
	InteractionUC interaction.UseCase
}

func NewInteractionHandler(euc interaction.UseCase) InteractionHandler {
	return &interactionHandler{
		InteractionUC: euc,
	}
}

// CreateNewInteraction godoc
// @Summary Create new interaction
// @Description Create new interaction
// @Tags interaction
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Router /v1/interactions [post]
// @Param user_id path string true "User ID"
// @Param post_id path string true "Post ID"
// @Param interaction_type path string true "Interaction Type"
// @Success 200 {object} string
// @Failure 500
func (h *interactionHandler) CreateNewInteraction(c *gin.Context) {
	var req *entity.UserInteractionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(400, err)
		return
	}

	if err := h.InteractionUC.CreateNewInteraction(c, req); err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	c.JSON(200, "Interaction created successfully")
	return
}

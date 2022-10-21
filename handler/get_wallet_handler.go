package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *handler) GetWalletHandler(c *gin.Context) {
	resp, err := h.uc.GetWalletUsecase()
	if err != nil {
		c.JSON(422, gin.H{"message": err})
	}

	c.JSON(200, resp)
}

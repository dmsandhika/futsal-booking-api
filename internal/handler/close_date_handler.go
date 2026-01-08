package handler

import (
	"time"

	"futsal-booking/internal/repository"

	"github.com/gin-gonic/gin"
)

type CloseDateHandler struct {
	Repo *repository.CloseDateRepository
}

func (h *CloseDateHandler) CreateCloseDate(c *gin.Context) {
	var req struct {
		Date   string `json:"date" binding:"required"`
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid date format"})
		return
	}
	closeDate, err := h.Repo.CreateCloseDate(date, req.Reason)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, closeDate)
}
func (h *CloseDateHandler) GetAllCloseDates(c *gin.Context) {
	closeDates, err := h.Repo.GetAllCloseDates()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success", "data": closeDates})
}

func (h *CloseDateHandler) DeleteCloseDate(c *gin.Context) {
	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.DeleteCloseDate(req.ID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "close date deleted"})
}
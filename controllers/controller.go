package controllers

import (
	"api-service/models"
	"api-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{s}
}

// POST /request
func (h *Handler) PostRequest(c *gin.Context) {
	var req models.RequestPayload

	// validate input
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	// apply rate limit + update stats
	err := h.s.HandleRequest(req.UserID)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "request created"})
}

// GET /stats
func (h *Handler) GetStats(c *gin.Context) {
	data, _ := h.s.Stats()
	c.JSON(http.StatusOK, data)
}

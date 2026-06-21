package video

import (
	"errors"
	"net/http"

	"coi/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoHandler struct {
	service VideoService
}

func NewVideoHandler(service VideoService) *VideoHandler {
	return &VideoHandler{service: service}
}

func (h *VideoHandler) CreateUploadSession(c *gin.Context) {
	var req model.CreateVideoUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(int)
	resp, err := h.service.CreateUploadSession(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(statusFromError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *VideoHandler) ConfirmUpload(c *gin.Context) {
	var uri struct {
		ID int `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(int)
	video, err := h.service.ConfirmUpload(c.Request.Context(), userID, uri.ID)
	if err != nil {
		c.JSON(statusFromError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": video})
}

func (h *VideoHandler) HandleUploadCompleted(c *gin.Context) {
	var req struct {
		StorageKey string `json:"storage_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	video, err := h.service.HandleUploadCompleted(c.Request.Context(), req.StorageKey)
	if err != nil {
		c.JSON(statusFromError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": video})
}

func (h *VideoHandler) DeleteVideo(c *gin.Context) {
	var uri struct {
		ID int `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(int)
	if err := h.service.DeleteVideo(c.Request.Context(), userID, uri.ID); err != nil {
		c.JSON(statusFromError(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func statusFromError(err error) int {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound
	}
	if err.Error() == "you do not own this channel" {
		return http.StatusForbidden
	}
	return http.StatusBadRequest
}

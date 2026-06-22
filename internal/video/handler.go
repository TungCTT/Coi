package video

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

func (h *VideoHandler) ListVideos(c *gin.Context) {
	limit := queryInt(c, "limit", 20)
	offset := queryInt(c, "offset", 0)
	query := c.Query("q")

	videos, err := h.service.ListPublicVideos(c.Request.Context(), query, limit, offset)
	if err != nil {
		c.JSON(statusFromError(err), gin.H{"error": err.Error()})
		return
	}
	attachStreamURLs(c, videos)

	c.JSON(http.StatusOK, gin.H{"data": videos})
}

func (h *VideoHandler) GetVideoByID(c *gin.Context) {
	var uri struct {
		ID int `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	video, err := h.service.GetPublicVideo(c.Request.Context(), uri.ID)
	if err != nil {
		c.JSON(statusFromError(err), gin.H{"error": err.Error()})
		return
	}
	attachStreamURL(c, video)

	c.JSON(http.StatusOK, gin.H{"data": video})
}

func (h *VideoHandler) StreamVideo(c *gin.Context) {
	var uri struct {
		ID int `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	video, body, info, err := h.service.OpenPublicVideo(c.Request.Context(), uri.ID, c.GetHeader("Range"))
	if err != nil {
		c.JSON(statusFromError(err), gin.H{"error": err.Error()})
		return
	}
	defer body.Close()

	contentType := info.ContentType
	if contentType == "" {
		contentType = video.ContentType
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", contentType)
	if info.ETag != "" {
		c.Header("ETag", info.ETag)
	}
	if info.Size > 0 {
		c.Header("Content-Length", strconv.FormatInt(info.Size, 10))
	}
	if info.ContentRange != "" {
		c.Header("Content-Range", info.ContentRange)
	}
	if c.GetHeader("Range") != "" {
		c.Status(http.StatusPartialContent)
	} else {
		c.Status(http.StatusOK)
	}

	_, _ = io.Copy(c.Writer, body)
}

func (h *VideoHandler) ListVideosByChannel(c *gin.Context) {
	var uri struct {
		ID int `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := queryInt(c, "limit", 20)
	offset := queryInt(c, "offset", 0)
	videos, err := h.service.ListPublicVideosByChannel(c.Request.Context(), uri.ID, limit, offset)
	if err != nil {
		c.JSON(statusFromError(err), gin.H{"error": err.Error()})
		return
	}
	attachStreamURLs(c, videos)

	c.JSON(http.StatusOK, gin.H{"data": videos})
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

func queryInt(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func attachStreamURLs(c *gin.Context, videos []model.Video) {
	for i := range videos {
		attachStreamURL(c, &videos[i])
	}
}

func attachStreamURL(c *gin.Context, video *model.Video) {
	video.VideoURL = fmt.Sprintf("%s/api/v1/videos/%d/stream", requestBaseURL(c), video.ID)
}

func requestBaseURL(c *gin.Context) string {
	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	return scheme + "://" + c.Request.Host
}

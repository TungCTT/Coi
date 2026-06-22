package user

import (
	"net/http"

	"coi/internal/model"

	"github.com/gin-gonic/gin"
)

// AuthHandler nhận HTTP request, gọi service, và trả HTTP response.
type AuthHandler struct {
	service AuthService
}

// NewAuthHandler là constructor của handler.
func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register godoc
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.CreateUserRequest

	// c.ShouldBindJSON() thực hiện 2 việc:
	//  1. Deserialize JSON body thành struct
	//  2. Validate theo các tag `binding:"..."` đã khai báo trong model
	// Nếu lỗi (thiếu field, sai format email, password < 6 ký tự), tự trả 400.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.service.Register(&req)
	if err != nil {
		// Lỗi từ service (email trùng, username trùng, lỗi DB...) → 400 hoặc 500
		// Ở đây đơn giản hóa: trả 400 cho mọi lỗi business logic
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 201 Created — đúng HTTP semantics cho việc tạo resource mới
	c.JSON(http.StatusCreated, gin.H{
		"message": "Đăng ký thành công",
		"data":    resp,
	})
}

// Login godoc
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		// 401 Unauthorized — đúng cho lỗi xác thực sai credentials
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đăng nhập thành công",
		"data":    resp,
	})
}

// RefreshToken godoc
// POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Nhận refresh_token từ JSON body
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token là bắt buộc"})
		return
	}

	tokenPair, err := h.service.RefreshToken(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cấp token mới thành công",
		"data":    tokenPair,
	})
}

// Me godoc
// GET /api/v1/auth/me
// Route này được bảo vệ bởi AuthMiddleware — chỉ đến được đây nếu token hợp lệ.
func (h *AuthHandler) Me(c *gin.Context) {
	// Middleware đã decode JWT và set "userID" vào context.
	// c.MustGet() panic nếu không có key — dùng được vì middleware đảm bảo luôn có.
	userID := c.MustGet("userID").(int)

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (h *AuthHandler) CreateMediaUploadSession(c *gin.Context) {
	var req model.CreateUserMediaUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(int)
	resp, err := h.service.CreateMediaUploadSession(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

func (h *AuthHandler) ConfirmMediaUpload(c *gin.Context) {
	var req model.ConfirmUserMediaUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(int)
	resp, err := h.service.ConfirmMediaUpload(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"invest-mate/internal/users/models/domain"
	"invest-mate/internal/users/services"
	"invest-mate/pkg/handlers"
)

type UserHandler struct {
	userService services.UserService
}

// Создание нового хендлера
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Регистрация маршрутов
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("/register", h.Register)
		users.POST("/login", h.Login)

		// TODO: Защищенные маршруты (потребуют middleware авторизации)
		users.GET("/profile", h.GetProfile)
		users.PUT("/profile", h.UpdateProfile)
		users.DELETE("/profile", h.DeleteUser)
	}

	// Admin routes
	admin := router.Group("/admin/users")
	{
		admin.GET("/", handlers.HandleListRequest(h.userService.GetListUsers))
		admin.GET("/:id", h.GetUserByID)
		admin.PUT("/:id", h.UpdateUser)
		admin.DELETE("/:id", h.DeleteUser)
	}
}

// Обработчик регистрации нового пользователя
func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userResponse, err := h.userService.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		errorMsg := err.Error()

		if err.Error() == "email already exists" ||
			err.Error() == "username already taken" ||
			err.Error() == "password must be at least 8 characters" {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"error": errorMsg})
		return
	}

	response := handlers.BuildResponse(userResponse)
	c.JSON(http.StatusOK, response)
}

// Обработчик авторизации
func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userResponse, err := h.userService.LoginUser(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "invalid credentials" ||
			err.Error() == "account is deactivated" {
			status = http.StatusUnauthorized
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// TODO: Генерация JWT токена
	// token, err := generateJWT(userResponse.ID)

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Login successful",
	// 	"user":    userResponse,
	// 	// "token":   token,
	// })

	response := handlers.BuildResponse(userResponse)
	c.JSON(http.StatusOK, response)
}

// Обработчик получения профиля
func (h *UserHandler) GetProfile(c *gin.Context) {
	// TODO: Получать ID пользователя из JWT токена
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := handlers.BuildResponse(userResponse)
	c.JSON(http.StatusOK, response)
}

// Обработчик изменения профиля (пользователя)
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var updates domain.User
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userResponse, err := h.userService.UpdateUser(c.Request.Context(), userID, &updates)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "email already exists" ||
			err.Error() == "username already taken" {
			status = http.StatusBadRequest
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := handlers.BuildResponse(userResponse)
	c.JSON(http.StatusOK, response)
}

// Обработчик получения профиля по идентификатору
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	response := handlers.BuildResponse(userResponse)
	c.JSON(http.StatusOK, response)
}

// Обработчик изменения пользователя (админом)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var updates domain.User
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userResponse, err := h.userService.UpdateUser(c.Request.Context(), userID, &updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := handlers.BuildResponse(userResponse)
	c.JSON(http.StatusOK, response)
}

// TODO: Сделать защиту, чтобы пользователь мог удалить только свой аккаунт, а админ - любой
// Обработчик удаления пользователя
func (h *UserHandler) DeleteUser(c *gin.Context) {
	var req domain.DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.userService.DeleteUser(c.Request.Context(), req.ID)
	if err != nil {
		status := http.StatusInternalServerError
		errorMsg := err.Error()
		c.JSON(status, gin.H{"error": errorMsg})
		return
	}

	response := handlers.BuildResponse(result)
	c.JSON(http.StatusOK, response)
}

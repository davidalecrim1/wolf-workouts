package handler

import (
	"errors"
	"net/http"

	"github.com/davidalecrim1/wolf-workouts/internal/users/app"
	"github.com/davidalecrim1/wolf-workouts/internal/users/handler/generated"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *app.UserService
}

func NewUserHandler(svc *app.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/users", h.CreateUser)
	router.POST("/users/login", h.LoginUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var userRequest generated.CreateUserRequest
	if err := c.ShouldBind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := app.NewUser(userRequest.Name, userRequest.Email, userRequest.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.svc.CreateUser(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req generated.LoginUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.LoginUser(c.Request.Context(), req.Email, req.Password)
	if errors.Is(err, app.ErrInvalidEmailOrPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

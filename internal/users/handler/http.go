package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/davidalecrim1/wolf-workouts/internal/users/app"
	gen "github.com/davidalecrim1/wolf-workouts/internal/users/handler/generated"
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
	var userRequest gen.CreateUserRequest
	if err := c.ShouldBind(&userRequest); err != nil {
		slog.Error("failed to create user", "error", err)
		c.JSON(http.StatusBadRequest, gen.ResponseError{Message: err.Error()})
		return
	}

	u, err := app.NewUser(userRequest.Name, userRequest.Email, userRequest.Password)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		c.JSON(http.StatusBadRequest, gen.ResponseError{Message: err.Error()})
		return
	}

	err = h.svc.CreateUser(c.Request.Context(), u)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gen.ResponseError{Message: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req gen.LoginUserRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Error("failed to login user", "error", err)
		c.JSON(http.StatusBadRequest, gen.ResponseError{Message: err.Error()})
		return
	}

	token, err := h.svc.LoginUser(c.Request.Context(), req.Email, req.Password)
	if errors.Is(err, app.ErrInvalidEmailOrPassword) {
		slog.Debug("invalid email or password", "email", req.Email, "error", err)
		c.JSON(http.StatusUnauthorized, gen.ResponseError{Message: err.Error()})
		return
	}

	if err != nil {
		slog.Error("failed to login user", "error", err)
		c.JSON(http.StatusInternalServerError, gen.ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gen.LoginUserResponse{AccessToken: &token})
}

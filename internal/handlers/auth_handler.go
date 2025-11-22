package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ronak4195/personal-assistant/internal/middleware"
	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/services"
)

type AuthHandler struct {
	svc services.AuthService
}

func NewAuthHandler(svc services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(c echo.Context) error {
	var req signupRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid payload")
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)

	if req.Email == "" || req.Password == "" || len(req.Password) < 6 {
		return respondError(c, http.StatusBadRequest, "invalid email or password")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	user, token, err := h.svc.Signup(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			return respondError(c, http.StatusConflict, "email already exists")
		}
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	resp := map[string]any{
		"user": map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"token": token,
	}
	return c.JSON(http.StatusCreated, models.SingleResponse[map[string]any]{Data: resp})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid payload")
	}

	if req.Email == "" || req.Password == "" {
		return respondError(c, http.StatusBadRequest, "email and password required")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	user, token, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, "invalid credentials")
	}

	resp := map[string]any{
		"user": map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"token": token,
	}
	return c.JSON(http.StatusOK, models.SingleResponse[map[string]any]{Data: resp})
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return respondError(c, http.StatusUnauthorized, "unauthorized")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	user, err := h.svc.GetUser(ctx, userID)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return respondError(c, http.StatusNotFound, "user not found")
	}

	resp := map[string]any{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	}
	return c.JSON(http.StatusOK, models.SingleResponse[map[string]any]{Data: resp})
}

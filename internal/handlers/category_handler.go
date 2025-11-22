package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ronak4195/personal-assistant/internal/middleware"
	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/services"
)

type CategoryHandler struct {
	svc services.CategoryService
}

func NewCategoryHandler(svc services.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

type categoryRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parentId"`
}

func (h *CategoryHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	var req categoryRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Name == "" {
		return respondError(c, http.StatusBadRequest, "name is required")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	cat, err := h.svc.Create(ctx, userID, req.Name, req.ParentID)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	resp := map[string]any{
		"id":        cat.ID,
		"name":      cat.Name,
		"parentId":  cat.ParentID,
		"createdAt": cat.CreatedAt,
	}
	return c.JSON(http.StatusCreated, models.SingleResponse[map[string]any]{Data: resp})
}

func (h *CategoryHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)
	parentID := c.QueryParam("parentId")
	var pid *string
	if parentID != "" {
		pid = &parentID
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	cats, err := h.svc.List(ctx, userID, pid)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	resp := make([]map[string]any, 0, len(cats))
	for _, cat := range cats {
		resp = append(resp, map[string]any{
			"id":       cat.ID,
			"name":     cat.Name,
			"parentId": cat.ParentID,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{"data": resp})
}

func (h *CategoryHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	cat, err := h.svc.Get(ctx, userID, id)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	if cat == nil {
		return respondError(c, http.StatusNotFound, "category not found")
	}

	resp := map[string]any{
		"id":       cat.ID,
		"name":     cat.Name,
		"parentId": cat.ParentID,
	}
	return c.JSON(http.StatusOK, models.SingleResponse[map[string]any]{Data: resp})
}

func (h *CategoryHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")
	var req categoryRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Name == "" {
		return respondError(c, http.StatusBadRequest, "name is required")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	cat, err := h.svc.Update(ctx, userID, id, req.Name, req.ParentID)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	resp := map[string]any{
		"id":       cat.ID,
		"name":     cat.Name,
		"parentId": cat.ParentID,
	}
	return c.JSON(http.StatusOK, models.SingleResponse[map[string]any]{Data: resp})
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	if err := h.svc.Delete(ctx, userID, id); err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

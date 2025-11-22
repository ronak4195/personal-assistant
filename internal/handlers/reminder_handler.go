package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ronak4195/personal-assistant/internal/middleware"
	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/repositories"
	"github.com/ronak4195/personal-assistant/internal/services"
)

type ReminderHandler struct {
	svc services.ReminderService
}

func NewReminderHandler(svc services.ReminderService) *ReminderHandler {
	return &ReminderHandler{svc: svc}
}

type reminderRequest struct {
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	DueAt          string  `json:"dueAt"`
	RepeatInterval string  `json:"repeatInterval"`
	IsActive       *bool   `json:"isActive"`
}

func (h *ReminderHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	var req reminderRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Title == "" {
		return respondError(c, http.StatusBadRequest, "title is required")
	}
	dueAt, err := time.Parse(time.RFC3339, req.DueAt)
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid dueAt")
	}

	r := &models.Reminder{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueAt:       dueAt,
	}
	if req.RepeatInterval != "" {
		r.RepeatInterval = models.RepeatInterval(req.RepeatInterval)
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	created, err := h.svc.Create(ctx, r)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, models.SingleResponse[*models.Reminder]{Data: created})
}

func (h *ReminderHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)

	isActiveParam := c.QueryParam("isActive")
	var isActive *bool
	if isActiveParam != "" {
		val, err := strconv.ParseBool(isActiveParam)
		if err != nil {
			return respondError(c, http.StatusBadRequest, "invalid isActive")
		}
		isActive = &val
	}
	from, err := parseTimeParam(c, "from")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid from")
	}
	to, err := parseTimeParam(c, "to")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid to")
	}

	filter := repositories.ReminderFilter{
		UserID:   userID,
		IsActive: isActive,
		From:     from,
		To:       to,
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	rems, err := h.svc.List(ctx, filter)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"data": rems})
}

func (h *ReminderHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	rem, err := h.svc.Get(ctx, userID, id)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	if rem == nil {
		return respondError(c, http.StatusNotFound, "reminder not found")
	}
	return c.JSON(http.StatusOK, models.SingleResponse[*models.Reminder]{Data: rem})
}

func (h *ReminderHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req reminderRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid payload")
	}

	r := &models.Reminder{
		ID:     id,
		UserID: userID,
	}

	if req.Title != "" {
		r.Title = req.Title
	}
	if req.Description != nil {
		r.Description = req.Description
	}
	if req.DueAt != "" {
		d, err := time.Parse(time.RFC3339, req.DueAt)
		if err != nil {
			return respondError(c, http.StatusBadRequest, "invalid dueAt")
		}
		r.DueAt = d
	}
	if req.RepeatInterval != "" {
		r.RepeatInterval = models.RepeatInterval(req.RepeatInterval)
	}
	if req.IsActive != nil {
		r.IsActive = *req.IsActive
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	updated, err := h.svc.Update(ctx, userID, r)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, models.SingleResponse[*models.Reminder]{Data: updated})
}

func (h *ReminderHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	if err := h.svc.Delete(ctx, userID, id); err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

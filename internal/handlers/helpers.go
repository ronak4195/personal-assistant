package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ronak4195/personal-assistant/internal/models"
)

func respondError(c echo.Context, status int, msg string, code ...string) error {
	var ccode *string
	if len(code) > 0 {
		ccode = &code[0]
	}
	return c.JSON(status, models.ErrorResponse{
		Error: models.ErrorBody{
			Message: msg,
			Code:    ccode,
		},
	})
}

func parsePagination(c echo.Context, defaultLimit int64) (limit, offset int64) {
	limit = defaultLimit
	offset = 0

	if lStr := c.QueryParam("limit"); lStr != "" {
		if v, err := strconv.ParseInt(lStr, 10, 64); err == nil && v > 0 {
			limit = v
		}
	}
	if oStr := c.QueryParam("offset"); oStr != "" {
		if v, err := strconv.ParseInt(oStr, 10, 64); err == nil && v >= 0 {
			offset = v
		}
	}
	return
}

func parseTimeParam(c echo.Context, key string) (*time.Time, error) {
	val := c.QueryParam(key)
	if val == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseBodyTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, s)
}

func ok(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

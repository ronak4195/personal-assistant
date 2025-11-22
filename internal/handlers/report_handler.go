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

type ReportHandler struct {
	svc services.ReportService
}

func NewReportHandler(svc services.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) Summary(c echo.Context) error {
	userID := middleware.GetUserID(c)

	periodStr := c.QueryParam("period")
	if periodStr == "" {
		periodStr = string(services.PeriodMonthly)
	}
	period := services.SummaryPeriod(periodStr)

	groupByStr := c.QueryParam("groupBy")
	if groupByStr == "" {
		groupByStr = string(services.GroupNone)
	}
	groupBy := services.GroupBy(groupByStr)

	var start, end *time.Time
	if period == services.PeriodCustom {
		sStr := c.QueryParam("start")
		eStr := c.QueryParam("end")
		if sStr == "" || eStr == "" {
			return respondError(c, http.StatusBadRequest, "start and end are required for custom period")
		}
		s, err := time.Parse("2006-01-02", sStr)
		if err != nil {
			return respondError(c, http.StatusBadRequest, "invalid start date")
		}
		e, err := time.Parse("2006-01-02", eStr)
		if err != nil {
			return respondError(c, http.StatusBadRequest, "invalid end date")
		}
		start = &s
		end = &e
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 15*time.Second)
	defer cancel()

	report, err := h.svc.GetSummary(ctx, userID, period, start, end, groupBy)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, models.SingleResponse[*services.SummaryReport]{Data: report})
}

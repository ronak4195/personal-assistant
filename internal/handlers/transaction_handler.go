package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ronak4195/personal-assistant/internal/middleware"
	"github.com/ronak4195/personal-assistant/internal/models"
	"github.com/ronak4195/personal-assistant/internal/repositories"
	"github.com/ronak4195/personal-assistant/internal/services"
)

type TransactionHandler struct {
	svc services.TransactionService
}

func NewTransactionHandler(svc services.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

type transactionCreateRequest struct {
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	CategoryID    *string `json:"categoryId"`
	SubcategoryID *string `json:"subcategoryId"`
	Description   *string `json:"description"`
	Date          string  `json:"date"`
}

func (h *TransactionHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	var req transactionCreateRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Amount <= 0 || req.Currency == "" || req.Type == "" {
		return respondError(c, http.StatusBadRequest, "type, amount and currency are required")
	}

	var date time.Time
	if req.Date != "" {
		var err error
		date, err = parseBodyTime(req.Date)
		if err != nil {
			return respondError(c, http.StatusBadRequest, "invalid date format")
		}
	}

	tx := &models.Transaction{
		UserID:        userID,
		Type:          models.TransactionType(req.Type),
		Amount:        req.Amount,
		Currency:      req.Currency,
		CategoryID:    req.CategoryID,
		SubcategoryID: req.SubcategoryID,
		Description:   req.Description,
		Date:          date,
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	created, err := h.svc.Create(ctx, tx)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, models.SingleResponse[*models.Transaction]{Data: created})
}

func (h *TransactionHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)

	limit, offset := parsePagination(c, 20)

	var ttype *models.TransactionType
	if tStr := c.QueryParam("type"); tStr != "" {
		tt := models.TransactionType(tStr)
		ttype = &tt
	}

	from, err := parseTimeParam(c, "from")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid from param")
	}
	to, err := parseTimeParam(c, "to")
	if err != nil {
		return respondError(c, http.StatusBadRequest, "invalid to param")
	}

	catID := c.QueryParam("categoryId")
	var catPtr *string
	if catID != "" {
		catPtr = &catID
	}
	subCatID := c.QueryParam("subcategoryId")
	var subCatPtr *string
	if subCatID != "" {
		subCatPtr = &subCatID
	}

	sortParam := c.QueryParam("sort")
	sortAsc := sortParam == "date_asc"

	filter := repositories.TransactionFilter{
		UserID:        userID,
		Type:          ttype,
		From:          from,
		To:            to,
		CategoryID:    catPtr,
		SubcategoryID: subCatPtr,
		Limit:         limit,
		Offset:        offset,
		SortDateAsc:   sortAsc,
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	txs, total, err := h.svc.List(ctx, filter)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}

	pagination := models.Pagination{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
	return c.JSON(http.StatusOK, models.ListResponse[models.Transaction]{
		Data:       txs,
		Pagination: pagination,
	})
}

func (h *TransactionHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	tx, err := h.svc.Get(ctx, userID, id)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	if tx == nil {
		return respondError(c, http.StatusNotFound, "transaction not found")
	}
	return c.JSON(http.StatusOK, models.SingleResponse[*models.Transaction]{Data: tx})
}

type transactionUpdateRequest struct {
	Type          *string  `json:"type"`
	Amount        *float64 `json:"amount"`
	Currency      *string  `json:"currency"`
	CategoryID    *string  `json:"categoryId"`
	SubcategoryID *string  `json:"subcategoryId"`
	Description   *string  `json:"description"`
	Date          *string  `json:"date"`
}

func (h *TransactionHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	var req transactionUpdateRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, "invalid payload")
	}

	tx := &models.Transaction{
		ID:     id,
		UserID: userID,
	}

	if req.Type != nil {
		tx.Type = models.TransactionType(*req.Type)
	}
	if req.Amount != nil {
		tx.Amount = *req.Amount
	}
	if req.Currency != nil {
		tx.Currency = *req.Currency
	}
	if req.CategoryID != nil {
		tx.CategoryID = req.CategoryID
	}
	if req.SubcategoryID != nil {
		tx.SubcategoryID = req.SubcategoryID
	}
	if req.Description != nil {
		tx.Description = req.Description
	}
	if req.Date != nil && *req.Date != "" {
		d, err := time.Parse(time.RFC3339, *req.Date)
		if err != nil {
			return respondError(c, http.StatusBadRequest, "invalid date format")
		}
		tx.Date = d
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	updated, err := h.svc.Update(ctx, userID, tx)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, models.SingleResponse[*models.Transaction]{Data: updated})
}

func (h *TransactionHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	if err := h.svc.Delete(ctx, userID, id); err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

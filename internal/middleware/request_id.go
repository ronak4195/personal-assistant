package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const HeaderRequestID = "X-Request-ID"

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header.Get(HeaderRequestID)
			if reqID == "" {
				reqID = uuid.NewString()
			}
			c.Response().Header().Set(HeaderRequestID, reqID)
			return next(c)
		}
	}
}

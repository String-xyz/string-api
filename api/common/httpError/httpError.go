package httpError

import (
	"net/http"

	validator "github.com/String-xyz/string-api/api/validator"
	"github.com/labstack/echo/v4"
)

type JSONError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Details any    `json:"details"`
}

func InvalidPayloadError(c echo.Context, err error) error {
	errorParams := validator.ExtractErrorParams(err)
	return c.JSON(http.StatusBadRequest, JSONError{Message: "Invalid Payload", Code: "INVALID_PAYLOAD", Details: errorParams})
}

func InternalError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, JSONError{Message: "Something went wrong", Code: "INTERNAL_SERVER"})
}

func BadRequestError(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, JSONError{Message: "Bad Request", Code: "BAD_REQUEST"})
}

func NotFoundError(c echo.Context) error {
	return c.JSON(http.StatusNotFound, JSONError{Message: "Resource Not Found", Code: "NOT_FOUND"})
}

func NotAllowedError(c echo.Context) error {
	return c.JSON(http.StatusMethodNotAllowed, JSONError{Message: "Not Allowed", Code: "NOT_ALLOWED"})
}

package middlewares

import (
	error2 "github.com/aliarbak/ethereum-connector/errors"
	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, ctx echo.Context) {
	errorResponse, ok := err.(*error2.Error)
	if ok {
		ctx.Logger().Error(err)
		ctx.JSON(errorResponse.StatusCode, errorResponse)
		return
	}

	ctx.JSON(500, error2.Error{
		StatusCode: 500,
		Message:    err.Error(),
	})
	return
}

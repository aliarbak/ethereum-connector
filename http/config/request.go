package httpconfig

import (
	"encoding/json"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ReadValidatedBody(ctx echo.Context, data interface{}) error {
	err := json.NewDecoder(ctx.Request().Body).Decode(&data)
	if err != nil {
		return errors.FromInvalidInput(err)
	}

	validate := validator.New()
	err = validate.Struct(data)
	if err != nil {
		return errors.FromInvalidInput(err)
	}

	return nil
}

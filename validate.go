package gox

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var v = validator.New()

type validationErrorField struct {
	Field      string `json:"field"`
	Tag        string `json:"tag"`
	Constraint string `json:"constraint"`
}

type validationError struct {
	Errors []validationErrorField `json:"errors"`
}

func Validate(to any, c echo.Context) bool {
	if err := c.Bind(to); err != nil {
		c.JSON(400, validationError{
			Errors: []validationErrorField{},
		})

		return false
	}

	err := v.Struct(to)
	if err != nil {
		fieldErrors := []validationErrorField{}

		for _, err := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, validationErrorField{
				Field:      err.Field(),
				Tag:        err.Tag(),
				Constraint: err.Param(),
			})
		}

		c.JSON(400, validationError{
			Errors: fieldErrors,
		})

		return false
	}

	return true
}

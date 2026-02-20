package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ParseValidatorError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, e := range ve {
			switch e.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", e.Field())
			case "email":
				return fmt.Sprintf("%s is not valid email", e.Field())
			case "min":
				return fmt.Sprintf("%s must be at lease %s characters", e.Field(), e.Param())
			}
		}
	}

	return err.Error()
}
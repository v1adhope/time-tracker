package v1

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("v1: registerCustomValidations: engine not found")
	}

	if err := v.RegisterValidation("filterstring", filterString); err != nil {
		return fmt.Errorf("v1: registerCustomValidations: registerValidation: %w", err)
	}

	return nil
}

func filterString(fl validator.FieldLevel) bool {
	parts := strings.Split(fl.Field().String(), ":")

	if len(parts) != 2 {
		return false
	}

	if parts[0] != "eq" && parts[0] != "ilike" {
		return false
	}

	return true
}

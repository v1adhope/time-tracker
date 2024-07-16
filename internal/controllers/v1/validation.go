package v1

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("v1: registerCustomValidations: engine not found")
	}

	if err := v.RegisterValidation("filterstring", filterString); err != nil {
		return fmt.Errorf("v1: registerCustomValidations: registerValidation: filterString: %w", err)
	}

	if err := v.RegisterValidation("sorttime", sortTime); err != nil {
		return fmt.Errorf("v1: registerCustomValidations: registerValidation: sortTime: %w", err)
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

func sortTime(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	if err != nil {
		return false
	}

	return true
}

package validator

import (
	"context"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"os"
)

func ValidatorLifecycleHook() fx.Hook {
	return fx.Hook{
		OnStart: func(context.Context) error {
			if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
				v.RegisterValidation("notempty", validateNotEmpty)
				v.RegisterValidation("requiredIfEmpty", validateConditionalRequired)
			}
			return nil
		},
		OnStop: func(context.Context) error {
			os.Exit(1)
			return nil
		},
	}
}

func validateConditionalRequired(fl validator.FieldLevel) bool {
	field := fl.Field()
	paramValue := fl.Parent().FieldByName(fl.Param())

	if paramValue.IsZero() && !field.IsZero() {
		return true
	} else if !paramValue.IsZero() && field.IsZero() {
		return true
	}

	return false
}

func validateNotEmpty(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value != ""
}

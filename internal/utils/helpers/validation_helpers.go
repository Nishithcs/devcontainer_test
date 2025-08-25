package helpers

import (
	"clusterix-code/internal/utils/errors"
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidate[T any](c *gin.Context, obj T) *errors.AppError {
	validationErrors := make(map[string][]string)

	if err := c.ShouldBindJSON(&obj); err != nil {
		if err.Error() == "EOF" {
			return errors.NewValidationError("Validation error", map[string][]string{
				"general": {"Request body is empty"},
			})
		}

		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				field := e.Field()
				message := fmt.Sprintf("The %s field is invalid: %s", field, e.Tag())
				validationErrors[field] = append(validationErrors[field], message)
			}
		} else {
			validationErrors["general"] = append(validationErrors["general"], err.Error())
		}
	}

	if validator, ok := any(obj).(interface{ Validate() map[string][]string }); ok {
		if errs := validator.Validate(); len(errs) > 0 {
			for field, messages := range errs {
				validationErrors[field] = append(validationErrors[field], messages...)
			}
		}
	}

	if len(validationErrors) > 0 {
		return errors.NewValidationError("Validation error", validationErrors)
	}

	return nil
}

// ValidateStruct dynamically validates any struct and returns validation errors.
func ValidateStruct(obj interface{}) map[string][]string {
	validate := validator.New()
	err := validate.Struct(obj)
	validationErrors := make(map[string][]string)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			param := err.Param()

			// Customize error messages based on the validation tag
			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s is a required field", field)
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", field, param)
			default:
				message = fmt.Sprintf("%s is invalid", field)
			}

			// Convert field name to snake_case
			fieldName := snakeCase(field)
			validationErrors[fieldName] = append(validationErrors[fieldName], message)
		}
	}

	// Handle nested structs and slices recursively
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		fieldName := fieldType.Name

		// Handle nested structs
		if field.Kind() == reflect.Struct || (field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct) {
			nestedErrors := ValidateStruct(field.Interface())
			for nestedField, messages := range nestedErrors {
				validationErrors[snakeCase(fieldName+"."+nestedField)] = messages
			}
		}

		// Handle slices/arrays of structs
		if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.Struct || (elem.Kind() == reflect.Ptr && elem.Elem().Kind() == reflect.Struct) {
					nestedErrors := ValidateStruct(elem.Interface())
					for nestedField, messages := range nestedErrors {
						validationErrors[snakeCase(fmt.Sprintf("%s.%d.%s", fieldName, j, nestedField))] = messages
					}
				}
			}
		}
	}

	return validationErrors
}

func snakeCase(s string) string {
	var res string
	for i, r := range s {
		if 'A' <= r && r <= 'Z' {
			if i > 0 {
				res += "_"
			}
			res += string(r + 32)
		} else {
			res += string(r)
		}
	}
	return res
}

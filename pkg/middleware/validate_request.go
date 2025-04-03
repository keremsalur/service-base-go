package middleware

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

var requestMap = map[string]interface{}{
	// request structs
}

func DynamicDTOValidationMiddleware(c *fiber.Ctx) error {

	key := fmt.Sprintf("%s %s", c.Method(), c.Route().Path)
	dto, exists := requestMap[key]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No DTO configuration found for this route",
		})
	}
	return validateDynamicDTO(c, dto)
}

func validateDynamicDTO(c *fiber.Ctx, req any) error {

	// Yeni bir DTO nesnesi oluştur
	dtoType := reflect.TypeOf(req)
	if dtoType.Kind() == reflect.Ptr {
		dtoType = dtoType.Elem()
	}
	dtoValue := reflect.New(dtoType).Interface()

	// İstek türüne göre parse et
	if c.Method() == fiber.MethodGet {
		// GET isteklerinde sadece parametreleri DTO'ya eşle
		for i := 0; i < dtoType.NumField(); i++ {
			field := dtoType.Field(i)
			paramName := field.Tag.Get("params")
			if paramName != "" {
				// Parametreyi al ve DTO'ya eşle
				paramValue := c.Params(paramName)
				if paramValue != "" {
					reflect.ValueOf(dtoValue).Elem().FieldByName(field.Name).SetString(paramValue)
				}
			}

			// Query parametrelerini al
			queryName := field.Tag.Get("query")
			if queryName != "" {
				// Query parametreyi al
				queryValue := c.Query(queryName)
				if queryValue != "" {
					fieldValue := reflect.ValueOf(dtoValue).Elem().FieldByName(field.Name)
					if err := setFieldValue(fieldValue, queryValue, field.Type.Kind()); err != nil {
						fmt.Printf("Error setting field value: %v\n", err)
					}
				}
			}
		}
	} else {
		// Diğer isteklerde body parse işlemi
		if err := c.BodyParser(dtoValue); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input format",
			})
		}
	}

	// Validasyonu çalıştır
	if err := validate.Struct(dtoValue); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}

	// DTO'yu context'e ekle
	c.Locals("validatedDTO", dtoValue)
	return c.Next()
}

// Genel dönüştürme fonksiyonu
func setFieldValue(fieldValue reflect.Value, queryValue string, fieldType reflect.Kind) error {
	switch fieldType {
	case reflect.String:
		fieldValue.SetString(queryValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(queryValue, 10, 64)
		if err != nil {
			return fmt.Errorf("error converting %s to int: %w", queryValue, err)
		}
		fieldValue.SetInt(intValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(queryValue, 64)
		if err != nil {
			return fmt.Errorf("error converting %s to float: %w", queryValue, err)
		}
		fieldValue.SetFloat(floatValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(queryValue)
		if err != nil {
			return fmt.Errorf("error converting %s to bool: %w", queryValue, err)
		}
		fieldValue.SetBool(boolValue)
	default:
		return fmt.Errorf("unsupported field type: %s", fieldType)
	}
	return nil
}

package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"golang.org/x/net/context"
)

type ErrValidation []ResponseValidation

func (e ErrValidation) Error() string {
	return "This form doesn't correct"
}

type ResponseValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(ctx context.Context, in interface{}) error {
	validate := validator.New()

	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.String()
		}
		return nil
	}, decimal.Decimal{})

	if err := validate.RegisterValidation("dgt", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		baseValue, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}
		return value.GreaterThan(baseValue)
	}); err != nil {
		return err
	}

	if err := validate.RegisterValidation("dgte", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		baseValue, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}
		return value.GreaterThanOrEqual(baseValue)
	}); err != nil {
		return err
	}

	err := validate.StructCtx(ctx, in)
	if err != nil {
		errList := validateStruct(err)
		if len(errList) > 0 {
			return errList
		}
		return err
	}

	return nil
}

func validateStruct(errStruct error) ErrValidation {
	var validationMessages ErrValidation
	var validationErrors validator.ValidationErrors
	if errors.As(errStruct, &validationErrors) {
		for _, errParam := range validationErrors {
			field := errParam.Field()
			tag := errParam.Tag()
			message := generateValidationMessage(field, tag)
			validationMessages = append(validationMessages, ResponseValidation{field, message})
		}

		return validationMessages
	}

	return validationMessages
}

func generateValidationMessage(field, tag string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least 3 characters long"
	default:
		return field + " is not valid"
	}
}

func BindStringSlices(obj any) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("obj must be a pointer to a struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("form")

		if tag == "" {
			continue
		}

		switch field.Kind() {
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				params := []string{}
				for j := 0; j < field.Len(); j++ {
					param := fmt.Sprintf("%s", field.Index(j).Interface())
					params = append(params, strings.Split(param, ",")...)
				}
				field.Set(reflect.ValueOf(params))
			}
		}
	}

	return nil
}

type Request interface {
	Validate(ctx context.Context) error
}

func MakeValidateBody[T Request](ctx *gin.Context) (*T, error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}

	var req T
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	if err := req.Validate(ctx); err != nil {
		return nil, err
	}
	return &req, nil
}

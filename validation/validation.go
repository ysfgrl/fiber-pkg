package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/ysfgrl/gerror"
)

var myValidator = validator.New()

func Validate(schema interface{}) *gerror.Error {
	err := myValidator.Struct(schema)
	if err != nil {
		var vErrors []Error
		for _, err := range err.(validator.ValidationErrors) {
			var el Error
			el.Field = err.Field()
			el.Value = err.Value()
			el.Tag = err.Tag()
			el.Param = err.Param()
			vErrors = append(vErrors, el)
		}

		return &gerror.Error{
			Function: "Add",
			File:     "Controller",
			Detail:   vErrors[0].Field + " " + vErrors[0].Tag,
			Code:     "api.validation_error",
		}
	}
	return nil
}

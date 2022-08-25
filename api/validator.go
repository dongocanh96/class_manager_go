package api

import (
	"github.com/dongocanh96/class_manager_go/util"
	"github.com/go-playground/validator/v10"
)

var validSubject validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if subject, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedSubject(subject)
	}
	return false
}

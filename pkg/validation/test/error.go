package test

import (
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateError(path *field.Path, obj *v1alpha1.Error) field.ErrorList {
	var errs field.ErrorList
	if obj != nil {
		errs = append(errs, ValidateFileRefOrCheck(path, obj.FileRefOrCheck)...)
	}
	return errs
}

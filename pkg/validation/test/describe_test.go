package test

import (
	"testing"

	v1alpha1 "github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestValidateDescribe(t *testing.T) {
	tests := []struct {
		name      string
		input     *v1alpha1.Describe
		expectErr bool
		errMsg    string
	}{{
		name:      "No resource provided",
		input:     &v1alpha1.Describe{},
		expectErr: true,
		errMsg:    "a resource must be specified",
	}, {
		name: "Neither Name nor Selector provided",
		input: &v1alpha1.Describe{
			Resource: "pods",
		},
		expectErr: false,
	}, {
		name: "Both Name and Selector provided",
		input: &v1alpha1.Describe{
			Resource: "pods",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name:     "example-name",
				Selector: "example-selector",
			},
		},
		expectErr: true,
		errMsg:    "a name or label selector must be specified (found both)",
	}, {
		name: "Only Name provided",
		input: &v1alpha1.Describe{
			Resource: "pods",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Name: "example-name",
			},
		},
		expectErr: false,
	}, {
		name: "Only Selector provided",
		input: &v1alpha1.Describe{
			Resource: "pods",
			ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
				Selector: "example-selector",
			},
		},
		expectErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateDescribe(field.NewPath("testPath"), tt.input)
			if tt.expectErr {
				assert.NotEmpty(t, errs)
				assert.Contains(t, errs.ToAggregate().Error(), tt.errMsg)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

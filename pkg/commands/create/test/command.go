package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func Command() *cobra.Command {
	save := false
	force := false
	description := false
	cmd := &cobra.Command{
		Use:          "test",
		Short:        "Create a Chainsaw test",
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			for _, path := range args {
				abs, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				name := filepath.Base(abs)
				test := v1alpha1.Test{
					TypeMeta: metav1.TypeMeta{
						APIVersion: v1alpha1.SchemeGroupVersion.String(),
						Kind:       "Test",
					},
					Spec: v1alpha1.TestSpec{
						Description: getDescription(description, "test description"),
						Steps:       sampleSteps(description),
					},
				}
				test.SetName(strings.ToLower(strings.ReplaceAll(name, "_", "-")))
				data, err := yaml.Marshal(&test)
				if err != nil {
					return err
				}
				if save {
					path := filepath.Join(path, "chainsaw-test.yaml")
					fmt.Fprintf(out, "Saving file %s ...\n", path)
					if err := os.WriteFile(path, data, os.ModePerm); err != nil {
						return err
					}
				} else {
					fmt.Fprintln(out, string(data))
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&save, "save", false, "If set, created test will be saved")
	cmd.Flags().BoolVar(&force, "force", false, "If set, existing test will be deleted if needed")
	cmd.Flags().BoolVar(&description, "description", true, "If set, adds description when applicable")
	return cmd
}

func sampleSteps(description bool) []v1alpha1.TestSpecStep {
	return []v1alpha1.TestSpecStep{{
		Name: "step 1",
		TestStepSpec: v1alpha1.TestStepSpec{
			Description: getDescription(description, "sample step 1"),
			Try: []v1alpha1.Operation{{
				Description: getDescription(description, "sample apply operation"),
				Apply: &v1alpha1.Apply{
					FileRefOrResource: v1alpha1.FileRefOrResource{
						FileRef: v1alpha1.FileRef{
							File: "resources.yaml",
						},
					},
				},
			}, {
				Description: getDescription(description, "sample assert operation"),
				Assert: &v1alpha1.Assert{
					FileRefOrCheck: v1alpha1.FileRefOrCheck{
						FileRef: v1alpha1.FileRef{
							File: "assert.yaml",
						},
					},
				},
			}, {
				Description: getDescription(description, "sample error operation"),
				Error: &v1alpha1.Error{
					FileRefOrCheck: v1alpha1.FileRefOrCheck{
						FileRef: v1alpha1.FileRef{
							File: "error.yaml",
						},
					},
				},
			}, {
				Description: getDescription(description, "sample delete operation"),
				Delete: &v1alpha1.Delete{
					ObjectReference: v1alpha1.ObjectReference{
						ObjectSelector: v1alpha1.ObjectSelector{
							Name: "foo",
						},
						APIVersion: "v1",
						Kind:       "Pod",
					},
				},
			}, {
				Description: getDescription(description, "sample script operation"),
				Script: &v1alpha1.Script{
					Content: `echo "test namespace = $NAMESPACE"`,
				},
			}},
			Catch: []v1alpha1.Catch{{
				Description: getDescription(description, "sample events collector"),
				Events: &v1alpha1.Events{
					ObjectLabelsSelector: v1alpha1.ObjectLabelsSelector{
						Name: "foo",
					},
				},
			}, {
				Description: getDescription(description, "sample pod logs collector"),
				PodLogs: &v1alpha1.PodLogs{
					Selector: "app=foo",
				},
			}},
			Finally: []v1alpha1.Finally{{
				Description: getDescription(description, "sample sleep operation"),
				Sleep: &v1alpha1.Sleep{
					Duration: metav1.Duration{Duration: 5 * time.Second},
				},
			}},
		},
	}}
}

func getDescription(enabled bool, desc string) string {
	if !enabled {
		return ""
	}
	return desc
}

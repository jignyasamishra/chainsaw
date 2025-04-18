package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type ReportFormatType string

const (
	JSONFormat ReportFormatType = "JSON"
	XMLFormat  ReportFormatType = "XML"
	NoReport   ReportFormatType = ""
)

// ConfigurationSpec contains the configuration used to run tests.
type ConfigurationSpec struct {
	// Global timeouts configuration. Applies to all tests/test steps if not overridden.
	// +optional
	Timeouts Timeouts `json:"timeouts"`

	// If set, do not delete the resources after running the tests (implies SkipClusterDelete).
	// +optional
	SkipDelete bool `json:"skipDelete,omitempty"`

	// Template determines whether resources should be considered for templating.
	// +optional
	Template *bool `json:"template,omitempty"`

	// FailFast determines whether the test should stop upon encountering the first failure.
	// +optional
	FailFast bool `json:"failFast,omitempty"`

	// The maximum number of tests to run at once.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	Parallel *int `json:"parallel,omitempty"`

	// ReportFormat determines test report format (JSON|XML|nil) nil == no report.
	// maps to report.Type, however we don't want generated.deepcopy to have reference to it.
	// +optional
	// +kubebuilder:validation:Enum:=JSON;XML;
	ReportFormat ReportFormatType `json:"reportFormat,omitempty"`

	// ReportPath defines the path.
	// +optional
	ReportPath string `json:"reportPath,omitempty"`

	// ReportName defines the name of report to create. It defaults to "chainsaw-report".
	// +optional
	// +kubebuilder:default:="chainsaw-report"
	ReportName string `json:"reportName,omitempty"`

	// Namespace defines the namespace to use for tests.
	// If not specified, every test will execute in a random ephemeral namespace
	// unless the namespace is overridden in a the test spec.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// NamespaceTemplate defines a template to create the test namespace.
	// +optional
	NamespaceTemplate *Any `json:"namespaceTemplate,omitempty"`

	// FullName makes use of the full test case folder path instead of the folder name.
	// +optional
	FullName bool `json:"fullName,omitempty"`

	// ExcludeTestRegex is used to exclude tests based on a regular expression.
	// +optional
	ExcludeTestRegex string `json:"excludeTestRegex,omitempty"`

	// IncludeTestRegex is used to include tests based on a regular expression.
	// +optional
	IncludeTestRegex string `json:"includeTestRegex,omitempty"`

	// RepeatCount indicates how many times the tests should be executed.
	// +kubebuilder:validation:Format:=int
	// +kubebuilder:validation:Minimum:=1
	// +optional
	RepeatCount *int `json:"repeatCount,omitempty"`

	// TestFile is the name of the file containing the test to run.
	// +kubebuilder:default:="chainsaw-test.yaml"
	// +optional
	TestFile string `json:"testFile,omitempty"`

	// ForceTerminationGracePeriod forces the termination grace period on pods, statefulsets, daemonsets and deployments.
	// +optional
	ForceTerminationGracePeriod *metav1.Duration `json:"forceTerminationGracePeriod,omitempty"`

	// DelayBeforeCleanup adds a delay between the time a test ends and the time cleanup starts.
	// +optional
	DelayBeforeCleanup *metav1.Duration `json:"delayBeforeCleanup,omitempty"`

	// Clusters holds a registry to clusters to support multi-cluster tests.
	// +optional
	Clusters map[string]Cluster `json:"clusters,omitempty"`
}

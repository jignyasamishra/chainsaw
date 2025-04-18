package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodLogs defines how to collect pod logs.
type PodLogs struct {
	// Timeout for the operation. Overrides the global timeout set in the Configuration.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// Cluster defines the target cluster (default cluster will be used if not specified and/or overridden).
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	Name string `json:"name,omitempty"`

	// Selector defines labels selector.
	// +optional
	Selector string `json:"selector,omitempty"`

	// Container in pod to get logs from else --all-containers is used.
	// +optional
	Container string `json:"container,omitempty"`

	// Tail is the number of last lines to collect from pods. If omitted or zero,
	// then the default is 10 if you use a selector, or -1 (all) if you use a pod name.
	// This matches default behavior of `kubectl logs`.
	// +optional
	Tail *int `json:"tail,omitempty"`
}

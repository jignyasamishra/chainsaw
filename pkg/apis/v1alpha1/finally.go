package v1alpha1

// Finally defines actions to be executed at the end of a test.
type Finally struct {
	// Description contains a description of the operation.
	// +optional
	Description string `json:"description,omitempty"`

	// PodLogs determines the pod logs collector to execute.
	// +optional
	PodLogs *PodLogs `json:"podLogs,omitempty"`

	// Events determines the events collector to execute.
	// +optional
	Events *Events `json:"events,omitempty"`

	// Describe determines the resource describe collector to execute.
	// +optional
	Describe *Describe `json:"describe,omitempty"`

	// Get determines the resource get collector to execute.
	// +optional
	Get *Get `json:"get,omitempty"`

	// Delete represents a deletion operation.
	// +optional
	Delete *Delete `json:"delete,omitempty"`

	// Command defines a command to run.
	// +optional
	Command *Command `json:"command,omitempty"`

	// Script defines a script to run.
	// +optional
	Script *Script `json:"script,omitempty"`

	// Sleep defines zzzz.
	// +optional
	Sleep *Sleep `json:"sleep,omitempty"`
}

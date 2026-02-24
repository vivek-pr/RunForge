/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AIJobSpec defines the desired state of AIJob
type AIJobSpec struct {
	// image is the container image used to run this job.
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// command is the container entrypoint command.
	// +optional
	Command []string `json:"command,omitempty"`

	// args are arguments passed to the container command.
	// +optional
	Args []string `json:"args,omitempty"`

	// env is a list of name/value environment variables.
	// +optional
	Env []AIJobEnvVar `json:"env,omitempty"`

	// envFrom references external sources for environment variables.
	// +optional
	EnvFrom []AIJobEnvFromSource `json:"envFrom,omitempty"`

	// resources defines compute resource requests and limits.
	// +optional
	Resources *AIJobResourceRequirements `json:"resources,omitempty"`

	// restartPolicy controls retry behavior for failed pods.
	// +kubebuilder:default:=Never
	// +kubebuilder:validation:Enum=Never;OnFailure
	RestartPolicy string `json:"restartPolicy,omitempty"`

	// backoffLimit is the number of retries before marking the job failed.
	// +kubebuilder:validation:Minimum=0
	// +optional
	BackoffLimit *int32 `json:"backoffLimit,omitempty"`

	// activeDeadlineSeconds is the max runtime in seconds before termination.
	// +kubebuilder:validation:Minimum=1
	// +optional
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty"`

	// ttlSecondsAfterFinished is the retention period for completed jobs.
	// +kubebuilder:validation:Minimum=0
	// +optional
	TTLSecondsAfterFinished *int32 `json:"ttlSecondsAfterFinished,omitempty"`

	// nodeSelector constrains scheduling onto specific nodes.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// tolerations allow scheduling onto tainted nodes.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// affinity defines affinity/anti-affinity scheduling preferences.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// serviceAccountName is the pod service account identity.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
}

// AIJobStatus defines the observed state of AIJob.
type AIJobStatus struct {
	// observedGeneration reflects the generation last processed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// phase is a coarse-grained summary of current execution state.
	// +kubebuilder:validation:Enum=Pending;Running;Succeeded;Failed
	// +optional
	Phase AIJobPhase `json:"phase,omitempty"`

	// jobName references the backing Kubernetes Job resource name.
	// +optional
	JobName string `json:"jobName,omitempty"`

	// conditions represent the current state transitions of AIJob.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []AIJobCondition `json:"conditions,omitempty"`

	// startTime is when execution began.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// completionTime is when execution reached a terminal phase.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// lastError stores the latest terminal or transient error message.
	// +optional
	LastError string `json:"lastError,omitempty"`
}

// AIJobPhase defines the execution phase of an AIJob.
type AIJobPhase string

const (
	AIJobPhasePending   AIJobPhase = "Pending"
	AIJobPhaseRunning   AIJobPhase = "Running"
	AIJobPhaseSucceeded AIJobPhase = "Succeeded"
	AIJobPhaseFailed    AIJobPhase = "Failed"
)

// AIJobEnvVar is an explicit name/value environment variable.
type AIJobEnvVar struct {
	// name is the environment variable key.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// value is the environment variable value.
	// +optional
	Value string `json:"value,omitempty"`
}

// AIJobLocalObjectReference references an object in the same namespace.
type AIJobLocalObjectReference struct {
	// name is the referenced object name.
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

// AIJobEnvFromSource references a ConfigMap or Secret for env injection.
type AIJobEnvFromSource struct {
	// configMapRef references a ConfigMap source.
	// +optional
	ConfigMapRef *AIJobLocalObjectReference `json:"configMapRef,omitempty"`

	// secretRef references a Secret source.
	// +optional
	SecretRef *AIJobLocalObjectReference `json:"secretRef,omitempty"`
}

// AIJobResourceList captures cpu/memory and optional gpu values.
type AIJobResourceList struct {
	// cpu is a CPU quantity (for example "500m" or "2").
	// +optional
	CPU string `json:"cpu,omitempty"`

	// memory is a memory quantity (for example "512Mi" or "2Gi").
	// +optional
	Memory string `json:"memory,omitempty"`

	// gpu is an optional accelerator quantity (for example "1").
	// +optional
	GPU string `json:"gpu,omitempty"`
}

// AIJobResourceRequirements describes requests and limits.
type AIJobResourceRequirements struct {
	// requests are minimum required resources.
	// +optional
	Requests *AIJobResourceList `json:"requests,omitempty"`

	// limits are max allowed resources.
	// +optional
	Limits *AIJobResourceList `json:"limits,omitempty"`
}

// AIJobCondition records a status transition detail.
type AIJobCondition struct {
	// type identifies the condition category.
	// +kubebuilder:validation:MinLength=1
	Type string `json:"type"`

	// reason is a machine-readable condition reason.
	// +kubebuilder:validation:MinLength=1
	Reason string `json:"reason"`

	// message is a human-readable condition message.
	// +optional
	Message string `json:"message,omitempty"`

	// lastTransitionTime is when this condition value last changed.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AIJob is the Schema for the aijobs API
type AIJob struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of AIJob
	// +required
	Spec AIJobSpec `json:"spec"`

	// status defines the observed state of AIJob
	// +optional
	Status AIJobStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// AIJobList contains a list of AIJob
type AIJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []AIJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AIJob{}, &AIJobList{})
}

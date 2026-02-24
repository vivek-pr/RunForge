package jobfactory

import (
	"fmt"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runforgev1alpha1 "github.com/vivekpradhan/runforge/api/v1alpha1"
)

const gpuResourceName corev1.ResourceName = "nvidia.com/gpu"

// BuildJob builds a Kubernetes Job from an AIJob spec without mutating inputs.
func BuildJob(aijob *runforgev1alpha1.AIJob, jobName string) (*batchv1.Job, error) {
	if aijob == nil {
		return nil, fmt.Errorf("aijob must not be nil")
	}
	if strings.TrimSpace(jobName) == "" {
		return nil, fmt.Errorf("jobName must not be empty")
	}

	restartPolicy := corev1.RestartPolicyNever
	if aijob.Spec.RestartPolicy != "" {
		restartPolicy = corev1.RestartPolicy(aijob.Spec.RestartPolicy)
	}

	container := corev1.Container{
		Name:    "worker",
		Image:   aijob.Spec.Image,
		Command: append([]string{}, aijob.Spec.Command...),
		Args:    append([]string{}, aijob.Spec.Args...),
	}

	for _, e := range aijob.Spec.Env {
		container.Env = append(container.Env, corev1.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		})
	}
	for _, from := range aijob.Spec.EnvFrom {
		src := corev1.EnvFromSource{}
		if from.ConfigMapRef != nil {
			src.ConfigMapRef = &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: from.ConfigMapRef.Name},
			}
		}
		if from.SecretRef != nil {
			src.SecretRef = &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: from.SecretRef.Name},
			}
		}
		container.EnvFrom = append(container.EnvFrom, src)
	}

	if aijob.Spec.Resources != nil {
		resources, err := buildResourceRequirements(aijob.Spec.Resources)
		if err != nil {
			return nil, err
		}
		container.Resources = resources
	}

	labels := map[string]string{
		"app.kubernetes.io/name":      "runforge",
		"app.kubernetes.io/component": "worker",
		"runforge.io/aijob":           aijob.Name,
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: aijob.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            aijob.Spec.BackoffLimit,
			ActiveDeadlineSeconds:   aijob.Spec.ActiveDeadlineSeconds,
			TTLSecondsAfterFinished: aijob.Spec.TTLSecondsAfterFinished,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					RestartPolicy:      restartPolicy,
					Containers:         []corev1.Container{container},
					NodeSelector:       aijob.Spec.NodeSelector,
					Tolerations:        aijob.Spec.Tolerations,
					Affinity:           aijob.Spec.Affinity,
					ServiceAccountName: aijob.Spec.ServiceAccountName,
				},
			},
		},
	}, nil
}

func buildResourceRequirements(spec *runforgev1alpha1.AIJobResourceRequirements) (corev1.ResourceRequirements, error) {
	requests, err := buildResourceList(spec.Requests)
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}
	limits, err := buildResourceList(spec.Limits)
	if err != nil {
		return corev1.ResourceRequirements{}, err
	}
	return corev1.ResourceRequirements{
		Requests: requests,
		Limits:   limits,
	}, nil
}

func buildResourceList(in *runforgev1alpha1.AIJobResourceList) (corev1.ResourceList, error) {
	if in == nil {
		return nil, nil
	}
	out := corev1.ResourceList{}
	if in.CPU != "" {
		q, err := resource.ParseQuantity(in.CPU)
		if err != nil {
			return nil, fmt.Errorf("invalid cpu quantity %q: %w", in.CPU, err)
		}
		out[corev1.ResourceCPU] = q
	}
	if in.Memory != "" {
		q, err := resource.ParseQuantity(in.Memory)
		if err != nil {
			return nil, fmt.Errorf("invalid memory quantity %q: %w", in.Memory, err)
		}
		out[corev1.ResourceMemory] = q
	}
	if in.GPU != "" {
		q, err := resource.ParseQuantity(in.GPU)
		if err != nil {
			return nil, fmt.Errorf("invalid gpu quantity %q: %w", in.GPU, err)
		}
		out[gpuResourceName] = q
	}
	return out, nil
}

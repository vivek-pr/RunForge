package jobfactory

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runforgev1alpha1 "github.com/vivekpradhan/runforge/api/v1alpha1"
)

func TestBuildJob(t *testing.T) {
	aijob := &runforgev1alpha1.AIJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo",
			Namespace: "default",
		},
		Spec: runforgev1alpha1.AIJobSpec{
			Image:         "busybox:1.36",
			Command:       []string{"sh", "-c"},
			Args:          []string{"echo ok"},
			RestartPolicy: "Never",
			Resources: &runforgev1alpha1.AIJobResourceRequirements{
				Requests: &runforgev1alpha1.AIJobResourceList{
					CPU:    "100m",
					Memory: "128Mi",
				},
				Limits: &runforgev1alpha1.AIJobResourceList{
					CPU:    "500m",
					Memory: "256Mi",
					GPU:    "1",
				},
			},
		},
	}

	job, err := BuildJob(aijob, "demo-job")
	if err != nil {
		t.Fatalf("BuildJob returned error: %v", err)
	}
	if job.Name != "demo-job" {
		t.Fatalf("unexpected job name: %s", job.Name)
	}
	if job.Spec.Template.Spec.Containers[0].Image != "busybox:1.36" {
		t.Fatalf("unexpected image: %s", job.Spec.Template.Spec.Containers[0].Image)
	}
	if _, ok := job.Spec.Template.Spec.Containers[0].Resources.Limits["nvidia.com/gpu"]; !ok {
		t.Fatalf("expected nvidia.com/gpu resource in limits")
	}
}

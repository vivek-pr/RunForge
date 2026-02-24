package status

import (
	"testing"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runforgev1alpha1 "github.com/vivekpradhan/runforge/api/v1alpha1"
)

func TestApplyFromJobFailed(t *testing.T) {
	now := metav1.Now()
	start := metav1.NewTime(now.Add(-1 * time.Second))
	completion := metav1.NewTime(now.Time)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-job"},
		Status: batchv1.JobStatus{
			StartTime:      &start,
			CompletionTime: &completion,
			Conditions: []batchv1.JobCondition{
				{
					Type:               batchv1.JobFailed,
					Status:             corev1.ConditionTrue,
					Reason:             "BackoffLimitExceeded",
					Message:            "job has reached backoff limit",
					LastTransitionTime: completion,
				},
			},
		},
	}
	aijob := &runforgev1alpha1.AIJob{}

	ApplyFromJob(aijob, job, now)

	if aijob.Status.Phase != runforgev1alpha1.AIJobPhaseFailed {
		t.Fatalf("expected Failed phase, got %s", aijob.Status.Phase)
	}
	if aijob.Status.LastError == "" {
		t.Fatalf("expected lastError to be set")
	}
	if aijob.Status.CompletionTime == nil {
		t.Fatalf("expected completionTime to be set")
	}
	if len(aijob.Status.Conditions) != 1 || aijob.Status.Conditions[0].Type != "Failed" {
		t.Fatalf("expected Failed condition")
	}
}

func TestApplyFromJobRunningThenSucceeded(t *testing.T) {
	now := metav1.Now()
	start := metav1.NewTime(now.Add(-2 * time.Second))

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{Name: "demo-job"},
		Status: batchv1.JobStatus{
			StartTime: &start,
			Active:    1,
		},
	}
	aijob := &runforgev1alpha1.AIJob{}

	ApplyFromJob(aijob, job, now)
	if aijob.Status.Phase != runforgev1alpha1.AIJobPhaseRunning {
		t.Fatalf("expected Running phase, got %s", aijob.Status.Phase)
	}

	completion := metav1.NewTime(now.Add(1 * time.Second))
	job.Status.Active = 0
	job.Status.CompletionTime = &completion
	job.Status.Conditions = []batchv1.JobCondition{
		{
			Type:               batchv1.JobComplete,
			Status:             corev1.ConditionTrue,
			Reason:             "Completed",
			Message:            "job completed successfully",
			LastTransitionTime: completion,
		},
	}
	ApplyFromJob(aijob, job, now)
	if aijob.Status.Phase != runforgev1alpha1.AIJobPhaseSucceeded {
		t.Fatalf("expected Succeeded phase, got %s", aijob.Status.Phase)
	}
	if len(aijob.Status.Conditions) < 2 {
		t.Fatalf("expected lifecycle condition history to include Running and Succeeded")
	}
}

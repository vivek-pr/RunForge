package status

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runforgev1alpha1 "github.com/vivekpradhan/runforge/api/v1alpha1"
)

// ApplyFromJob maps Job lifecycle state into AIJob status fields.
func ApplyFromJob(aijob *runforgev1alpha1.AIJob, job *batchv1.Job, now metav1.Time) {
	phase, reason, message, transitionTime := phaseAndCondition(job, now)

	aijob.Status.Phase = phase
	aijob.Status.StartTime = copyTime(job.Status.StartTime)
	aijob.Status.CompletionTime = copyTime(job.Status.CompletionTime)
	if aijob.Status.CompletionTime == nil &&
		(phase == runforgev1alpha1.AIJobPhaseSucceeded || phase == runforgev1alpha1.AIJobPhaseFailed) {
		t := transitionTime
		aijob.Status.CompletionTime = &t
	}

	if phase == runforgev1alpha1.AIJobPhaseFailed {
		aijob.Status.LastError = buildLastError(reason, message)
	} else {
		aijob.Status.LastError = ""
	}

	cond := runforgev1alpha1.AIJobCondition{
		Type:               string(phase),
		Reason:             reason,
		Message:            message,
		LastTransitionTime: transitionTime,
	}
	aijob.Status.Conditions = upsertCondition(aijob.Status.Conditions, cond)
}

func phaseAndCondition(job *batchv1.Job, now metav1.Time) (runforgev1alpha1.AIJobPhase, string, string, metav1.Time) {
	for _, cond := range job.Status.Conditions {
		if cond.Type == batchv1.JobComplete && cond.Status == corev1.ConditionTrue {
			reason := cond.Reason
			if reason == "" {
				reason = "JobCompleted"
			}
			return runforgev1alpha1.AIJobPhaseSucceeded, reason, cond.Message, cond.LastTransitionTime
		}
		if cond.Type == batchv1.JobFailed && cond.Status == corev1.ConditionTrue {
			reason := cond.Reason
			if reason == "" {
				reason = "JobFailed"
			}
			return runforgev1alpha1.AIJobPhaseFailed, reason, cond.Message, cond.LastTransitionTime
		}
	}

	if job.Status.Active > 0 {
		transitionTime := now
		if job.Status.StartTime != nil {
			transitionTime = *job.Status.StartTime
		}
		return runforgev1alpha1.AIJobPhaseRunning, "JobRunning", fmt.Sprintf("active pods: %d", job.Status.Active), transitionTime
	}

	return runforgev1alpha1.AIJobPhasePending, "JobPending", "waiting for job to start", now
}

func buildLastError(reason, message string) string {
	if message == "" {
		return reason
	}
	return fmt.Sprintf("%s: %s", reason, message)
}

func upsertCondition(conditions []runforgev1alpha1.AIJobCondition, incoming runforgev1alpha1.AIJobCondition) []runforgev1alpha1.AIJobCondition {
	for i := range conditions {
		if conditions[i].Type != incoming.Type {
			continue
		}
		if conditions[i].Reason == incoming.Reason && conditions[i].Message == incoming.Message {
			return conditions
		}
		conditions[i].Reason = incoming.Reason
		conditions[i].Message = incoming.Message
		conditions[i].LastTransitionTime = incoming.LastTransitionTime
		return conditions
	}
	return append(conditions, incoming)
}

func copyTime(in *metav1.Time) *metav1.Time {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}

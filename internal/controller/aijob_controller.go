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

package controller

import (
	"context"
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	runforgev1alpha1 "github.com/vivekpradhan/runforge/api/v1alpha1"
	"github.com/vivekpradhan/runforge/internal/jobfactory"
	statusutil "github.com/vivekpradhan/runforge/internal/status"
)

// AIJobReconciler reconciles a AIJob object
type AIJobReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=runforge.runforge.io,resources=aijobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=runforge.runforge.io,resources=aijobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=runforge.runforge.io,resources=aijobs/finalizers,verbs=update
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AIJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.1/pkg/reconcile
func (r *AIJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var aijob runforgev1alpha1.AIJob
	if err := r.Get(ctx, req.NamespacedName, &aijob); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	jobName := desiredJobName(&aijob)
	var job batchv1.Job
	err := r.Get(ctx, client.ObjectKey{Namespace: aijob.Namespace, Name: jobName}, &job)
	if err != nil && !apierrors.IsNotFound(err) {
		return ctrl.Result{}, err
	}
	if apierrors.IsNotFound(err) {
		newJob, buildErr := jobfactory.BuildJob(&aijob, jobName)
		if buildErr != nil {
			return ctrl.Result{}, buildErr
		}
		if setErr := controllerutil.SetControllerReference(&aijob, newJob, r.Scheme); setErr != nil {
			return ctrl.Result{}, setErr
		}
		if createErr := r.Create(ctx, newJob); createErr != nil {
			return ctrl.Result{}, createErr
		}
		if r.Recorder != nil {
			r.Recorder.Eventf(&aijob, corev1.EventTypeNormal, "JobCreated", "Created Job %s", newJob.Name)
		}
		log.Info("created Job for AIJob", "aijob", req.NamespacedName, "job", newJob.Name)

		if statusErr := r.updateAIJobStatus(ctx, &aijob, newJob); statusErr != nil {
			return ctrl.Result{}, statusErr
		}
		return ctrl.Result{}, nil
	}

	if statusErr := r.updateAIJobStatus(ctx, &aijob, &job); statusErr != nil {
		return ctrl.Result{}, statusErr
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AIJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&runforgev1alpha1.AIJob{}).
		Owns(&batchv1.Job{}).
		Named("aijob").
		Complete(r)
}

func (r *AIJobReconciler) updateAIJobStatus(ctx context.Context, aijob *runforgev1alpha1.AIJob, job *batchv1.Job) error {
	original := aijob.DeepCopy()
	aijob.Status.JobName = job.Name
	aijob.Status.ObservedGeneration = aijob.Generation
	statusutil.ApplyFromJob(aijob, job, metav1.Now())
	if statusUnchanged(original.Status, aijob.Status) {
		return nil
	}
	return r.Status().Patch(ctx, aijob, client.MergeFrom(original))
}

func desiredJobName(aijob *runforgev1alpha1.AIJob) string {
	if aijob.Status.JobName != "" {
		return aijob.Status.JobName
	}
	return fmt.Sprintf("%s-job", aijob.Name)
}

func statusUnchanged(oldStatus, newStatus runforgev1alpha1.AIJobStatus) bool {
	if oldStatus.JobName != newStatus.JobName ||
		oldStatus.Phase != newStatus.Phase ||
		oldStatus.ObservedGeneration != newStatus.ObservedGeneration ||
		oldStatus.LastError != newStatus.LastError ||
		!timeEqual(oldStatus.StartTime, newStatus.StartTime) ||
		!timeEqual(oldStatus.CompletionTime, newStatus.CompletionTime) ||
		len(oldStatus.Conditions) != len(newStatus.Conditions) {
		return false
	}
	for i := range oldStatus.Conditions {
		a := oldStatus.Conditions[i]
		b := newStatus.Conditions[i]
		if a.Type != b.Type ||
			a.Reason != b.Reason ||
			a.Message != b.Message ||
			!a.LastTransitionTime.Equal(&b.LastTransitionTime) {
			return false
		}
	}
	return true
}

func timeEqual(a, b *metav1.Time) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equal(b)
}

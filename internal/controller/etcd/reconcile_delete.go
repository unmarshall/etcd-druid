package etcd

import (
	"context"
	"fmt"
	"time"

	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"github.com/gardener/etcd-druid/internal/common"
	ctrlutils "github.com/gardener/etcd-druid/internal/controller/utils"
	"github.com/gardener/etcd-druid/internal/resource"
	"github.com/gardener/gardener/pkg/controllerutils"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// deleteStepFn is a step in the deletion flow. Every deletion step must have this signature.
type deleteStepFn func(ctx resource.OperatorContext, logger logr.Logger, etcdObjectKey client.ObjectKey) ctrlutils.ReconcileStepResult

// triggerDeletionFlow is the entry point for the deletion flow triggered for an etcd resource which has a DeletionTimeStamp set on it.
func (r *Reconciler) triggerDeletionFlow(ctx context.Context, logger logr.Logger, etcd *druidv1alpha1.Etcd) (ctrl.Result, error) {
	operatorCtx := resource.NewOperatorContext(ctx, r.client, r.logger, etcd.GetNamespaceName())
	deleteStepFns := []deleteStepFn{
		r.recordDeletionStartOperation,
		r.deleteEtcdResources,
		r.verifyNoResourcesAwaitCleanUp,
		r.removeFinalizer,
	}
	for _, fn := range deleteStepFns {
		if stepResult := fn(operatorCtx, logger, etcd.GetNamespaceName()); ctrlutils.ShortCircuitReconcile(stepResult) {
			return r.recordIncompleteDeletionOperation(operatorCtx, logger, etcd.GetNamespaceName(), stepResult).ReconcileResult()
		}
	}
	// if we are here this means that all deletion steps have been successful
	return r.recordDeletionSuccessOperation(operatorCtx, logger, etcd.GetNamespaceName()).ReconcileResult()
}

func (r *Reconciler) deleteEtcdResources(ctx resource.OperatorContext, logger logr.Logger, etcdObjKey client.ObjectKey) ctrlutils.ReconcileStepResult {
	etcd := &druidv1alpha1.Etcd{}
	if result := r.getLatestEtcd(ctx, etcdObjKey, etcd); ctrlutils.ShortCircuitReconcile(result) {
		return result
	}
	operators := r.operatorRegistry.AllOperators()
	deleteTasks := make([]resource.OperatorTask, len(operators))
	for kind, operator := range operators {
		operator := operator
		deleteTasks = append(deleteTasks, resource.OperatorTask{
			Name: fmt.Sprintf("triggerDeletionFlow-%s-operator", kind),
			Fn: func(ctx resource.OperatorContext) error {
				return operator.TriggerDelete(ctx, etcd)
			},
		})
	}
	logger.Info("triggering triggerDeletionFlow operators for all resources")
	if errs := resource.RunConcurrently(ctx, deleteTasks); len(errs) > 0 {
		return ctrlutils.ReconcileWithError(errs...)
	}
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) verifyNoResourcesAwaitCleanUp(ctx resource.OperatorContext, logger logr.Logger, etcdObjKey client.ObjectKey) ctrlutils.ReconcileStepResult {
	etcd := &druidv1alpha1.Etcd{}
	if result := r.getLatestEtcd(ctx, etcdObjKey, etcd); ctrlutils.ShortCircuitReconcile(result) {
		return result
	}
	operators := r.operatorRegistry.AllOperators()
	resourceNamesAwaitingCleanup := make([]string, 0, len(operators))
	for _, operator := range operators {
		existingResourceNames, err := operator.GetExistingResourceNames(ctx, etcd)
		if err != nil {
			return ctrlutils.ReconcileWithError(err)
		}
		resourceNamesAwaitingCleanup = append(resourceNamesAwaitingCleanup, existingResourceNames...)
	}
	if len(resourceNamesAwaitingCleanup) > 0 {
		logger.Info("Cleanup of all resources has not yet been completed", "resourceNamesAwaitingCleanup", resourceNamesAwaitingCleanup)
		return ctrlutils.ReconcileAfter(5*time.Second, "Cleanup of all resources has not yet been completed. Skipping removal of Finalizer")
	}
	logger.Info("All resources have been cleaned up")
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) removeFinalizer(ctx resource.OperatorContext, logger logr.Logger, etcdObjKey client.ObjectKey) ctrlutils.ReconcileStepResult {
	etcd := &druidv1alpha1.Etcd{}
	if result := r.getLatestEtcd(ctx, etcdObjKey, etcd); ctrlutils.ShortCircuitReconcile(result) {
		return result
	}
	logger.Info("Removing finalizer", "finalizerName", common.FinalizerName)
	if err := controllerutils.RemoveFinalizers(ctx, r.client, etcd, common.FinalizerName); client.IgnoreNotFound(err) != nil {
		return ctrlutils.ReconcileWithError(err)
	}
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) recordDeletionStartOperation(ctx resource.OperatorContext, logger logr.Logger, etcdObjKey client.ObjectKey) ctrlutils.ReconcileStepResult {
	etcd := &druidv1alpha1.Etcd{}
	if result := r.getLatestEtcd(ctx, etcdObjKey, etcd); ctrlutils.ShortCircuitReconcile(result) {
		return result
	}
	if err := r.lastOpErrRecorder.RecordStart(ctx, etcd, druidv1alpha1.LastOperationTypeDelete); err != nil {
		logger.Error(err, "failed to record etcd deletion start operation")
		return ctrlutils.ReconcileWithError(err)
	}
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) recordIncompleteDeletionOperation(ctx resource.OperatorContext, logger logr.Logger, etcdObjKey client.ObjectKey, exitReconcileStepResult ctrlutils.ReconcileStepResult) ctrlutils.ReconcileStepResult {
	etcd := &druidv1alpha1.Etcd{}
	if result := r.getLatestEtcd(ctx, etcdObjKey, etcd); ctrlutils.ShortCircuitReconcile(result) {
		return result
	}
	if err := r.lastOpErrRecorder.RecordError(ctx, etcd, druidv1alpha1.LastOperationTypeDelete, exitReconcileStepResult.GetDescription(), exitReconcileStepResult.GetErrors()...); err != nil {
		logger.Error(err, "failed to record last operation and last errors for etcd deletion")
		return ctrlutils.ReconcileWithError(err)
	}
	return exitReconcileStepResult
}

func (r *Reconciler) recordDeletionSuccessOperation(ctx resource.OperatorContext, logger logr.Logger, etcdObjKey client.ObjectKey) ctrlutils.ReconcileStepResult {
	etcd := &druidv1alpha1.Etcd{}
	if result := r.getLatestEtcd(ctx, etcdObjKey, etcd); ctrlutils.ShortCircuitReconcile(result) {
		return result
	}
	if err := r.lastOpErrRecorder.RecordSuccess(ctx, etcd, druidv1alpha1.LastOperationTypeDelete); err != nil {
		logger.Error(err, "failed to record last operation for successful etcd deletion")
		return ctrlutils.ReconcileWithError(err)
	}
	return ctrlutils.DoNotRequeue()
}

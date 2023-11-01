package etcd

import (
	"context"
	"errors"
	"fmt"
	"time"

	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	ctrlutils "github.com/gardener/etcd-druid/internal/controller/utils"
	"github.com/gardener/etcd-druid/internal/resource"
	"github.com/gardener/etcd-druid/pkg/common"
	"github.com/gardener/gardener/pkg/controllerutils"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type deleteStepFn func(ctx resource.OperatorContext, logger logr.Logger, etcd *druidv1alpha1.Etcd) ctrlutils.ReconcileStepResult

func (r *Reconciler) delete(ctx context.Context, logger logr.Logger, etcd *druidv1alpha1.Etcd) (ctrl.Result, error) {
	operatorCtx := resource.OperatorContext{Context: ctx}
	deleteStepFns := []deleteStepFn{
		r.deleteEtcdResources,
		r.verifyNoResourcesAwaitCleanUp,
		r.removeFinalizer,
	}
	for _, fn := range deleteStepFns {
		if stepResult := fn(operatorCtx, logger, etcd); ctrlutils.ShortCircuitReconcile(stepResult) {
			return stepResult.ReconcileResult()
		}
	}
	return ctrlutils.DoNotRequeue().ReconcileResult()
}

func (r *Reconciler) deleteEtcdResources(ctx resource.OperatorContext, logger logr.Logger, etcd *druidv1alpha1.Etcd) ctrlutils.ReconcileStepResult {
	operators := r.operatorRegistry.AllOperators()
	deleteTasks := make([]resource.OperatorTask, len(operators))
	for kind, operator := range operators {
		operator := operator
		deleteTasks = append(deleteTasks, resource.OperatorTask{
			Name: fmt.Sprintf("delete-%s-operator", kind),
			Fn: func(ctx resource.OperatorContext) error {
				return operator.TriggerDelete(ctx, etcd)
			},
		})
	}
	logger.Info("triggering delete operators for all resources")
	if err := errors.Join(resource.RunConcurrently(ctx, deleteTasks)...); err != nil {
		return ctrlutils.ReconcileWithError(err)
	}
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) verifyNoResourcesAwaitCleanUp(ctx resource.OperatorContext, logger logr.Logger, etcd *druidv1alpha1.Etcd) ctrlutils.ReconcileStepResult {
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
		return ctrlutils.ReconcileAfter(5 * time.Second)
	}
	logger.Info("All resources have been cleaned up")
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) removeFinalizer(ctx resource.OperatorContext, logger logr.Logger, etcd *druidv1alpha1.Etcd) ctrlutils.ReconcileStepResult {
	logger.Info("Removing finalizer", "finalizerName", common.FinalizerName)
	if err := controllerutils.RemoveFinalizers(ctx, r.client, etcd, common.FinalizerName); client.IgnoreNotFound(err) != nil {
		return ctrlutils.ReconcileWithError(err)
	}
	return ctrlutils.ContinueReconcile()
}

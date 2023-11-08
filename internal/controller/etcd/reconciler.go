package etcd

import (
	"context"

	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	ctrlutils "github.com/gardener/etcd-druid/internal/controller/utils"
	"github.com/gardener/etcd-druid/internal/resource"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/utils/imagevector"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Reconciler struct {
	client            client.Client
	config            *Config
	recorder          record.EventRecorder
	imageVector       imagevector.ImageVector
	operatorRegistry  resource.OperatorRegistry
	lastOpErrRecorder ctrlutils.LastOperationErrorRecorder
	logger            logr.Logger
}

// NewReconciler creates a new reconciler for Etcd.
func NewReconciler(mgr manager.Manager, config *Config) (*Reconciler, error) {
	imageVector, err := ctrlutils.CreateImageVector()
	logger := log.Log.WithName(controllerName)
	if err != nil {
		return nil, err
	}
	operatorReg := resource.NewOperatorRegistry(mgr.GetClient(),
		logger,
		resource.OperatorConfig{
			DisableEtcdServiceAccountAutomount: config.DisableEtcdServiceAccountAutomount,
		},
	)
	lastOpErrRecorder := ctrlutils.NewLastOperationErrorRecorder(mgr.GetClient(), logger)
	return &Reconciler{
		client:            mgr.GetClient(),
		config:            config,
		recorder:          mgr.GetEventRecorderFor(controllerName),
		imageVector:       imageVector,
		logger:            logger,
		operatorRegistry:  operatorReg,
		lastOpErrRecorder: lastOpErrRecorder,
	}, nil
}

/*
	If deletionTimestamp set:
		triggerDeletionFlow(); if err then requeue
	Else:
		If ignore-reconciliation is set to true:
			skip reconcileSpec()
		Else If IgnoreOperationAnnotation flag is true:
			always reconcileSpec()
		Else If IgnoreOperationAnnotation flag is false and reconcile-op annotation is present:
			reconcileSpec()
			if err in getting etcd, return with requeue
			if err in deploying any of the components, then record pending requeue
	reconcileStatus()
	requeue after minimum of X seconds (EtcdStatusSyncPeriod) and previously recorded requeue request
*/

// TODO: where/how is this being used?
// +kubebuilder:rbac:groups=druid.gardener.cloud,resources=etcds,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=druid.gardener.cloud,resources=etcds/status,verbs=get;create;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;create;update;patch;triggerDeletionFlow
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;create;update;patch;triggerDeletionFlow
// +kubebuilder:rbac:groups="",resources=serviceaccounts;services;configmaps,verbs=get;list;create;update;patch;triggerDeletionFlow
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;create;update;patch;triggerDeletionFlow
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;create;update;patch;triggerDeletionFlow
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;create;update;patch;triggerDeletionFlow
// +kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;watch
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;get;list

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	type reconcileFn func(ctx context.Context, objectKey client.ObjectKey) ctrlutils.ReconcileStepResult
	reconcileFns := []reconcileFn{
		r.reconcileEtcdDeletion,
		r.reconcileSpec,
		r.reconcileStatus,
	}
	for _, fn := range reconcileFns {
		if result := fn(ctx, req.NamespacedName); ctrlutils.ShortCircuitReconcile(result) {
			return result.ReconcileResult()
		}
	}

	return ctrlutils.DoNotRequeue().ReconcileResult()

	//if r.shouldReconcileSpec(etcd) {
	//	rLog := r.logger.WithValues("etcd", req.NamespacedName, "operation", "reconcile-spec")
	//	if result := r.reconcileSpec(ctx, etcd, rLog);
	//}
	//
	//r.reconcileStatus(ctx, req.NamespacedName)
	//
	//if reconcileSpecErr != nil {
	//	return ctrl.Result{Requeue: true}, reconcileSpecErr
	//}
	//return ctrl.Result{RequeueAfter: r.config.EtcdStatusSyncPeriod}, nil
}

func (r *Reconciler) reconcileEtcdDeletion(ctx context.Context, etcdObjectKey client.ObjectKey) ctrlutils.ReconcileStepResult {
	etcd := &druidv1alpha1.Etcd{}
	if result := r.getLatestEtcd(ctx, etcdObjectKey, etcd); ctrlutils.ShortCircuitReconcile(result) {
		return result
	}
	if etcd.IsMarkedForDeletion() {
		dLog := r.logger.WithValues("etcd", etcdObjectKey, "operation", "delete")
		return r.triggerDeletionFlow(ctx, dLog, etcdObjectKey)
	}
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) reconcileSpec(ctx context.Context, etcdObjectKey client.ObjectKey) ctrlutils.ReconcileStepResult {

	return ctrlutils.DoNotRequeue()
}

func (r *Reconciler) reconcileStatus(ctx context.Context, etcdNamespacedName types.NamespacedName) ctrlutils.ReconcileStepResult {
	/*
		fetch EtcdMember resources
		fetch member leases
		status.condition checks
		status.members checks
		fetch latest Etcd resource
		update etcd status
	*/

	return ctrlutils.DoNotRequeue()
}

func (r *Reconciler) getLatestEtcd(ctx context.Context, objectKey client.ObjectKey, etcd *druidv1alpha1.Etcd) ctrlutils.ReconcileStepResult {
	if err := r.client.Get(ctx, objectKey, etcd); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrlutils.DoNotRequeue()
		}
		return ctrlutils.ReconcileWithError(err)
	}
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) handleReconcilePause(etcd *druidv1alpha1.Etcd) ctrlutils.ReconcileStepResult {
	//TODO: Once no one uses IgnoreReconciliationAnnotation annotation, then we can simplify this code.
	var annotationKey string
	if metav1.HasAnnotation(etcd.ObjectMeta, druidv1alpha1.SuspendEtcdSpecReconcileAnnotation) {
		annotationKey = druidv1alpha1.SuspendEtcdSpecReconcileAnnotation
	} else if metav1.HasAnnotation(etcd.ObjectMeta, druidv1alpha1.IgnoreReconciliationAnnotation) {
		annotationKey = druidv1alpha1.IgnoreReconciliationAnnotation
	}
	if len(annotationKey) > 0 {
		r.recorder.Eventf(
			etcd,
			corev1.EventTypeWarning,
			"SpecReconciliationSkipped",
			"spec reconciliation of %s/%s is skipped by etcd-druid due to the presence of annotation %s on the etcd resource",
			etcd.Namespace,
			etcd.Name,
			annotationKey,
		)
		return ctrlutils.DoNotRequeue()
	}
	return ctrlutils.ContinueReconcile()
}

func (r *Reconciler) shouldReconcileSpec(etcd *druidv1alpha1.Etcd) bool {
	// TODO: replace v1beta1constants.GardenerOperation with own `druid.gardener.cloud/operation: reconcile`
	// and make gardener use druid's constant instead of druid use gardener's constant
	return r.config.IgnoreOperationAnnotation ||
		(!r.config.IgnoreOperationAnnotation && metav1.HasAnnotation(etcd.ObjectMeta, v1beta1constants.GardenerOperation))
}

package etcd

import (
	"context"
	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	ctrlutils "github.com/gardener/etcd-druid/internal/controller/utils"
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

const (
	// IgnoreReconciliationAnnotation is an annotation set by an operator in order to stop reconciliation.
	// TODO: move to api constants?
	IgnoreReconciliationAnnotation = "druid.gardener.cloud/ignore-reconciliation"
)

type Reconciler struct {
	client      client.Client
	config      *Config
	recorder    record.EventRecorder
	imageVector imagevector.ImageVector
	logger      logr.Logger
}

// NewReconciler creates a new reconciler for Etcd.
func NewReconciler(mgr manager.Manager, config *Config) (*Reconciler, error) {
	imageVector, err := ctrlutils.CreateImageVector()
	if err != nil {
		return nil, err
	}

	return &Reconciler{
		client:      mgr.GetClient(),
		config:      config,
		recorder:    mgr.GetEventRecorderFor(controllerName),
		imageVector: imageVector,
		logger:      log.Log.WithName(controllerName),
	}, nil
}

// TODO: where/how is this being used?
// +kubebuilder:rbac:groups=druid.gardener.cloud,resources=etcds,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=druid.gardener.cloud,resources=etcds/status,verbs=get;create;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts;services;configmaps,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;watch
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;get;list

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	/*
		If deletionTimestamp set:
			delete(); if err then requeue
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

	r.logger.WithValues("Etcd", req.NamespacedName)
	r.logger.Info("Etcd-controller reconciliation started")

	etcd := &druidv1alpha1.Etcd{}
	if err := r.client.Get(ctx, req.NamespacedName, etcd); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{Requeue: false}, nil
		}
		return ctrl.Result{Requeue: true}, err
	}

	if !etcd.DeletionTimestamp.IsZero() {
		if err := r.delete(ctx, etcd); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	}

	if metav1.HasAnnotation(etcd.ObjectMeta, IgnoreReconciliationAnnotation) {
		r.recorder.Eventf(
			etcd,
			corev1.EventTypeWarning,
			"ReconciliationIgnored",
			"reconciliation of %s/%s is ignored by etcd-druid due to the presence of annotation %s on the etcd resource",
			etcd.Namespace,
			etcd.Name,
			IgnoreReconciliationAnnotation,
		)
		return ctrl.Result{Requeue: false}, nil
	}

	var reconcileSpecErr error
	if r.shouldReconcileSpec(etcd) {
		reconcileSpecErr = r.reconcileSpec(ctx, etcd)
	}

	r.reconcileStatus(ctx, req.NamespacedName)

	if reconcileSpecErr != nil {
		return ctrl.Result{Requeue: true}, reconcileSpecErr
	}
	return ctrl.Result{RequeueAfter: r.config.EtcdStatusSyncPeriod}, nil
}

func (r *Reconciler) reconcileSpec(ctx context.Context, etcd *druidv1alpha1.Etcd) error {
	/*
		update status to reflect observedGeneration, lastError, etc (fields specific to spec reconciliation)
	*/

	return nil
}

func (r *Reconciler) reconcileStatus(ctx context.Context, etcdNamespacedName types.NamespacedName) {
	/*
		fetch EtcdMember resources
		fetch member leases
		status.condition checks
		status.members checks
		fetch latest Etcd resource
		update etcd status
	*/

}

func (r *Reconciler) delete(ctx context.Context, etcd *druidv1alpha1.Etcd) error {
	/*
		components [];
		for component in components:
			get component
			if component exists:
				delete component; if error then requeue
			if component does not exist, then record skip
		if all components have been recorded as non-existent, then remove finalizer and exit
	*/

	return nil
}

func (r *Reconciler) shouldReconcileSpec(etcd *druidv1alpha1.Etcd) bool {
	// TODO: replace v1beta1constants.GardenerOperation with own `druid.gardener.cloud/operation: reconcile`
	// and make gardener use druid's constant instead of druid use gardener's constant
	return r.config.IgnoreOperationAnnotation ||
		(!r.config.IgnoreOperationAnnotation && metav1.HasAnnotation(etcd.ObjectMeta, v1beta1constants.GardenerOperation))
}

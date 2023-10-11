package etcd

import (
	"context"
	ctrlutils "github.com/gardener/etcd-druid/internal/controller/utils"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/gardener/gardener/pkg/utils/imagevector"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	// IgnoreReconciliationAnnotation is an annotation set by an operator in order to stop reconciliation.
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
			Else If IgnoreOperationAnnotation flag is false:
				reconcileSpec()
				if err in getting etcd, return with requeue
				if err in deploying any of the components, then record pending requeue
		reconcileStatus()
		requeue after minimum of X seconds (EtcdStatusSyncPeriod) and previously recorded requeue request
	*/

	return ctrl.Result{}, nil
}

func (r *Reconciler) reconcileSpec() {
	/*

	 */
}

func (r *Reconciler) reconcileStatus() {
	/*

	 */
}

func (r *Reconciler) delete() {
	/*
		components [];
		for component in components:
			get component
			if component exists:
				delete component; if error then requeue
			if component does not exist, then record skip
		if all components have been recorded as non-existent, then remove finalizer and exit
	*/
}

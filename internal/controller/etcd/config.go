package etcd

import (
	"errors"
	"flag"
	"time"

	"github.com/gardener/etcd-druid/internal/controller/utils"
	"github.com/gardener/etcd-druid/internal/features"
	"k8s.io/component-base/featuregate"
)

// Flag names
const (
	workersFlagName                            = "etcd-workers"
	ignoreOperationAnnotationFlagName          = "ignore-operation-annotation"
	disableEtcdServiceAccountAutomountFlagName = "disable-etcd-serviceaccount-automount"
	etcdStatusSyncPeriodFlagName               = "etcd-status-sync-period"
)

const (
	defaultWorkers                            = 3
	defaultIgnoreOperationAnnotation          = false
	defaultDisableEtcdServiceAccountAutomount = false
	defaultEtcdStatusSyncPeriod               = 15 * time.Second
)

// featureList holds the feature gate names that are relevant for the Etcd Controller.
var featureList = []featuregate.Feature{
	features.UseEtcdWrapper,
}

// Config defines the configuration for the Etcd Controller.
type Config struct {
	// Workers is the number of workers concurrently processing reconciliation requests.
	Workers int
	// IgnoreOperationAnnotation specifies whether to ignore or honour the operation annotation on resources to be reconciled.
	// TODO: better name please, or deprecate and use new flag
	IgnoreOperationAnnotation bool
	// DisableEtcdServiceAccountAutomount controls the auto-mounting of service account token for ETCD StatefulSets.
	DisableEtcdServiceAccountAutomount bool
	// EtcdStatusSyncPeriod is the duration after which an event will be re-queued ensuring ETCD status synchronization.
	EtcdStatusSyncPeriod time.Duration
	// FeatureGates contains the feature gates to be used by Etcd Controller.
	FeatureGates map[featuregate.Feature]bool
}

// InitFromFlags initializes the config from the provided CLI flag set.
func InitFromFlags(fs *flag.FlagSet, cfg *Config) {
	fs.IntVar(&cfg.Workers, workersFlagName, defaultWorkers,
		"Number of workers spawned for concurrent reconciles of etcd spec and status changes. If not specified then default of 3 is assumed.")
	flag.BoolVar(&cfg.IgnoreOperationAnnotation, ignoreOperationAnnotationFlagName, defaultIgnoreOperationAnnotation,
		"Specifies whether to ignore or honour the operation annotation on resources to be reconciled.")
	fs.BoolVar(&cfg.DisableEtcdServiceAccountAutomount, disableEtcdServiceAccountAutomountFlagName, defaultDisableEtcdServiceAccountAutomount,
		"If true then .automountServiceAccountToken will be set to false for the ServiceAccount created for etcd StatefulSets.")
	fs.DurationVar(&cfg.EtcdStatusSyncPeriod, etcdStatusSyncPeriodFlagName, defaultEtcdStatusSyncPeriod,
		"Period after which an etcd status sync will be attempted.")
}

// Validate validates the config.
func (cfg *Config) Validate() error {
	var errs []error
	if err := utils.MustBeGreaterThan(workersFlagName, 0, cfg.Workers); err != nil {
		errs = append(errs, err)
	}
	if err := utils.MustBeGreaterThan(etcdStatusSyncPeriodFlagName, 0, cfg.EtcdStatusSyncPeriod); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

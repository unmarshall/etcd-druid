package monitor

import (
	"context"
	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"github.com/gardener/etcd-druid/internal/monitor/condition"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Prober probes the etcd cluster and runs probes to enquire conditions, etcd-member status, etc.
// Note: It should be ensured that no Prober implementation should mutate the etcd object.
type Prober[T any] interface {
	// Probe probes the etcd cluster and returns a specific part of the status based on the probe results.
	Probe(ctx context.Context, etcd *druidv1alpha1.Etcd) []T
}

func ProbeConditions(ctx context.Context, cl client.Client, etcd *druidv1alpha1.Etcd) []druidv1alpha1.Condition {
	probers := []Prober[druidv1alpha1.Condition]{
		condition.NewBackupReadyProber(cl),
		condition.NewDataVolumesReadyProber(cl),
	}
	conditions := make([]druidv1alpha1.Condition, 0, 4)
	for _, p := range probers {
		conditions = append(conditions, p.Probe(ctx, etcd)...)
	}
	return conditions
}

package condition

import (
	"context"
	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"github.com/gardener/etcd-druid/internal/monitor"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type _dataVolReadyProbe struct {
	cl client.Client
}

func NewDataVolumesReadyProber(cl client.Client) monitor.Prober[druidv1alpha1.Condition] {
	return &_dataVolReadyProbe{
		cl: cl,
	}
}

func (_ _dataVolReadyProbe) Probe(ctx context.Context, etcd *druidv1alpha1.Etcd) []druidv1alpha1.Condition {
	//TODO implement me
	panic("implement me")
}

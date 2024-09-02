package condition

import (
	"context"
	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"github.com/gardener/etcd-druid/internal/monitor"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type _etcdClusterReadyProbe struct {
	cl client.Client
}

func NewEtcdClusterReadyProber(client client.Client) monitor.Prober[druidv1alpha1.Condition] {
	return &_etcdClusterReadyProbe{
		cl: client,
	}
}

func (_ _etcdClusterReadyProbe) Probe(ctx context.Context, etcd *druidv1alpha1.Etcd) []druidv1alpha1.Condition {
	panic("implement me")
}

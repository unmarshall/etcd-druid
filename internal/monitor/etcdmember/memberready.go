package etcdmember

import (
	"context"
	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"github.com/gardener/etcd-druid/internal/monitor"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type _probe struct {
	cl client.Client
}

func NewMemberReadyProber(client client.Client) monitor.Prober[druidv1alpha1.EtcdMemberStatus] {
	return &_probe{
		cl: client,
	}
}

func (_ _probe) Probe(ctx context.Context, etcd *druidv1alpha1.Etcd) []druidv1alpha1.EtcdMemberStatus {
	//TODO implement me
	panic("implement me")
}

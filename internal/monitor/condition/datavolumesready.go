package condition

import (
	"context"
	"fmt"
	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"github.com/gardener/etcd-druid/internal/monitor"
	"github.com/gardener/etcd-druid/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	reasonStatefulSetNotFound      = "StatefulSetNotFound"
	reasonErrorFetchingStatefulSet = "ErrorFetchingStatefulSet"
	reasonErrorFetchingPVCEvents   = "ErrorFetchingPVCEvents"
	reasonWarningsFoundForPVCs     = "WarningsFoundForPVCs"
	reasonNoWarningsFoundForPVCs   = "NoWarningsFoundForPVCs"
)

type _dataVolReadyProbe struct {
	cl client.Client
}

func NewDataVolumesReadyProber(cl client.Client) monitor.Prober[druidv1alpha1.Condition] {
	return &_dataVolReadyProbe{
		cl: cl,
	}
}

func (p _dataVolReadyProbe) Probe(ctx context.Context, etcd *druidv1alpha1.Etcd) []druidv1alpha1.Condition {
	cond := druidv1alpha1.Condition{
		Type:   druidv1alpha1.ConditionTypeDataVolumesReady,
		Status: druidv1alpha1.ConditionUnknown,
	}
	stsName := druidv1alpha1.GetStatefulSetName(etcd.ObjectMeta)
	sts, err := utils.GetStatefulSet(ctx, p.cl, etcd)
	if err != nil {
		cond.Reason = reasonErrorFetchingStatefulSet
		cond.Message = fmt.Sprintf("Error fetching StatefulSet: %s: err: %s", stsName, err.Error())
		return []druidv1alpha1.Condition{cond}
	}
	if sts == nil {
		cond.Reason = reasonStatefulSetNotFound
		cond.Message = fmt.Sprintf("StatefulSet %s not found", stsName)
		return []druidv1alpha1.Condition{cond}
	}

	pvcWarnMsgs, err := utils.FetchPVCWarningMessagesForStatefulSet(ctx, p.cl, sts)
	if err != nil {
		cond.Reason = reasonErrorFetchingPVCEvents
		cond.Message = fmt.Sprintf("Error fetching PVC warning events for StatefulSet: %s: err: %s", stsName, err.Error())
		return []druidv1alpha1.Condition{cond}
	}
	if !utils.IsEmptyString(pvcWarnMsgs) {
		cond.Status = druidv1alpha1.ConditionFalse
		cond.Reason = reasonWarningsFoundForPVCs
		cond.Message = pvcWarnMsgs
		return []druidv1alpha1.Condition{cond}
	}

	cond.Status = druidv1alpha1.ConditionTrue
	cond.Reason = reasonNoWarningsFoundForPVCs
	cond.Message = fmt.Sprintf("No warning events found for PVCs used by StatefulSet %s", stsName)
	return []druidv1alpha1.Condition{cond}
}

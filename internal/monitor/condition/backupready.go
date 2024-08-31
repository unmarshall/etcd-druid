package condition

import (
	"context"
	"fmt"
	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"github.com/gardener/etcd-druid/internal/monitor"
	"github.com/gardener/etcd-druid/internal/utils"
	coordinationv1 "k8s.io/api/coordination/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// BackupReadyUnknown is a constant that means that the etcd backup status is currently not known
	BackupReadyUnknown    string = "Unknown"
	SnapshotLeaseNotFound string = "SnapshotLeaseNotFound"
)

type _backupReadyProbe struct {
	cl client.Client
}

func NewBackupReadyProber(client client.Client) monitor.Prober[druidv1alpha1.Condition] {
	return &_backupReadyProbe{
		cl: client,
	}
}

func (p *_backupReadyProbe) Probe(ctx context.Context, etcd *druidv1alpha1.Etcd) []druidv1alpha1.Condition {
	backupConditions := make([]druidv1alpha1.Condition, 0, 3)
	if !etcd.IsBackupStoreEnabled() {
		return backupConditions
	}
	backupConditions = append(backupConditions, p.probeFullSnapshotBackup(ctx, etcd))
	backupConditions = append(backupConditions, p.probeDeltaSnapshotBackup(ctx, etcd))
	backupConditions = append(backupConditions, deriveBackupReadyCondition(backupConditions))
	return backupConditions
}

func (p *_backupReadyProbe) probeFullSnapshotBackup(ctx context.Context, etcd *druidv1alpha1.Etcd) druidv1alpha1.Condition {
	fullSnapLease := &coordinationv1.Lease{}
	if err := p.cl.Get(ctx, client.ObjectKey{Name: druidv1alpha1.GetFullSnapshotLeaseName(etcd.ObjectMeta)}, fullSnapLease); err != nil {
		return createConditionOnError(err, etcd, druidv1alpha1.ConditionTypeFullSnapshotBackupReady)
	}
	panic("implement me")
}

func (p *_backupReadyProbe) probeDeltaSnapshotBackup(ctx context.Context, etcd *druidv1alpha1.Etcd) druidv1alpha1.Condition {
	panic("implement me")
}

// deriveBackupReadyCondition derives the backup ready condition from the constituent conditions (full and delta snapshot backup ready conditions).
// TODO: Remove this once BackupReady condition has been removed (wait for deprecation period).
func deriveBackupReadyCondition(constituentConditions []druidv1alpha1.Condition) druidv1alpha1.Condition {
	panic("implement me")
}

func createConditionOnError(err error, etcd *druidv1alpha1.Etcd, conditionType druidv1alpha1.ConditionType) druidv1alpha1.Condition {
	now := metav1.Now()
	errCondition := druidv1alpha1.Condition{
		Type:           conditionType,
		LastUpdateTime: now,
	}
	if apierrors.IsNotFound(err) {
		errCondition.Reason = SnapshotLeaseNotFound
		errCondition.Message = utils.IfConditionOr(conditionType == druidv1alpha1.ConditionTypeFullSnapshotBackupReady,
			"Full snapshot lease not found",
			"Delta snapshot lease not found")
	} else {
		errCondition.Reason = BackupReadyUnknown
		errCondition.Message = utils.IfConditionOr(conditionType == druidv1alpha1.ConditionTypeFullSnapshotBackupReady,
			fmt.Sprintf("Error getting full snapshot lease: %v", err),
			fmt.Sprintf("Error getting delta snapshot lease: %v", err))
	}
	if hasConditionChanged(etcd, errCondition) {
		errCondition.LastTransitionTime = now
	}
	return errCondition
}

func hasConditionChanged(etcd *druidv1alpha1.Etcd, newCondition druidv1alpha1.Condition) bool {
	if existingCondition := GetConditionByType(etcd, newCondition.Type); existingCondition != nil {
		return existingCondition.Status != newCondition.Status || existingCondition.Reason != newCondition.Reason
	}
	return true
}

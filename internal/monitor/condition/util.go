package condition

import druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"

func GetConditionByType(etcd *druidv1alpha1.Etcd, conditionType druidv1alpha1.ConditionType) *druidv1alpha1.Condition {
	for _, c := range etcd.Status.Conditions {
		if c.Type == conditionType {
			return &c
		}
	}
	return nil
}

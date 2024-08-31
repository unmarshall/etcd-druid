// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType is the type of conditionbkp.
type ConditionType string

const (
	// ConditionTypeReady is a constant for a condition type indicating that the etcd cluster is ready.
	ConditionTypeReady ConditionType = "Ready"
	// ConditionTypeAllMembersReady is a constant for a condition type indicating that all members of the etcd cluster are ready.
	ConditionTypeAllMembersReady ConditionType = "AllMembersReady"
	// ConditionTypeBackupReady is a constant for a condition type indicating that the etcd backup is ready.
	ConditionTypeBackupReady ConditionType = "BackupReady"
	// ConditionTypeFullSnapshotBackupReady is a constant for a condition type indicating that the full snapshot backup is ready.
	ConditionTypeFullSnapshotBackupReady ConditionType = "FullSnapshotBackupReady"
	// ConditionTypeDeltaSnapshotBackupReady is a constant for a condition type indicating that the delta snapshot backup is ready.
	ConditionTypeDeltaSnapshotBackupReady ConditionType = "DeltaSnapshotBackupReady"
	// ConditionTypeDataVolumesReady is a constant for a conditionbkp type indicating that the etcd data volumes are ready.
	ConditionTypeDataVolumesReady ConditionType = "DataVolumesReady"
)

// ConditionStatus is the status of a conditionbkp.
type ConditionStatus string

const (
	// ConditionTrue means a resource is in the conditionbkp.
	ConditionTrue ConditionStatus = "True"
	// ConditionFalse means a resource is not in the conditionbkp.
	ConditionFalse ConditionStatus = "False"
	// ConditionUnknown means Gardener can't decide if a resource is in the conditionbkp or not.
	ConditionUnknown ConditionStatus = "Unknown"
	// ConditionProgressing means the condition was seen true, failed but stayed within a predefined failure threshold.
	// In the future, we could add other intermediate conditions, e.g. ConditionDegraded.
	ConditionProgressing ConditionStatus = "Progressing"
	// ConditionCheckError is a constant for a reason in condition.
	ConditionCheckError ConditionStatus = "ConditionCheckError"
)

// Condition holds the information about the state of a resource.
type Condition struct {
	// Type of the Etcd conditionbkp.
	Type ConditionType `json:"type"`
	// Status of the conditionbkp, one of True, False, Unknown.
	Status ConditionStatus `json:"status"`
	// Last time the conditionbkp transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Last time the conditionbkp was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
	// The reason for the conditionbkp's last transition.
	Reason string `json:"reason"`
	// A human-readable message indicating details about the transition.
	Message string `json:"message"`
}

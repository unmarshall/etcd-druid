// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	corev1alpha1 "github.com/gardener/etcd-druid/api/core/v1alpha1"
	scheme "github.com/gardener/etcd-druid/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// EtcdListsGetter has a method to return a EtcdListInterface.
// A group's client should implement this interface.
type EtcdListsGetter interface {
	EtcdLists(namespace string) EtcdListInterface
}

// EtcdListInterface has methods to work with EtcdList resources.
type EtcdListInterface interface {
	Create(ctx context.Context, etcdList *corev1alpha1.EtcdList, opts v1.CreateOptions) (*corev1alpha1.EtcdList, error)
	Update(ctx context.Context, etcdList *corev1alpha1.EtcdList, opts v1.UpdateOptions) (*corev1alpha1.EtcdList, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*corev1alpha1.EtcdList, error)
	List(ctx context.Context, opts v1.ListOptions) (*corev1alpha1.EtcdListList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1alpha1.EtcdList, err error)
	EtcdListExpansion
}

// etcdLists implements EtcdListInterface
type etcdLists struct {
	*gentype.ClientWithList[*corev1alpha1.EtcdList, *corev1alpha1.EtcdListList]
}

// newEtcdLists returns a EtcdLists
func newEtcdLists(c *DruidV1alpha1Client, namespace string) *etcdLists {
	return &etcdLists{
		gentype.NewClientWithList[*corev1alpha1.EtcdList, *corev1alpha1.EtcdListList](
			"etcdlists",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *corev1alpha1.EtcdList { return &corev1alpha1.EtcdList{} },
			func() *corev1alpha1.EtcdListList { return &corev1alpha1.EtcdListList{} },
		),
	}
}

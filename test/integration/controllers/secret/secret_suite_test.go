// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"testing"

	druidconfigv1alpha1 "github.com/gardener/etcd-druid/api/config/v1alpha1"
	"github.com/gardener/etcd-druid/internal/controller/secret"
	"github.com/gardener/etcd-druid/test/integration/controllers/assets"
	"github.com/gardener/etcd-druid/test/integration/setup"

	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	intTestEnv *setup.IntegrationTestEnv
	k8sClient  client.Client
	namespace  string
)

const (
	testNamespacePrefix = "secret-"
)

func TestSecretController(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(
		t,
		"Secret Controller Suite",
	)
}

var _ = BeforeSuite(func() {
	crdPaths := []string{assets.GetEtcdCrdPath()}
	intTestEnv = setup.NewIntegrationTestEnv(testNamespacePrefix, "secret-int-tests", crdPaths)
	intTestEnv.RegisterReconcilers(func(mgr manager.Manager) {
		reconciler := secret.NewReconciler(mgr, druidconfigv1alpha1.SecretControllerConfiguration{
			ConcurrentSyncs: ptr.To(5),
		})
		Expect(reconciler.RegisterWithManager(context.TODO(), mgr)).To(Succeed())
	}).StartManager()
	k8sClient = intTestEnv.K8sClient
	namespace = intTestEnv.TestNs.Name
})

package main

import (
	"context"
	"fmt"
	"time"

	druidv1alpha1 "github.com/gardener/etcd-druid/api/v1alpha1"
	"gopkg.in/yaml.v3"
	coordinationv1 "k8s.io/api/coordination/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type etcdConfig struct {
	Name                    string                       `yaml:"name"`
	DataDir                 string                       `yaml:"data-dir"`
	Metrics                 druidv1alpha1.MetricsLevel   `yaml:"metrics"`
	SnapshotCount           int                          `yaml:"snapshot-count"`
	EnableV2                bool                         `yaml:"enable-v2"`
	QuotaBackendBytes       resource.Quantity            `yaml:"quota-backend-bytes"`
	AutoCompactionMode      druidv1alpha1.CompactionMode `yaml:"auto-compaction-mode"`
	AutoCompactionRetention time.Duration                `yaml:"auto-compaction-retention"`
}

func main() {
	//testDeleteAllOf()
	ctx := context.Background()
	testWithCtx1(ctx)
	testWithCtx2(ctx)
	testWithCtx3(ctx)
	fmt.Printf("context result: %s\n", ctx.Value("result"))
}

func testWithCtx1(ctx context.Context) {
	fmt.Printf("testWithCtx1 called")
	ctx = context.WithValue(ctx, "prevFuncName", "testWithCtx1")
}

func testWithCtx2(ctx context.Context) {
	fmt.Printf("testWithCtx2 called")
	val := ctx.Value("prevFuncName")
	fmt.Printf("value prevFuncName: %s\n", val)
	ctx = context.WithValue(ctx, "prevFuncName", "testWithCtx2")
}

func testWithCtx3(ctx context.Context) {
	fmt.Printf("testWithCtx3 called")
	val := ctx.Value("prevFuncName")
	fmt.Printf("value prevFuncName: %s\n", val)
	ctx = context.WithValue(ctx, "result", "success")
}

func testDeleteAllOf() {
	ctx := context.Background()
	config, err := clientcmd.BuildConfigFromFlags("", "/var/folders/p1/4ggwp2hn69n7b4ldt0kfvbch0000gn/T/garden/34e938c2-f22b-4ba8-8cbc-2377c4b66907/kubeconfig.yaml")
	if err != nil {
		fmt.Printf("error getting config: %v\n", err)
		return
	}
	cli, err := client.New(config, client.Options{})
	if err != nil {
		fmt.Printf("error creating client: %v\n", err)
		return
	}
	err = cli.DeleteAllOf(ctx, &coordinationv1.Lease{}, client.InNamespace("default"), client.MatchingLabels{"bingo": "tringo"})
	if err != nil {
		fmt.Printf("error deleting lease: %v\n", err)
		return
	}
}

func testMarshal() {
	cfg := &etcdConfig{
		Name:                    "bingo",
		DataDir:                 "/var/data/bingo",
		Metrics:                 druidv1alpha1.Basic,
		SnapshotCount:           10,
		EnableV2:                false,
		QuotaBackendBytes:       *resource.NewQuantity(5*1024*1024*1024, resource.BinarySI),
		AutoCompactionMode:      druidv1alpha1.Periodic,
		AutoCompactionRetention: 30 * time.Second,
	}
	cfgYaml, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Printf("error in marshalling etcd config %v\n", err)
		return
	}
	fmt.Printf("marshalled yaml: \n%s\n", cfgYaml)
}

package etcd

import (
	"github.com/aws/aws-sdk-go/aws/client"
)

type reconciler struct {
	client client.Client
	config *Config
}

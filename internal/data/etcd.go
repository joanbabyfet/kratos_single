package data

import (
	"kratos_single/internal/conf"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewEtcdClient(c *conf.Data) (*clientv3.Client, error) {

	return clientv3.New(clientv3.Config{
		Endpoints:   c.Etcd.Endpoints,
		DialTimeout: c.Etcd.DialTimeout.AsDuration(),
	})
}
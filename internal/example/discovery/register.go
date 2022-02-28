package discovery

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"

	"google.golang.org/grpc"
)

var dis *Discovery

type Discovery struct {
	cli *clientv3.Client
	em  endpoints.Manager
}

func init() {
	dis, _ = New("")
}

func New(etcdaddress string) (*Discovery, error) {
	d := &Discovery{}
	var err error
	if etcdaddress == "" {
		etcdaddress = "http://localhost:2379"
	}
	d.cli, err = clientv3.NewFromURL(etcdaddress)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func Register(ctx context.Context, serviceID, instanceID, endpoint string) error {
	var err error
	if dis.em == nil {
		if dis.em, err = endpoints.NewManager(dis.cli, serviceID); err != nil {
			return err
		}
	}
	lease := clientv3.NewLease(dis.cli)
	tick, err := lease.Grant(ctx, 30)
	if err != nil {
		return err
	}
	lease.KeepAlive(ctx, tick.ID)
	return dis.em.AddEndpoint(ctx, instanceID, endpoints.Endpoint{Addr: endpoint}, clientv3.WithLease(tick.ID))
}

func DeleteRegister(ctx context.Context, instanceID string) error {
	if dis.em != nil {
		return dis.em.DeleteEndpoint(ctx, instanceID)
	}
	return nil

}

func NewConn(serviceID string) (*grpc.ClientConn, error) {
	resolver, err := resolver.NewBuilder(dis.cli)
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("etcd:///%s", serviceID), grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), grpc.WithInsecure(), grpc.WithResolvers(resolver))
}

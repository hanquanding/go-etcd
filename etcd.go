package getcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

type ETCD struct {
	cfg clientv3.Config
	cli *clientv3.Client
	kvc *clientv3.KV
}

func New(endpoints ...string) (*ETCD, error) {
	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	etcd := &ETCD{
		cfg: cfg,
		cli: cli,
	}

	return etcd, nil
}

func (etcd *ETCD) Set(name, value string) (*clientv3.PutResponse, error) {
	return etcd.cli.Put(context.Background(), name, value)
}

func (etcd *ETCD) SetWithTTL(name, value string, ttl int64) (*clientv3.PutResponse, error) {
	lease, _ := etcd.cli.Grant(context.TODO(), ttl)
	return etcd.cli.Put(context.Background(), name, value, clientv3.WithLease(lease.ID))
}

func (etcd *ETCD) GetValue(name string) (string, error) {
	resp, err := etcd.cli.Get(context.Background(), name, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}

	if resp.Count == 0 {
		return "", fmt.Errorf("not found service by name [%s]", name)
	}
	return string(resp.Kvs[0].Value), nil
}

func (etcd *ETCD) Get(name string) (*clientv3.GetResponse, error) {
	return etcd.cli.Get(context.Background(), name)
}

func (etcd *ETCD) GetPrefix(name string) (*clientv3.GetResponse, error) {
	return etcd.cli.Get(context.Background(), name, clientv3.WithPrefix())
}

func (etcd *ETCD) GetPrefixFirstValue(name string) (string, error) {
	resp, err := etcd.cli.Get(context.Background(), name, clientv3.WithPrefix())
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", fmt.Errorf("not found service by name [%s]", name)
	}
	return string(resp.Kvs[0].Value), nil
}

func (etcd *ETCD) Remove(name string) (*clientv3.DeleteResponse, error) {
	return etcd.cli.Delete(context.Background(), name)
}

func (etcd *ETCD) RemovePrefix(name string) (*clientv3.DeleteResponse, error) {
	return etcd.cli.Delete(context.Background(), name, clientv3.WithPrefix())
}

func (etcd *ETCD) NewMutex(name string) (*concurrency.Mutex, error) {
	key := fmt.Sprintf("/lock/%s", name)
	session, err := concurrency.NewSession(etcd.cli)
	if err != nil {
		return nil, err
	}
	mutex := concurrency.NewMutex(session, key)
	return mutex, nil
}

func (etcd *ETCD) Close() error {
	return etcd.cli.Close()
}

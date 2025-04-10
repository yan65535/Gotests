package config

import (
	"context"
	"errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"log"
	"time"
)

type etcdClient struct {
	client      *clientv3.Client
	configKey   string
	lastVersion int64
	lastConf    Conf
}

func newEtcdClient() (*etcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdHosts,
		DialTimeout: time.Duration(5) * time.Second,
	})
	if err != nil {

		return nil, err
	}
	return &etcdClient{
		client:    client,
		configKey: etcdKey,
	}, nil
}

func (cc *etcdClient) fetchData() (kvData *mvccpb.KeyValue, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	resp, err := cc.client.Get(ctx, cc.configKey)
	cancel()

	if resp.Count == 0 {
		err = errors.New("etcd client is empty")
		return nil, err
	}
	return resp.Kvs[0], nil
}
func (cc *etcdClient) GetRemoteConfig() (conf Conf, err error) {
	data, err := cc.fetchData()
	if err != nil {
		log.Println(err)
	}

	if data.Version < cc.lastVersion {
		return cc.lastConf, nil
	}
	cc.handleUpdateConfig(data)
	return cc.lastConf, nil
}

func (cc *etcdClient) updateConf(conf Conf) {
	cc.lastConf = conf
}

func (cc *etcdClient) decodeConf(data []byte) (conf Conf, err error) {
	if err = yaml.Unmarshal(data, &conf); err != nil {
		log.Println(err)
	}
	return conf, err
}

func (cc *etcdClient) WatchConfig() {
	var err error
	for {
		watchChan := cc.client.Watch(context.Background(), cc.configKey)
		for w := range watchChan {
			err = w.Err()
			if err != nil {
				log.Println(err)
			}
			for _, event := range w.Events {
				if string(event.Kv.Key) == cc.configKey {
					cc.handleUpdateConfig(event.Kv)
					break
				}
			}
		}
	}
}

func (cc *etcdClient) handleUpdateConfig(kvData *mvccpb.KeyValue) {
	var conf Conf
	if kvData.Version <= cc.lastVersion {
		return
	}
	conf, err := cc.decodeConf(kvData.Value)
	if err != nil {
		zap.L().Error("decode conf failed", zap.Error(err))
		return
	}

	cc.updateConf(conf)
}

package config

import (
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
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
	if err != nil {
		trace.GlobalLogger.For(context.TODO()).Error("fetchData get config failed", zap.Error(err))
	}
	if resp.Count == 0 {
		err = errors.New("etcd client is empty")
		return nil, err
	}
	return resp.Kvs[0], nil
}
func (cc *etcdClient) GetRemoteConfig() (conf Conf, err error) {
	data, err := cc.fetchData()
	if err != nil {
		trace.GlobalLogger.For(context.TODO()).Error("GetRemoteConfig fetchData config failed", zap.Error(err))
		return conf, err
	}
	if data == nil {
		trace.GlobalLogger.For(context.TODO()).Error("GetRemoteConfig fetchData data is nil")
		return conf, errors.New("got data is nil")
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
		trace.GlobalLogger.For(context.TODO()).Error("watchSystemConfig failed", zap.Error(err), zap.String("val", string(data)))
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
				trace.GlobalLogger.For(context.TODO()).Error("WatchConfig failed", zap.Error(err), zap.String("configKey", cc.configKey))
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
		trace.GlobalLogger.For(context.TODO()).Error("watchSystemConfig failed", zap.Error(err), zap.String("val", string(kvData.Value)))
		return
	}
	for _, update := range updateFuncs {
		update(cc.lastConf, conf)
	}

	cc.updateConf(conf)
}

func RegisterUpdateFunc(updateFunc ...UpdateConfFunc) {
	updateFuncs = updateFunc
}

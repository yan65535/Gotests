package config

import "fmt"

const (
	ServerName = "latex"
)

var (
	etcdKey         = fmt.Sprintf("/configs/%s/system", ServerName)
	etcdAddr        string
	etcdHosts       []string
	localConfigPath string
)

type Conf struct {
	Server    ServerConf    `yaml:"server"`
	Etcd      EtcdConf      `yaml:"etcd"`
	Zookeeper ZookeeperConf `yaml:"zookeeper"`
	RpcServer RpcServerConf `yaml:"rpcServer"`
	RpcClient RpcClientConf `yaml:"rpcClient"`
}

type ServerConf struct {
}

type EtcdConf struct {
}

type ZookeeperConf struct {
}

type RpcServerConf struct {
}

type RpcClientConf struct {
}

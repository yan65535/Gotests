package config

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	ServerName     = "basic"
	ServerFullName = "latexOcrService"
)

var (
	etcdKey         = fmt.Sprintf("/configs/%s/system", ServerFullName)
	etcdAddr        string
	etcdHosts       []string
	localConfigPath string
	zookeeperHost   string
)

// 定义配置结构体
type Conf struct {
	Port   string `yaml:"port"`
	Server struct {
		Env     string  `yaml:"env"`
		RpcConf RpcConf `yaml:"rpc_conf"`
	} `yaml:"server"`
	// 可以添加更多配置项
}

type RpcConf struct {
	Name      string   `yaml:"name"`
	ListenOn  string   `yaml:"listen_on"`
	EtcdHosts []string `yaml:"etcd_hosts"`
	EtcdKey   string   `yaml:"etcd_key"`
}

type RedisConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func init() {
	flag.StringVar(&localConfigPath, "c", ServerName+"_local.yml", "default config path")
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_HOSTS"), "default etcd address")
	flag.StringVar(&zookeeperHost, "z", os.Getenv("ZOOKEEPER_HOSTS"), "default zookeeper address")
	flag.Parse()
}

// InitConfig 初始化配置，解析 YAML 文件
func Init() (conf Conf, err error) {
	flag.Parse()
	//etcdAddr = "127.0.0.1:2379"  for test
	if etcdAddr != "" {
		etcdHosts = strings.Split(etcdAddr, ",")
		conf, err = getFromRemoteAndWatchUpdate()
	} else {
		conf, err = getFromLocal()
	}

	fmt.Printf("static config => [%#v]\n", conf)
	return
}

func getFromRemoteAndWatchUpdate() (conf Conf, err error) {
	client, err := newEtcdClient()
	if err != nil {
		return conf, err
	}
	conf, err = client.GetRemoteConfig()
	if err != nil {
		return conf, err
	}
	go client.WatchConfig()
	return conf, err
}

func getFromLocal() (confObj Conf, err error) {
	content, err := os.ReadFile(localConfigPath)
	if err != nil {
		return confObj, err
	}

	err = yaml.Unmarshal(content, &confObj)
	if err != nil {
		return confObj, err
	}
	err = conf.LoadFromYamlBytes(content, &confObj)
	if err != nil {
		return Conf{}, err
	}
	return confObj, err
}

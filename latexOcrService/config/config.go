package config

type Conf struct {
	Port   string
	Server struct {
		Env       string
		RpcConf   RpcConf
		RedisConf RedisConf
	}
}

type Server struct {
}

type RpcConf struct {
	Name      string
	ListenOn  string
	EtcdHosts []string
	EtcdKey   string
}

type RedisConf struct {
	Host string
	Port string
}

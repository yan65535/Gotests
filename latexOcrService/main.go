package latexOcrService

import (
	"github.com/samuel/go-zookeeper/zk"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"latexOcrService/adaptor"
	"latexOcrService/config"
	latexOcr "latexOcrService/latexOcrProto"
	"latexOcrService/rpc"
	"time"
)

//tkRpc "e.coding.net/g-mneg1542/block/block.proto/token/protoc/token"
//"e.coding.net/g-mneg1542/block/block.utils/utrace/tredis"
//trace "e.coding.net/g-mneg1542/block/block.utils/utrace/tutil"

func main() {
	conf := initConfig()

	// 初始化Zookeeper连接
	zkConn := initZookeeper(conf)
	defer zkConn.Close()

	// 初始化gRPC服务器
	grpcServer := initGrpcServer(conf, adaptor.Adaptors{})

	// 注册服务到Zookeeper
	registerToZookeeper(zkConn, conf.Zookeeper.ServicePath, conf.Server.RpcConf.ListenOn)

	// 启动服务
	grpcServer.Start()

}

func initConfig() config.Conf {
	conf, err := config.Init()
	handleErr(err)

	return conf
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initGrpcServer(conf config.Conf, adaptor adaptor.Adaptors) *zrpc.RpcServer {
	rpcConf := zrpc.RpcServerConf{
		ServiceConf: service.ServiceConf{
			Name: conf.Server.RpcConf.Name,
			Mode: conf.Server.Env,
		},
		ListenOn: conf.Server.RpcConf.ListenOn,
		Etcd: discov.EtcdConf{
			Hosts: conf.Server.RpcConf.EtcdHosts,
			Key:   conf.Server.RpcConf.EtcdKey,
		},
		Middlewares: zrpc.ServerMiddlewaresConf{
			Trace:      true,
			Recover:    true,
			Stat:       true,
			Prometheus: true,
		},
	}

	return zrpc.MustNewServer(rpcConf, func(grpcServer *grpc.Server) {
		latexOcr.RegisterLatexServiceServer(grpcServer, rpc.NewOcrGrpcService(adaptor, conf))

		if rpcConf.Mode == service.DevMode || rpcConf.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
}

//func initApp(conf config.Conf, router *route.Router) *uweb.App {
//
//	app := uweb.NewApp(conf.Server.HttpPort, "", tracer, router)
//	return app
//}

// 初始化Zookeeper连接
func initZookeeper(conf config.Conf) *zk.Conn {
	conn, _, err := zk.Connect(
		conf.Zookeeper.Hosts,
		time.Duration(conf.Zookeeper.Timeout)*time.Second,
	)
	handleErr(err)
	return conn
}

// Zookeeper服务注册
func registerToZookeeper(conn *zk.Conn, servicePath string, listenOn string) {
	// 创建持久节点（如果不存在）
	exists, _, err := conn.Exists(servicePath)
	handleErr(err)

	if !exists {
		_, err = conn.Create(servicePath, []byte{}, 0, zk.WorldACL(zk.PermAll))
		handleErr(err)
	}

	// 创建临时顺序节点
	nodePath := servicePath + "/node-"
	_, err = conn.CreateProtectedEphemeralSequential(
		nodePath,
		[]byte(listenOn),
		zk.WorldACL(zk.PermAll))
	handleErr(err)
}

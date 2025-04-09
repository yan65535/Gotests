package latexOcrService

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"latexOcrService/adaptor"
	"latexOcrService/config"
	latexOcr "latexOcrService/latexOcrProto"
	"latexOcrService/rpc"
)

//tkRpc "e.coding.net/g-mneg1542/block/block.proto/token/protoc/token"
//"e.coding.net/g-mneg1542/block/block.utils/utrace/tredis"
//trace "e.coding.net/g-mneg1542/block/block.utils/utrace/tutil"

func main() {

}

func initConfig() config.Conf {
	conf, err := config.Init()
	handleErr(err)
	config.RegisterUpdateFunc(config.GetUpdateLogLevelFun())
	return conf
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

package route

import (
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"latexOcrService/adaptor"
	"latexOcrService/config"
	"strings"
)

// 定义500错误处理函数
func err500(ctx iris.Context) {
	_, _ = ctx.WriteString("CUSTOM 500 ERROR")
}

// 定义404错误处理函数
func err404(ctx iris.Context) {
	_, _ = ctx.WriteString("CUSTOM 404 ERROR")
}

type Router struct {
	FullPPROF   bool
	rootPath    string
	conf        config.Conf
	checkFunc   func() error
	adaptors    adaptor.Adaptors
	controllers *Controller
}

func NewRouter(conf config.Conf, adaptors adaptor.Adaptors, checkFunc func() error) *Router {
	rootPath := "/api/content"
	return &Router{
		rootPath:    rootPath,
		conf:        conf,
		checkFunc:   checkFunc,
		adaptors:    adaptors,
		controllers: NewController(adaptors, conf),
	}
}

func (r *Router) Register(app *iris.Application) {
	app.OnErrorCode(iris.StatusInternalServerError, err500)
	app.OnErrorCode(iris.StatusNotFound, err404)
	root := app.Party(r.rootPath)
	r.route(root)
	// 注入 pprof
	// 加入 prometheus
	root.Get("/metrics", iris.FromStd(promhttp.Handler()))
	// 加入 ping检查接口
	root.Any("/ping", func(ctx iris.Context) {
		msg := "success"
		if err := r.checkFunc(); err != nil {
			msg = err.Error()
		}
		_, _ = ctx.WriteString(msg)
	})

}
func (r *Router) route(app *iris.Application) {
	//todo
}

//白名单暂不需设置
//func (r *Router) SpanFilter(req *http.Request) bool {
//	path := strings.ReplaceAll(req.URL.Path, r.rootPath, "")
//
//	checkMap := getNoSpanList()
//	if _, ok := checkMap[path]; !ok {
//		return false
//	}
//	return true
//}

//b
//func (r *Router) AccessRecordFilter(ctx iris.Context) bool {
//	path := strings.ReplaceAll(ctx.Path(), r.rootPath, "")
//	checkMap := getNoAccessLogList()
//	if _, ok := checkMap[path]; !ok {
//		return false
//	}
//	return true
//}

func (r *Router) AccessRecordFilter(ctx iris.Context) bool {
	path := strings.ReplaceAll(ctx.Path(), r.rootPath, "")
	checkMap := getNoAccessLogList()
	if _, ok := checkMap[path]; !ok {
		return false
	}
	return true
}

func (r *Router) adminRoute(adminRoot iris.Party) {
	adminRoot.Post("/v1/basic/toTex", r.controllers.basic.ToTex)

}

func (r *Router) customerRoute(cstRoot iris.Party) {

}

func (r *Router) businessRoute(bizRoot iris.Party) {

}

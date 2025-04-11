package v2

import (
	"context"
	sentryiris "github.com/getsentry/sentry-go/iris"
	"github.com/kataras/iris/v12"
	"github.com/opentracing/opentracing-go"
	"net/http"
	"strconv"
	"time"

	"e.coding.net/g-mneg1542/block/block.utils/utrace/tweb"
	"e.coding.net/g-mneg1542/block/block.utils/uweb"
)

type IRouter interface {
	Register(app *iris.Application)
	SpanFilter(r *http.Request) bool
	AccessRecordFilter(ctx iris.Context) bool
}

type App struct {
	server *iris.Application
	addr   iris.Runner
	config iris.Configurator
}

func NewApp(port int, sentryDsn string, tracer opentracing.Tracer, router IRouter) *App {
	app := iris.New()
	if sentryDsn != "" {
		app.Use(sentryiris.New(sentryiris.Options{})) // sentry中间件会对panic进行recover，因此如果在uweb.NewRecoverMdw后面，uweb.NewRecoverMdw就捕获不到异常
	}
	// 全局异常
	app.Use(uweb.NewRecoverMdw())
	// 全局链路追踪
	app.Use(tweb.Middleware(tracer, tweb.MWSpanFilter(router.SpanFilter)))
	// 全局日志
	app.Use(uweb.NewAccessLogMdw(router.AccessRecordFilter))
	// 注册路由
	router.Register(app)

	// server配置
	config := iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           true,
		EnablePathEscape:                  false,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: true,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "2006-01-02 15:04:05",
		Charset:                           "UTF-8",
		IgnoreServerErrors:                []string{iris.ErrServerClosed.Error()},
		RemoteAddrHeaders:                 []string{"X-Real-Ip", "X-Forwarded-For"},
	})
	return &App{
		server: app,
		addr:   iris.Addr(":" + strconv.Itoa(port)),
		config: config,
	}
}

func (app *App) Run() {
	_ = app.server.Run(app.addr, app.config)
}

func (app *App) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_ = app.server.Shutdown(ctx)
}

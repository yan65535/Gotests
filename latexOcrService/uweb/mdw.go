package uweb

import (
	"context"
	"e.coding.net/g-mneg1542/block/block.utils/ulog"
	trace "e.coding.net/g-mneg1542/block/block.utils/utrace/tutil"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/gogf/gf/util/gconv"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
	"strings"
	"time"
)

// 获取iris的请求体
func RequestBody(ctx iris.Context) string {
	body, err := ctx.GetBody()
	if err != nil {
		return err.Error()
	}
	return string(body)
}

// 获取iris的get参数
func RequestQueries(ctx iris.Context) string {
	var requestQuery string
	for k, v := range ctx.URLParams() {
		requestQuery += k + "=" + v + "&"
	}
	requestQuery = strings.Trim(requestQuery, "&")

	return requestQuery
}

// 统一异常处理
// 依赖：zlog.LogSystem、sentry
func NewRecoverMdw() iris.Handler {
	return func(ctx iris.Context) {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 3)
				if ctx.IsStopped() {
					return
				}

				body, _ := ctx.GetBody()
				fields := []zap.Field{
					zap.String("Recovered from a route's Handler", ctx.HandlerName()),
					zap.Int("status code", ctx.GetStatusCode()),
					zap.String("path", ctx.Path()),
					zap.String("method", ctx.Method()),
					zap.Any("path param", ctx.URLParams()),
					zap.ByteString("body param", body),
					zap.String("remote addr", ctx.RemoteAddr()),
					zap.Stack("stack"),
					zap.Any("err", err),
				}
				if trace.GlobalLogger.HasInit() {
					trace.GlobalLogger.For(ctx.Request().Context()).Error("recovered", fields...)
				} else {
					ulog.Error("recovered", fields...)
				}

				ctx.StatusCode(500)
				ctx.StopExecution()
			}
		}()

		ctx.Next()
	}
}

// 请求日志记录
// 依赖：zlog.LogAccess
func NewAccessLogMdw(filters ...func(ctx iris.Context) bool) iris.Handler {
	return func(ctx iris.Context) {
		for _, filter := range filters {
			if filter(ctx) {
				ctx.Next()
				return
			}
		}

		begin := time.Now()
		fields := []zap.Field{
			zap.String("ip", ctx.RemoteAddr()),
			zap.String("method", ctx.Method()),
			zap.String("path", ctx.Path()),
		}
		if ctx.Method() == iris.MethodPost {
			fields = append(fields, zap.String("req_body", RequestBody(ctx)))
		} else if ctx.Method() == iris.MethodGet {
			fields = append(fields, zap.Any("params", RequestQueries(ctx)))
		}
		if trace.GlobalLogger.HasInit() {
			trace.GlobalLogger.For(ctx.Request().Context()).With(fields...).Info("api begin")
		} else {
			ulog.Info("api begin", fields...)
		}

		recorder := ctx.Recorder()

		ctx.Next()

		b := recorder.Body()
		if len(b) > 1024 {
			b = b[:1024]
		}
		fields = []zap.Field{}
		fields = append(fields, zap.Duration("duration", time.Since(begin)))
		if ctx.Method() == iris.MethodPost {
			fields = append(fields, zap.ByteString("resp_body", b))
		} else if ctx.Method() == iris.MethodGet {
			fields = append(fields, zap.ByteString("resp_body", b))
		}
		if trace.GlobalLogger.HasInit() {
			trace.GlobalLogger.For(ctx.Request().Context()).With(fields...).Info("api end")
		} else {
			ulog.Info("api end", fields...)
		}
	}
}

func checkPass(noCheckPaths map[string]bool, rootPath, path string) bool {
	path = strings.ReplaceAll(path, rootPath, "")
	if v, ok := noCheckPaths[path]; ok && v {
		return true
	}
	return false
}

// 用户信息
type User struct {
	ID         int64  `json:"id"` // 全局user_id
	NickName   string `json:"nick_name"`
	MobileBind int32  `json:"mobile_bind"` // 默认0：未绑定 1：已绑定
	WechatBind int32  `json:"wechat_bind"` // 默认0：未绑定 1：已绑定
}
type AppUser struct {
	UserID  int64  `json:"user_id"`  // 全局user_id
	AppCode int32  `json:"app_code"` // 1000=公众号，1001=小程序， 其他进行扩展，固定值不能变更
	OpenID  string `json:"open_id"`  // 对应微信应用下的openid
}

type MobileUser struct {
	UserID int64  `json:"user_id"` // 全局user_id
	Mobile string `json:"mobile"`  // 手机号码
}

type WeChatUser struct {
	UserID        int64  `json:"user_id"`         // 全局user_id
	UnionID       string `json:"union_id"`        // 微信唯一用户ID
	NickName      string `json:"nick_name"`       // 微信昵称
	WechatIconURL string `json:"wechat_icon_url"` // 128长度 微信头像
}

type UserProfile struct {
	UserID         int64  `json:"user_id"`          // 全局user_id
	StudentNo      string `json:"student_no"`       // 学号
	SchoolSystemID int64  `json:"school_system_id"` // 系统内的学校ID
	GradeEnum      int32  `json:"grade_enum"`       // 1--12分别为一年级到高三
	ClassName      string `json:"class_name"`
}

type ParentUser struct {
	UserID       int64 `json:"user_id"`
	ParentUserID int64 `json:"parent_user_id"`
}

type AuthorUser struct {
	UserID     int64  `json:"user_id"`
	AuthorID   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
}

type UserInfo struct {
	User        User        `json:"user"`
	MobileUser  MobileUser  `json:"mobile_user"`
	WechatUser  WeChatUser  `json:"wechat_user"`
	UserProfile UserProfile `json:"user_profile"`
	ParentUser  ParentUser  `json:"parent_user"`
	AuthorUser  AuthorUser  `json:"author_user"`
	AppUsers    []AppUser   `json:"app_users"`
}

type AdminUser struct {
	UserID int64  `json:"id"`
	Name   string `json:"name"`
}

const (
	TokenKey    = "token"
	UserInfoKey = "user_info"
)

type TokenFun func(ctx context.Context, token string) (string, error)

// 统一token检查处理
func NewAuthorityMdwV2(noCheckPaths map[string]bool, rootPath string, getTokenFun TokenFun) iris.Handler {
	return func(ctx iris.Context) {
		if checkPass(noCheckPaths, rootPath, ctx.Path()) {
			ctx.Next()
			return
		}

		token := ctx.GetHeader(TokenKey)
		if len(token) == 0 {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.StopExecution()
		} else {
			tokenStr, err := getTokenFun(ctx.Request().Context(), token)
			if err != nil {
				trace.GlobalLogger.For(ctx.Request().Context()).Error("NewAuthorityMdwV2 getTokenFun error", zap.Error(err), zap.String("token", token))
			} else {
				doUser := &UserInfo{}
				err = json.Unmarshal(gconv.Bytes(tokenStr), doUser)
				if err != nil {
					trace.GlobalLogger.For(ctx).Error("NewAuthorityMdwV2 Unmarshal", zap.Any("token", token), zap.Error(err))
				} else {
					ctx.Values().Set(UserInfoKey, doUser)
				}
			}
			if err != nil {
				ctx.StatusCode(iris.StatusUnauthorized)
				ctx.StopExecution()
			}
		}
		ctx.Next()
	}
}

// 管理后台-统一token检查处理
func NewAuthorityAdminMdw(noCheckPaths map[string]bool, rootPath string, getTokenFun TokenFun) iris.Handler {
	return func(ctx iris.Context) {
		if checkPass(noCheckPaths, rootPath, ctx.Path()) {
			ctx.Next()
			return
		}

		token := ctx.GetHeader(TokenKey)
		if len(token) == 0 {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.StopExecution()
		} else {
			tokenStr, err := getTokenFun(ctx.Request().Context(), token)
			if err != nil {
				trace.GlobalLogger.For(ctx.Request().Context()).Error("NewAuthorityAdminMdw getTokenFun error", zap.Error(err), zap.String("token", token))
			} else {
				doUser := &AdminUser{}
				err = json.Unmarshal(gconv.Bytes(tokenStr), doUser)
				if err != nil {
					trace.GlobalLogger.For(ctx).Error("NewAuthorityAdminMdw Unmarshal", zap.Any("token", token), zap.Error(err))
				} else {
					ctx.Values().Set(UserInfoKey, doUser)
				}
			}
			if err != nil {
				ctx.StatusCode(iris.StatusUnauthorized)
				ctx.StopExecution()
			}
		}
		ctx.Next()
	}
}

func GetUserInfoFromCtx(ctx iris.Context) *UserInfo {
	get := ctx.Values().Get(UserInfoKey)
	if get == nil {
		return nil
	}
	return get.(*UserInfo)
}

func GetAdminUserFromCtx(ctx iris.Context) *AdminUser {
	get := ctx.Values().Get(UserInfoKey)
	if get == nil {
		return nil
	}
	return get.(*AdminUser)
}

package auth

import (
	"context"
	"net"
	"strings"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

const userKey string = "user"
const roleKey string = "role"
const langKey string = "lang"

func SetUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func GetUser(ctx context.Context) string {
	if v, ok := ctx.Value(userKey).(string); ok {
		return v
	}
	return "0"
}

func GetLang(ctx context.Context) string {
	if v, ok := ctx.Value(langKey).(string); ok {
		return v
	}
	return "en"
}

func GetClientIp(ctx context.Context) string {
	tr, ok := khttp.RequestFromServerContext(ctx)
	if !ok {
		return ""
	}

	// 1. 代理转发头（优先）
	ip := tr.Header.Get("X-Forwarded-For")
	if ip != "" {
		arr := strings.Split(ip, ",")
		return strings.TrimSpace(arr[0])
	}

	// 2. Nginx 常用头
	ip = tr.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 3. RemoteAddr
	host, _, err := net.SplitHostPort(tr.RemoteAddr)
	if err == nil {
		return host
	}

	return tr.RemoteAddr
}
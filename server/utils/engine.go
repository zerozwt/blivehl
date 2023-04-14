package utils

import (
	"net/http"
	"strings"

	"github.com/zerozwt/blivehl/server/engine"
)

const (
	COOKIE_KEY = "blivehl_sess"
	USER_KEY   = "blivehl_user"
	ADMIN_KEY  = "blivehl_admin"
)

func GetCtxCookie(ctx *engine.Context) (string, bool) {
	cookie, err := ctx.RawRequest.Cookie(COOKIE_KEY)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

func PutCtxCookie(ctx *engine.Context, cookie string) {
	cookieObj := &http.Cookie{
		Name:     COOKIE_KEY,
		Value:    cookie,
		Path:     "/",
		Domain:   trimPort(ctx.RawRequest.Host),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	ctx.RawResponse.Header().Set("Set-Cookie", cookieObj.String())
}

func trimPort(host string) string {
	if idx := strings.LastIndex(host, ":"); idx >= 0 {
		return host[:idx]
	}
	return host
}

func GetCtxUser(ctx *engine.Context) (string, bool) {
	return engine.CtxValue[string](ctx, USER_KEY)
}

func PutCtxUser(ctx *engine.Context, user string) {
	ctx.PutValue(USER_KEY, user)
}

func CheckCtxAdmin(ctx *engine.Context) bool {
	value, ok := engine.CtxValue[bool](ctx, ADMIN_KEY)
	return ok && value
}

func SetCtxAdmin(ctx *engine.Context, isAdmin bool) {
	ctx.PutValue(ADMIN_KEY, isAdmin)
}

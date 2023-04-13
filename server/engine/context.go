package engine

import (
	"net/http"
)

type Context struct {
	RawRequest  *http.Request
	RawResponse http.ResponseWriter

	values map[any]any
	chain  []func(*Context)
}

type HandlerFunc func(*Context)

func makeContext(w http.ResponseWriter, r *http.Request, handlers ...HandlerFunc) *Context {
	tmp := []func(*Context){}
	for _, f := range handlers {
		tmp = append(tmp, f)
	}
	return &Context{
		RawRequest:  r,
		RawResponse: w,
		values:      make(map[any]any),
		chain:       tmp,
	}
}

func (ctx *Context) Next() {
	if len(ctx.chain) == 0 {
		panic("context have no handlers")
	}
	handler := ctx.chain[0]
	ctx.chain = ctx.chain[1:]
	handler(ctx)
}

func (ctx *Context) PutValue(key, value any) {
	ctx.values[key] = ctx.values[value]
}

func (ctx *Context) GetValue(key any) any {
	value, ok := ctx.values[key]
	if ok {
		return value
	}
	return nil
}

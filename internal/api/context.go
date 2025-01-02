package api

import (
	"context"
	"encoding/json"
	"net/http"
)

type Context struct {
	context.Context
	cancel context.CancelFunc
	w      http.ResponseWriter
	r      *http.Request
}

type Validator interface {
	Validate(ctx *Context) error
}

func (ctx *Context) SetCancellableCtx(baseCtx context.Context, cancel context.CancelFunc) {
	ctx.Context = baseCtx
	ctx.cancel = cancel
}

func (ctx *Context) SetHTTPWriter(w http.ResponseWriter) {
	ctx.w = w
}

func (ctx *Context) SetHTTPRequest(r *http.Request) {
	ctx.r = r
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.w.Header().Set(key, value)
}

func (ctx *Context) GetHeader(key string) string {
	return ctx.r.Header.Get(key)
}

func (ctx *Context) StopChain() {
	ctx.cancel()
}

func (ctx *Context) GetMethod() string {
	return ctx.r.Method
}

func (ctx *Context) Decode(dest interface{}) error {
	err := json.NewDecoder(ctx.r.Body).Decode(dest)
	if err != nil {
		return err
	}

	return dest.(Validator).Validate(ctx)
}

func (ctx *Context) WriteResponse(statusCode int, resp interface{}) error {
	ctx.cancel()

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if ctx.w.Header().Get("Content-Type") == "" {
		ctx.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	ctx.w.WriteHeader(statusCode)

	_, err = ctx.w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

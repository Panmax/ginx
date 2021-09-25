package ginx

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"io/ioutil"
)

type IRequest interface {
	QueryInt(key string, def int) (int, bool)

	ParamInt(key string, def int) (int, bool)

	FormInt(key string, def int) (int, bool)

	BindJson(obj interface{}) error

	Uri() string
	Method() string
	Host() string
	ClientIp() string

	Headers() map[string][]string
	Header(key string) (string, bool)

	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request != nil {
		return ctx.request.URL.Query()
	}
	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	if vals, ok := ctx.QueryAll()[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) Param(key string) interface{} {
	if ctx.params != nil {
		if val, ok := ctx.params[key]; ok {
			return val
		}
	}
	return nil
}

func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	if val := ctx.Param(key); val != nil {
		return cast.ToInt(val), true
	}
	return def, false
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		ctx.request.ParseForm()
		return ctx.request.PostForm
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) (int, bool) {
	if vals, ok := ctx.FormAll()[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}
	return def, false
}

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		errors.New("ctx.request empty")
	}
	return nil
}

func (ctx *Context) Uri() string {
	return ctx.request.RequestURI
}

func (ctx *Context) Method() string {
	return ctx.request.Method
}

func (ctx *Context) Host() string {
	return ctx.request.Host
}

func (ctx *Context) ClientIp() string {
	r := ctx.request
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}

func (ctx *Context) Headers() map[string][]string {
	return ctx.request.Header
}

func (ctx *Context) Header(key string) (string, bool) {
	vals := ctx.request.Header.Values(key)
	if len(vals) > 0 {
		return vals[0], true
	}
	return "", false
}

func (ctx *Context) Cookies() map[string]string {
	cookies := ctx.request.Cookies()
	ret := make(map[string]string, len(cookies))
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}

func (ctx *Context) Cookie(key string) (string, bool) {
	cookies := ctx.Cookies()
	if val, ok := cookies[key]; ok {
		return val, true
	}
	return "", false
}

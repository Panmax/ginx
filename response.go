package ginx

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
)

type IResponse interface {
	Json(obj interface{}) IResponse

	JsonP(object interface{}) IResponse

	SetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse

	SetHeader(key, val string) IResponse

	SetStatus(code int) IResponse

	SetOkStatus() IResponse
}

func (ctx *Context) Json(obj interface{}) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.SetStatus(http.StatusInternalServerError)
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.responseWriter.Write(byt)
	return ctx
}

func (ctx *Context) JsonP(obj interface{}) IResponse {
	callbackFunc, _ := ctx.QueryString("callback", "callback_function")
	ctx.SetHeader("Content-Type", "application/javascript")
	callback := template.JSEscapeString(callbackFunc)

	_, err := ctx.responseWriter.Write([]byte(callback))
	if err != nil {
		return ctx
	}

	_, err = ctx.responseWriter.Write([]byte("("))
	if err != nil {
		return ctx
	}

	content, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.responseWriter.Write(content)
	if err != nil {
		return ctx
	}

	_, err = ctx.responseWriter.Write([]byte(")"))
	if err != nil {
		return ctx
	}
	return ctx
}

func (ctx *Context) SetHeader(key, val string) IResponse {
	ctx.responseWriter.Header().Add(key, val)
	return ctx
}

func (ctx *Context) SetCookie(key, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.responseWriter, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: 1,
	})
	return ctx
}

func (ctx *Context) SetStatus(code int) IResponse {
	ctx.responseWriter.WriteHeader(code)
	return ctx
}

func (ctx *Context) SetOkStatus() IResponse {
	ctx.responseWriter.WriteHeader(http.StatusOK)
	return ctx
}

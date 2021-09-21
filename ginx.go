package ginx

import (
	"log"
	"net/http"
	"strings"
)

type Ginx struct {
	router map[string]*Tree
}

func NewGinx() *Ginx {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Ginx{router: router}
}

func (g *Ginx) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := NewContext(request, response)

	handler := g.FindRouteByRequest(request)
	if handler == nil {
		ctx.Json(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if err := handler(ctx); err != nil {
		ctx.Json(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (g *Ginx) FindRouteByRequest(request *http.Request) ControllerHandler {
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	if methodTree, ok := g.router[upperMethod]; ok {
		return methodTree.FindHandler(uri)
	}
	return nil
}

func (g *Ginx) Get(url string, handler ControllerHandler) {
	if err := g.router["GET"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (g *Ginx) Post(url string, handler ControllerHandler) {
	if err := g.router["POST"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (g *Ginx) Put(url string, handler ControllerHandler) {
	if err := g.router["PUT"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (g *Ginx) Delete(url string, handler ControllerHandler) {
	if err := g.router["DELETE"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (g *Ginx) Group(prefix string) IGroup {
	return NewGroup(g, prefix)
}

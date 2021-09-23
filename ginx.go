package ginx

import (
	"log"
	"net/http"
	"strings"
)

type Ginx struct {
	router      map[string]*Tree
	middlewares []ControllerHandler
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

	handlers := g.FindRouteByRequest(request)
	if handlers == nil {
		ctx.Json(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	ctx.SetHandlers(handlers)

	if err := ctx.Next(); err != nil {
		ctx.Json(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (g *Ginx) FindRouteByRequest(request *http.Request) []ControllerHandler {
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	if methodTree, ok := g.router[upperMethod]; ok {
		return methodTree.FindHandler(uri)
	}
	return nil
}

func (g *Ginx) Use(middirewares ...ControllerHandler) {
	g.middlewares = middirewares
}

func (g *Ginx) Get(url string, handlers ...ControllerHandler) {
	allHandlers := append(g.middlewares, handlers...)
	if err := g.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (g *Ginx) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(g.middlewares, handlers...)
	if err := g.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (g *Ginx) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(g.middlewares, handlers...)
	if err := g.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (g *Ginx) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(g.middlewares, handlers...)
	if err := g.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (g *Ginx) Group(prefix string) IGroup {
	return NewGroup(g, prefix)
}

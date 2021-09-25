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

	node := g.FindRouteNodeByRequest(request)
	if node == nil {
		ctx.SetStatus(http.StatusNotFound).Json(http.StatusText(http.StatusNotFound))
		return
	}

	ctx.SetHandlers(node.handlers)
	params := node.parseParamsFromEndNode(request.URL.Path)
	ctx.SetParams(params)

	if err := ctx.Next(); err != nil {
		ctx.SetStatus(http.StatusInternalServerError).Json(http.StatusText(http.StatusInternalServerError))
		return
	}
}

func (g *Ginx) FindRouteNodeByRequest(request *http.Request) *node {
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	if methodTree, ok := g.router[upperMethod]; ok {
		return methodTree.root.matchNode(uri)
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

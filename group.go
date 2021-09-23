package ginx

type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)
	Use(middlewares ...ControllerHandler)
}

type Group struct {
	ginx        *Ginx
	prefix      string
	middlewares []ControllerHandler
}

func NewGroup(ginx *Ginx, prefix string) *Group {
	return &Group{
		ginx:   ginx,
		prefix: prefix,
	}
}

func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = middlewares
}

func (g *Group) Get(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.middlewares, handlers...)
	g.ginx.Get(uri, allHandlers...)
}

func (g *Group) Post(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.middlewares, handlers...)
	g.ginx.Post(uri, allHandlers...)
}

func (g *Group) Put(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.middlewares, handlers...)
	g.ginx.Put(uri, allHandlers...)
}

func (g *Group) Delete(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.middlewares, handlers...)
	g.ginx.Delete(uri, allHandlers...)
}

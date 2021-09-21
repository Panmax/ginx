package ginx

type IGroup interface {
	Get(string, ControllerHandler)
	Post(string, ControllerHandler)
	Put(string, ControllerHandler)
	Delete(string, ControllerHandler)
}

type Group struct {
	ginx   *Ginx
	prefix string
}

func NewGroup(ginx *Ginx, prefix string) *Group {
	return &Group{
		ginx:   ginx,
		prefix: prefix,
	}
}

func (g *Group) Get(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.ginx.Get(uri, handler)
}

func (g *Group) Post(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.ginx.Post(uri, handler)
}

func (g *Group) Put(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.ginx.Put(uri, handler)
}

func (g *Group) Delete(uri string, handler ControllerHandler) {
	uri = g.prefix + uri
	g.ginx.Delete(uri, handler)
}

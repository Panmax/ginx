package framework

type NewInstance func(...interface{}) (interface{}, error)

type ServiceProvider interface {
	Register(Container) NewInstance

	Boot(Container) error

	IsDefer() bool

	Params(Container) []interface{}

	Name() string
}

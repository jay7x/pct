package template

const DefaultEngineName = "go"

type Engine interface {
	Name() string
	Render(content string, vars map[string]interface{}) (string, error)
}

var engines = map[string]Engine{
	DefaultEngineName: &goEngine{},
}

func Get(name string) (Engine, bool) {
	e, ok := engines[name]
	return e, ok
}

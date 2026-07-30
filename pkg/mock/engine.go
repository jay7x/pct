package mock

type Engine struct {
	NameFn   func() string
	RenderFn func(content string, vars map[string]interface{}) (string, error)
}

func (e *Engine) Name() string {
	if e.NameFn == nil {
		return "mock"
	}
	return e.NameFn()
}

func (e *Engine) Render(content string, vars map[string]interface{}) (string, error) {
	if e.RenderFn == nil {
		return content, nil
	}
	return e.RenderFn(content, vars)
}

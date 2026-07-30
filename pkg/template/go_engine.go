package template

import (
	"bytes"
	gotemplate "text/template"

	"github.com/jay7x/pct/pkg/utils"
)

var goFuncs = gotemplate.FuncMap{
	"toClassName": utils.ToClassName,
	"ns2path":     utils.Ns2Path,
}

type goEngine struct{}

func (e *goEngine) Name() string { return DefaultEngineName }

func (e *goEngine) Render(content string, vars map[string]interface{}) (string, error) {
	tmpl, err := gotemplate.New("").Funcs(goFuncs).Parse(content)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}

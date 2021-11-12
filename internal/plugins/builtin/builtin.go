package builtin

import (
	"html/template"

	"github.com/drewstinnett/go-output-format/formatter"
	"github.com/drewstinnett/labdoc/pkg/labdoc"
)

type plug struct{}

func (p *plug) TemplateFunctions() (template.FuncMap, error) {
	templ := template.FuncMap{
		"debug": debug,
	}
	return templ, nil
}

func debug(item interface{}) (string, error) {
	c := &formatter.Config{
		Format: "plain",
	}
	out, err := formatter.OutputData(item, c)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func init() {
	labdoc.Add("builtin", func() labdoc.Plugin {
		return &plug{}
	})
}

package labdoc

import (
	"html/template"
)

type Plugin interface {
	TemplateFunctions() (template.FuncMap, error)
}

type Creator func() Plugin

var Plugins = map[string]Creator{}

func Add(name string, creator Creator) {
	Plugins[name] = creator
}

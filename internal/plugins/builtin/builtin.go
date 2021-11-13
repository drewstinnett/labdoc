package builtin

import (
	"fmt"
	"html/template"
	"time"

	"github.com/apex/log"
	"github.com/araddon/dateparse"
	"github.com/drewstinnett/go-output-format/formatter"
	"github.com/drewstinnett/labdoc/pkg/labdoc"
	"github.com/dustin/go-humanize"
)

type plug struct{}

func (p *plug) TemplateFunctions() (template.FuncMap, error) {
	templ := template.FuncMap{
		"debug": debug,
		"ago":   ago,
	}
	return templ, nil
}

func (p *plug) Examples() string {
	return `## Debug Item Output
{{ range rssListFeed "https://www.example.com/feed" 10}}
{{ . | builtinDebug }}
{{ end }}

## Time Ago Info
{{ range rssListFeed "https://www.example.com/feed" 10}}
{{ .Title }} - {{ .PublishedParsed | builtinAgo }}
{{ end }}`
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

func ago(t interface{}) string {
	// log.Warnf("%v is a %v", t, fmt.Sprintf("%T", t))
	switch v := t.(type) {
	case time.Time:
		return humanize.Time(v)
	case *time.Time:
		return humanize.Time(*v)
	case string:
		// Try to make a good guess on converting a string to a time...this feels gross...
		g, err := dateparse.ParseAny("3/1/2014")
		if err != nil {
			log.Warnf("Could not convert %v to an actual time", fmt.Sprintf("%v", v))
			return fmt.Sprintf("%v", v)
		}
		return humanize.Time(g)
	}
	log.Warnf("Could not convert %v to an actual time", fmt.Sprintf("%v", t))
	return fmt.Sprintf("%v", t)
}

func init() {
	labdoc.Add("builtin", func() labdoc.Plugin {
		return &plug{}
	})
}

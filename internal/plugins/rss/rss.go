package rss

import (
	"html/template"

	"github.com/drewstinnett/labdoc/pkg/labdoc"
	"github.com/mmcdole/gofeed"
)

type plug struct{}

func (p *plug) TemplateFunctions() (template.FuncMap, error) {
	templ := template.FuncMap{
		"listFeed": listFeed,
	}
	return templ, nil
}

func (p *plug) Examples() string {
	return `## Recent Reviews
{{ range rssListFeed "https://www.rogerebert.com/feed" 10}}
[{{ .Title }}]({{ .Link }}) {{ .PublishedParsed | builtinAgo }}
{{- end }}`
}

func listFeed(url string, limit int) ([]*gofeed.Item, error) {
	var items []*gofeed.Item
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}
	items = feed.Items[0:limit]
	return items, nil
}

func init() {
	labdoc.Add("rss", func() labdoc.Plugin {
		return &plug{}
	})
}

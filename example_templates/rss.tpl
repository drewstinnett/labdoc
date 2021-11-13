# RSS

## Recent Reviews
{{ range rssListFeed "https://www.rogerebert.com/feed" 10}}
[{{ .Title }}]({{ .Link }}) {{ .PublishedParsed | builtinAgo }}
{{- end }}
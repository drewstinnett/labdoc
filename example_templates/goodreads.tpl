## Recently Read Books
{{ range goodreadsRecentlyRead 10 }}
- [{{ .Title }}]({{ .Link }}) by {{ .Author }} {{ .PubDate | builtinAgo }}
{{- end }}
# LabDoc

Generate a README.md for GitLab similar to the way
[markscribe](https://github.com/muesli/markscribe) does in GitHub.

## Builtin Template Functions

### builtinDebug

Oh dude, would love do document this better. For now, you can use the
`builtinDebug` filter though. Would love some help on making this more clear

```go
{{ . | builtinDebug }}
```

This will show a json representation of the object

### builtinAgo

Convert a time.Time or string to 'time ago'. If it can't convert the string to a
time.Time, it'll just return the original string

```go
{{ .Title }} - {{ .PublishedParsed | builtinAgo }}
```

## Extending with more Plugins

We would love to have more plugins!

Plugins should have their own directory under
[internal/plugins](internal/plugins), and have the following interface:

```go
type plug struct{}

func (p *plug) TemplateFunctions() (template.FuncMap, error) {
  templ := template.FuncMap{
    "templateFunctionFoo":       templateFunctionFoo,
  }
  return templ, nil
}

func (p *plug) Examples() string {
    return `## Example Entries
{{ range rssListFeed "https://www.example.com/feed" 10}}
[{{ .Title }}]({{ .Link }}) on {{ .Published }}
{{- end }}
`
}

func templateFunctionFoo(limit int) ([]someOject, error){
  ...
  return objects, nil
}

func init() {
  labdoc.Add("myplugin", func() labdoc.Plugin {
    return &plug{}
  })
}
```

Once your plugin is ready, register it in
[plugins/all/all.go](plugins/all/all.go)

Note that in the above example, the template function will be namespaced by the
plugin name, and can be called using `mypluginTemplateFunctionFoo`.

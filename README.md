# LabDoc

Generate a README.md for GitLab similar to the way
[markscribe](https://github.com/muesli/markscribe) does in GitHub.

## Extending with more Plugins

We would love to have more plugins!

Plugins should have their own directory under [internal/plugins](internal/plugins), and have the following interface:

```go
type plug struct{}

func (p *plug) TemplateFunctions() (template.FuncMap, error) {
  templ := template.FuncMap{
    "templateFunctionFoo":       templateFunctionFoo,
  }
  return templ, nil
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

Note that in the above example, the template function will be namespaced by the
plugin name, and can be called using `mypluginTemplateFunctionFoo`.

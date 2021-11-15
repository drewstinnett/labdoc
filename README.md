# LabDoc

Generate a README.md for GitLab similar to the way
[markscribe](https://github.com/muesli/markscribe) does in GitHub.

## Usage

### Generate README.md

Currently the rendered version of your template will be printed to stdout. Just
direct it in to the target file like so:

```bash
labdoc generate example_templates/rss.tpl > README.md
```

### Show Examples

Print out example usage, based on plugin data:

```bash
labdoc examples
```

### Run in Gitlab CI

Create a new personal access token from your profile page on Gitlab.com (or your
locally hosted URL). It will need the following permissions: `read_user`,
`read_api`, and `write_repository`. Save this token as `GITLAB_TOKEN` in your
runner variable section, along with any other plugin variables that you might
need.

Set up your `.gitlab-ci.yml` file with something like the following:

```yaml
---
stages:
  - release

release:
  stage: release
  image:
    name: brewerdrewer/labdoc:latest
    entrypoint: ['']
  variables:
    GIT_DEPTH: 0
  script:
    - /labdoc generate templates/README.md.tpl > README.md
    - git config --global user.email "${GITLAB_USER_EMAIL}"
    - git config --global user.name "${GITLAB_USER_NAME}"
    - CHANGED=$(git diff)
    - |
      if [ -n "${CHANGED}" ]; then
        git commit -a -m 'Updating readme'
        git push https://${GITLAB_USER_LOGIN}:${GITLAB_TOKEN}@${CI_SERVER_HOST}/${CI_PROJECT_ROOT_NAMESPACE}/${CI_PROJECT_NAME}.git HEAD:main
      fi
```

## Plugins

### Builtin Template

This is for generic template functions across multiple plugins

#### builtinDebug

Oh dude, would love do document this better. For now, you can use the
`builtinDebug` filter though. Would love some help on making this more clear

```go
{{ . | builtinDebug }}
```

This will show a json representation of the object

#### builtinAgo

Convert a time.Time or string to 'time ago'. If it can't convert the string to a
time.Time, it'll just return the original string

```go
{{ .Title }} - {{ .PublishedParsed | builtinAgo }}
```

### RSS Plugin

Simple plugin for just loading an rss file and returning the entries

#### rssListFeed

No additional configuration needed, just pass in a feed URL and a limit

```go
{{ range rssListFeed "https://www.rogerebert.com/feed" 10}}
```

### Letterboxd Plugin

Uses the letterboxd rss feeds to query your recent events. Requires
`LETTERBOXD_USER` environment variable to be set

#### letterboxdRecentlyWatched

```go
{{ range letterboxdRecentlyWatched 5 }}
```

### Goodreads Plugin

Uses the goodreads rss feed, since the API has been deprecated. You'll need to
set both of these variables in order for this to work:

`GOODREADS_RSSUSERID` - You can find this in your profile link. For example, if
you see `https://www.goodreads.com/user/show/111216449-drew-stinnett`, then
`111216449` is what `GOODREADS_RSSUSERID` should be set to

`GOODREADS_RSSKEY` - Unsure how private this should be, but if it's in an RSS
url, assuming it's not too private. Pull this by clicking on the rss icon in the
goodreads web interface, and selecting the long string after `key=` and before
`&`in the url

#### goodreadsRecentlyRead

```go
{{ range goodreadsRecentlyRead 10 }}
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

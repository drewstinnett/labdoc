## Watched List
{{ range letterboxdRecentlyWatched 5 }}
* {{ .Verb }} [{{ .TitleWithRating }})]({{ .Link }})
{{- end }}

## Watched Posters
{{ range letterboxdRecentlyWatched 5 }}
![{{ .FilmTitle }}]({{ .Poster}})
{{- end }}
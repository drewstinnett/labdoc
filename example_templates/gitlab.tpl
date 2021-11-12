# Gitlab

## Recent MRs
{{ range gitlabRecentlyAcceptedMergeRequests 10}}
{{ .Event.TargetTitle }} for {{ .Project.Name }}
{{- end }}


## Recently Created Projects
{{ range gitlabRecentlyCreatedProjects 10 }}
[{{ .Name }}]({{ .WebURL }})
{{- end }}
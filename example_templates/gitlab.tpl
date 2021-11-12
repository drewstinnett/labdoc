# Gitlab

## Recent MRs
{{ range gitlabRecentlyAcceptedMergeRequests 10}}
{{ .Event.TargetTitle }} for {{ .Project.Name }}
{{- end }}


## Recent Events 
{{ range gitlabRecentlyCreatedEvents 10 }}
{{ .Event.ActionName }}
{{- end }}

## Recently Created Projects
{{ range gitlabRecentlyCreatedProjects 10 }}
[{{ .Name }}]({{ .WebURL }})
{{- end }}
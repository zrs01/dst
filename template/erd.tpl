@startuml

skinparam {
	linetype ortho
	arrowFontSize 10
}
{{- range .Schemas }}
  {{- range .Tables }}
entity "{{ .Name }}{{ if .Title != "" }}\n<size:11>({{ .Title }})</size>{{ end }}" as {{ .Name }} {
  |= |= <size:11>name</size> |= <size:11>type</size> |
    {{- range .Columns }}
  |  | <size:11>{{ .Name }}</size> | <size:11>{{ .DataType }}</size> |
    {{- end }}
}
  {{- end }}
{{- end }}

{{- range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns }}
      {{- if .ForeignKey != "" }}
        {{- parts := split(.ForeignKey, ".") }}
{{ table }} }|-- {{ parts[0] }} #text:FireBrick : {{ .Name }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}

@enduml
@startuml

skinparam {
	linetype ortho
	arrowFontSize 10
}

{{ range .Schemas }}
  {{- range .Tables }}
entity "{{ .Name }}{{ if .Title != "" }}\n<size:11>({{ .Title }})</size>{{ end }}" as {{ .Name }} {
  |= |= <size:11>name</size> |= <size:11>type</size> |
    {{- range .Columns }}
      {{- if .Identity == "Y" || .ForeignKey != "" }}
  | {{ if .Identity == "Y" }}<size:11>PK</size>{{ end }}{{ if .ForeignKey != "" }}<size:11>FK</size>{{ end }} | <size:11>{{ .Name }}</size> | <size:11>{{ .DataType }}</size> |
      {{- end }}
    {{- end }}
}
  {{- end }}
{{- end }}

{{ range .Schemas }}
  {{- range .Tables }}
    {{- range .Columns }}
      {{- if .ForeignKey != "" }}
        {{- parts := split(.ForeignKey, ".") }}
entity {{ parts[0] }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}

{{ range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns }}
      {{- if .ForeignKey != "" }}
        {{- parts := split(.ForeignKey, ".") }}
{{ table }} }o--{{ if .NotNull == "Y" }}||{{ end }}{{ if .NotNull != "Y" }}o|{{ end }} {{ parts[0] }} #text:FireBrick : "{{ .Name }}"
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}

@enduml
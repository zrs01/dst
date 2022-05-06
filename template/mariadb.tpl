{* ---------------------------------- tables ---------------------------------- *}
{{- fixed := .Fixed }}
{{ range .Schemas }}
  {{- range .Tables }}
CREATE TABLE IF NOT EXISTS {{ .Name }} (
    {{- range i := .Columns}}
      {{ .Name }} {{ .DataType }}
      {{- if .NotNull == "Y" }} NOT NULL {{- end }}
      {{- if .Value != "" }} DEFAULT '{{ .Value }}' {{- end }}
      {{- if .Identity == "Y" }} AUTO_INCREMENT PRIMARY KEY {{- end }}
      {{- if .Desc != "" }} COMMENT '{{ .Desc }}' {{- end }}
      {{- "," }}
    {{- end }}
    {{- fixedCount := len(fixed) }}
    {{- range i := fixed }}
      {{ .Name }} {{ .DataType }}
      {{- if .NotNull == "Y" }} NOT NULL {{- end }}
      {{- if .Value != "" }} DEFAULT '{{ .Value }}' {{- end }}
      {{- if .Identity == "Y" }} AUTO_INCREMENT PRIMARY KEY {{- end }}
      {{- if .Desc != "" }} COMMENT '{{ .Desc }}' {{- end }}
      {{- if i < fixedCount - 1 }},{{- end }}
    {{- end }}
);
  {{- end }}
{{ end }}
{* ------------------------------- foreign keys ------------------------------- *}
{{ range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns}}
      {{- parts := split(.ForeignKeyHint, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE {{ table }} ADD CONSTRAINT fk_{{ table }}_{{ .Name }} FOREIGN KEY ({{ parts[0] }}) REFERENCE {{ parts[0] }} ({{ parts[1] }})
      {{- end }}
    {{- end }}
    {{- range fixed }}
      {{- parts := split(.ForeignKeyHint, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE {{ table }} ADD CONSTRAINT fk_{{ table }}_{{ .Name }} FOREIGN KEY ({{ parts[0] }}) REFERENCE {{ parts[0] }} ({{ parts[1] }})
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}

{* syntax: https://github.com/CloudyKit/jet/blob/master/docs/syntax.md *}

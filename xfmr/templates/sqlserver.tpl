{* ---------------------------------- tables ---------------------------------- *}
-- up
{{- fixed := .Fixed }}
{{ range .Schemas }}
  {{- range .Tables }}
CREATE TABLE {{ .Name }} (
    {{- idc := "" }}
    {{- range i := .Columns}}
      {{- if .Identity == "Y" }}
        {{- idc = . }}
      {{- end }}
      {{ .Name }} {{ .DataType }}
      {{- if .Identity == "Y" }} IDENTITY(1,1){{- end }}
      {{- if .Value != "" }} DEFAULT '{{ .Value }}' {{- end }}
      {{- if .NotNull == "Y" }} NOT NULL {{- end }}
      {{- "," }}
    {{- end }}
    {{- fixedCount := len(fixed) }}
    {{- range i := fixed }}
      {{ .Name }} {{ .DataType }}
      {{- if .Value != "" }} DEFAULT '{{ .Value }}' {{- end }}
      {{- if .NotNull == "Y" }} NOT NULL {{- end }}
      {{- if i < fixedCount - 1 }},{{- end }}
    {{- end }}
    {{- if idc != "" }},
      CONSTRAINT pk{{ .Name }}{{ idc.Name }} PRIMARY KEY ({{ idc.Name }})
    {{- end }}
);
  {{- end }}
{{ end }}
{* ------------------------------- foreign keys ------------------------------- *}
{{ range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns}}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE {{ table }} ADD CONSTRAINT fk{{ table }}{{ .Name }} FOREIGN KEY ({{ .Name }}) REFERENCES {{ parts[0] }} ({{ parts[1] }});
      {{- end }}
    {{- end }}
    {{- range fixed }}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE {{ table }} ADD CONSTRAINT fk{{ table }}{{ .Name }} FOREIGN KEY ({{ .Name }}) REFERENCES {{ parts[0] }} ({{ parts[1] }});
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}



-- down
{* ------------------------------- foreign keys ------------------------------- *}
{{ range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns}}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE {{ table }} DROP CONSTRAINT fk{{ table }}{{ .Name }};
      {{- end }}
    {{- end }}
    {{- range fixed }}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE {{ table }} DROP CONSTRAINT fk{{ table }}{{ .Name }};
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
{* ---------------------------------- tables ---------------------------------- *}
-- up
{{ range .Schemas }}
  {{- range .Tables }}
DROP TABLE {{ .Name }};
  {{- end }}
{{ end }}


{* syntax: https://github.com/CloudyKit/jet/blob/master/docs/syntax.md *}

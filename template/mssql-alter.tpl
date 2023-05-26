-- +goose Up
{* -- ------------------------------ ADD COLUMNS ------------------------------- -- *}
{{- range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns}}
IF COL_LENGTH(N'{{ table }}', N'{{ .Name }}') IS NULL ALTER TABLE {{ table }} ADD {{ .Name }} {{ .DataType }} {{- if .NotNull == "Y" }} NOT NULL {{- end }};
    {{- end }}
  {{- end }}
{{- end }}
{* ------------------------------- foreign keys ------------------------------- *}
{{- range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns}}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
IF OBJECT_ID(N'fk{{ table }}{{ .Name }}', N'F') IS NULL ALTER TABLE {{ table }} ADD CONSTRAINT fk{{ table }}{{ .Name }} FOREIGN KEY ({{ .Name }}) REFERENCES {{ parts[0] }} ({{ parts[1] }});

      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}


-- +goose Down
{* ------------------------------- foreign keys ------------------------------- *}
{{- range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns}}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
IF OBJECT_ID(N'fk{{ table }}{{ .Name }}', N'F') IS NOT NULL ALTER TABLE {{ table }} DROP CONSTRAINT IF EXISTS fk{{ table }}{{ .Name }};
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
{* -- ------------------------------ DROP COLUMNS ------------------------------ -- *}
{{- range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns}}
IF COL_LENGTH(N'{{ table }}', N'{{ .Name }}') IS NOT NULL ALTER TABLE {{ table }} DROP COLUMN {{ .Name }};
    {{- end }}
  {{- end }}
{{- end }}


{* syntax: https://github.com/CloudyKit/jet/blob/master/docs/syntax.md *}

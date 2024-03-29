-- +goose Up
-- +goose NO TRANSACTION
{* ---------------------------------- tables ---------------------------------- *}
{{- fixed := .Fixed }}
{{- range .Schemas }}
  {{- range .Tables }}
IF OBJECT_ID(N'{{ .Name }}', N'U') IS NULL
BEGIN
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
)
END;
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
    {{- range fixed }}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
IF OBJECT_ID(N'fk{{ table }}{{ .Name }}', N'F') IS NULL ALTER TABLE {{ table }} ADD CONSTRAINT fk{{ table }}{{ .Name }} FOREIGN KEY ({{ .Name }}) REFERENCES {{ parts[0] }} ({{ parts[1] }});
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}

{* -- ------------------------------- CREATE INDEX ------------------------------- -- *}
{{- range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns}}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx{{ table }}{{ .Name }}' AND object_id = OBJECT_ID('{{ table }}')) CREATE INDEX idx{{ table }}{{ .Name }} ON {{ table }} ({{ .Name }});
      {{- end }}
      {{- if .Index == "Y" }}
IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx{{ table }}{{ .Name }}' AND object_id = OBJECT_ID('{{ table }}')) CREATE{{ if .Unique == "Y" }} UNIQUE{{ end }} INDEX idx{{ table }}{{ .Name }} ON {{ table }} ({{ .Name }});
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
    {{- range fixed }}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
IF OBJECT_ID(N'fk{{ table }}{{ .Name }}', N'F') IS NOT NULL ALTER TABLE {{ table }} DROP CONSTRAINT IF EXISTS fk{{ table }}{{ .Name }};
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
{* ---------------------------------- tables ---------------------------------- *}

{{- range .Schemas }}
  {{- range .Tables }}
DROP TABLE IF EXISTS {{ .Name }};
  {{- end }}
{{- end }}


{* syntax: https://github.com/CloudyKit/jet/blob/master/docs/syntax.md *}

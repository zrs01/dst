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

{* -------------------------- foreign keys (create) ------------------------- *}
{{ range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns }}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE IF EXISTS {{ table }} ADD CONSTRAINT fk_{{ table }}_{{ .Name }} FOREIGN KEY ({{ .Name }}) REFERENCES {{ parts[0] }} ({{ parts[1] }});
      {{- end }}
    {{- end }}
    {{- range fixed }}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE IF EXISTS {{ table }} ADD CONSTRAINT fk_{{ table }}_{{ .Name }} FOREIGN KEY ({{ .Name }}) REFERENCES {{ parts[0] }} ({{ parts[1] }});
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}

{* ------------------------------ alter columns ----------------------------- *}
{{ range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- lastColumn := "" }}
    {{- range i := .Columns }}
ALTER TABLE {{ table }} ADD COLUMN IF NOT EXISTS {{ .Name }} {{ .DataType }}
      {{- if i != 0 }} AFTER {{ lastColumn }}{{ end }}
      {{- ";" }}
      {{- lastColumn = .Name }}
    {{- end }}
  {{- end }}
{{- end }}

{* -------------------------- foreign key (remove) -------------------------- *}
{{ range .Schemas }}
  {{- range .Tables }}
    {{- table := .Name }}
    {{- range .Columns }}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE IF EXISTS {{ table }} DROP CONSTRAINT IF EXISTS fk_{{ table }}_{{ .Name }};
      {{- end }}
    {{- end }}
    {{- range fixed }}
      {{- parts := split(.ForeignKey, ".") }}
      {{- if len(parts) > 1 }}
ALTER TABLE IF EXISTS {{ table }} DROP CONSTRAINT fk_{{ table }}_{{ .Name }};
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}

{* -------------------------------- drop table -------------------------------- *}
{{ range .Schemas }}
  {{- range .Tables }}
DROP TABLE IF EXISTS {{ .Name }};
    {{- table := .Name }}
  {{- end }}
{{- end }}


{* syntax: https://github.com/CloudyKit/jet/blob/master/docs/syntax.md *}

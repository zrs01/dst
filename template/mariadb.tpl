{{ range $schema := .Schemas }}
  {{ $schema.Name }}
  {{ range $table := $schema.Tables }}
    CREATE TABLE {{ $table.Name }} {
      {{- range $i, $column := $table.Columns }}
        {{- $colCount := len $table.Columns }}
        {{ $column.Name }} {{ $column.DataType }}
        {{- if eq $column.NotNull "Y" }} NOT NULL{{ end }}
        {{- if eq $column.Identity "Y" }} AUTO_INCREMENT PRIMARY KEY{{ end }}
        {{- if ne $column.Desc "" }} {{$column.Desc}}{{ end }}
        {{- if lt $i $colCount }} {{ $i }} - {{  $colCount  }},{{ end }}
      {{- end }}
    }
  {{ end }}
{{ end }}
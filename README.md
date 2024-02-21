# Database Schema Tool

Export to different format from definition file.

## Usage

### Example

Create a schema file in yaml format (e.g. schema.yml)

```yml
# fixed columns:
# the columns will be appended at the end of each table
fixed:
  - { na: deleted, ty: int, va: 0, nu: Y, dc: "0:activated, others:deleted" }
  - { na: create_user_id, ty: int, nu: Y, dc: "record create user ID" }
  - { na: create_time, ty: datetime, nu: Y, dc: "record create time" }
  - { na: update_user_id, ty: int, nu: Y, dc: "record update user ID" }
  - { na: update_time, ty: datetime, nu: Y, dc: "record update time" }
  - { na: recver, ty: int, nu: Y, va: 0, dc: "record version" }

schemas:
  - name: General
    tables:
      - name: doc
        title: document
        desc: External document/files
        columns:
          - { na: doc_id, ty: INT, id: Y, nu: Y, dc: "unique identifier" }
          - { na: ver, ty: INT, nu: Y, dc: "version, starts from zero" }
          - { na: ref, ty: VARCHAR(50), nu: Y, un: Y, dc: "reference to locate the external file" }

      - name: doc_tag
        title: document tag
        desc: "Relationship table between doc and tag"
        columns:
          - { na: doc_tag_id, ty: INT, id: Y, nu: Y, dc: "unique identifier" }
          - { na: doc_id, ty: INT, nu: Y, fk: doc.doc_id, cd: "0..*:1", dc: "doc id" }
          - { na: tag_id, ty: INT, nu: Y, fk: tag.tag_id, cd: "0..*:1", dc: "tag id" }

      - name: tag
        desc: "Code table for tag"
        columns:
          - { na: tag_id, ty: INT, id: Y, nu: Y, dc: "unique identifier" }
          - { na: code, ty: VARCHAR(10), dc: "tag code. eg. premarket md application index A1" }
          - { na: name, ty: VARCHAR(20), nu: Y, dc: "tag name" }
          - { na: priority, ty: VARCHAR(50), dc: "tag priority" }
          - { na: descr, ty: VARCHAR(50), dc: "tag description" }
          - { na: source, ty: CHAR(3), nu: Y, dc: "Source; COM: common, PRE: pre, POS: post" }
          - { na: type, ty: CHAR(3), nu: Y, dc: "Tag Type; ATT: attachment, APP: application, COM: company" }
```

```sh
# if the schema.yml is in the current directory, you can skip the -i option
# if -o option does not exist, result will be printed to stdout

# -- YAML to SQL schema file
# generate using template.tpl
dst convert text -i schema.yml -o sample.sql -t template.tpl
# generate using template.tpl, select the tables start with 'tag' only
dst convert text -i schema.yml -o sample.sql -t template.tpl --table 'tag%'

# -- YAML to ER diagram definition file
dst convert diagram -i schema.yml -o sample.puml

# -- YAML to ER diagram (.png)
# download plantuml.jar from https://plantuml.com/download

# plantuml.jar can be skiped, it try to find the .jar from the PATH variable
dst convert diagram -i schema.yml -o sample.puml
# specify plantuml.jar
dst convert diagram -i schema.yml -o sample.puml -j plantuml.jar
# you may convert some of tables, select the tables starts with 'tag' only
dst convert diagram -i schema.yml -o sample.puml --table 'tag%'
# template file is used for the ER diagram
dst convert diagram -i schema.yml -t template.tpl -o sample.puml

```

## Template Syntax
https://github.com/CloudyKit/jet/blob/master/docs/syntax.md

## Data Model

### Root
| Field Name | Type     | YAML Tag | Description                         |
| ---------- | -------- | -------- | ----------------------------------- |
| Schemas    | []Schema | schemas  | Represents a list of schemas.       |

### Schema
| Field Name | Type    | YAML Tag | Description                                    |
| ---------- | ------- | -------- | ---------------------------------------------- |
| Name       | string  | name     | Represents the name of the schema.             |
| Desc       | string  | desc     | Represents the description of the schema.      |
| Tables     | []Table | tables   | Represents a list of tables within the schema. |

### Table
| Field Name | Type     | YAML Tag   | Description                                       |
| ---------- | -------- | ---------- | ------------------------------------------------- |
| Name       | string   | name       | Represents the name of the table.                 |
| Title      | string   | title      | Represents the title of the table.                |
| Desc       | string   | desc       | Represents the description of the table.          |
| Version    | bool     | version    | Indicates whether the table has a version.        |
| Columns    | []Column | columns    | Represents a list of columns within the table.    |
| References | []Ref    | references | Represents a list of references within the table. |

### Column
| Field Name  | Type   | YAML Tag | Description                                              |
| ----------- | ------ | -------- | -------------------------------------------------------- |
| Name        | string | na       | Represents the name of the column.                       |
| DataType    | string | ty       | Represents the data type of the column.                  |
| Identity    | string | id       | Represents the identity of the column.                   |
| NotNull     | string | nu       | Indicates whether the column is not null (default: "N"). |
| Unique      | string | un       | Indicates whether the column is unique.                  |
| Value       | string | va       | Represents the value of the column.                      |
| ForeignKey  | string | fk       | Represents the foreign key of the column.                |
| Cardinality | string | cd       | Represents the cardinality of the column.                |
| Title       | string | tt       | Represents the title of the column.                      |
| Index       | string | in       | Represents the index of the column.                      |
| Desc        | string | dc       | Represents the description of the column.                |
| Compute     | string | cm       | Represents the compute column.                           |

### Ref
| Field Name | Type           | YAML Tag | Description                        |
| ---------- | -------------- | -------- | ---------------------------------- |
| ColumnName | string         | column   | Represents the name of the column. |
| ForeignKey | []ForeignTable | foreign  | Represents the foreign keys.       |

### ForeignTable
| Field Name | Type   | YAML Tag | Description                        |
| ---------- | ------ | -------- | ---------------------------------- |
| Table      | string | table    | Represents the name of the table.  |
| Column     | string | column   | Represents the name of the column. |

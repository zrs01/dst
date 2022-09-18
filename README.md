# Database Schema Tool

Export to different format from definition file.

## Usage

### Command

```sh
NAME:
   dst - Database schema tool

USAGE:
   dst [global options] command [command options] [arguments...]

VERSION:
   development

COMMANDS:
   convert, c  Convert to other format
   verify, v   Verify the foreign key
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d    Debug mode (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

Use `--help` after command to show addition information, e.g.

```sh
$ dst convert --help

NAME:
   dst convert - Convert to other format

USAGE:
   dst convert command [command options] [arguments...]

COMMANDS:
   yaml, y
   excel, e
   text, t
   diagram, d
   help, h     Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
```

### Example

Create a definition file (e.g. sample.yml)

```yml
# column definition:
#   na: column-name
#   ty: data-type
#   nu: not null (Y/N)
#   id: identity (Y/N)
#   un: unique
#   va: default value
#   fk: foreign key hint
#   cd: cardinality
#   tt: title
#   dc: description

# fixed columns:
# the columns will be appended to each table
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
# -- YAML to Excel
$ dst convert excel -i sample.yml -o sample.xlsx

# -- Excel to YAML
$ dst convert yaml -i sample.xlsx -o sample.yml

# -- YAML to ER diagram definition file
$ dst convert diagram -i sample.yml -o sample.puml

# -- YAML to ER diagram (.png)
# download plantuml.jar from https://plantuml.com/download
$ dst convert diagram -i sample.yml -o sample.puml -j plantuml.jar
# you may convert some of tables only, below select the tables has 'tag' prefix only
$ dst convert diagram -i sample.yml -o sample.puml -j plantuml.jar -p tag
# for simple mode, only show PK and FK in the diagram
$ dst convert diagram -i sample.yml -o sample.puml -j plantuml.jar -p tag --simple

# -- YAML to SQL schema file
# MariaDB
$ dst convert text -i sample.yml -o sample.sql -t mariadb
# MsSQL
$ dst convert text -i sample.yml -o sample.sql -t mssql

# -- Verify foreign key
# make sure the foreign table and key exist
$ dst verify -i sample.yml
```

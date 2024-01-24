# Database Schema Tool

Export to different format from definition file.

## Usage

### Example

Create a definition file (e.g. sample.yml)

```yml
# column definition:
#   na: column-name
#   ty: data-type
#   nu: not null (Y/N)
#   id: identity (Y/N)
#   in: index (Y/N)
#   un: unique (Y/N)
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

# -- YAML to ER diagram definition file
$ dst convert diagram -i sample.yml -o sample.puml

# -- YAML to ER diagram (.png)
# download plantuml.jar from https://plantuml.com/download

# plantuml.jar can be skiped, it try to find the .jar from the PATH variable
$ dst convert diagram -i sample.yml -o sample.puml
# specify plantuml.jar
$ dst convert diagram -i sample.yml -o sample.puml -j plantuml.jar
# you may convert some of tables, select the tables starts with 'tag' only
$ dst convert diagram -i sample.yml -o sample.puml --table 'tag%'
# template file is used for the ER diagram
$ dst convert diagram -i sample.yml -t template.tpl -o sample.puml

# -- YAML to SQL schema file
# generate using template.tpl
$ dst convert text -i sample.yml -o sample.sql -t template.tpl
# generate using template.tpl, select the tables start with 'tag' only
$ dst convert text -i sample.yml -o sample.sql -t template.tpl --table 'tag%'

```

# Database Schema Tool

Export to different format from definition file.

## Usage

### Create a definition file (e.g. sample.yml)

```yml
# column definition:
#   na: column-name
#   ty: data-type
#   nu: is not null (y/n)
#   id: is identity (y/n)
#   va: default value
#   fk: foreign key hint
#   dc: description

# fixed columns:
# the columns will be appended to all tables
fixed:
  - {na: createUser, ty: varchar(20), fk: users.id, dc: create user ID}
  - {na: createTime, ty: datetime, fk: users.id, dc: date time of creation}
  - {na: updateUser, ty: varchar(20), fk: user.id, dc: update user ID}
  - {na: updateTime, ty: timestamp, fk: user.id, dc: date time of update}

schemas:
- name: "the name of schema"
  tables:
  - name: "the name of table"
    desc: "the description of table"
    columns:
    - [...]

  - name: "other table"
  ...
```

### Yaml to Excel

```
$ dst -i sample.yml -o sample.xslx
```

### Excel to Yaml

```
$ dst -i sample.xslx -o sample.yml
```

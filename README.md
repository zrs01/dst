# Database Schema Tool
Export to different format from definition file.

## Usage
### Create a definition file (e.g. sample.yml)

```yml
# column definition:
# column name, data type, is primary key, is unique value, is nullable, foreign hint, comment

# repeat columns for every table
fixColumns: &fc
  - [createUser, varchar(20), n, n, n, users.id, create user ID]
  - [createTime, datetime, n, n, n, users.id, date time of creation]
  - [updateUser, varchar(20), n, n, n, user.id, update user ID]
  - [updateTime, timestamp, n, n, n, user.id, date time of update]

schemas:
  - name: "the name of schema"
    tables:
      - name: "the name of table"
        description: "the description of table"
        columns:
          - [column name, data type, n, n, n, foreign key hint, comment]
          - *fc

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

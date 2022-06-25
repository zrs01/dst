# Database Schema Tool

Export to different format from definition file.

## Usage

### Create a definition file (e.g. sample.yml)

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
#   dc: description

# fixed columns:
# the columns will be appended to all tables
fixed:
  - { na: createUser, ty: varchar(20), fk: users.id, dc: create user ID }
  - { na: createTime, ty: datetime, fk: users.id, dc: date time of creation }
  - { na: updateUser, ty: varchar(20), fk: user.id, dc: update user ID }
  - { na: updateTime, ty: timestamp, fk: user.id, dc: date time of update }

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
### Command

```sh
NAME:
   dst - Database schema tool

USAGE:
   dst [global options] command [command options] [arguments...]

VERSION:
   0.0.1-202204

COMMANDS:
   convert  Convert to other format
   verify   Verify the foreign key
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d    Debug mode (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)


# dst convert - Convert to other format
USAGE:
   dst convert [command options] [arguments...]

OPTIONS:
   --input value, -i value     Input file (source)
   --output value, -o value    Output file
   --template value, -t value  Template file
   --help, -h                  show help (default: false)

# tranformation supported:
- yml to xlsx
  e.g. dst convert -i sample.yml -o sample.xslx

- xlsx to yml
  e.g. dst convert -i sample.xslx -o sample.yml

- yml to txt (using template file), build-in templates: 'mariadb.tpl' and 'sqlserver.tpl'
  e.g. dst convert -i sample.yml -t mariadb.tpl -o sample.sql

# dst verify - Verify the foreign key
Verify the table and column of foreign key
```

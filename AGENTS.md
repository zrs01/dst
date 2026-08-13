# AGENTS.md

## Repository structure detected

**Single-app repository.** This is one Go CLI application (`dst`, module `github.com/zrs01/dst`) with no sub-packages that are independent apps. A single `AGENTS.md` is therefore produced at the repository root. Evidence: one `go.mod`, one `main.go`, and all `internal/`/`model/` packages belong to the same binary.

---

## What this is

`dst` ("Database Schema Tool") is a Go CLI that converts database schema definitions between several formats: a YAML schema definition file, generated SQL DDL, PlantUML ER diagrams, Excel workbooks, and arbitrary text/code via Jet templates. It can also reverse-engineer a live MariaDB/MySQL database back into a YAML schema definition (`export`) and compute a diff between a live database and a schema file (`ddl diff`). It includes a terminal UI (`edit`) for editing a schema file.

The primary tech stack is Go (module `go 1.24.0`), `urfave/cli/v3` for the CLI, `CloudyKit/jet/v6` for templating, `jinzhu/configor` for config loading, `logrus` + `nested-logrus-formatter` + `lumberjack.v2` for logging, `tracerr` for wrapped errors, `samber/lo` for utilities, `excelize/v2` for Excel, `tview`/`tcell` for the TUI, and `go-sql-driver/mysql` for database access. It serves developers who need to maintain schema definitions as source-of-truth YAML and generate DB artifacts from them.

---

## Application entry points

The only entry point is `main.go`, which sets `config.Version` (from build ldflags `-X main.version=...`), configures logrus formatter/level, and builds a single root `cli.Command` named `dst`. It then registers six subcommands from `cmd/dst/*.go` and calls `cmd.Run(context.Background(), os.Args)`. Each subcommand's `Action` calls `config.LoadConfig(...)` to merge config.yml + command-specific settings + CLI flags into the global `config.MergeSetting`, then dispatches to a service function:

- `ddl` (cmd-ddl.go) → `internal/service/ddl` (`CreateDatabase`, `CreateTable`, `DropTable`, `Diff`); subcommands `cd`, `ct`, `dt`, `diff`.
- `erd` (cmd-erd.go) → `internal/service/erd.Generate`.
- `xls` (cmd-xls.go) → `internal/service/xls.Generate`.
- `export` (cmd-exp.go) → `internal/service/exp.Generate`.
- `txt` (cmd-txt.go) → `internal/service/txt.Generate`.
- `edit` (cmd-edit.go) → `internal/service/edit.Launch` (TUI).

There are no background workers, scheduled jobs, or message consumers. On error, `main.go` prints via `tracerr` (source-colored stack only when log level ≥ debug) and exits 1.

---

## Commands

Commands defined by the project's build tooling (`Makefile`):

- `make` / `make default` — builds a host binary to `bin/dst` (`.exe` on Windows). Uses `-ldflags="-X main.version=$(VERSION)"` where `VERSION` comes from `git describe --tags --always --long --dirty`. `CGO_ENABLED=0`.
- `make windows` / `linux` / `darwin` / `arm` — cross-compile to `build/dst_<os>_<arch>` with `-trimpath -ldflags="-s -w -X main.version=..."`.
- `make build` — builds all four platform binaries (windows, linux, darwin, arm) and prints the version.
- `make clean` — removes `build/` binaries and the host `bin/dst`.
- `make help` — lists targets by parsing `##` comments.

Other tooling:
- `modd.conf` — `modd` watches `**/*.go **/*.jet` and runs `make` on change (live-rebuild dev loop; invoked via `support/layout.kdl` inside a zellij layout).
- `.goreleaser.yaml` — release automation config (goreleaser). No CI workflow files exist (no `.github/`).

Runtime CLI commands (registered in `main.go`):
- `dst ddl cd` — create database DDL. Requires `--sdf` (or config `sdf`).
- `dst ddl ct [--alter] [--table]` — create table DDL; `--alter` emits `ALTER TABLE ... ADD COLUMN` instead of a direct create.
- `dst ddl dt [--table]` — drop table DDL (drops indexes + constraints first).
- `dst ddl diff [--dsn] [--sdf] [--table]` — DDL to sync a live DB with the schema file.
- `dst erd -o <file> [--table] [--template] [--plantuml]` — PlantUML output; `-o` is required and must end in `.puml` or `.png`; `.png` requires `--plantuml` (path to `plantuml.jar`) and Java.
- `dst xls [--sdf] [--ccf] [-o] [--table]` — YAML → `.xlsx`; appends common columns.
- `dst export [--dsn] [--ccf] [-o] [--table]` — live DB → YAML; removes common columns.
- `dst txt [--sdf] [-o] [--table] [--template]` — template → text; if no `--template`, outputs YAML.
- `dst edit [--sdf]` — interactive TUI editor.

Common flags: `--sdf`/`-i` (schema definition file), `--table` (wildcard filter `*`/`%`, comma-separated), `--template`/`-t`, `--output`/`-o`, `--dsn`, `--ccf`/`-cc` (common column file). If `-o` is omitted, output goes to stdout.

There are **no tests, lint, or format commands** defined (Makefile `test`/`all` are commented out; no `go test`, `golangci-lint`, or `gofmt` targets).

---

## Required runtime environment

- **MariaDB/MySQL database** — required for `ddl diff` and `export`, and for `ddl cd`/`ct`/`dt` only when the DSN is used to connect (the DDL generators themselves only need a schema file). Connection is via a `mariadb://` DSN (see `bin/config.yml` example: `mariadb://root:pass@tcp(host:port)/dbname`). The `internal/db` service panics if the DSN does not start with `mariadb://`. Missing/unreachable DB breaks these commands.
- **`config.yml`** in the working directory — loaded by `config.LoadConfig` via `config.Filename = "config.yml"`. Missing file is silently tolerated (`configor.Silent: true`), so it is optional at runtime, but without it DSN/SDF defaults must come from CLI flags.
- **Java + `plantuml.jar`** — required only for `erd` with `.png` output; `srv-erd.go` shells out to `java -jar <plantuml>`. For `.puml` output no Java is needed.
- **Interactive terminal** — required for `edit` (tview TUI).
- **Go toolchain** — `go.mod` declares `go 1.24.0`; the installed toolchain on this machine is `go1.26.5` (build succeeded with it). No Docker, Redis, Kafka, or cloud dependencies exist.

Startup order: none beyond the CLI itself; each subcommand loads config and connects to the DB lazily.

---

## Project structure

```
AGENTS.md                     this file
main.go                       CLI entry point; registers subcommands
go.mod / go.sum               module deps (go 1.24.0)
Makefile                      build targets (host + cross-compile)
modd.conf                     modd live-rebuild watch config (**/*.go **/*.jet)
.goreleaser.yaml              release automation
.gitignore                    ignores bin/, build/, dist/, vendor/, etc.
README.md                     usage + data-model docs
LICENSE                       (Apache-2.0 style, 2022)

cmd/dst/                      CLI command registrations (one file per subcommand)
  cmd-ddl.go                  ddl cd/ct/dt/diff (large commented-out blocks)
  cmd-erd.go, cmd-xls.go, cmd-exp.go, cmd-txt.go, cmd-edit.go

config/                       config loading + logging setup
  config.go                   ConfigDef/SettingDef, LoadConfig, merge-struct
  logging.go                  logrus formatter/level/file output; LOG_LEVEL env

model/                        schema data model (YAML tags, custom MarshalYAML)
  model.go                    getContent() reflection helper
  model-root.go, model-schema.go, model-table.go, model-column.go,
  model-index.go, model-reference.go, model-routine.go,
  model-foreign-table.go      data structures
  model-utils.go              Verify() integrity checks; isNumeric/toInt

internal/
  db/                         DB service abstraction
    service.go                factory: NewService() → dbcm.Service (MariaDB only)
    dbcm/                     constants (DSN prefixes) + Service/DDL interfaces
    mariadb/                  MariaDB implementation
      service.go              Load() reads live DB → model.Schema
      builder-ddl.go          DDLBuilder: CREATE/ALTER/DROP statements
      builder-export.go       ExportBuilder: DB → schema model
      def-*.go                INFORMATION_SCHEMA managers (table, column, index,
                              constraint, database, routine, view, trigger, sequence)
  dbm/                        DEAD CODE: old MySQL/MSSQL DataService (see Deprecated)
  flagbuilder/                builder wrappers for cli flags (string/bool)
  ddwriter/                   YAML writer (OutputYml); XLSX reader is dead code
  service/
    loader.go                 LoadBuilder: load YAML, filter, build references
    common-column.go          RemoveCommonColumns / AppendCommonColumns
    ddl/srv-ddl.go            DDL generation + diff logic
    erd/srv-erd.go            ER diagram (.puml/.png via plantuml)
    exp/srv-exp.go            DB → YAML export
    txt/srv-txt.go            Jet template rendering + global template funcs
    xls/srv-xls.go            YAML → Excel workbook
    edit/                     TUI editor (ui-main, ui-table, ui-column,
                              dialog/ui-dialog.go, app/ui-common.go)
    sql/                      DEAD CODE: legacy Jet SQL templates (see Deprecated)

template/                     user-facing example Jet templates
  mariadb.tpl, mssql-create.tpl, mssql-alter.tpl, dto.jet, model.jet

util/                         helpers (wildcard match, YAML unmarshal, DB open,
                              condition builder, IsYes); datautil.go is dead code

support/layout.kdl            zellij layout running modd (dev loop)
example/                      sample schema files + generated outputs (sample.yml,
                              sample2.yml, sample.puml, sample.sql, sample.xlsx)
bin/  build/  dist/           gitignored build/artifact outputs (do not edit)
```

Omitted (generated/ignored): `bin/`, `build/`, `dist/`, `.git/`, `node_modules`-equivalents. `bin/config.yml` is a local, gitignored dev config file, not tracked.

---

## Architecture notes that aren't obvious from filenames

**Config merging.** `config.LoadConfig(configType, source)` (config/config.go) loads `config.yml` into `config.Default` via `jinzhu/configor`, then uses `geraldo-labs/merge-struct` (`mp.Struct`) to merge three layers into the global `config.MergeSetting`: (1) global `Dsn`/`Sdf`/`Ccf`, (2) the type-specific settings block (`DataDef`, `Erd`, `Xls`, `Export`, `Text`), (3) the CLI flag struct passed as `source`. CLI flags therefore override config.yml, which overrides globals. `EmptyConf` (used by `edit`) skips the type-specific merge.

**Loader/filter pipeline.** `internal/service/loader.go` uses a functional-options builder (`NewLoadBuilder(WithTablePattern, WithColumnPattern, ...)`) to read a YAML schema file into `model.Schema`, then `model.Schema.Filter(tablePattern, columnPattern)` (model/model-schema.go) applies wildcard matching via `util.WildCardMatchWithCommaPattern` (supports `*`, `%`, and comma-separated lists). `schema.Filter` returns an error when no tables/columns match — callers treat that as fatal.

**Reference graph.** `LoadBuilder.updateReferenceTables` scans every column's `ForeignKey` (format `table.column`), looks up the referenced table, and appends a `Reference` entry to that table's `References` list (reverse mapping). This is what feeds `References` in templates. `model.Verify` (model/model-utils.go) validates that column names/data types exist and that each `fk` target exists.

**DB abstraction is MariaDB-only in practice.** `internal/db/service.go` switches on `config.MergeSetting.Dsn`; only the `mariadb://` branch is wired (`mariadb.NewMariadbService()`). MySQL/MSSQL branches are commented out, and `db.NewService()` panics for any other prefix. The `dbcm.Service` interface exposes `Load(tableFilter)` and `DDLBuilder()`; `dbcm.DDL` declares many methods (create/alter/drop table, add/drop column, index, constraint) — only a subset is implemented by `mariadb.DDLBuilder`.

**DDL generation.** `internal/db/mariadb/builder-ddl.go` builds MariaDB DDL strings directly (no templating). Notable behaviors: `CREATE TABLE ... IF NOT EXISTS`, `ALTER TABLE ... ADD COLUMN IF NOT EXISTS ... AFTER <prev>`, `CREATE INDEX IF NOT EXISTS`, foreign keys named `fk_<table>_<column>`, indexes named `idx_<table>_<column>` (auto-named `idx_<table>_<NNN>` for multi-column indexes), and `MODIFY COLUMN IF EXISTS` for diffs. `--alter` mode creates only identity columns in the base table and adds the rest via `ADD COLUMN`.

**Diff algorithm.** `internal/service/ddl/srv-ddl.go` `Diff()` loads the live DB schema (via `db.NewService().Load`) and the schema file, then per table: creates new tables, adds missing columns (with an index if `column.Index` is yes), modifies columns when data type / not-null / desc / unique differ, and drops tables absent from the schema file. Column comparison uses `lo.If(dbColumn.NotNull == "N", "").Else(...)` to normalize the `N` sentinel.

**Templating.** `internal/service/txt/srv-txt.go` registers global Jet functions used by all templates: `toCamel`, `toLowerCamel` (via `iancoleman/strcase`), `toPlural`/`toSingular` (via `gertd/go-pluralize`), `toJavaType` and `toTypescriptType` (string-matching on data type substrings). `erd` embeds `templates/default.jet` via `embed.FS`; `sql` embeds its Jet templates the same way. Output goes to stdout when `-o` is empty.

**YAML schema format.** Column fields use short YAML tags (`na`, `ty`, `id`, `nu`, `un`, `va`, `fk`, `cd`, `tt`, `in`, `dc`, `cm`, `computetype`). `model.Column.MarshalYAML` uses `getContent` (model/model.go) to emit only non-empty fields, drops `N`-valued `Identity`/`Unique`/`NotNull`, and renders `compute` expressions as YAML literal style. `ddwriter.restoreFixColumns` (writer-yml.go) detects columns identical across all tables and hoists them into `Root.Fixed`. `expandFixColumns` in the loader is currently a no-op (body commented out), so fixed columns are not actually expanded on load.

**Logging.** `config/logging.go` sets a `nested.Formatter` (colors disabled at Info level), reads `LOG_LEVEL` env to set the logrus level, and — if `logging.output` is set in config.yml — writes to both stderr and a `lumberjack` rolling file (`io.MultiWriter`). `main.go` prints full source-colored stack traces only when level ≥ debug; otherwise `tracerr.Sprint(err)` (message only).

**TUI editor.** `internal/service/edit/` uses `tview` with a `Pages` container: a table widget lists tables, F2 opens an edit dialog, selecting a table opens a column widget, Esc pops pages. It edits the in-memory `model.Schema` only; there is no save-to-file path in the code inspected.

**Error handling.** All service functions return errors wrapped with `tracerr.Wrap`; the CLI action handlers wrap again with `tracerr.Wrap`. Validation errors from `model.Verify` are printed then returned as `tracerr.Errorf("invalid data")`.

---

## Environment variables

Only one environment variable is read at runtime:

- **`LOG_LEVEL`** — sets the logrus level (e.g. `debug`, `trace`, `info`). Optional; default is `logrus.InfoLevel`. Consumed in `config/logging.go` `SetLogLevel()` (line ~43, `os.Getenv("LOG_LEVEL")`). Gotcha: if the value fails `logrus.ParseLevel`, the error is logged and Info is kept; at Trace level `logrus.SetReportCaller(true)` is enabled, which adds caller info to logs.

The cli flag builders have `WithEnvVars` methods defined but **commented out**, so no flags read env vars. `config.Version` is set from build ldflags (`-X main.version=...`), not an env var. No other `os.Getenv`/`os.LookupEnv` calls exist in the codebase (verified by grep).

---

## Conventions to preserve

**Naming conventions.** Files are prefixed by role: `cmd-*` (CLI commands), `srv-*` (service entry), `def-*` (INFORMATION_SCHEMA managers), `builder-*` (DDL/export builders), `db-*` (DB services/models), `writer-*` (output writers), `flag-*`/`*-flag-builder.go` (flag builders), `ui-*` (TUI widgets). Packages live under `internal/` or top-level `model/`, `config/`, `util/`, `cmd/`.

**Builder pattern.** Several layers use a fluent builder with `With*` methods returning the receiver: `flagbuilder` flag builders, `service.NewLoadBuilder` (functional options), `mariadb` `ExportBuilder`/manager types (`WithDriverName`, `WithDataSourceName`, `WithTableSchema`, ...). New builders should follow this shape.

**YAML short-tag schema.** The `model.Column` YAML tags (`na`, `ty`, `id`, `nu`, `un`, `va`, `fk`, `cd`, `tt`, `in`, `dc`, `cm`) are a fixed convention used by README, templates, and the export writer. `Y`/`N` string sentinels represent booleans (`Identity`, `NotNull`, `Unique`); `util.IsYes` interprets "contains Y".

**Flag building.** All CLI flags are constructed via `flagbuilder` builders with `WithDestination(&options.Field)` so that `config.LoadConfig` receives a populated `SettingDef`. Aliases are set with `WithAliases` (e.g. `-i` for `--sdf`, `-o` for `--output`, `-t` for `--template`).

**Templating.** Jet templates are the standard for SQL/ERD/code generation; template global functions are registered centrally in `txt.setJetFunc`. Embedded templates are loaded via `embed.FS` (`//go:embed`); user templates via `jet.NewOSFileSystemLoader`.

**Error wrapping.** Errors are consistently wrapped with `tracerr.Wrap` and re-wrapped at the CLI boundary; `main.go` decides stack-trace verbosity from the log level.

**Logging.** Logrus is the single logging library; `config/logging.go` centralizes formatter/level/output. No `fmt.Println`-based debug logging in service code is used for diagnostics (though `dbm` dead code does use `spew.Dump`).

---

## Guidance for future agents

**Dead-code packages (`internal/dbm`, `internal/service/sql`, `internal/ddwriter/writer-xlsx.go`, `util/datautil.go`).** These compile but are unreachable: `dbm` is referenced only in comments in `loader.go`; `sql` only in commented-out blocks in `cmd-ddl.go`; `ReadXlsx` has no callers; `datautil.go` is fully commented out. Do not treat them as live features. If a task appears to require them, first confirm the caller exists — it likely does not. Avoid "fixing" commented-out code; it is intentionally disabled.

**Large commented-out CLI blocks in `cmd/dst/cmd-ddl.go`.** Roughly half the file is commented-out subcommands (`ac`, `dc`, `rc`, `mc`, `ci`, `di`) that reference the dead `sql` package. Do not uncomment them blindly; they reference APIs (`service.NewLoadBuilder().Filter`, `sql.AddColumn`, etc.) that no longer match current signatures.

**DB-dependent commands need a live MariaDB.** `ddl diff` and `export` connect to a real DB via `--dsn`/config `dsn`. They cannot be verified without a reachable `mariadb://` instance. `db.NewService()` panics if the DSN doesn't start with `mariadb://` — don't expect graceful handling for `mysql://`/`sqlserver://` DSNs even though the flag help mentions them.

**`erd` with `.png` needs plantuml.jar + Java.** `srv-erd.go` shells out with `sh.Command("java", "-jar", ...)` and writes a temp `output.puml` in the CWD (deleted after). Missing jar/Java yields a confusing failure; the code only checks the jar path exists, not that Java is installed.

**`edit` is an interactive TUI.** It requires a terminal and blocks until the app exits. It edits only the in-memory model — there is no file-save path in the code inspected, so don't claim it persists changes.

**Config file gotchas.** `config.yml` is loaded from the current working directory (not `bin/`). `bin/config.yml` is a gitignored local dev file with a hard-coded DB DSN — don't commit it or rely on it existing. `configor.Silent: true` means a missing config.yml is silently ignored; defaults come only from CLI flags.

**`make` version ldflags.** `Makefile` uses `git describe --tags --always --long --dirty`; on a repo with no tags this produces a long hash. Cross-builds use `-trimpath -s -w`. `modd`/`support/layout.kdl` trigger `make` on every `.go`/`.jet` change — a full rebuild on each edit.

**No tests exist.** There are no unit/integration tests, so behavior is verified by building (`go build ./...`) and manual CLI runs. When changing logic, build first; expect to validate against a real DB or schema file manually.

**`go.mod` vs toolchain.** `go.mod` declares `go 1.24.0`; the installed toolchain may be newer (verified `go1.26.5` here). Keep language features compatible with the declared version.

---

## Deprecated, stale, or unused components

Evidence-based; speculation is flagged where applicable.

- **`internal/dbm/`** — dead. `dbm.go`, `db-mysql*.go`, `db-mssql*.go` define `DataService` implementations for MySQL/MSSQL, but grep shows the only references are commented-out lines in `internal/service/loader.go` (lines ~135, ~158). No live caller.
- **`internal/service/sql/` + its Jet templates** — dead. `sql.go` and the `templates/mariadb|mssql/*.jet` files are referenced only by commented-out CLI blocks in `cmd/dst/cmd-ddl.go`. No live caller (grep confirmed zero non-comment references).
- **`internal/ddwriter/writer-xlsx.go` (`ReadXlsx`)** — dead. No callers anywhere in the codebase; the constants `CName`..`CDesc` are used only by that function.
- **`util/datautil.go`** — dead. The entire file body is commented out (`FilterData`).
- **`internal/ddwriter/writer-option.go`** — an empty `options` struct with no fields; vestigial.
- **`flagbuilder.SchemaNameFlag`** — defined but only referenced in commented-out code in `cmd-ddl.go`; no active `--schema` flag exists.
- **`internal/db/dbcm/constants.go` `MYSQL_PREFIX` / `MSSQL_PREFIX`** — defined but unused by the active switch in `internal/db/service.go`, which only handles `MARIADB_PREFIX`; the MySQL/MSSQL branches are commented out.
- **`model/model-table.go` and `model/model-root.go` `MarshalYAML`** — commented out; only `Column` and `ForeignTable` implement custom marshaling via `getContent`.
- **`service.LoadBuilder.expandFixColumns`** — a no-op; its body is commented out, so `Root.Fixed` columns are not expanded into tables on load (the `mariadb.tpl` template still references `fixed`).
- **Commented-out `LoadFromDB`/`loadFromMysql`/`loadFromMssql`/`loadFromMariadb` in `loader.go`** — disabled DB-loading paths superseded by `internal/db` + `internal/db/mariadb`.
- **`dist/`, `bin/`, `build/`** — generated/ignored artifact directories (gitignored), not source.

Nothing was found that is actively broken; the above are disabled or unreachable code paths rather than failing features.

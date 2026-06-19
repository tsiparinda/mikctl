# mikctl

MikroTik fleet management utility.

Features:

* Router inventory
* RouterOS version tracking
* Remote command execution
* RouterOS script import
* Neighbor discovery
* Router grouping
* SQLite inventory database

## Directory structure

```text
backup/                 Router backups
db/                     SQLite database
scripts/                RouterOS scripts and templates
src/                    Source code
~/.config/mikctl/       Configuration
```

Configuration file:

```text
~/.config/mikctl/config.yaml
```

## Build

```bash
go build -o mikctl ./src
```

## Requirements

FreeBSD:

```bash
pkg install sqlite3 sqlitestudio
```

## Configuration

Example:

```yaml
db_path: db/mikrotik.db
workers: 10
ssh_user: admin
password_id: 0
```

## Commands

Show help:

```bash
./mikctl --help
./mikctl update --help
```

Show router count:

```bash
./mikctl count
```

List routers:

```bash
./mikctl list
```

Update router information:

```bash
./mikctl update
./mikctl update -r mt24_gippo03
./mikctl update -g v6
```

Run commands:

```bash
./mikctl run scripts/commands/show_version.txt
```

Import RouterOS script:

```bash
./mikctl import scripts/commands/backup/backup2ftp_v6.rsc
./mikctl import scripts/commands/backup/backup2ftp_v6.rsc -V 6
```

Discover MikroTik neighbors (check neighbors of main routers only!):

```bash
./mikctl discovery
```

## Global flags

```text
-g, --group string     Router group
-r, --router string    Router name
-s, --site string      Site
-v, --verbose          Verbose output
-V, --ros string       RouterOS major version (6 or 7)
-w, --workers int      Parallel workers
-m, --main int         1=main, 0=slave
```

## Examples

Update all RouterOS 6 devices:

```bash
./mikctl update -V 6
```

Import script to all main routers:

```bash
./mikctl import scripts/example.rsc -m 1
```

Import script to all slave routers:

```bash
./mikctl import scripts/example.rsc -m 0
```

Update one router:

```bash
./mikctl update -r mt24_gippo03
```
## License

MIT License. See LICENSE file for details.
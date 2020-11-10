![pg-mig master](https://github.com/djordjev/pg-mig/workflows/pg-mig%20master%20test%20&%20build/badge.svg)

# pg-mig
Database-centric tool for running migrations against PostgreSQL database.

This tool is somehow similar to git but for a database. Each database state is described as a serie of
revisions (called migration). So it allows user to change between those revisions (similar to 
`checkout` previous commit in git). It's compiled to executable to it doesn't require any runtime environment.
Migrations are identified by creation timestamp, and the current state of applied migrations is stored
in the database itself, so it's safe for developers to create multiple migrations in different git branches
and later merge them together preserving an order of execution. 

## Basic principles
Each incremental database upgrade is stored in 2 `.sql` files. A file with suffix `_up.sql` is used
to upgrade database to a next revision. The file with suffix `_down.sql` is used to downgrade to a previous
revision. That way each database state can be reached as sequence of consecutive upgrades/downgrades.

The current database state is stored in the meta-table `__pg_mig_meta`. The table represents a list of
migrations that are applied to current database state. Thus, `pg-mig` can determine which migrations should
be executed by comparing its content with files present in a workspace folder.

To print existing commands you can run `./pg-mig help`. For each particular command you can get additional
information with a list of command flags by running `./pg-mig command -h`.

## Commands

### init
Before using `pg-mig` a user has to initialize it first. Command `init` takes a form:

```shell
./pg-mig init -db=localhost -name=main_db -credentials=postgres:pg_pass -path="~/workspace_folder" -ssl=disable -port=5432 
```

As a result it outputs config file in `json` format in the current working directory. This file will be named
`pgmig.config.json` and can be edited manually if something in database setup changes. Most of the commands 
will use this file to establish connection to the database, and it must be always accessible (have right access
permissions).

**Available flags for `init` command:**
- *path* - Path to a workspace directory. It's a directory where all migration files should be stored. If omitted
defaults to current directory.
- *db* - Address of the database server. If omitted, defaults to `localhost`.
- *name* - The name of database against which migrations will run. This is the required flag.
- *credentials* - Also required flag. It expects credentials to connect on the database in form `username:password`.
It will be used in building postgres connection string so all PostgreSQL connection credentials are supported.
- *ssl* - ssl flag from PostgreSQL connection string. Defaults to `disabled`.
- *port* - Port on which PostgreSQL server is running. If omitted default PostgreSQL port `5432` will be used.
- *nocolor* - By default `pg-mig` uses different colors and emojis to print different types of messages. Setting
this flag will force pure textual output.

### add
Command creates two empty files for up and down migrations. Files need to be in particular format so this is 
a suggested way for creating a new migration.

```shell
./pg-mig add -name="migration name"
```
**Available flags for `add` command:**
- *name* - This is the only argument `add` command accepts. If passed it will be included in file names. Its 
purpose is only visual, to more easily detect what is migration supposed to do. It can be safely omitted.

### run
The main command in `pg-mig` as it executes migrations until provided time. So depending on current state in 
the database it can execute up migrations or down migrations (or even both of them) to bring the database
in the state for given timestamp. Since database state is stored in the database itself, it's completely valid
to have two developers working in a separate branches, adding their own migrations. Once branches are merged
on next `run` command execution, all migration files found on a filesystem that has not been executed will 
be executed. It's similar to git `checkout` command.
```shell
./pg-mig run -time="2010-09-20T15:04:05Z"
```


**Available flags for `run` command:**
- *time* - Accepts time on which database needs to be reset or updated. It expects time in various formats.
If not provided it will execute all available `up` until current time. There are 2 special values: `push` and
`pop`. Push command will execute the next `up` migration after the last one that already has been applied to
the database. Pop command will do the opposite, execute `down` migration for the last one that has been applied
to the database.
- *dry-run* - This flag is used only for a testing purposes. Use it to print which migrations would be executed
without applying them. 

Formats accepted for *time* flag:
- *2006-01-02T15:04:05Z07:00* - RFC3339 format.
- *2006-01-02T15:04:05* - date without a timezone, UTC timezone will be used.
- *2006-01-02* - only a date 00:00 time will be used (also UTC).
- *3:04PM* - time only (in UTC), for date it assumes current date.
- *push* - push command described in the paragraph above.
- *pop* - pop command described in the paragraph above.
- *1604752594* - unix timestamp

`run` command can take a simple form
```shell
./pg-mig run
```
The command in this form will execute all available up migrations and bring the database to the latest state.
It's useful to run it to ensure the database is up to date with all existing migrations.

### squash
This command is similar to git squash. During the time it's possible that there will be a lot of migration
files. After some time there might be no need for a fine-grained moving between some of them. Such migrations
can be merged into one with just 2 files (up and down migration files). **Use with caution. Once migrations
are squashed there's no way to bring them back to original form, and further on they will be considered a
single migration.**

```shell
./pg-mig squash -from="2010-09-10T15:04:05Z" -to="2010-09-20T15:04:05Z"
```

**Available flags for `squash` command:**
- *from* - start date for squash. 
- *to* - end date for squash.

Note: for squash command both *from* and *to* values are inclusive (meaning if there's a migration with
exact the same time as in the flag it will be included in squash). 

### log
Similar to git log command. Prints migrations present on filesystem and those that are already applied
to the database.

```shell
./pg-mig log
```

This command does not accept any flags.

## Usage with docker
When running PostgreSQL in docker container it can be handy to have `pg-mig` installed directly in container.
That way it's not needed to have `pg-mig` installed on development machine. Docker multi-stage builds come 
handy in situations like this.

First create a `Dockerfile`

```dockerfile
FROM golang:1.15.4-alpine3.12 AS builder

RUN apk update && apk add git

WORKDIR "/"
RUN ["git", "clone", "https://github.com/djordjev/pg-mig"]

WORKDIR "/pg-mig"
RUN ["go",  "build",  "-o", "./build/pg-mig", "./cmd/pg-mig/main.go"]

FROM postgres:13.0-alpine

COPY ./workspace/ /usr/pg-mig/workspace/

COPY --from=builder /pg-mig/build/pg-mig /usr/pg-mig/pg-mig

RUN chmod -R 777 /usr/pg-mig

EXPOSE 5432
```

Then in `docker-compose`:
```yaml
version: "3.7"

services:
  db:
    build:
      context: ./db
    container_name: db
    ports:
      - 5432:5432
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=pg_pass
      - POSTGRES_DB=main_db
    volumes:
      - './db/workspace:/usr/pg-mig/workspace'

```

Few things to note:
1. We use golang container to build `pg-mig` and copy executable into final `postgres` container. That way
final container is becoming more lightweight, since everything form `builder` container is thrown out except
the main executable.
2. Binding local `workspace` folder to `/usr/pg-mig/workspace` will make all files in container accessible
on the local filesystem. This is important with creating new migrations.

After a container has been run with 
```shell script
docker-compose up db
```
you can exec into it
```shell script
docker exec -it db sh
```

move to `usr/pg-mig`
```shell script
cd /usr/pg-mig
```

run init
```shell script
./pg-mig init -name=main_db -credentials=postgres:pg_pass -path=./workspace
```

run add 
```shell script
./pg-mig add -name="some name for migration"
```

execute migrations
```shell script
./pg-mig run
```
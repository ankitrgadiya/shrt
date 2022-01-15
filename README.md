# Shrt - A "go" short-link service

This is my fork of the [kellegous/go](https://github.com/kellegous/go) server
[[Original
Readme](https://github.com/kellegous/go/blob/043936a95042fbbdcec4ff5263325607797ea54a/README.md)]

## Differences

* SQLite based backend
* Uses Go's builtin embedding for static assets
* Command-line interface for CRUD to the server as well as database directly 

## Installation

Shrt can directly be installed using the Go toolchain by running the following
command.

``` sh
go install argc.in/shrt@latest
```

The SQLite library [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
supports different compilation strategies. By default, Go will try to compile
the the SQLite library using the available C Compiler toolchain. On most linux
systems that would be GLibC and GCC. This can be configured by build and link
flags described [here](https://github.com/mattn/go-sqlite3#compilation).

I personally prefer a statically linked binary that can be compiled using Musl
and GCC. On Alpine Linux, it can be done by running the following command. On
other Linux distributions, the CC variable can be configured to point to the
Musl GCC binary. Also note the "-s -w" flags to strip off the symbol tables and
debugging information to reduce the size.

``` sh
CC=/usr/bin/x86_64-alpine-linux-musl-gcc go build --ldflags '-linkmode external -extldflags "-static" -s -w' -o shrt main.go
```

Alternatively, I also provide a Docker image through Github Packages with
precompiled binary.

``` sh
docker run -it \
    -v /path/to/data:/data \
    -p 8080:8080 \
    ghcr.io/ankitrgadiya/shrt@master serve --database /data/routes.db --addr 0.0.0.0:8080
```

## Backup

SQLite provides backup API that is available through its command-line tool. This
can be used to backup the database.

``` sh
sqlite3 /data/routes.db ".backup /path/to/backup.db"
```

Alternatively, Shrt can be backed up remotely by using the Shrt command-line
tool. It generates a tab-separated list of all the links configured. SQLite's
command-line tool can import the TSV data directly for restoring.

``` sh
shrt list --server https://SERVER_ADDRESS
```


webdir      := join(justfile_directory(), "web")
staticdir   := join(justfile_directory(), "internal", "web", "static")

default:
    @just --list

assets:
    #!/usr/bin/env bash
    set -euo pipefail

    tmp=$(mktemp)
    tsc --out $tmp {{webdir}}/edit.ts
    closure-compiler --js $tmp --js_output_file {{staticdir}}/edit.js

    sass --no-source-map --style=compressed {{webdir}}/edit.scss {{staticdir}}/edit.css
    sass --no-source-map --style=compressed {{webdir}}/links.scss {{staticdir}}/links.css

    cp {{webdir}}/*.html {{staticdir}}/
    cp {{webdir}}/*.svg {{staticdir}}/

build: assets
    go build -o golinks main.go

serve: assets
    go run main.go serve

clean:
    rm -f {{justfile_directory()}}/golinks

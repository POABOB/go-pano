#!/bin/sh

sm=$(which "sql-migrate")
if [ -z sql-migrate ]; then
    go get -u github.com/go-sql-driver/mysql
    go get github.com/rubenv/sql-migrate/...@latest
fi

option="-config=config.yml -env=database"

echo $sm up $option

case "$1" in
    "new")
    $sm new $option $2
    ;;
    "up")
    $sm up $option
    ;;
    "redo")
    $sm redo $option
    ;;
    "status")
    $sm status $option
    ;;
    "down")
    $sm down $option
    ;;
    *)
    echo "Usage: $(basename "$0") new {name}/up/status/down"
    exit 1
esac
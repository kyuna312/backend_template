#!/bin/bash

case $1 in
start)
  go run main.go
  ;;
doc)
  swag init
  ;;
start:doc)
  swag init && go run main.go
  ;;
build)
  go build
  ;;
init)
  go mod init
  ;;
*)
  echo "start, doc or start:doc"
  ;;
esac

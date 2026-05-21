#!/usr/bin/env bash

go build -o dest/web cmd/web/main.go

./dest/web

#! /bin/bash

go generate ./...
go test ./... -gcflags=all=-l

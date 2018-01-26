#!/bin/bash

docker run --rm -v $(pwd):/opt -w /opt golang:latest /bin/sh -c "\
go get github.com/olekukonko/tablewriter &&\
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/isearch-amd64 . &&\
CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o build/isearch-darwin .
"
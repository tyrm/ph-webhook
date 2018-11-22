#!/bin/bash

rm -Rf ph-warehouse

go get github.com/gobuffalo/packr
go get github.com/google/jsonapi
go get github.com/gorilla/mux
go get github.com/juju/loggo
go get github.com/lib/pq
go get github.com/rubenv/sql-migrate

CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o ph-warehouse .

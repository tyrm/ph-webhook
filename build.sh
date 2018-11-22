#!/bin/bash

rm -Rf ph-warehouse

/usr/local/go/bin/go get github.com/gobuffalo/packr
/usr/local/go/bin/go get github.com/google/jsonapi
/usr/local/go/bin/go get github.com/gorilla/mux
/usr/local/go/bin/go get github.com/juju/loggo
/usr/local/go/bin/go get github.com/lib/pq
/usr/local/go/bin/go get github.com/rubenv/sql-migrate

CGO_ENABLED=0 GOOS=linux /usr/local/go/bin/go build -v -a -installsuffix cgo -o ph-warehouse .

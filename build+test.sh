#!/bin/bash -e
PATH=$HOME/gopath/bin:$GOPATH/bin:$PATH

if ! type -p goveralls; then
  echo go get github.com/mattn/goveralls
  go get github.com/mattn/goveralls
fi

go test -v -race .
go test -v -covermode=count -coverprofile=date.out .
go tool cover -func=date.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=date.out -service=travis-ci -repotoken $COVERALLS_TOKEN

go vet ./...

## Copyright 2020 Tokenomy Technologies Ltd. All rights reserved.
## Use of this source code is governed by a MIT-style license that can be
## found in the LICENSE file.

.PHONY: all lint test

all: build lint test

build:
	go build ./...
	go test -run=noop ./...

lint:
	-golangci-lint run ./...

test:
	CGO_ENABLED=1 go test -race -p=1 ./...

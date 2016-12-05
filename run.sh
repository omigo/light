#!/bin/bash

set -e

# Install yan
go-bindata template.txt
go install

# Use persist to generate code
go generate example/mapper/model.go

# Test
go test -v example/mapper/*.go

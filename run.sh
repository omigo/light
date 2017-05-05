#!/bin/bash

echo 'Install light:'
go install

echo 'Generate implementation file:'
go generate example/mapper/model.go

echo 'Run unit test:'
go test -v example/mapper/*.go

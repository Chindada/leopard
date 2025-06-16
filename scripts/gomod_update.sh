#!/bin/bash

rm go.mod
rm go.sum

go mod init github.com/chindada/leopard
go mod tidy

git add go.mod go.sum

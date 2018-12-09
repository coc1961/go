#!/bin/bash
rm /tmp/cover.out
rm /tmp/cover.html
go test -coverprofile /tmp/cover.out
go tool cover -html=/tmp/cover.out -o /tmp/cover.html
/opt/google/chrome/chrome /tmp/cover.html
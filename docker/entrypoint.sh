#!/bin/sh
set -x

go mod download && reflex -s -r .go go run "$1"

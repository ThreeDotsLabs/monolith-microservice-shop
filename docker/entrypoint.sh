#!/bin/sh
set -x

dep ensure -v && reflex -s -r .go go run "$1"

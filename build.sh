#!/bin/bash

docker run --rm -it -v $(pwd):/go -w /go/src golang:1.9.3 "${@}"

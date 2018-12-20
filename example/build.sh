#!/usr/bin/env bash

set -e

# 基本环境变量
SCRIPT=`dirname $0`
cd ${SCRIPT}
ROOT=`pwd`
DIST=${ROOT}/example.exe

CGO_ENABLED=1 go build -tags debug -ldflags="-H=windowsgui" -o ${DIST} blink/example/cmd

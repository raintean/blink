#!/usr/bin/env bash

set -e

# 基本环境变量
SCRIPT=`dirname $0`
cd ${SCRIPT}
ROOT=`pwd`

cp blink32.dll blink.dll
upx blink.dll
go-bindata --nocompress -o blink_dll_386.go -pkg dll ./blink.dll

cp blink64.dll blink.dll
upx blink.dll
go-bindata --nocompress -o blink_dll_amd64.go -pkg dll ./blink.dll

rm blink.dll
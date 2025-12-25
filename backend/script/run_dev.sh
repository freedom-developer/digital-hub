#!/bin/bash

cur_dir=$(cd $(dirname $0); pwd)
src_dir=$cur_dir/../src

set -e -x

function run()
{
    cd $src_dir
    go run main.go
}

run
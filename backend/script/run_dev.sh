#!/bin/bash

cur_dir=$(cd $(dirname $0); pwd)
src_dir=$cur_dir/../src

function run()
{
    cd $src_dir
    go run main.go
}

run
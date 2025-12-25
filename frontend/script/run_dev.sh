#!/bin/bash

cur_dir=$(cd $(dirname $0); pwd)
src_dir=$cur_dir/../src

set -e
set -x

# nvm最新版本

# node最新稳定版

function run()
{
    cd $src_dir
    [ ! -d node_modules ] && npm install && npm install axios
    npm run dev
}

run
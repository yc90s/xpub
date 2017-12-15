#!/bin/sh

command -v go > /dev/null 2>&1
if [ $? != 0 ]; then
    command -v apt-get > /dev/null 2>&1
    if [ $? == 0 ]; then
        sudo apt-get install golang
    else
        command -v yum > /dev/null 2>&1
        if [ $? == 0 ]; then
            sudo yum install golang
        fi 
    fi 
fi 

# check golang 
command -v go > /dev/null 2>&1
if [ $? != 0 ]; then
    echo "golang install failed"
    exit 1
fi

DIRNAME=$(pwd)
export GOPATH=$DIRNAME

if [ ! -d $GOPATH/bin ]; then
    mkdir $GOPATH/bin
fi
export GOBIN=$GOPATH/bin

set -x
go get github.com/mattn/go-colorable
go get github.com/peterh/liner
go get golang.org/x/crypto/ssh
go get golang.org/x/crypto/ssh/terminal

cd src && go build -o $GOBIN/xpub
chmod u+x $GOBIN/xpub

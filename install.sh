#!/bin/sh
DIR=$HOME/.gitdo
mkdir $DIR
cp -r ./plugins $DIR
cp -r ./hooks $DIR

GOARCH=`uname -p`
GOOS=`uname`
VERSIONTOCP="gitdo_${GOOS}_${GOARCH}"
cp $VERSIONTOCP /usr/local/bin/gitdo

#!/bin/sh
DIR=$HOME/.gitdo
mkdir $DIR
cp -r ./plugins $DIR
cp -r ./hooks $DIR
cp ./secrets.json $DIR

#!/usr/bin/env bash

HOME=`pwd`

echo "Running service to SYNC HDFS to Redis"

echo $HOME

if [ ! -d "$HOME/bin/" ]; then
  # Control will enter here if $DIRECTORY doesn't exist.
  mkdir $HOME/bin
fi

cd $HOME/syncer/
go build
cp syncer $HOME/bin/

cp config/config* $HOME/bin/

cd $HOME/bin/

nohup ./syncer &
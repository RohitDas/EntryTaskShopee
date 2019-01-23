#!/usr/bin/env bash
HOME=`pwd`
echo $HOME

if [ ! -d "$HOME/bin/" ]; then
  # Control will enter here if $DIRECTORY doesn't exist.
  mkdir $HOME/bin
fi

cd $HOME
go build
cp httpapp $HOME/bin/

cd $HOME/LoadBalancer/; go build
cp LoadBalancer $HOME/bin/

cd $HOME/StressTest/; go build
cp StressTest $HOME/bin/

cd $HOME/bin/

nohup ./httpapp 8081 &
nohup ./httpapp 8082 &
nohup ./httpapp 8083 &

nohup ./LoadBalancer &

time ./StressTest
#!/usr/bin/env bash
HERE=`pwd`

echo $HERE

if [ ! -d "$HERE/bin/" ]; then
  # Control will enter here if $DIRECTORY doesn't exist.
  mkdir $HERE/bin
fi

echo -e "\n  BUILD MAPPER AND REDUCER \n"

cd $HERE/cmd/mapper && go build
cp $HERE/cmd/mapper/mapper $HERE/bin/
cd $HERE/cmd/reducer && go build
cp $HERE/cmd/reducer/reducer $HERE/bin/
cd $HERE/cmd/mapper2 && go build
cp $HERE/cmd/mapper2/mapper2 $HERE/bin/
cd $HERE/cmd/reducer2 && go build
cp $HERE/cmd/reducer2/reducer2 $HERE/bin/
cd $HERE/bin

echo -e "\n  RUN HADOOP JOB \n"

hadoop  jar /opt/hadoop-streaming/hadoop-streaming.jar -input /ads_tracking/2019-01-01/ -output /user/ld-rohitangsu_das/map1 -mapper "mapper" -file "mapper" -reducer "reducer" -file "reducer"
hadoop  jar /opt/hadoop-streaming/hadoop-streaming.jar -input /user/ld-rohitangsu_das/map1/ -output /user/ld-rohitangsu_das/reduce1/ -mapper "mapper2" -file "mapper2" -reducer "reducer2" -file "reducer2"



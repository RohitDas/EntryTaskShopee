HERE=$(pwd)
MINDEGREE=3

echo -e "\n  CHECKING INSTALLED SOFTWARE \n"

GO=${GO:-$(which go)}

if [ -z "$GO" ]; then
    echo "    Please install Go (golang)"
    exit 1
else
  echo "    GO=" $GO
fi


echo -e "\n  BUILD MAPPER AND REDUCER \n"

cd $HERE/cmd/mapper && go build
cd $HERE/cmd/reducer && go build
cd $HERE/cmd/mapper2 && go build
cd $HERE/cmd/reducer2 && go build
cd $HERE
echo -e "\n  RUN HADOOP JOB \n"

OUTPUT="./output"

if [ -d "$OUTPUT" ]; then
  echo "    Please delete $OUTPUT directory"
  exit 1
else
  echo "    hadoop  jar $HADOOP_HOME/share/hadoop/tools/lib/hadoop-streaming-$HADOOP_VERSION.jar -input ./input -output ./output -mapper ./mapper/mapper -reducer ./reducer/reducer"

  hadoop  jar /usr/local/Cellar/hadoop/3.1.1/libexec/share/hadoop/tools/lib/hadoop-streaming-3.1.1.jar -input /rohit/test.data -output /rohit/out1 -mapper ./cmd/mapper/mapper -file ./cmd/mapper/mapper -reducer "./cmd/reducer/reducer" -file ./cmd/reducer/reducer

  echo "    Check output directory to see the results"
fi

hadoop  jar /usr/local/Cellar/hadoop/3.1.1/libexec/share/hadoop/tools/lib/hadoop-streaming-3.1.1.jar -input /rohit/out1/ -output /rohit/out2 -mapper ./cmd/mapper2/mapper2 -file ./cmd/mapper2/mapper2 -reducer ./cmd/reducer2/reducer2 -file ./cmd/reducer2/reducer2




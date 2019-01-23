#EntryTaskShopee

The Task involves mainly 3 components:

1. MapReduce: This component is delegated the responsibility to consume input from hdfs and schedule map/reduce job that computes intermediate data to cater the Http component later. It is completely an offline computation task.
2. HdfsRedisSyncer: This component syncs the output of the MapReduce Component to a in-memory Redis Database layer. All the http requests are served from the Redis Layer.
3. HttpServer: A single-instance of a go-based server that handles request and uses REDIS connection pool to fetch data.
4. LoadBalancer: With enormous requests per sec, a single service may be overwhelmed. LoadBalancers takes as an input a config.yml that contains the list of up and running servers and acts as a proxy to balance the requests across the servers.


A. MapReduce Component

Source Code: MapReduceComputation

Aim: To understand the Hadoop's Mapreduce framework and get hands on experience in submitting jobs to a cluster of commodity machines using the Hadoop Streaming API. Learn Golang.

Understanding the log structure:

The log contains mainly user-based actions like Clicks, Add to Carts etc. The main focus of the task is on operation type 2, which are clicks. An rough json structure is like the following:
{
  userId: 1,
  sessionId: 1,
  operation: 2
  timestamp: 182384432112,
  items: [{
    shopid: 123
    itemId: 321
  }.{
    shopid: 1234
    itemId: 4321
  }]
}

Steps Taken:

Input Directory: /ads_tracking/2019-01-01/

1. Map Task 1: mapper.go reads lines from stdin and outputs (key, value) in the following format: 
  (userId<TAB>shopid,itemid,timestamp)
2. Reduce Task 1: reducer.go reads the output from the mapper task, aggrgrates the <shopid,itemid,timestamp> for each userId, and for each combination <itemA, itemB> pair, calculates whether diff(tsA, tsB) is within 1 hour. If this is the case, then it outputs (<itemA, itemB> \t 1) pair.
  
  Note: In the above task, a user might have click the same item many times, however, the timestamp considered here is the latest click. This is done to simplify the computation. In the next version, this would be changed to consider all the clicks.
  
  Note: itemA is lexicographically smaller than itemB, this is done on purpose considering the symmetry, i,e 
  degree(itemA, itemB) == degree(itemB, itemA)
3. Map Task 2: Identity mapper than just takes the output of the previous phase and outputs the same thing.
4. Reduce Task 2: Reduce task uses the aggregration function SUM, to reduce by key, thereby calculating the degree of 2 items.

Output Directory /user/ld-rohitangsu_das/reduce1

Total runtime of the 2 map reduce tasks together was about 25~30 minutes.

B. HdfsRedisSyncer

SourceCode: HDFSRedisSyncer

Version1: In the first version, A synchronous process is followed. 

1. Get all the output files from the path hdfs:/user/ld-rohitangsu_das/reduce1
2. For each output file:
    For each line of the output:
      Update keys in the redis

Important things to note:
1. The output of a single file is not streamed, it is read from the hdfs completely to memory. The size of the files are about 1 GB each. Although, not the best approach in terms of memory usage. 
2. Each line of the hdfs file is of the form <shopid:itemid DEGREE>. In redis, there is a key for each shopid:itemid. This key is computed by hashing the <shopid:itemid> pair. The Value for each Key is a Map which is of the form 
{
  shopid: 123
  itemid: 122
  degreeInfo: "shopid1:itemid1=degree1|shopid2:itemid2=degree2"
}

Total time taken in updating the Redis: Not tracked yet

C. HttpApp
  a. LoadBalancer: Simply takes a config.yml file and forwards request to a server which is chosen on the basis of simple algorithm based on the number of open connections currently handled by each server. The request is passed to a server, which is least busy. This algorithm might be made more mature and complex to consider many factors such as network congestion.
  
  b. HttpServer: Each httpServer, consumes a request and uses one of 100 open redis connections to fetch the value for a <shopid:itemid> key. However, at this step, there is still some computation that needs to be done.
  
  For each key:
    val := Get the value from Redis
    itemsAndDegree := Split val by "|"
    items, Degree := Split itemsAndDegree by Degree
    Filter by minDegree
    Sort items based on Degree
    return the top 50 degree items in the form of [{"shopid": 123, "item": 246, degree:456}]
  
  
  c. StressTest:
     It sends 500 simultaneous requests in goroutines and calculates the time taken in completing all the requests.
     The time complexity of a particular request is dependent on the number of other items that share same user clicks that are close in time.
     
     For shopid=49969001&itemid=1042094244&min_degree=1, with 2 items in response, the total time taken is the 
     following:
     real	0m0.313s
     user	0m2.216s
     sys	0m1.096s
     
     More advance benchmarking: To be done
     
     
Deployment:




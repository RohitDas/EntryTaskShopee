#!/usr/bin/env bash

nophup ./httpapp 8081 &
nohup ./httpapp 8082 &
nohup ./httpapp 8083 &

nohup ./LoadBalancer &

time ./StressTest
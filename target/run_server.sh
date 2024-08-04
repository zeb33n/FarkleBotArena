#!bin/bash

python ../pysrc/main.py & 
main_id=$! # id of most recent background process
./server
kill $main_id
killall server #should probably rename the server to something this seems a bit sketch 
rm ../pipes/*

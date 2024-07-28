#!bin/bash

python ../pysrc/main.py & 
main_id=$! # id of most recent background process
./server
kill $main_id
rm pipe
trap kill $main_id

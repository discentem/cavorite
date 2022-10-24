#!/bin/bash

# Start the first process
minio server /data --address :9000 --console-address :9001 &
  
# Start the second process
/usr/bin/create_bucket &
  
# Wait for any process to exit
wait -n
  
# Exit with status of process that exited first
exit $?
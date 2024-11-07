#!/bin/bash

source ./venv/bin/activate
echo "Launching development server"
fastapi dev server.py &  #
SERVER_PID=$!  

# Wait for the server to start
sleep 2  

# Check if the server is up by sending a request
until curl -s -o /dev/null -w "%{http_code}" http://localhost:8000/ > /dev/null; do
    echo "Waiting for server to start..."
    sleep 1  # Wait and check again
done

echo "Server is up, running endpoint tests"
python ./tests/test_endpoints.py


kill $SERVER_PID

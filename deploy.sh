#!/bin/bash

# builds, copys, and runs clashbot on remote system

set -e

FILE=clashbot
ENV_FILE=clash_env

echo "Building binary for Raspberry Pi..."
GOOS=linux GOARCH=arm GOARM=5 go build -o $FILE
echo "$FILE built..."

if [ ! -f "$FILE" ]; then
  echo "$FILE does not exist"
  exit 1
fi

echo "Killing existing process..."
ssh pi@raspberrypi "pkill $FILE || true"

echo "Copying binary and ENV to remote server..."
scp -i ~/.ssh/id_rsa $FILE pi@raspberrypi:~/dev
echo "Binary copied."
scp -i ~/.ssh/id_rsa $ENV_FILE pi@raspberrypi:~/dev
echo "ENV copied."

ssh pi@raspberrypi /bin/bash << EOF
  cd dev;
  source $ENV_FILE
  echo "Sourced $ENV_FILE"
  nohup ./test >/dev/null 2>&1 &
  echo "Running $FILE"
  exit
EOF

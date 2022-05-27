#!/bin/sh
if [[ ! $(pgrep -f script.sh) ]]; then
  nats-streaming-server
fi
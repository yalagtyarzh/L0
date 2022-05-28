#!/bin/bash
if [[ ! $(pgrep -f nats-streaming-server) ]]; then
  nats-streaming-server
fi
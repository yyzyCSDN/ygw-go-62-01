#!/bin/sh
set -eu
cd "$(dirname "$0")"
docker build -f benzhi.Dockerfile -t columnstore:latest .

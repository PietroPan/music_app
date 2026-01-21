#!/usr/bin/env bash

set -e

BASE_URL="http://localhost:8080"

METHOD=$1
ENDPOINT=$2

shift 2

case "$METHOD" in
  get)
    /usr/bin/curl -s "$BASE_URL$ENDPOINT" ;;
  post)
    DATA="$1"
    /usr/bin/curl -s -X POST -H "Content-Type: application/json" -d "$DATA" "$BASE_URL$ENDPOINT" ;;
  *)
    echo "Usage: $0 {get|post} /endpoint [data]"
    exit 1
    ;;
esac

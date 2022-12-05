#!/bin/bash

# args: $1 = number between 1 and 9999
function pad {
  printf "%04d" $1
}
ID=$(pad $1)

if [ -z "$1" ]; then
  echo "Usage: $0 <id>"
  exit 1
fi

# check if ID dir already exists
if [ -d "$ID" ]; then
  echo "problem $ID already exists"
  exit 1
fi

echo "initialising problem $ID/ directory"
mkdir "$ID"
echo "package p$ID" >"$ID/code.go"

echo "paste problem stub code:"
# while STUB is empty, read multiple lines into STUB

STUB=""
while [ -z "$STUB" ]; do
  STUB=$(cat)
done

echo "" >>"$ID/code.go"
echo "$STUB" >>"$ID/code.go"

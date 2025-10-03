#!/bin/bash
if [ "$#" -lt 2 ]; then
  echo "Usage: $0 input.mmd output.png"
  exit 1
fi
mmdc -i "$1" -o "$2"

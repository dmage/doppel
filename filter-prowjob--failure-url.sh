#!/bin/sh
jq -rc '.items[].status | select(.state == "failure") | .url | select(. != null)'

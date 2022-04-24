#!/bin/bash
# create_blog.sh

set -e #exit on most errors

podman pod rm -f vendingmachine

podman rmi vending-machine:latest

podman build . --tag vending-machine:latest

./pods.sh
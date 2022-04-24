#!/bin/bash
# create_blog.sh

set -e #exit on most errors

podman pod create --name vendingmachine -p 8080:8080 -p 3306:3306 -p 8083:8083

podman run -dt --pod vendingmachine --name mysql_vm -e MYSQL_DATABASE=vending_machine_db \
    -e MYSQL_ROOT_PASSWORD=r00t --volume /var/lib/mysql -d docker.io/library/mysql:latest

podman unshare chown 1000:1000 -R realm

podman run -dt --pod vendingmachine --name keycloak_vm -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin \
    -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin --user 1000 \
    -v realm:/realm:Z -e KEYCLOAK_IMPORT=/realm/vendingmachine-realm-export.json \
    -d quay.io/keycloak/keycloak:15.0.2

podman cp realm/vendingmachine-realm-export.json keycloak_vm:/realm

podman run -dt --pod vendingmachine --name rest_api_vm --restart "on-failure" -d vending-machine:latest

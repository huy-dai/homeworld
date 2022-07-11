#!/bin/sh

set -e

mkdir -p /data/nginx/nix-cache-info/temp
mkdir -p /data/nginx/nix-cache-info/store
mkdir -p /data/nginx/nar/temp
mkdir -p /data/nginx/nar/store
chown -R nginx.nginx /data/nginx/
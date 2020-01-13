#!/bin/sh

find /app/main*.js -type f -exec sed -i -e "s|{{!API_URL!}}|${API_URL}|" -e "s|{{!WS_URL!}}|${WS_URL}|" {} \;

nginx -g 'daemon off;'

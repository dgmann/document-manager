#!/bin/sh

find /app/main*.js -type f -exec sed -i -e "s|{{!API_URL!}}|${API_URL}|" -e "s|{{!WS_URL!}}|${WS_URL}|" -e "s|{{!BUGSNAG_KEY!}}|${BUGSNAG_KEY}|" {} \;

nginx -g 'daemon off;'

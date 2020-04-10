#!/bin/sh
echo "{\"host\": \"$API_HOST\", \"useSSL\": $ENABLE_SSL }" > /app/assets/config.json

exec nginx -g 'daemon off;'

#!/bin/bash
set -e
echo "127.0.1.1 ${HOSTNAME}" >> /etc/hosts
exec "$@"
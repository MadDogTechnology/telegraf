#!/bin/sh
set -e

if [ "${1:0:1}" = '-' ]; then
    set -- telegraf "$@"
fi

envsubst < /telegraf-template.conf > /etc/telegraf/telegraf.conf

exec "$@"

#!/bin/sh
set -e

if [ "${1:0:1}" = '-' ]; then
    set -- telegraf "$@"
fi

export METRIC_BATCH_SIZE=${METRIC_BATCH_SIZE:-5000}
export METRIC_BUFFER_LIMIT=${METRIC_BUFFER_LIMIT:-1000000}

envsubst < /telegraf-template.conf > /etc/telegraf/telegraf.conf

exec "$@"

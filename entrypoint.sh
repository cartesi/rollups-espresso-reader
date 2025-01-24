#!/bin/sh

cartesi-rollups-node &

if [ "${CARTESI_FEATURE_ESPRESSO_READER_ENABLED:-false}" != "false" ]; then
    cartesi-rollups-espresso-reader
fi

wait

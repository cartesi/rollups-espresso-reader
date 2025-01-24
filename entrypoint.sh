#!/bin/sh

cartesi-rollups-node &

if [ "$CARTESI_FEATURE_INPUT_READER_ENABLED" != "true" ]; then
    cartesi-rollups-espresso-reader
fi

wait

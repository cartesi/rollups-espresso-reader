#!/bin/sh
set -e

/migrate-db.sh

if [ "${CARTESI_FEATURE_ESPRESSO_READER_ENABLED:-false}" != "false" ]; then
    /espresso-dev-node-health.sh
    cartesi-rollups-node & cartesi-rollups-espresso-reader
else
    cartesi-rollups-node
fi

wait

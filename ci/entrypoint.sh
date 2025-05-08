#!/bin/sh
set -e

/migrate-db.sh

/espresso-dev-node-health.sh

cartesi-rollups-node &

# TODO: remove in the future when we stop test the node v2 without espresso in this repo.
if [ "${CARTESI_FEATURE_ESPRESSO_READER_ENABLED:-false}" != "false" ]; then
    cartesi-rollups-espresso-reader
fi

wait

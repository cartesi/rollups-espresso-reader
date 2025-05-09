#!/bin/bash
set -e
echo "Waiting for Espresso network to be available..."
for i in {1..300}; do
    if curl -sSf $ESPRESSO_BASE_URL/v1/availability/header/1 > /dev/null; then
        echo "Espresso network is available."
        exit 0
    else
        echo "Still waiting... attempt $i"
        sleep 1
    fi
done
echo "Espresso network did not become available in time."
docker logs espresso-dev-node || true
exit 1
#!/bin/bash
set -e
echo "Waiting for Espresso dev node to start..."
for i in {1..60}; do
    if curl -sSf http://localhost:24000/v0/availability/header/1 > /dev/null; then
        echo "Espresso dev node is available."
        exit 0
    else
        echo "Still waiting... attempt $i"
        sleep 1
    fi
done
echo "Espresso dev node did not become available in time."
docker logs espresso-dev-node || true
exit 1

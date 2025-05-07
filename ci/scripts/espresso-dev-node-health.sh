#!/bin/bash
set -e
echo "Waiting for Espresso dev node to start..."
for i in {1..60}; do
    # https://github.com/EspressoSystems/espresso-network/blob/985e0fb1bfdfddab509d02fb55508caa635ca550/sequencer/src/api.rs#L1926
    if curl -sSf http://localhost:24000/healthcheck > /dev/null; then
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

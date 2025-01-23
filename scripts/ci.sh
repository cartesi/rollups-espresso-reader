#!/bin/bash
set -e

docker stop $(docker ps -q)

docker buildx prune --all --force && docker system prune --volumes --force

docker run -d --rm --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=rollupsdb postgres:16-alpine

echo "Migrate DB node v2"
cd rollups-node
eval $(make env)
export CGO_CFLAGS="-D_GNU_SOURCE -D__USE_MISC"
go run dev/migrate/main.go
cd -

echo "Migrate DB Espresso"
eval $(make env)
make migrate
make generate-db

echo "Build image"
#docker build -t espresso .
docker build -t espresso -f Dockerfile-espresso .

echo "Run Anvil"
cd rollups-node
make devnet
make run-devnet
cd -

docker run --env-file env.nodev2-local --rm --network=host -v ./rollups-node:/var/lib/cartesi-rollups-node/src --name c_espresso espresso

docker exec c_espresso cartesi-rollups-cli app deploy -n echo-dapp -t applications/echo-dapp/ -v

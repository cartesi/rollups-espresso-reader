#!/bin/bash
set -e

rm -rf ./rollups-node

git clone -b feature/new-build --recurse-submodules https://github.com/cartesi/rollups-node.git

docker stop $(docker ps -q) || true

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
docker build -t espresso -f Dockerfile .

echo "Run Anvil"
cd rollups-node
make devnet
make run-devnet
cd -

docker run --env-file env.nodev2-local --rm --network=host --name c_espresso espresso
# docker run --env-file env.nodev2-sepolia --rm --network=host --name c_espresso espresso

exit 0

docker exec c_espresso cartesi-rollups-cli app deploy -n echo-dapp -t applications/echo-dapp/ -v

export ACCOUNT=0xB5C1674c0527b6C31A5019fD04a6C1529396DA37
export PRIVATE_KEY=ad03...e462
export RPC_URL=https://eth-sepolia.g.alchemy.com/v2/<key>
docker exec c_espresso cartesi-rollups-cli app deploy -v -n echo-dapp -o $ACCOUNT -O $ACCOUNT -k $PRIVATE_KEY -r $RPC_URL -t applications/echo-dapp/

# Milton
docker exec c_espresso cartesi-rollups-cli app register -v -c 0xe82D9ebc0c2773516914a1285F4492Ad5f5Ab9F6 -a 0x2fBe606e211b1BFD0ffE53aa7e15d299824a9478 -n echo-dapp -t applications/echo-dapp/

# Oshiro
docker exec c_espresso cartesi-rollups-cli app register -v -c 0x48F176733DBdEc10EC4d1692e98403E0927E869C -a 0x5a205Fcb6947e200615B75C409ac0aa486D77649 -n echo-dapp -t applications/echo-dapp/

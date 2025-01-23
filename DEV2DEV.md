# Dev 2 Dev

Run the database

```sh
cd rollups-node
docker run --rm --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_DB=rollupsdb -v ./test/postgres/init-test-db.sh:/docker-entrypoint-initdb.d/init-test-db.sh postgres:16-alpine
```

Migrate the database

```sh
cd rollups-node
eval $(make env)
export CGO_CFLAGS="-D_GNU_SOURCE -D__USE_MISC"
go run dev/migrate/main.go
```

The 
`export CGO_CFLAGS="-D_GNU_SOURCE -D__USE_MISC"`
fixes:

```sh
runtime/cgo
# runtime/cgo
In file included from libcgo.h:7,
                 from gcc_context.c:7:
/usr/include/stdio.h:205:27: error: 'L_tmpnam' undeclared here (not in a function)
  205 | extern char *tmpnam (char[L_tmpnam]) __THROW __wur;
      |                           ^~~~~~~~
/usr/include/stdio.h:210:33: error: 'L_tmpnam' undeclared here (not in a function); did you mean 'tmpnam'?
  210 | extern char *tmpnam_r (char __s[L_tmpnam]) __THROW __wur;
      |                                 ^~~~~~~~
      |                                 tmpnam
```

Build the Espresso image

Make sure that you are on the project root folder

```sh
docker build -t espresso .
```

Run the image

```sh
docker run --rm --network=host -v ./rollups-node:/var/lib/cartesi-rollups-node/src --name c_espresso espresso
```

docker run --env-file env.nodev2-local --rm --network=host -v ./rollups-node:/var/lib/cartesi-rollups-node/src --name c_espresso espresso

docker exec -it c_espresso /bin/bash

mkdir -p applications
cartesi-machine --ram-length=128Mi --store=applications/echo-dapp --final-hash -- ioctl-echo-loop --vouchers=1 --notices=1 --reports=1 --verbose=1

cartesi-rollups-cli app deploy -n echo-dapp -t applications/echo-dapp/ -v

Output:

```bash
Transaction submitted: 0x65fde97551978c587378a195bfbf20807d738883a05ec5c62c4ed9d5060a9ea5
Transaction successful!
New Authority contract deployed at address: 0xd121f8aE5Ab0d5F472687AF19E393D18fD3e140c
Transaction submitted: 0xe6aee52ea9c921940f49ad544b80e6237da35a78a174c82368a7820963debaba
Transaction successful!
New Application contract deployed at address: 0x36B9E60ACb181da458aa8870646395CD27cD0E6E
Application 0x36b9e60acb181da458aa8870646395cd27cd0e6e successfully deployed
```


INPUT=0xdeadbeef; \
INPUT_BOX_ADDRESS=0x593E5BCf894D6829Dd26D0810DA7F064406aebB6; \
APPLICATION_ADDRESS=0x36B9E60ACb181da458aa8870646395CD27cD0E6E; \
cast send \
    --mnemonic "test test test test test test test test test test test junk" \
    --rpc-url "http://localhost:8545" \
    $INPUT_BOX_ADDRESS "addInput(address,bytes)(bytes32)" $APPLICATION_ADDRESS $INPUT
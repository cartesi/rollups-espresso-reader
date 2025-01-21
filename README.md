# Cartesi Rollups Espresso Reader

To build:

```bash
go build
```

To run:

```bash
eval $(make env)
make migrate
./espresso-reader
```

To run automated integrated tests, 
```
(cartesi-rollups-node)
make run-postgres && make migrate
./cartesi-rollups-node

(rollups-espresso-reader)
make migrate
go test github.com/cartesi/rollups-espresso-reader/internal/espressoreader -v
```

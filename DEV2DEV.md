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
docker run --rm --network=host espresso
```

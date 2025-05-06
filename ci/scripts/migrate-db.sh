#!/bin/bash
set -e

GO_MIGRATE_VERSION=4.18.2
ARCH=$(dpkg --print-architecture || uname -m)
case "$ARCH" in
    amd64|x86_64)
        MIGRATE_URL="https://github.com/golang-migrate/migrate/releases/download/v${GO_MIGRATE_VERSION}/migrate.linux-amd64.deb"
        ;;
    arm64|aarch64)
        MIGRATE_URL="https://github.com/golang-migrate/migrate/releases/download/v${GO_MIGRATE_VERSION}/migrate.linux-arm64.deb"
        ;;
    armhf|armv7l)
        MIGRATE_URL="https://github.com/golang-migrate/migrate/releases/download/v${GO_MIGRATE_VERSION}/migrate.linux-armv7.deb"
        ;;
    i386)
        MIGRATE_URL="https://github.com/golang-migrate/migrate/releases/download/v${GO_MIGRATE_VERSION}/migrate.linux-386.deb"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac
curl -fsSL "$MIGRATE_URL" -o /tmp/migrate.deb
dpkg -i /tmp/migrate.deb
rm /tmp/migrate.deb
migrate --version

echo "Run rollups node migration"
migrate -source file:///rollups-node/migrations -database "$CARTESI_DATABASE_CONNECTION" up

echo "Run espresso reader migration"
migrate -source file:///internal/repository/postgres/schema/migrations -database "$CARTESI_DATABASE_CONNECTION&x-migrations-table=espresso_schema_migrations" up

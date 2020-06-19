#!/bin/sh

rm -rf /tmp/scorch.bleve || /bin/true # wipeout the temporary index on start
restore-central-db
rollback-rocksdb
/stackrox/bin/migrator
rocksdb-migration

RESTART_EXE="$(readlink -f "$0")" exec /stackrox/central "$@"

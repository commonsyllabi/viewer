# internal

## migrations

the cli tool relies on the env variable set `POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/example?sslmode=disable'`


to create a new migration:

```bash
env POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/example?sslmode=disable' migrate create -ext sql -dir ./migrations/ -seq create_contributors_table
```

when there's a migration conflict, go back down to the last working version, then fix things if necessary (for instance, bulk adding fields to a new column that has a NOT NULL attribute), then doing the migrations up again
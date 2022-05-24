# internal

## migrations

the cli tool relies on the env variable set `POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/example?sslmode=disable'`


to create a new migration:

```bash
env POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/example?sslmode=disable' migrate create -ext sql -dir ./migrations/ -seq create_contributors_table
```
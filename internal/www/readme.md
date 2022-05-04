# CC Viewer Frontend

To start the frontend (_note:_ manual reload needed):

```
yarn watch
```

This would start site on localhost:2046

To start the server, from root of project:

```
go run internal/main.go
```

### To Test

Open cypress:

```
yarn run cypress open
```

Start frontend on `localhost:8080` (_note:_ this server does not work for reloading):

```
yarn start
```

Run the tests:

```
yarn test
```

Or alternatively, `yarn autotest` combines `start` and `test`

### Database

The API connects to a Postgres database called `cosyl`, with username `cosyl` and password `cosyl`. These are being set in the .env file, along with the `DB_HOST`.

`DB_HOST` is the hostname of where Postgres runs. Locally, this is `localhost`. When running with Docker, it's `db` (the `name` set in `docker-compose.yml`).

If you change the username, password or db name, you need to make sure the new users and databases are created.  You can do it with `psql` locally, and by removing the data volume of the docker image, with

```
docker-compose down --volumes
docker-compose build
```
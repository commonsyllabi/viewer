# CC Viewer Frontend

### To start the backend

From root of project:

```
godotenv -f .env go run cmd/api/main.go
```

This runs the Go backend, loading environment variables from the `.env` file with the [godotenv](https://github.com/joho/godotenv) binary. Environment variables are:

- `PORT`, starts the Go process on the given port, default `3046`
- `DATABASE_URL`, full string to connect to database, default to `postgres://cosyl:cosyl@localhost:5432/cosyl`, varies depending on local dev, docker or deployed.
- `DATABASE_TEST_URL`, full string to connect to test database, default to `postgres://cosyl:cosyl@localhost:5432/test`, only the database name changes to `test`.
- `DEBUG`, set to `true` or `false`.
- `GIN_MODE`, set to `debug`, `production`, `test`


### To start the frontend:

There are two cases: the pages rendered client-side (such as `www/src/cartridge.html`with Vue embedded, and pages rendered server-side (such as `api/templates/syllabus.tmpl`).

To work on the client-side pages, you can run:

```
yarn dev
```

This would start site on `localhost:3000`, and has hot-reload enabled.

To work on the server-side pages, you can run 

```
yarn watch
```

This would make the files available to `localhost:{server_port}`, but doesn't have hot-reload enabled, since the HTML templates are rendered by the go app.


### To Test

Tests are written in the `integration` folder, and can use dummy data as request responses (ajax, fetch, axios, etc.), from the `fixtures` folder.

To work on the tests interactively, open cypress:

```
yarn run cypress open
```

Serve the frontend on `localhost:3046` (_note:_ this server does not work for reloading, and is only used for testing). From root of project:

```
godotenv -f .env go run cmd/api/main.go
```

You can also run the tests in headless mode (no UI). While the backend is running, run the frontend tests:

```
yarn test
```

Or alternatively, `yarn autotest` combines `start` and `test`. It is used the `pre-commit` hook—a script in the `.git/hooks/` folder.This is in the process of being migrated to a docker image, for consistency. WIP:

```
docker-compose -f docker-compose.test.yml run api_test --abort-on-container-exit    
```

### Database

The API connects to a Postgres database called `cosyl`, with username `cosyl` and password `cosyl`. These are being set in the .env file, along with the `DB_HOST`.

`DB_HOST` is the hostname of where Postgres runs. Locally, this is `localhost`. When running with Docker, it's `db` (the `name` set in `docker-compose.yml`).

If you change the username, password or db name, you need to make sure the new users and databases are created.  You can do it with `psql` locally, and by removing the data volume of the docker image, with

```
docker-compose down --volumes
docker-compose build # this rebuilds the db image with the env user, then the volume with the env name
```
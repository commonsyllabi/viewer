# CC Viewer Frontend

### To start the server

From root of project:

```
go run internal/main.go
```

This runs the Go backend on a port defined by `config.go`, either read from `config.yml` or set to `2046` in the `defaults` method.


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

Open cypress:

```
yarn run cypress open
```

Serve the frontend on `localhost:8080` (_note:_ this server does not work for reloading, and is only used for testing):

```
yarn start
```

Run the tests:

```
yarn test
```

Or alternatively, `yarn autotest` combines `start` and `test`. It is used the `pre-commit` hook—a script in the `.git/hooks/` folder.

### Database

The API connects to a Postgres database called `cosyl`, with username `cosyl` and password `cosyl`. These are being set in the .env file, along with the `DB_HOST`.

`DB_HOST` is the hostname of where Postgres runs. Locally, this is `localhost`. When running with Docker, it's `db` (the `name` set in `docker-compose.yml`).

If you change the username, password or db name, you need to make sure the new users and databases are created.  You can do it with `psql` locally, and by removing the data volume of the docker image, with

```
docker-compose down --volumes
docker-compose build # this rebuilds the db image with the env user, then the volume with the env name
```
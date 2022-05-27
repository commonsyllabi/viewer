# CC Viewer Frontend

### To start the backend

You can either run it as a Docker container, or as a local process. The advantage of docker containers is that we mirror the deploy environment and that we completely control database operations. Locally, it's faster, but idiosyncratic. Run both commands from the root of project:

```bash
# running locally
godotenv -f ".env,.secrets" go run cmd/api/main.go

# running as docker container
docker compose up
```

You can pass `--build` to the `docker compose` command to recreate the image, if there are any modifications to the source code that have made.

This runs the Go backend, loading environment variables from the `.env` file with the [godotenv](https://github.com/joho/godotenv) binary. Environment variables are:

- `PORT`, starts the Go process on the given port, default `3046`
- `DB_PORT`, starts the database on the given port, and exposes it to other docker containers, default `5433`.
- `DB_USER`, default `cosyl`
- `DB_PASSWORD`, default `cosyl`
- `DB_HOST`, the host the database, either `db` when run in postgres or `localhost` when run locally.
- `DB_NAME`, the name of the database to connect to, default `cosyl` (for testing, only the database name changes to `test`)

Ultimately, these are used to compose the database connection URL, which resembles: `postgres://cosyl:cosyl@localhost:5432/cosyl`. Additionally, you can set:

- `DEBUG`, set to `true` or `false`.
- `GIN_MODE`, set to `debug`, `production`, `test`

#### Secrets

Necessary secrets to run the app, loaded from `.env.secrets`:

- `MAILGUN_PRIVATE_API_KEY` for sending emails.

### To start the frontend:

There are two cases: the pages rendered client-side (such as `www/src/cartridge.html`with Vue embedded, and pages rendered server-side (such as `api/templates/syllabus.tmpl`).

To work on the client-side pages, you can run:

```bash
yarn dev
```

This would start site on `localhost:3000`, and has hot-reload enabled.

To work on the server-side pages, you can run 

```bash
yarn watch
```

This would make the files available to `localhost:{server_port}`, but doesn't have hot-reload enabled, since the HTML templates are rendered by the go app.

### To Test

Tests are written in the `integration` folder, and can use dummy data as request responses (ajax, fetch, axios, etc.), from the `fixtures` folder.

#### interactive

To work on the tests interactively, open cypress:

```
yarn run cypress open
```

Serve the frontend on `localhost:3046` (_note:_ this server does not work for reloading, and is only used for testing). From root of project:

```bash
godotenv -f ".env,.secrets" go run cmd/api/main.go
```

You can also run the tests in headless mode (no UI). While the backend is running, run the frontend tests:

```
yarn test
```

#### automated

An automated, end-to-end testing against a test database can be found in `docker-compose.yml` in the `tests` folder. It is used the `pre-commit` hook—a script in the `.git/hooks/` folder.

```bash
docker compose -f docker compose.test.yml up --build backend_test
docker compose -f docker compose.test.yml up --build frontend_test

```

The full `pre-commit` script is as follows:

```bash
current_branch=$(git rev-parse --abbrev-ref HEAD)

if [ "$current_branch" = "main" ]
then
	docker compose -f docker-compose.test.yml --remove-orphans run backend_test 
	docker compose -f docker-compose.test.yml run docker compose -f docker compose.test.yml run frontend_test --remove-orphans frontend_test 
else
	echo "skipping tests... (not on main)"
fi

```

### Database

#### Connection

The API connects to a Postgres database called `cosyl`, with username `cosyl` and password `cosyl`. These are being set in the .env file, along with the `DB_HOST`.

`DB_HOST` is the hostname of where Postgres runs. Locally, this is `localhost`. When running with Docker, it's `db` (the `name` set in `docker-compose.yml`).

If you change the username, password or db name, you need to make sure the new users and databases are created.  You can do it with `psql` locally, and by removing the data volume of the docker image, with

```bash
docker compose down --volumes
docker compose build
```

#### Migrations

Migrating databases is the process of reaching a desired database state in a step-wise fashion. To create a new migration, run:

```bash
export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/test?sslmode=disable'
migrate create -ext sql -dir ${OUTPUT_DIR} -seq create_users_table
```

I should also look into uniqueness and non-nullability and other column attributes.

The table creation should also take into account inter-table dependencies (has-a, belongs-to), and look into foreign keys
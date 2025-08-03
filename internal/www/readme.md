# CC Viewer Frontend

## To start the backend

You will need some way to provide environment variables (here we use godotenv)

```bash
godotenv go run cmd/api/main.go
```

## To start the frontend

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

```bash
yarn run cypress open
```

Serve the frontend on `localhost:3046` (_note:_ this server does not work for reloading, and is only used for testing). From root of project:

```bash
godotenv -f ".env,.secrets" go run cmd/api/main.go
```

You can also run the tests in headless mode (no UI). While the backend is running, run the frontend tests:

```bash
yarn test
```

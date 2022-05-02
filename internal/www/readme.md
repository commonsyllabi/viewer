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

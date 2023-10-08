# RubiksSolver

A WIP rubiks cube solver. Uses precomputed cube solutions and IDFS to run a meet in the middle style solve.
All generated solutions are optimal (but it might take a while)

Full write up on my [Website](https://www.swannyscode.com/projects/6)

## Precomputing database
```
go run rubiks.go generate -db "path/to/database/file.db"
```

## Building the frontend
```
cd frontEnd
npm install
npm run build
```
## Starting web app / server
(requires a built frontend)
```
go run rubiks.go server -db "path/to/database/file.db" -port 3000
```

## Running tests
```
go test -v -p 1 ./... -db "path/to/database/file.db"
```
(using the flag `-p 1` ensures the output of the tests are live, this is useful as some of the tests can take a long time to run)

# Development

## Build from source

### Get playground (frontend) submodule if required

```bash
git submodule update --init
```

### Generate swagger

```bash
# install swagger cli (only once)
go install github.com/swaggo/swag/cmd/swag@latest

swag init -g "server/server.go" -o "server/docs"
```

### Run Tests
```
go test -cover ./...
```

### Build binary

```bash
go build -o pdf-turtle

# run binary
./pdf-turtle
```

### Build Docker image

```bash
docker build -t pdf-turtle .

# run docker image
docker run -d -p 8000:8000 --name pdf-turtle pdf-turtle
```

## Run Tests including race and coverage

```
go test -race -cover ./...
```

## Update all dependencies

```
go get -u
go mod tidy

# verify functionality
go test ./...
```

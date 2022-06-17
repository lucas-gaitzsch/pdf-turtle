# PdfTurtle 🐢

[![build and test](https://github.com/lucas-gaitzsch/pdf-turtle/actions/workflows/pipeline.yml/badge.svg)](https://github.com/lucas-gaitzsch/pdf-turtle/actions/workflows/pipeline.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucas-gaitzsch/pdf-turtle)](https://goreportcard.com/report/github.com/lucas-gaitzsch/pdf-turtle)

### A painless html to pdf rendering service

[PdfTurtle](https://github.com/lucas-gaitzsch/pdf-turtle) generates PDF reports and documents from HTML templates or raw HTML.

## How to run

### With docker

```bash
docker pull lucasgaitzsch/pdf-turtle:latest

docker run -d \
    -p 8000:8000 \
    --name pdf-turtle \
    lucasgaitzsch/pdf-turtle:latest
```

### With prebuilt binaries

_...COMING SOON_

<!-- TODO:!! -->

### Config

| command line argument | environment variable | type    | default | description                                             |
| --------------------- | -------------------- | ------- | ------- | ------------------------------------------------------- |
| --help                | -                    | -       | -       | Show help                                               |
| --logDebug            | LOG_LEVEL_DEBUG      | boolean | false   | Debug log level active                                  |
| --logJsonOutput       | LOG_JSON_OUTPUT      | boolean | false   | Json log output                                         |
| --renderTimeout       | RENDER_TIMEOUT       | integer | 30      | Render timeout in seconds                               |
| --workerInstances     | WORKER_INSTANCES     | integer | 30      | Count of worker instances                               |
| --port                | PORT                 | integer | 8000    | Server port                                             |
| --maxBodySize         | MAX_BODY_SIZE        | integer | 32      | Max body size in megabyte                               |
| --servePlayground     | SERVE_PLAYGROUND     | boolean | false   | Serve playground from path "./static-files/playground/" |
| --secret              | SECRET               | string  | ""      | Secret used as bearer token                             |

## How to use

### Swagger

Use Swagger-UI under _/swagger/index.html_ as API documentation.

You can use the swagger description (_/swagger/doc.json_ or [./server/docs/swagger.json](./server/docs/swagger.json)) to generate a API client for the language of your choice.

### Postman

_...COMING SOON_

<!-- TODO:!! -->

### PdfTurtle Playground

_...COMING SOON_

<!-- TODO:!! -->

## Included template engines

| Template style                               | Package       | PdfTurtle key  | URL                                 |
| -------------------------------------------- | ------------- | -------------- | ----------------------------------- |
| Golang                                       | html/template | **golang**     | https://pkg.go.dev/html/template    |
| Handlebars-syntax like                       | raymond       | **handlebars** | https://github.com/aymerick/raymond |
| Django-syntax like (require _model._ prefix) | pongo2        | **django**     | https://github.com/flosch/pongo2    |

## Development / Build from source

### Generate swagger

```bash
# install swagger cli (only once)
go install github.com/swaggo/swag/cmd/swag@latest

swag init -g "server/server.go" -o "server/docs"
```

### Build binary

```bash
go build -o pdf-turtle

# run binary
./pdf-turtle
```

### Build Docker image
The docker image is hosted on [Docker Hub](https://hub.docker.com/r/lucasgaitzsch/pdf-turtle).

With the tag *:playground* you get a bundled image with the web playground.

```bash
docker build -t pdf-turtle .

# run docker image
docker run -d -p 8000:8000 --name pdf-turtle pdf-turtle
```

### Test

<!-- `go test -race ./...` -->

```
go test -cover ./...
```

<!-- `go test -coverprofile coverage ./...` -->

## Build with

- [go](https://github.com/golang/go)
- [chromedp (golang chromium driver)](https://github.com/chromedp/chromedp)
- [chromium (render engine)](https://github.com/chromium/chromium)
- [raymond (handlebars template engine)](https://github.com/aymerick/raymond)
- [pongo2 (django template engine)](https://github.com/flosch/pongo2)
- [zerolog](https://github.com/rs/zerolog)
- [go-arg](https://github.com/alexflint/go-arg)
<!-- TODO:!! -->

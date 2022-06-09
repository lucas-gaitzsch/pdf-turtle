# PdfTurtle 🐢
### A painless html to pdf rendering service
Generate PDF reports from HTML templates or raw HTML.


## Included template engines
|Template style|Package|PdfTurtle key|URL|
|-|-|-|-|
|Golang|html/template|**golang**|https://pkg.go.dev/html/template|
|Handlebars-syntax like|raymond|**handlebars**|https://github.com/aymerick/raymond|
|Django-syntax like| pongo2|**django**|https://github.com/flosch/pongo2|

## Development / Build from source
### Test
`go test ./...`
`go test -race ./...`
`go test -cover ./...`
`go test -coverprofile coverage ./...`

### Generate swagger
`go install github.com/swaggo/swag/cmd/swag@latest`
`swag init -g "server/server.go" -o "server/docs"`


# scapo

Simple webapi scaffold project.

- go
- openapi3
- chi
- sqlite3
- air
- delve
- docker

## How to Run

1. build golang debug enviroment.

```shell
$docker-compose -f docker-compose.base.yml build
```

2. build and run project.

```shell
$docker-compose up
```

3. check

http://localhost:18080/pets

```shell
$curl http://localhost:18080/pets
$curl -X POST -H "Content-Type: application/json" -d '{"name":"foo", "tag":"bar"}' localhost:18080/pets
$curl -X DELETE localhost:18080/pets/21
$curl localhost:18080/pets/1
```

## Generate Source Code

```shell
$go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0 -generate types -package openapi petstore-expanded.yaml > petstore/openapi/oapi_types.gen.go

$go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0 -generate chi-server -package openapi petstore-expanded.yaml > petstore/openapi/oapi_server.gen.go

$go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.0 -generate spec -package openapi petstore-expanded.yaml > petstore/openapi/oapi_spec.gen.go
```

The generated source uses `oapi-codegen` v2 and requires Go 1.26.6 or later.
The request validation middleware and runtime are provided by the maintained
`oapi-codegen` modules in `go.mod`. Do not use the deprecated
`github.com/deepmap/oapi-codegen` v1 module when regenerating source code.

## Dependency and Security Maintenance

Run the following checks after changing dependencies or generated source code.

```shell
$go mod tidy
$go mod verify
$go test ./...
$go vet ./...
$go run golang.org/x/vuln/cmd/govulncheck@v1.7.0 ./...
```

Keep the Go toolchain and dependencies on versions that include the latest
available security fixes. Review each vulnerability by its CVE or GHSA,
severity, fixed version, and dependency path before upgrading. Prefer the
smallest compatible update, then regenerate the OpenAPI source and rerun all
checks.

## Debug

edit `.air.toml`.

```toml
full_bin = "APP_ENV=dev APP_USER=air ./tmp/main"
#full_bin = "APP_ENV=dev APP_USER=air /go/bin/dlv exec ./tmp/main --headless=true --listen=:18081 --api-version=2 --accept-multiclient"
```

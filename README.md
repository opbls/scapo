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
$oapi-codegen -generate types -package openapi petstore-expanded.yaml > petstore/openapi/oapi_types.gen.go

$oapi-codegen -generate chi-server -package openapi petstore-expanded.yaml > petstore/openapi/oapi_server.gen.go

$oapi-codegen -generate spec -package openapi petstore-expanded.yaml > petstore/openapi/oapi_spec.gen.go
```

## Debug

edit `.air.toml`.

```toml
full_bin = "APP_ENV=dev APP_USER=air ./tmp/main"
#full_bin = "APP_ENV=dev APP_USER=air /go/bin/dlv exec ./tmp/main --headless=true --listen=:18081 --api-version=2 --accept-multiclient"
```

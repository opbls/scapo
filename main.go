package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/opbls/scapo/petstore/delivery"
	"github.com/opbls/scapo/petstore/openapi"
	"github.com/opbls/scapo/petstore/repository"
	"github.com/opbls/scapo/petstore/usecase"
)

func init() {
	dbConfig = databaseConfig{}

	buf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(buf, &dbConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func main() {

	// address and port
	port := flag.Int("port", 18080, "Port for test HTTP server")
	flag.Parse()
	addr := fmt.Sprintf("0.0.0.0:%d", *port)

	// router swagger
	router := chi.NewRouter()
	swagger, err := openapi.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	router.Use(middleware.OapiRequestValidator(swagger))

	// database
	db, err := sqlx.Connect(dbConfig.getDbDriver(), dbConfig.getDbDataSource())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database\n: %s", err)
	}
	defer db.Close()

	// handlres
	repo := repository.NewPetStoreRepository(db)
	usecase := usecase.NewPetStoreUsecase(repo)
	handler := delivery.NewPetStoreDelivery(usecase)

	// server
	server := http.Server{}
	server.Handler = openapi.HandlerFromMux(handler, router)

	log.Fatal(http.ListenAndServe(addr, router))
}

type databaseConfig struct {
	DbDriver     string `yaml:"DbDriver"`
	DbDataSource string `yaml:"DbDataSource"`
}

var dbConfig databaseConfig

func (dbConfig databaseConfig) getDbDriver() string {
	return dbConfig.DbDriver
}

func (dbConfig databaseConfig) getDbDataSource() string {
	return dbConfig.DbDataSource
}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/baiv84/chiapi/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run() error {
	pgConn, err := pgx.Connect(context.Background(), "postgres://postgres:daemon@172.19.0.2:5432/testdb")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pgConn.Close(context.Background())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	routes.RegisterUserRoutes(pgConn, r)

	//log.Println("Сервер запущен на :3000")
	ok := http.ListenAndServe(":3000", r)
	return ok
}

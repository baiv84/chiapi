package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

type CrudItem struct {
	Id    int
	Name  string
	Descr string
	Year  int
}

type PGConnector struct {
	conn *pgx.Conn
}

func RegisterUserRoutes(pgConnect *pgx.Conn, r chi.Router) {
	pg := &PGConnector{conn: pgConnect}
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	r.Route("/films", func(r chi.Router) {
		r.Get("/", pg.ReadAllFilms)
		r.Post("/", pg.CreateFilm)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", pg.ReadOneFilm)
			r.Delete("/", pg.DeleteFilm)
			r.Patch("/", pg.UpdateFilm)
		})
	})
}

/*
* Create film
 */
func (pg *PGConnector) CreateFilm(w http.ResponseWriter, r *http.Request) {
	var item CrudItem
	err_json := json.NewDecoder(r.Body).Decode(&item)
	if err_json != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err_json.Error()))
		return

	}
	pg.conn.Exec(context.Background(), "INSERT INTO public.films (id, name, descr, year) VALUES(nextval('generate_film_id'::regclass), $1, $2, $3)", item.Name, item.Descr, item.Year)
}

/*
* Read all films
 */
func (pg *PGConnector) ReadAllFilms(w http.ResponseWriter, r *http.Request) {
	rows, err := pg.conn.Query(context.Background(), "SELECT id, name, descr, year FROM films")
	if err != nil {
		fmt.Println("error")
	}
	defer rows.Close()

	items := []CrudItem{}
	for rows.Next() {
		var item CrudItem
		err = rows.Scan(&item.Id, &item.Name, &item.Descr, &item.Year)
		if err != nil {
			fmt.Println("!!!")
		}
		fmt.Printf("ID: %d, Name: %s, Description: %s, Year - %d\n", item.Id, item.Name, item.Descr, item.Year)
		items = append(items, item)
	}

	resultJson, err := json.Marshal(items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(resultJson)

}

/*
* Read one film
 */
func (pg *PGConnector) ReadOneFilm(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var item CrudItem
	err = pg.conn.QueryRow(context.Background(), "select id, name, descr, year from films where id=$1", id).Scan(&item.Id, &item.Name, &item.Descr, &item.Year)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	jsonResult, _ := json.Marshal(&item)
	w.Write(jsonResult)

}

/*
* Update film
 */
func (pg *PGConnector) UpdateFilm(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var item CrudItem
	err_json := json.NewDecoder(r.Body).Decode(&item)
	if err_json != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err_json.Error()))
		return
	}
	pg.conn.Exec(context.Background(), "UPDATE public.films SET name=$1, descr=$2, year=$3 WHERE id=$4", item.Name, item.Descr, item.Year, id)

}

/*
* Delete film
 */
func (pg *PGConnector) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}
	pg.conn.Exec(context.Background(), "DELETE FROM public.films WHERE id=$1", id)

}

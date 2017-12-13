package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type movies struct {
	ID       int
	Movie    string
	Director string
	Genre    string
	Year     int
}

type movies2 struct {
	ID       int    `json:"ID, omitempty"`
	Movie    string `json:"Movie, omitempty"`
	Director string `json:"Director, omitempty"`
	Genre    string `json:"Genre, omitempty"`
	Year     string `json:"Year, omitempty"`
}

func main() {
	port := 13131

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "home.html")
	})

	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "get.html")
	})

	http.HandleFunc("/getall/", func(w http.ResponseWriter, r *http.Request) {
		GetAll(w, r)
	})

	http.HandleFunc("/getmovies/", func(w http.ResponseWriter, r *http.Request) {
		GetMovies(w, r)
	})

	http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "post.html")
	})

	http.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case "POST":
			Insert(w, r)
			break
		}
	})
	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/movies")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	film := movies{}

	rows, err := db.Query("select ID, Movie, Director, Genre, Year from MOVIES")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&film.ID, &film.Movie, &film.Director, &film.Genre, &film.Year)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(&film)
	}

	err = rows.Err()
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/movies")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	film := movies{}

	rows, err := db.Query("select Movie from MOVIES")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&film.Movie)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(&film)
	}

	err = rows.Err()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	var film movies2
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&film)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/movies")
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("INSERT INTO MOVIES (ID, Movie, Director, Genre, Year) VALUES (?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(film.ID, film.Movie, film.Director, film.Genre, film.Year)
}

func Delete(w http.ResponseWriter, r *http.Request, id string) {
	idmovie, _ := strconv.Atoi(id)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/movies")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("DELETE FROM MOVIES WHERE ID=?", idmovie)

	defer rows.Close()
}

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"html/template"
	"net/http"
	"os"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/static").Handler(
		http.StripPrefix(
			"/static",
			http.FileServer(
				http.Dir("static/"),
			),
		),
	)

	r.HandleFunc("/post/{id}", ViewHandler)

	r.HandleFunc("/", HomeHandler)

	fmt.Println(http.ListenAndServe(":8080", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	db := initDb()

	items, err := ListPosts(db)
	checkError(err)

	t := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/list.html",
		),
	)

	if err := t.ExecuteTemplate(w, "layout.html", items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic recovered in ViewHandler", r)
		}
	}()

	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, "You must provide a post id", http.StatusBadRequest)
	}

	db := initDb()

	t := template.Must(
		template.ParseFiles(
			"templates/layout.html",
			"templates/view.html",
		),
	)

	if err := t.ExecuteTemplate(w, "layout.html", GetPostById(id, db)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetPostById(id string, db *sql.DB) *Post {
	row := db.QueryRow("SELECT * FROM posts WHERE id=?", id)

	var post Post

	err := row.Scan(&post.Id, &post.Title, &post.Body)
	checkError(err)

	defer func(db *sql.DB) {
		closeErr := db.Close()
		checkError(closeErr)
	}(db)

	return &post
}

func initDb() *sql.DB {
	err := godotenv.Load(".env")

	checkError(err)

	var dbUser = os.Getenv("DB_USER")
	var dbPassword = os.Getenv("DB_PASSWORD")

	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@/db?charset=utf8",
		dbUser,
		dbPassword,
	))

	checkError(err)

	return db
}

func checkError(err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic recovered in checkError", r)
		}
	}()

	if err != nil {
		panic(err)
	}
}

func ListPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT * FROM posts")

	if err != nil {
		return nil, err
	}

	var items []Post

	for rows.Next() {
		post := Post{}
		scanError := rows.Scan(&post.Id, &post.Title, &post.Body)
		checkError(scanError)
		items = append(items, post)
	}

	defer func(db *sql.DB) {
		closeErr := db.Close()
		checkError(closeErr)
	}(db)

	return items, nil
}

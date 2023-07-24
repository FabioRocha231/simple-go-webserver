package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"html/template"
	"io"
	"net/http"
	"os"
)

type Post struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
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

	r.HandleFunc("/post", CreatePostHandler).Methods("POST")

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

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic recovered in ViewHandler", r)
		}
	}()

	reqBody, _ := io.ReadAll(r.Body)
	var post Post
	json.Unmarshal(reqBody, &post)

	if post.Title == "" || post.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "You must provide a title and body",
		})
		return
	}

	db := initDb()

	err := CreatePost(db, post)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": err.Error(),
		})
		return
	}

	postId := LastInsertId(db)

	post.Id = postId

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Post created",
		"post":    post,
	})
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

func LastInsertId(db *sql.DB) int {
	var id int

	err := db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&id)

	checkError(err)

	return id
}

func CreatePost(db *sql.DB, post Post) error {
	_, err := db.Exec("INSERT INTO posts (title, body) VALUES (?, ?)", post.Title, post.Body)

	return err
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

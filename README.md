# Golang Project with Gorilla Mux and MySQL Driver

This project is a basic Golang web application that uses the Gorilla Mux router and the MySQL driver to handle HTTP requests and interact with a MySQL database. The project allows you to view a list of posts and display individual post details.

## Prerequisites

Before you start, make sure you have the following installed:

1. Golang: https://golang.org/doc/install
2. Docker and Docker Compose: https://docs.docker.com/get-docker/

## Getting Started

1. Clone the repository:

```
git clone <repository-url>
cd <project-folder>
```

2. Set up the MySQL Database:

The project uses Docker to run a MySQL database. Open the `docker-compose.yml` file and customize the following environment variables:

- `MYSQL_DATABASE`: The name of the MySQL database to be created.
- `MYSQL_USER`: The username for the MySQL user.
- `MYSQL_PASSWORD`: The password for the MySQL user.
- `MYSQL_ROOT_PASSWORD`: The root password for the MySQL server.

Save the changes and start the MySQL database using Docker Compose:

```
docker-compose up -d
```

3. Environment Variables:

Create a `.env` file in the project root directory and set the following environment variables:

```
DB_USER=user
DB_PASSWORD=password
```

4. Initialize the database:

Run the `initDb` function to create the required `posts` table:

```
go run main.go
```

5. Run the Application:

Start the Golang web application:

```
go run main.go
```

The application will be accessible at `http://localhost:8080/`.

## Project Structure

```
- main.go
- templates/
  - layout.html
  - list.html
  - view.html
- static/
  - main.css
```

- `main.go`: The main Golang file that sets up the Gorilla Mux router and handles HTTP requests.
- `templates/`: Contains HTML templates for the application.
- `static/`: Contains static files like CSS.

## Endpoints

- `/`: Displays a list of posts.
- `/post/{id}`: Displays the details of a specific post.

## HTML Templates

The application uses Go templates to render HTML pages. There are three templates:

- `layout.html`: The base layout for all pages.
- `list.html`: Displays the list of posts.
- `view.html`: Displays the details of a specific post.

## Data Model

The `Post` struct represents the data model for a post with the following fields:

- `Id`: The unique identifier for the post.
- `Title`: The title of the post.
- `Body`: The body/content of the post.

## Database Operations

- `initDb()`: Initializes the database connection and sets up the required table (Run once to create the `posts` table).
- `ListPosts(db *sql.DB) ([]Post, error)`: Retrieves a list of all posts from the database.
- `GetPostById(id string, db *sql.DB) *Post`: Retrieves a specific post by its ID from the database.

## Dependencies

- Gorilla Mux: A powerful URL router and dispatcher for Golang.
- MySQL Driver: The official MySQL driver for Golang.

## License

This project is licensed under the [MIT License](LICENSE). Feel free to modify and use it for your needs.

Happy coding!
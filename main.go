package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Tasks struct {
	Itasks	[]Task
	Mtasks	[]Task
	Ltasks	[]Task
	Ctasks	[]Task
}

type Task struct {
	Id		int
	Content		string
	IsCompleted	bool
	CreatedAt	time.Time
	Importance	string
}

func deleteItemHandler(c *gin.Context) {
	_ = c.Param("importance")
	id := c.Param("id")

	fmt.Println("Here:", id)

	_, err := db.Exec("DELETE FROM todos WHERE todo_id = ?", id)
	if err != nil {
		log.Println("[ERROR] Error deleting from database -> ", err)
	}
}

func addTodoHandler(c *gin.Context) {
	// Add item to database
	// add item to list and refresh the whole page.
	// OR
	// Add item to database and respond with html as previously done.
	content := c.PostForm("content")
	importance := c.Param("importance")

	if importance == "important" {
		importance = "high"
	} else if importance == "minor" {
		importance = "medium"
	} else if importance == "later" {
		importance = "low"
	}

	task := Task{
		Content: content,
		IsCompleted: false,
		CreatedAt: time.Now(),
		Importance: importance,
	}

	_, err := db.Exec("INSERT INTO todos (title, description, importance, created_at) VALUES(?, ?, ?, ?)", "todo", task.Content, task.Importance, task.CreatedAt)
	if err != nil {
		log.Println("[ERROR] Error inserting into database -> ", err)
	}
}

func indexHandler(c *gin.Context) {
	var itasks []Task
	var mtasks []Task
	var ltasks []Task
	var ctasks []Task

	// All of these need to get tasks based on the username as well now
	itasks = getTasksByImportance("high")
	mtasks = getTasksByImportance("medium")
	ltasks = getTasksByImportance("low")

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"ImportantTasks": itasks,
		"MinorTasks": mtasks,
		"LaterTasks": ltasks,
		"CompletedTasks": ctasks,
	})
}

func loginHandler(c *gin.Context) {
	c.File("./static/login.html")
}

func loginFormHandler(c *gin.Context) {
	username := c.PostForm("username")
	if ok := getUsername(username); !ok {
		// Create a new user and insert in database

		// Create a new session for new user and insert in map

		// Redirect to index.
	} else {
		password := c.PostForm("password")
		if ok := checkPassword(username, password); !ok {
			// Return an error message for user in login page
		} else {
			// Create a new session and store in session map

			// Redirect to index (btw: Index will now need to grab tasks/todos based on the user as well)
		}
	}
}

func getUsername(username string) bool {
	var user string
	row := db.QueryRow("SELECT * FROM users WHERE username= ?", username)
	if err := row.Scan(&user); err != nil {
		return false
	} else {
		return true
	}
}

func checkPassword(username, password string) bool {
	var user string
	var pass string
	row := db.QueryRow("SELECT * FROM users WHERE username= ?", username)
	if err := row.Scan(&user, &pass); err != nil {
		return false
	} else {
		if pass == password {
			return true
		}
		return false
	}
}

func getTasksByImportance(importance string) []Task {	
	var tasks []Task

	rows, err :=  db.Query("SELECT * FROM todos WHERE importance= ?", importance)
	if err != nil {
		log.Println("[ERROR] Error querying row for important tasks -> ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var title string
		if err := rows.Scan(&task.Id, &title, &task.Content, &task.Importance, &task.CreatedAt); err != nil {
			log.Println("[ERROR] Could not scan row -> ", err)
		}

		tasks = append(tasks, task)
	}
	return tasks
}

// Middleware
func AuthFunc(c *gin.Context) {
	// Get the session passed so the unique identifier passed by sessionId
	// If it does not exists, send a forbidden http status and redirect to the login page.

	c.Next()
}

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "username",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "todos",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()
	router.LoadHTMLGlob("static/*")

	// Init the sessions map
	// Create the AuthGroup which will include everything except /login and /loginform

	router.POST("/add-item/:importance", addTodoHandler)
	router.DELETE("/delete/:importance/:id", deleteItemHandler)
	router.GET("/index", indexHandler)
	router.GET("/login", loginHandler)
	router.Run()
}

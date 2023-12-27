package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cs50-romain/MetisDeux/util/session"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var store *session.Store

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
	Username	string
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
	
	sessionid, _ := c.Cookie("session_id")
	session := store.Get(sessionid)

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
		Username: session.Username,
	}

	fmt.Println(task.Username)

	_, err := db.Exec("INSERT INTO todos (title, description, importance, created_at, username) VALUES(?, ?, ?, ?, ?)", "todo", task.Content, task.Importance, task.CreatedAt, task.Username)
	if err != nil {
		log.Println("[ERROR] Error inserting into database -> ", err)
	}
}

func indexHandler(c *gin.Context) {
	var tasks []Task
	var itasks []Task
	var mtasks []Task
	var ltasks []Task
	var ctasks []Task

	sessionid, _ := c.Cookie("session_id")
	session := store.Get(sessionid)
	log.Printf("[INFO] sessionid: %s, session username: %s\n", sessionid, session.Username)

	// All of these need to get tasks based on the username as well now
	tasks = getTasksByUsers(session.Username)
	itasks = filterTasksByImportance("High", tasks)
	mtasks = filterTasksByImportance("Medium", tasks)
	ltasks = filterTasksByImportance("Low", tasks)

	fmt.Println(tasks)
	fmt.Println(itasks)
	fmt.Println(mtasks)
	fmt.Println(ltasks)

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
	password := c.PostForm("password")
	// No usernmame found in database
	if ok := getUsername(username); !ok {
		// Create a new user and insert in database
		_, err := db.Exec("INSERT INTO users (username, password) VALUES(?, ?)", username, password)
		if err != nil {
			log.Println("[ERROR] Error inserting into database -> ", err)
		}
		// Create a new session for new user and insert in map
		// Set Cookie
		sessionid := store.Save(username)
		c.SetCookie("session_id", sessionid, 172800, "/", "localhost", false, true)
		// Redirect to index.
		c.Redirect(http.StatusFound, "/index")
	} else {
		if ok := checkPassword(username, password); !ok {
			// Return an error message for user in login page
			c.HTML(http.StatusForbidden, "login.html", nil)
			c.Abort()
			return
		} else {
			// Create a new session and store in session map
			sessionid := store.Save(username)
			c.SetCookie("session_id", sessionid, 172800, "/", "localhost", false, true)
			// Redirect to index.
			c.Redirect(http.StatusFound, "/index")
		}
	}
}

func filterTasksByImportance(importance string, original_tasks []Task) []Task {
	var tasks []Task
	for _, task := range original_tasks {
		if task.Importance == importance {
			tasks = append(tasks, task)
		}
	}
	return tasks
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

func getTasksByUsers(username string) []Task {	
	var tasks []Task

	rows, err :=  db.Query("SELECT * FROM todos WHERE username= ?", username)
	if err != nil {
		log.Println("[ERROR] Error querying row for tasks -> ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var title string
		if err := rows.Scan(&task.Id, &title, &task.Content, &task.Importance, &task.CreatedAt, &task.Username); err != nil {
			log.Println("[ERROR] Could not scan row -> ", err)
		}
		fmt.Println(task.Importance)

		tasks = append(tasks, task)
	}
	return tasks
}

// Middleware
func AuthFunc(c *gin.Context) {
	fmt.Println("Middleware starts")
	// Get the session_id passed so the unique identifier passed by session_id
	session_id, err := c.Cookie("session_id")
	if err != nil {
		c.HTML(http.StatusForbidden, "login.html", nil)
		c.Abort()
		return
	}
	// If it does not exists, send a forbidden http status and redirect to the login page.
	session := store.Get(session_id)
	fmt.Println(session.Username)
	if session.Username != " " {
		fmt.Println("Forbidden")
		c.HTML(http.StatusForbidden, "login.html", nil)
		c.Abort()
		return
	}

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

	authRouter := router.Group("/", AuthFunc)

	// Init the sessions map
	store = session.Init()
	// Create the AuthGroup which will include everything except /login and /loginform

	authRouter.POST("/add-item/:importance", addTodoHandler)
	authRouter.DELETE("/delete/:importance/:id", deleteItemHandler)
	authRouter.GET("/index", indexHandler)
	router.GET("/login", loginHandler)
	router.POST("/loginform", loginFormHandler)
	router.Run()
}

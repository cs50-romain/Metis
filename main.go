package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cs50-romain/Metis/util/session"
	ttask "github.com/cs50-romain/Metis/util/task"
)

const PATH_SEP_WINDOWS = '\\'
const PATH_SEP_LINUX = '/'	

func empty(content string) bool{
	if content == "" || content == " " {
		fmt.Println("Empty")
		return true 
	}
	return false
}

func sessionMiddleware(r *http.Request) *session.Session {
	// Do you have an exisiting session?
		// If not expired -> return the session
		// If expired, delete the session and return nil -> route will then redirect to the login page
	// If no existing session -> return nil and route will redirect to login page
	sessionID, err := getSessionIdCookie(r)
	if err != nil {
		log.Println("[ERROR] -> ", err)
	}

	if sessionID.String() != "" {
		session_pointer, ok := session.GetSession(strings.TrimLeft(sessionID.String(), "sessionID="))
		if ok && session_pointer.ExpiresAt.After(time.Now()) {
			return session_pointer
		} else {
			session.DeleteSession(strings.TrimLeft(sessionID.String(), "sessionID="))
			session.Save()
			return nil
		}
	} else {
		return nil
	}
}

func Route(w http.ResponseWriter, r *http.Request) {
	log.Printf("[REQUEST] for %s; Routing...\n", r.URL.Path)
	requestpath := strings.Split(r.URL.Path, "/")
	urlpath := requestpath[1]
	importance := requestpath[len(requestpath)-1]
	var session_pointer *session.Session

	if urlpath == "loginform" {
		loginFormHandler(w, r)
	} else {
		session_pointer = sessionMiddleware(r)
	}

	if session_pointer == nil {
		log.Println("[INFO] Redirecting to /login")
		loginHandler(w, r)
		return
	} else { 
		if r.URL.Path == "/" {
			index(w, r, session_pointer)
		} else if urlpath == "login" {
			loginHandler(w, r)
		}else if urlpath == "add-item" {
			AddItem(w, r, importance, session_pointer)
			session.Save()
		} else if urlpath == "delete" {
			DeleteItem(w, r, session_pointer)
			session.Save()
		} else if urlpath == "itemcompleted" {
			ItemCompleted(w, r, requestpath, session_pointer)
			session.Save()
		} else {
			log.Printf("Invalid Path Request: %s\n", r.URL.Path)
			http.Error(w, "Invalid Path Request", http.StatusBadRequest)
		}
	}
	fmt.Print("End of route -> showing sessions data:")
	session.Print()
	fmt.Println()
}

func AddItem(w http.ResponseWriter, r *http.Request, importance string, session *session.Session) {
	fmt.Println("Adding item to:", session)
	if r.Method == "POST" {
		var newid int
		content := r.PostFormValue("content")

		if empty(content) {
			http.Error(w, "No content", http.StatusBadRequest)
		} else {
			if importance == "important" {
				newid = len(session.UserData.Itasks)
				session.UserData.AddTaskToList(ttask.Task{len(session.UserData.Itasks), content, false, time.Now(), "important"}, importance)
			} else if importance == "minor" {
				newid = len(session.UserData.Mtasks)
				session.UserData.AddTaskToList(ttask.Task{len(session.UserData.Mtasks), content, false, time.Now(), "minor"}, importance)
			} else if importance == "later" {
				newid = len(session.UserData.Ltasks)
				session.UserData.AddTaskToList(ttask.Task{len(session.UserData.Ltasks), content, false, time.Now(), "later"}, importance)
			}

			htmlEl := fmt.Sprintf("<li id=%b class='todo-item' hx-trigger='change delay:2s' hx-target='#completed-box' hx-include='this' hx-post='/itemcompleted/%s/%b' hx-swap='beforeend'>\n<label>\n<input hx-delete='/delete/important/%b' hx-trigger='click delay:4s' hx-target='closest li' hx-swap='delete' type='checkbox'><span>%s</span>\n</label>\n<button hx-delete='/delete/%s/%b' hx-trigger='click' hx-confirm='Are you sure?' hx-target='closest li' hx-swap='delete'>Delete</button>\n</li>", newid, importance, newid, newid, content, importance, newid)
			tmpl, _ := template.New("t").Parse(htmlEl)
			tmpl.Execute(w, nil)
		}
		fmt.Println("Session's data:", session)
		saveFile(session)
	}
}

func DeleteItem(w http.ResponseWriter, r *http.Request, session *session.Session) {
	log.Print("[REQUEST] from Route to /delete")
	path := strings.Split(r.URL.Path, "/")
	importance := path[len(path)-2]
	taskId := path[len(path)-1]
	id, err := strconv.Atoi(taskId)

	log.Printf("INFO: Object information -> id:%b, importance:%s", id, importance)

	if err != nil {
		log.Print("ERROR (line 94): String Conversion Error:", err)
	}

	if importance == "important" {
		for idx, object := range session.UserData.Itasks {
			if object.Id == id {
				if len(session.UserData.Itasks) > 1 {
					session.UserData.Itasks = append(session.UserData.Itasks[:idx], session.UserData.Itasks[idx+1:]...)
				} else {
					session.UserData.Itasks = []ttask.Task{}
				}
			}
		}
	} else if importance == "minor" {
		for idx, object := range session.UserData.Mtasks {
			if object.Id == id {
				session.UserData.Mtasks = append(session.UserData.Mtasks[:idx], session.UserData.Mtasks[idx+1:]...)
			}
		}
	} else if importance == "later" {
		for idx, object := range session.UserData.Ltasks {
			if object.Id == id {
				session.UserData.Ltasks = append(session.UserData.Ltasks[:idx], session.UserData.Ltasks[idx+1:]...)
			}
		}
	} else {
		log.Printf("Error; Could not find importance type for: %s\n", importance)
	}

	// FIX ID OF TASKS
	ttask.FixId(session.UserData.Ltasks)
	ttask.FixId(session.UserData.Mtasks)
	ttask.FixId(session.UserData.Itasks)

	saveFile(session)
}

func ItemCompleted(w http.ResponseWriter, r *http.Request, path []string, session *session.Session) {
	log.Print("[REQUEST] from Route to /itemcompleted received")
	task_id_str := path[len(path)-1]
	task_id, _ := strconv.Atoi(task_id_str)
	importance := path[len(path)-2]

	task := ttask.GetTask(importance, task_id)
	task.IsCompleted = true
	task.Importance = "completed"
	task.CreatedAt = time.Now()
	//DeleteTask(importance, task_id)
	session.UserData.AddTaskToList(task, "completed")

	year, month, day := task.CreatedAt.Date()

	htmlEl := fmt.Sprintf("<li id=%b class='todo-item'>\n<label><span>%s</span><span style='font-size: 10px;'>  on:%v-%v-%v</span>\n</label>\n</li>", task.Id, task.Content, year, month, day)

	tmpl, _ := template.New("t").Parse(htmlEl)
	tmpl.Execute(w, nil)
	saveFile(session)
}

func index(w http.ResponseWriter, r *http.Request, session *session.Session) {
	log.Print("[REQUEST] from Route to index received.")
	tmpl := template.Must(template.ParseFiles("./static/index.html"))
	ttask.Prefill(session.UserData)

	// FIX ID OF TASKS
	ttask.FixId(session.UserData.Ltasks)
	ttask.FixId(session.UserData.Mtasks)
	ttask.FixId(session.UserData.Itasks)

	data := struct {
		ImportantTasks	[]ttask.Task
		MinorTasks	[]ttask.Task
		LaterTasks	[]ttask.Task
		CompletedTasks	[]ttask.Task
	}{
		ImportantTasks: session.UserData.Itasks,
		MinorTasks: session.UserData.Mtasks,
		LaterTasks: session.UserData.Ltasks,
		CompletedTasks: session.UserData.Ctasks,
	}

	tmpl.Execute(w, data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[REQUEST] from Route to /login")
	http.ServeFile(w, r, "./static/login.html")
}

func loginFormHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[REQUEST] from Route to /loginform")
	if r.Method == "POST" {
		username := r.PostFormValue("username")
		//userpass := r.PostFormValue("password")
		
		// check if session id alread exists, if so redirect
		session_pointer := sessionMiddleware(r)
		if session_pointer != nil {
			log.Println("[TESTING] found a common sessionid")
			http.Redirect(w, r, "/", 302)
		} else {
			sessionid := session.CreateSession(username)
			http.SetCookie(w, &http.Cookie{
				Name:	"sessionID",
				Value:	sessionid,
				Expires: time.Now().Add(time.Hour * 24 * 10),
			})
			session.Save()

			http.Redirect(w, r, "/", 302)
		}
	}
}


func getSessionIdCookie(r *http.Request) (*http.Cookie, error) {
	return r.Cookie("sessionID")
}

func saveFile(session *session.Session) {
	// Once sessions are created, check to see if the session is inactive, if so save the file and exit.
	tasks := ttask.Tasks{
		Itasks: session.UserData.Itasks,
		Mtasks: session.UserData.Mtasks,
		Ltasks: session.UserData.Ltasks,
		Ctasks: session.UserData.Ctasks,
	}

	b, err := json.Marshal(tasks)
	if err != nil {
		log.Print("[INFO] Couldn't save file, retrying in...")
	} else {
		fmt.Println("[INFO] File saved")
	}

	if os.PathSeparator == PATH_SEP_WINDOWS {
		err = os.WriteFile("C:\\Program Files\\Metis\\data.json", b, 0777)
		if err != nil {
			log.Print("[ERROR] Error saving to file")
		}
	} else {
		file_saving_location := "./data/" + session.UserName
		if _, err := os.Stat(file_saving_location); err != nil {
		    if os.IsNotExist(err) {
			direrr := os.Mkdir(file_saving_location, 0777)
			if direrr != nil {
				log.Print("[ERROR] Could not create directory: ", direrr)
			}
		    } else {
			log.Print("[ERROR] Path error for os.Stat")
		    }
		}

		err = os.WriteFile(file_saving_location + "/data.json", b, 0644)
		if err != nil {
			log.Print("[ERROR] Error saving to file -> ", err)
		}
	}
}

func main() {
	fmt.Println("[+] Decoding json information if any...")
	//prefill()
	session.PreFillSessions()
	session.Print()

	fmt.Println("[+] Starting web server...")

	http.HandleFunc("/", Route)
	log.Print(http.ListenAndServe(":8080", nil))
}

/*
TODO
1. Add ability to search for a youtube video. Make window a little bit bigger and in the future resizable - TOMORROW
2. Save json to file: import and export the json data. - 2023-12-04 -> DONE 
3. Session per client connection - 2023-12-09 -> 
4. Keyboard shortcuts - STARTED (functionality is there, just not pretty) -> 12/08/2023 - REMOVED FOR NOW
5. Draggable task items (within flex and from one flex, like from Later section to Important to another)
6. Move items from Later to Minor after 4 days - 2023-12-05 -> DONE 
7. Move items from Minor to Important after 7 days - 12/05/2023 -> DONE
8. Add ability to delete items/tasks - 12/06/2023 - FRONTEND DONE / WORKING ON BACKEND -> 12/08/2023 - DONE 
9. Add ability to mark task as completed (add sound effect) - TOMORROW -> 12/08/2023 - DONE
9. Publish website
*/

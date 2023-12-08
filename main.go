package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"net/http"
	"log"
	"strings"
	"strconv"
	"time"
)

var itasks []Task
var mtasks []Task
var ltasks []Task

const PATH_SEP_WINDOWS = '\\'
const PATH_SEP_LINUX = '/'	

type Tasks struct {
	Itasks	[]Task
	Mtasks	[]Task
	Ltasks	[]Task
}

type Task struct {
	Id		int
	Content		string
	isCompleted	bool
	CreatedAt	time.Time
	Importance	string
}

func empty(content string) bool{
	if content == "" || content == " " {
		fmt.Println("Empty")
		return true 
	}
	return false
}

func Route(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received for %s; Routing...\n", r.URL.Path)
	requestpath := strings.Split(r.URL.Path, "/")
	urlpath := requestpath[1]
	importance := requestpath[len(requestpath)-1]

	if r.URL.Path == "/" {
		index(w, r)
	} else if urlpath == "add-item" {
		AddItem(w, r, importance)
	} else if urlpath == "delete" {
		DeleteItem(w, r)
	} else {
		log.Printf("Invalid Path Request: %s\n", r.URL.Path)
		http.Error(w, "Invalid Path Request", http.StatusBadRequest)
	}
}

func AddItem(w http.ResponseWriter, r *http.Request, importance string) {
	if r.Method == "POST" {
		var newid int
		content := r.PostFormValue("content")

		if empty(content) {
			http.Error(w, "No content", http.StatusBadRequest)
		} else {
			if importance == "important" {
				newid = len(itasks)
				itasks = append(itasks, Task{len(itasks), content, false, time.Now(), "later"})
			} else if importance == "minor" {
				newid = len(mtasks)
				mtasks = append(mtasks, Task{len(mtasks), content, false, time.Now(), "later"})
			} else if importance == "later" {
				newid = len(ltasks)
				ltasks = append(ltasks, Task{len(ltasks), content, false, time.Now(), "later"})
			}

			go saveFile()
			htmlEl := fmt.Sprintf("<li id=%b class='todo-item'>\n<label>\n<input type='checkbox'>%s\n</label>\n<button hx-delete='/delete/%s/%b' hx-trigger='click' hx-confirm='Are you sure?' hx-target='closest li' hx-swap='delete'>Delete</button>\n</li>", newid, content, importance, newid)
			tmpl, _ := template.New("t").Parse(htmlEl)
			tmpl.Execute(w, nil)
		}
	}
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	log.Print("INFO: Delete request received")
	path := strings.Split(r.URL.Path, "/")
	importance := path[len(path)-2]
	taskId := path[len(path)-1]
	id, err := strconv.Atoi(taskId)

	log.Printf("INFO: Object information -> id:%b, importance:%s task:%s", id, importance)

	if err != nil {
		log.Print("ERROR (line 94): String Conversion Error:", err)
	}

	if importance == "important" {
		for idx, task := range itasks {
			if task.Id == id {
				log.Print("INFO: Found related task:", task)
				itasks = append(itasks[:idx], itasks[idx+1:]...)
				log.Print("INFO: New itasks:", itasks)
			}
		}
	} else if importance == "minor" {
		for idx, task := range mtasks {
			if task.Id == id {
				mtasks = append(mtasks[:idx], mtasks[idx+1:]...)
			}
		}
	} else if importance == "later" {
		for idx, task := range ltasks {
			if task.Id == id {
				ltasks = append(ltasks[:idx], ltasks[idx+1:]...)
			}
		}
	} else {
		log.Printf("Error; Could not find importance type for: %s\n", importance)
	}

	itasks = FixId(itasks)
	mtasks = FixId(mtasks)
	ltasks = FixId(ltasks)

	saveFile()
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Print("Index request received")
	tmpl := template.Must(template.ParseFiles("./static/index.html"))

	itasks = FixId(itasks)
	mtasks = FixId(mtasks)
	ltasks = FixId(ltasks)

	data := struct {
		ImportantTasks	[]Task
		MinorTasks	[]Task
		LaterTasks	[]Task
	}{
		ImportantTasks: itasks,
		MinorTasks: mtasks,
		LaterTasks: ltasks,
	}

	tmpl.Execute(w, data)
}

func saveFile() {
	// Once sessions are created, check to see if the session is inactive, if so save the file and exit.
	tasks := Tasks{
		Itasks: itasks,
		Mtasks: mtasks,
		Ltasks: ltasks,
	}

	b, err := json.Marshal(tasks)
	if err != nil {
		log.Print("{!] Couldn't save file, retrying in...")
	} else {
		fmt.Println("[+] File saved")
	}

	if os.PathSeparator == PATH_SEP_WINDOWS {
		err = os.WriteFile("C:\\Program Files\\Metis\\data.json", b, 0644)
		if err != nil {
			log.Print("[!] Error saving to file")
		}
	} else {
		err = os.WriteFile("./data.json", b, 0644)
		if err != nil {
			log.Print("[!] Error saving to file")
		}
	}
}

func prefill() {
	var tasks Tasks
	b, err := os.ReadFile("./data.json")
	err = json.Unmarshal(b, &tasks)
	if err != nil {
		fmt.Println(err)
	}

	itasks = tasks.Itasks
	mtasks = tasks.Mtasks
	ltasks = tasks.Ltasks

	itasks = FixId(itasks)
	mtasks = FixId(mtasks)
	ltasks = FixId(ltasks)

	// AFTER PREFILLING, CHECK DATES AND SEE IF AN OBJECT NEEDS TO SWITCH
	// CHECK EACH TASK LISTS AND MOVE TASK ACCORDINGLY IF NEEDED
	swapTasksBasedOnDate()
}

func swapTasksBasedOnDate() {
	for idx, object := range mtasks {
		//log.Printf("Object %s creation Date:%s\t Current Date:%s", object.Content, object.CreatedAt.String(), time.Now().String())
		if compareDate(object.CreatedAt) >= 7 {
			object.Importance = "important"
			itasks = append(itasks, object)
			mtasks = append(mtasks[:idx], mtasks[idx+1:]...)
		}
	}

	for idx, object := range ltasks {
		//log.Printf("Object %s creation Date:%s\t Current Date:%s", object.Content, object.CreatedAt.String(), time.Now().String())
		if compareDate(object.CreatedAt) >= 4 {
			object.Importance = "minor"
			mtasks = append(mtasks, object)
			ltasks = append(ltasks[:idx], ltasks[idx+1:]...)
		}
	}
}

// compareDate 2 DATES (DATE1 AND TIME.NOW()) AND RETURN AN INTEGER - INTEGER WILL BE HOW MANY DAYS DATE1 IS FROM TIME.NOW()
func compareDate(date time.Time) int{
	if date.Year() == time.Now().Year() {
		if date.Month() == time.Now().Month() {
			return time.Now().Day() - date.Day()
		} else {
			monthDiff := int(time.Now().Month()) - int(date.Month())
			startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
			var daysTillEndOfMonth int = 0 
			if monthDiff <= 1 {
				if date.Month() == time.January || date.Month() == time.March || date.Month() == time.May || date.Month() == time.July || date.Month() == time.August || date.Month() == time.October || date.Month() == time.December {
					daysTillEndOfMonth = 31 - date.Day()
				} else if date.Month() == time.February {
					if date.Year() / 2 == 0 {
						daysTillEndOfMonth = 29 - date.Day() + 1
					} else {
						daysTillEndOfMonth = 28 - date.Day() + 1
					}
					
				} else {
					daysTillEndOfMonth = 30 - date.Day() + 1
				}
				return time.Now().Day() - startOfMonth.Day() + daysTillEndOfMonth 
			} else {
				return -1
			}
		}
	} else {
		return -1
		//return time.Now().Year() - date.Year() + 365 
	}
	return -1
}

func FixId(array []Task) []Task {
	if array == nil {
		return array
	}

	for idx,_ := range array {
		array[idx].Id = idx
	}
	return array
}

func main() {
	fmt.Println("[+] Decoding json information if any...")
	prefill()

	fmt.Println("[+] Starting web server...")

	http.HandleFunc("/", Route)

	log.Print(http.ListenAndServe(":8080", nil))
}

/*
TODO
1. Add ability to search for a youtube video. Make window a little bit bigger and in the future resizable - TOMORROW
2. Save json to file: import and export the json data. - 2023-12-04 -> DONE 
3. Session per client connection - THURSDAY
4. Keyboard shortcuts - STARTED (functionality is there, just not pretty)
5. Draggable task items (within flex and from one flex, like from Later section to Important to another)
6. Move items from Later to Minor after 4 days - 2023-12-05 -> DONE 
7. Move items from Minor to Important after 7 days - 12/05/2023 -> DONE
8. Add ability to delete items/tasks - 12/06/2023 -> FRONTEND DONE / WORKING ON BACKEND 
9. Add ability to mark task as completed (add sound effect) - TOMORROW
9. Publish website
*/

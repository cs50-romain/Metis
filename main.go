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

var id int
var itasks []Task
var mtasks []Task
var ltasks []Task

const PATH_SEP_WINDOWS = '\\'
const PATHH_SEP_LINUX = '/'	

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

func addImportantItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		if empty(content) {

		}
			
		if len(itasks) <= 0 {
			htmlEl := fmt.Sprintf("<div class='flex m-auto text-center text-2xl'>\nImportant\n</div>\n<div class='flex flex-wrap space-x-2 flex-col p-4 w-2/6 h-32'>\n<div class='items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-important-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n<button hx-delete='/delete/minor/%b' hx-target='closest div' hx-swap='innerhtml' type='button' class='text-white bg-blue-700 hover:bg-red-800 focus:ring-2 focus:outline-none focus:ring-red-300 font-small rounded-full text-sm p-0.5 text-center inline-flex items-end me-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-black-800'>\n<svg xmlns='http://www.w3.org/2000/svg' width='15' height='15' viewBox='0 0 24 24'><path fill='currentColor' d='M16 1.75V3h5.25a.75.75 0 0 1 0 1.5H2.75a.75.75 0 0 1 0-1.5H8V1.75C8 .784 8.784 0 9.75 0h4.5C15.216 0 16 .784 16 1.75Zm-6.5 0V3h5V1.75a.25.25 0 0 0-.25-.25h-4.5a.25.25 0 0 0-.25.25ZM4.997 6.178a.75.75 0 1 0-1.493.144L4.916 20.92a1.75 1.75 0 0 0 1.742 1.58h10.684a1.75 1.75 0 0 0 1.742-1.581l1.413-14.597a.75.75 0 0 0-1.494-.144l-1.412 14.596a.25.25 0 0 1-.249.226H6.658a.25.25 0 0 1-.249-.226L4.997 6.178Z'/>\n<path fill='currentColor' d='M9.206 7.501a.75.75 0 0 1 .793.705l.5 8.5A.75.75 0 1 1 9 16.794l-.5-8.5a.75.75 0 0 1 .705-.793Zm6.293.793A.75.75 0 1 0 14 8.206l-.5 8.5a.75.75 0 0 0 1.498.088l.5-8.5Z'/>\n</svg>\n<span class='sr-only'>Icon description</span>\n</button>\n</div>\n</div>\n</div>\n<form id='minor-input-form' hx-post='/add-minor-item/' hx-target='#task-minor-list' hx-swap='afterend'>\n<div class='flex-none flex items-center border-b border-teal-500 py-2'>\n<input name='content' id='minor-input' class='appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none' type='text' placeholder='Add Item'>\n<button class='flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded' type='submit'>Add</button>\n</div>\n</form>\n", newid, content, newid)
			tmpl, _ := template.New("t").Parse(htmlEl)
			tmpl.Execute(w, nil)
		} else {

			htmlel := fmt.Sprintf("<div class='items-center bg-[#1da1f2] space-x-2 text-white rounded-lg mb-2' id='task-important-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewbox='0 0 20 20' fill='currentcolor'>\n<path fill-rule='evenodd' d='m16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414l8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n<button hx-delete='/delete/important/%b' hx-target='closest div' hx-swap='innerhtml' type='button' class='text-white bg-blue-700 hover:bg-red-800 focus:ring-2 focus:outline-none focus:ring-red-300 font-small rounded-full text-sm p-0.5 text-center inline-flex items-end me-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-black-800'>\n<svg xmlns='http://www.w3.org/2000/svg' width='15' height='15' viewBox='0 0 24 24'><path fill='currentColor' d='M16 1.75V3h5.25a.75.75 0 0 1 0 1.5H2.75a.75.75 0 0 1 0-1.5H8V1.75C8 .784 8.784 0 9.75 0h4.5C15.216 0 16 .784 16 1.75Zm-6.5 0V3h5V1.75a.25.25 0 0 0-.25-.25h-4.5a.25.25 0 0 0-.25.25ZM4.997 6.178a.75.75 0 1 0-1.493.144L4.916 20.92a1.75 1.75 0 0 0 1.742 1.58h10.684a1.75 1.75 0 0 0 1.742-1.581l1.413-14.597a.75.75 0 0 0-1.494-.144l-1.412 14.596a.25.25 0 0 1-.249.226H6.658a.25.25 0 0 1-.249-.226L4.997 6.178Z'/>\n<path fill='currentColor' d='M9.206 7.501a.75.75 0 0 1 .793.705l.5 8.5A.75.75 0 1 1 9 16.794l-.5-8.5a.75.75 0 0 1 .705-.793Zm6.293.793A.75.75 0 1 0 14 8.206l-.5 8.5a.75.75 0 0 0 1.498.088l.5-8.5Z'/>\n</svg>\n<span class='sr-only'>Icon description</span>\n</button>\n</div>\n</div>", newid, content, newid) 
			tmpl, _ := template.New("t").Parse(htmlel)
			tmpl.Execute(w, nil)
		}

		itasks = append(itasks, Task{newid, content, false, time.Now(), "important"})
		go saveFile()
	}
}

func addMinorItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		if empty(content) {

		}

		fmt.Println(len(mtasks))

		// The whole form is being re-rendered after we delete that last element in the list. Thisgives the effect of 2 of the same thing appearing. I need a way that this does not happen. Ideally I don't need to check the length and I can just render the htmlel that's in the else statement.
		if len(mtasks) <= 0 {
			fmt.Println("We are here")
			htmlEl := fmt.Sprintf("<div class='flex m-auto text-center text-2xl'>\nMinor\n</div>\n<div class='flex flex-wrap space-x-2 flex-col p-4 w-2/6 h-32'>\n<div class='items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-minor-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n<button hx-delete='/delete/minor/%b' hx-target='closest div' hx-swap='innerhtml' type='button' class='text-white bg-blue-700 hover:bg-red-800 focus:ring-2 focus:outline-none focus:ring-red-300 font-small rounded-full text-sm p-0.5 text-center inline-flex items-end me-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-black-800'>\n<svg xmlns='http://www.w3.org/2000/svg' width='15' height='15' viewBox='0 0 24 24'><path fill='currentColor' d='M16 1.75V3h5.25a.75.75 0 0 1 0 1.5H2.75a.75.75 0 0 1 0-1.5H8V1.75C8 .784 8.784 0 9.75 0h4.5C15.216 0 16 .784 16 1.75Zm-6.5 0V3h5V1.75a.25.25 0 0 0-.25-.25h-4.5a.25.25 0 0 0-.25.25ZM4.997 6.178a.75.75 0 1 0-1.493.144L4.916 20.92a1.75 1.75 0 0 0 1.742 1.58h10.684a1.75 1.75 0 0 0 1.742-1.581l1.413-14.597a.75.75 0 0 0-1.494-.144l-1.412 14.596a.25.25 0 0 1-.249.226H6.658a.25.25 0 0 1-.249-.226L4.997 6.178Z'/>\n<path fill='currentColor' d='M9.206 7.501a.75.75 0 0 1 .793.705l.5 8.5A.75.75 0 1 1 9 16.794l-.5-8.5a.75.75 0 0 1 .705-.793Zm6.293.793A.75.75 0 1 0 14 8.206l-.5 8.5a.75.75 0 0 0 1.498.088l.5-8.5Z'/>\n</svg>\n<span class='sr-only'>Icon description</span>\n</button>\n</div>\n</div>\n</div>\n<form id='minor-input-form' hx-post='/add-minor-item/' hx-target='#task-minor-list' hx-swap='afterend'>\n<div class='flex-none flex items-center border-b border-teal-500 py-2'>\n<input name='content' id='minor-input' class='appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none' type='text' placeholder='Add Item'>\n<button class='flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded' type='submit'>Add</button>\n</div>\n</form>\n", newid, content, newid)

			tmpl, _ := template.New("t").Parse(htmlEl)
			tmpl.Execute(w, nil)
		} else {

			htmlel := fmt.Sprintf("<div class='items-center bg-[#1da1f2] space-x-2 text-white rounded-lg mb-2' id='task-minor-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewbox='0 0 20 20' fill='currentcolor'>\n<path fill-rule='evenodd' d='m16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414l8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n<button hx-delete='/delete/minor/%b' hx-target='closest div' hx-swap='innerhtml' type='button' class='text-white bg-blue-700 hover:bg-red-800 focus:ring-2 focus:outline-none focus:ring-red-300 font-small rounded-full text-sm p-0.5 text-center inline-flex items-end me-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-black-800'>\n<svg xmlns='http://www.w3.org/2000/svg' width='15' height='15' viewBox='0 0 24 24'><path fill='currentColor' d='M16 1.75V3h5.25a.75.75 0 0 1 0 1.5H2.75a.75.75 0 0 1 0-1.5H8V1.75C8 .784 8.784 0 9.75 0h4.5C15.216 0 16 .784 16 1.75Zm-6.5 0V3h5V1.75a.25.25 0 0 0-.25-.25h-4.5a.25.25 0 0 0-.25.25ZM4.997 6.178a.75.75 0 1 0-1.493.144L4.916 20.92a1.75 1.75 0 0 0 1.742 1.58h10.684a1.75 1.75 0 0 0 1.742-1.581l1.413-14.597a.75.75 0 0 0-1.494-.144l-1.412 14.596a.25.25 0 0 1-.249.226H6.658a.25.25 0 0 1-.249-.226L4.997 6.178Z'/>\n<path fill='currentColor' d='M9.206 7.501a.75.75 0 0 1 .793.705l.5 8.5A.75.75 0 1 1 9 16.794l-.5-8.5a.75.75 0 0 1 .705-.793Zm6.293.793A.75.75 0 1 0 14 8.206l-.5 8.5a.75.75 0 0 0 1.498.088l.5-8.5Z'/>\n</svg>\n<span class='sr-only'>Icon description</span>\n</button>\n</div>\n</div>", newid, content, newid) 

			tmpl, _ := template.New("t").Parse(htmlel)
			tmpl.Execute(w, nil)
		}

		mtasks = append(mtasks, Task{newid, content, false, time.Now(), "minor"})
		go saveFile()
	}
}

func addLaterItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		if empty(content) {

		}

		if len(ltasks) <= 0 {
			htmlEl := fmt.Sprintf("<div class='flex m-auto text-center text-2xl'>\nFor Later\n</div>\n<div class='flex flex-wrap space-x-2 flex-col p-4 w-2/6 h-32'>\n<div class='items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-later-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>\n</div>\n<form id='later-input-form' hx-post='/add-later-item/' hx-target='#task-later-list' hx-swap='afterend'>\n<div class='flex-none flex items-center border-b border-teal-500 py-2'>\n<input name='content' id='later-input' class='appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none' type='text' placeholder='Add Item'>\n<button class='flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded' type='submit'>Add</button>\n</div>\n</form>\n", newid, content)

			tmpl, _ := template.New("t").Parse(htmlEl)
			tmpl.Execute(w, nil)
		} else {
			htmlel := fmt.Sprintf("<div class='items-center bg-[#1da1f2] space-x-2 text-white rounded-lg mb-2' id='task-later-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewbox='0 0 20 20' fill='currentcolor'>\n<path fill-rule='evenodd' d='m16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414l8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n<button hx-delete='/delete/later/%b' hx-target='closest div' hx-swap='innerhtml' type='button' class='text-white bg-blue-700 hover:bg-red-800 focus:ring-2 focus:outline-none focus:ring-red-300 font-small rounded-full text-sm p-0.5 text-center inline-flex items-end me-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-black-800'>\n<svg xmlns='http://www.w3.org/2000/svg' width='15' height='15' viewBox='0 0 24 24'><path fill='currentColor' d='M16 1.75V3h5.25a.75.75 0 0 1 0 1.5H2.75a.75.75 0 0 1 0-1.5H8V1.75C8 .784 8.784 0 9.75 0h4.5C15.216 0 16 .784 16 1.75Zm-6.5 0V3h5V1.75a.25.25 0 0 0-.25-.25h-4.5a.25.25 0 0 0-.25.25ZM4.997 6.178a.75.75 0 1 0-1.493.144L4.916 20.92a1.75 1.75 0 0 0 1.742 1.58h10.684a1.75 1.75 0 0 0 1.742-1.581l1.413-14.597a.75.75 0 0 0-1.494-.144l-1.412 14.596a.25.25 0 0 1-.249.226H6.658a.25.25 0 0 1-.249-.226L4.997 6.178Z'/>\n<path fill='currentColor' d='M9.206 7.501a.75.75 0 0 1 .793.705l.5 8.5A.75.75 0 1 1 9 16.794l-.5-8.5a.75.75 0 0 1 .705-.793Zm6.293.793A.75.75 0 1 0 14 8.206l-.5 8.5a.75.75 0 0 0 1.498.088l.5-8.5Z'/>\n</svg>\n<span class='sr-only'>Icon description</span>\n</button>\n</div>\n</div>", newid, content, newid) 

			tmpl, _ := template.New("t").Parse(htmlel)
			tmpl.Execute(w, nil)
		}

		ltasks = append(ltasks, Task{newid, content, false, time.Now(), "later"})
		go saveFile()
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	log.Print("Delete request received")
	path := strings.Split(r.URL.Path, "/")
	idString := path[len(path)-1]
	fmt.Println(idString)
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Print("Line 135", err)
	}
	importance := path[len(path)-2]

	if importance == "important" {
		for idx, task := range itasks {
			if task.Id == id {
				itasks = append(itasks[:idx], itasks[idx+1:]...)
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
	saveFile()
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Print("Index request received")
	tmpl := template.Must(template.ParseFiles("./static/index.html"))
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

	fmt.Println(itasks)
	fmt.Println(mtasks)
	fmt.Println(ltasks)

	// AFTER PREFILLING, CHECK DATES AND SEE IF AN OBJECT NEEDS TO SWITCH
	// CHECK EACH TASK LISTS AND MOVE TASK ACCORDINGLY IF NEEDED
	swapTasksBasedOnDate()
}

func swapTasksBasedOnDate() {
	for idx, object := range mtasks {
		log.Printf("Object %s creation Date:%s\t Current Date:%s", object.Content, object.CreatedAt.String(), time.Now().String())
		if compareDate(object.CreatedAt) >= 7 {
			object.Importance = "important"
			itasks = append(itasks, object)
			mtasks = append(mtasks[:idx], mtasks[idx+1:]...)
		}
	}

	for idx, object := range ltasks {
		log.Printf("Object %s creation Date:%s\t Current Date:%s", object.Content, object.CreatedAt.String(), time.Now().String())
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

func getId() int {
	newid := 0
	for _, obj := range itasks {
		if obj.Id > id {
			newid = obj.Id
		}
	}

	for _, obj := range mtasks {
		if obj.Id > id {
			newid = obj.Id
		}
	}

	for _, obj := range ltasks {
		if obj.Id > id {
			newid = obj.Id
		}
	}
	return newid + 1
}

func main() {
	id = getId()

	fmt.Println("[+] Decoding json information if any...")
	prefill()
	fmt.Println("[+] Starting web server...")

	http.HandleFunc("/", index)
	http.HandleFunc("/add-important-item/", addImportantItem)
	http.HandleFunc("/add-minor-item/", addMinorItem)
	http.HandleFunc("/add-later-item/", addLaterItem)
	http.HandleFunc("/delete/", deleteItem)
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

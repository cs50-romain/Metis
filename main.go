package main

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
)

var id int

type ImportantTask struct {
	Id		int
	Content		string
	isCompleted	bool
}

func addImportantItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	fmt.Println(r.Method)
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		htmlEl := fmt.Sprintf("<div class='flex-initial items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-10 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>", newid, content)
		tmpl, _ := template.New("t").Parse(htmlEl)
		tmpl.Execute(w, nil)

		fmt.Println(content, newid)
	}
}

func addMinorItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	fmt.Println(r.Method)
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		htmlEl := fmt.Sprintf("<div class='flex-initial items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-minor-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-10 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>", newid, content)
		tmpl, _ := template.New("t").Parse(htmlEl)
		tmpl.Execute(w, nil)

		fmt.Println(content, newid)
	}
}

func addLaterItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	fmt.Println(r.Method)
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		htmlEl := fmt.Sprintf("<div class='flex-initial items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-later-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-10 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>", newid, content)
		tmpl, _ := template.New("t").Parse(htmlEl)
		tmpl.Execute(w, nil)

		fmt.Println(content, newid)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./static/index.html"))
	tasks := map[string][]ImportantTask{
		"Tasks": {
			{Id: 1, Content: "Call consulate.", isCompleted: false},
			{Id: 2, Content: "D this.", isCompleted: false},
			{Id: 3, Content: "Call doctor.", isCompleted: false},
		},
	}
	tmpl.Execute(w, tasks)
}

func main() {
	id = 0
	fmt.Println("[+] Starting server...")
	http.HandleFunc("/", index)
	http.HandleFunc("/add-important-item/", addImportantItem)
	http.HandleFunc("/add-minor-item/", addMinorItem)
	http.HandleFunc("/add-later-item/", addLaterItem)
	log.Print(http.ListenAndServe(":8080", nil))
}

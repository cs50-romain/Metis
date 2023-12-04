package main

import (
	//"encoding/json"
	"fmt"
	"html/template"
	//"io"
	"net/http"
	"log"
	"time"
)

var id int
var itasks []ImportantTask
var mtasks []MinorTask
var ltasks []LaterTask

type Tasks struct {
	itasks	[]ImportantTask
	mtasks	[]MinorTask
	ltasks	[]LaterTask
}

type ImportantTask struct {
	Id		int
	Content		string
	isCompleted	bool
}

type MinorTask struct {
	Id		int
	Content		string
	isCompleted	bool
	CreatedAt	time.Time
}

type LaterTask struct {
	Id		int
	Content		string
	isCompleted	bool
	CreatedAt	time.Time
}

func addImportantItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		if len(itasks) <= 0 {
			htmlEl := fmt.Sprintf("<div class='flex m-auto text-center text-2xl'>\nImportant\n</div>\n<div class='flex flex-wrap space-x-2 flex-col p-4 w-2/6 h-32'>\n<div class='items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-important-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>\n</div>\n<form hx-post='/add-important-item/' hx-target='#task-important-list' hx-swap='afterend'>\n<div class='flex-none flex items-center border-b border-teal-500 py-2'>\n<input name='content' class='appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none' type='text' placeholder='Add Item'>\n<button class='flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded' type='submit'>Add</button>\n</div>\n</form>\n", newid, content)

			tmpl, _ := template.New("t").Parse(htmlEl)
			tmpl.Execute(w, nil)
		} else {
			htmlel := fmt.Sprintf("<div class='items-center bg-[#1da1f2] space-x-2 text-white rounded-lg mb-2' id='task-important-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewbox='0 0 20 20' fill='currentcolor'>\n<path fill-rule='evenodd' d='m16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414l8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>", newid, content)
			tmpl, _ := template.New("t").Parse(htmlel)
			tmpl.Execute(w, nil)
		}

		itasks = append(itasks, ImportantTask{newid, content, false})
		fmt.Println(itasks, mtasks, ltasks)
	}
}

func addMinorItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		if len(ltasks) <= 0 {
			htmlEl := fmt.Sprintf("<div class='flex m-auto text-center text-2xl'>\nMinor\n</div>\n<div class='flex flex-wrap space-x-2 flex-col p-4 w-2/6 h-32'>\n<div class='items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-minor-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>\n</div>\n<form hx-post='/add-minor-item/' hx-target='#task-minor-list' hx-swap='afterend'>\n<div class='flex-none flex items-center border-b border-teal-500 py-2'>\n<input name='content' class='appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none' type='text' placeholder='Add Item'>\n<button class='flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded' type='submit'>Add</button>\n</div>\n</form>\n", newid, content)

			tmpl, _ := template.New("t").Parse(htmlEl)
			tmpl.Execute(w, nil)
		} else {
			htmlel := fmt.Sprintf("<div class='items-center bg-[#1da1f2] space-x-2 text-white rounded-lg mb-2' id='task-minor-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewbox='0 0 20 20' fill='currentcolor'>\n<path fill-rule='evenodd' d='m16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414l8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>", newid, content)
			tmpl, _ := template.New("t").Parse(htmlel)
			tmpl.Execute(w, nil)
		}

		mtasks = append(mtasks, MinorTask{newid, content, false, time.Now()})
	}
}

func addLaterItem(w http.ResponseWriter, r *http.Request) {
	log.Print("HTMX request received")
	log.Print(r.Header.Get("HX-REQUEST"))
	if r.Method == "POST" {
		content := r.PostFormValue("content")
		newid := id
		id++

		if len(ltasks) <= 0 {
			htmlEl := fmt.Sprintf("<div class='flex m-auto text-center text-2xl'>\nFor Later\n</div>\n<div class='flex flex-wrap space-x-2 flex-col p-4 w-2/6 h-32'>\n<div class='items-center bg-[#1da1f2] text-white rounded-lg mb-2' id='task-later-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org2000/svg' viewBox='0 0 20 20' fill='currentColor'>\n<path fill-rule='evenodd' d='M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>\n</div>\n<form hx-post='/add-later-item/' hx-target='#task-later-list' hx-swap='afterend'>\n<div class='flex-none flex items-center border-b border-teal-500 py-2'>\n<input name='content' class='appearance-none bg-transparent border-none w-full text-gray-700 mr-3 py-1 px-2 leading-tight focus:outline-none' type='text' placeholder='Add Item'>\n<button class='flex-shrink-0 bg-teal-500 hover:bg-teal-700 border-teal-500 hover:border-teal-700 text-sm border-4 text-white py-1 px-2 rounded' type='submit'>Add</button>\n</div>\n</form>\n", newid, content)

			tmpl, _ := template.New("t").Parse(htmlEl)
			tmpl.Execute(w, nil)
		} else {
			htmlel := fmt.Sprintf("<div class='items-center bg-[#1da1f2] space-x-2 text-white rounded-lg mb-2' id='task-later-list'>\n<div>\n<input class='hidden' type='checkbox' id=%b name='item' checked>\n<label class='flex items-center h-7 px-2 rounded cursor-pointer' for='task_1'>\n<span class='flex items-center justify-center w-5 h-5 text-transparent border-2 border-gray-300 rounded-full'>\n<svg class='w-4 h-4 fill-current' xmlns='http://www.w3.org/2000/svg' viewbox='0 0 20 20' fill='currentcolor'>\n<path fill-rule='evenodd' d='m16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414l8 12.586l7.293-7.293a1 1 0 011.414 0z' clip-rule='evenodd' />\n</svg>\n</span>\n<span class='ml-4 text-sm'>%s</span>\n</label>\n</div>\n</div>", newid, content)
			tmpl, _ := template.New("t").Parse(htmlel)
			tmpl.Execute(w, nil)
		}

		ltasks = append(ltasks, LaterTask{newid, content, false, time.Now()})
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Print("Index request received")
	tmpl := template.Must(template.ParseFiles("./static/index.html"))
	/*
	importanttasks := map[string][]ImportantTask{
		"ImportantTasks": {
			{Id: 1, Content: "Call consulate.", isCompleted: false},
			{Id: 2, Content: "D this.", isCompleted: false},
			{Id: 3, Content: "Call doctor.", isCompleted: false},
		},
	}
	*/

	data := struct {
		ImportantTasks	[]ImportantTask
		MinorTasks	[]MinorTask
		LaterTasks	[]LaterTask
	}{
		ImportantTasks: itasks,
		MinorTasks: mtasks,
		LaterTasks: ltasks,
	}

	tmpl.Execute(w, data)
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

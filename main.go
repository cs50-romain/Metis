package main

import (
	"fmt"
	"net/http"
	"log"
)

func main() {
	fmt.Println("Hello")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Print(http.ListenAndServe(":8080", nil))
}

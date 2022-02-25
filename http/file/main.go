package main

import (
	"flag"
	"log"
	"net/http"
)

var choice int

func init() {
	flag.IntVar(&choice, "choice", 0, "1: static, 2: withoutHidingFile, 3: more")
}

func main() {
	flag.Parse()
	switch choice {
	case 1:
		staticHTTPServer()
	case 2:
		staticHTTPServerWithoutHidingFile()
	case 3:
		http.Handle("/abc", http.StripPrefix("/abc", http.FileServer(http.Dir("/home/go"))))
		log.Fatal(http.ListenAndServe(":8080", nil))
	default:
		log.Fatal("Invalid choice")
	}
}

func staticHTTPServer() {
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/home/go"))))
}

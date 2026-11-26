package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("página del Taller 13 disponible en http://localhost:5173")
	log.Fatal(http.ListenAndServe("127.0.0.1:5173", http.FileServer(http.Dir("."))))
}

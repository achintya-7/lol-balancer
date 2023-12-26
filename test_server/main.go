package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	port1 = ":3000"
	port2 = ":3001"
)

func server1() {
	server := http.Server{
		Addr: port1,
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Server 1")
			},
		),
	}

	log.Println("Server 1 started")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func server2() {
	server := http.Server{
		Addr: port2,
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Server 2")
			},
		),
	}

	log.Println("Server 2 started")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	go server1()
	go server2() 

	select {}
}
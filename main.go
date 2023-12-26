package main

import (
	"log"
	"lol-balancer/pkg/backend"
	serverpool "lol-balancer/pkg/server_pool"
	"net/http"
	"net/url"
	"time"
)

const (
	server1 = "http://localhost:3000"
	server2 = "http://localhost:3001"
	port    = ":2205"
)

func main() {
	server_urls := []string{server1, server2}
	backnedServers := []*backend.Backend{}

	for _, server_url := range server_urls {
		url, err := url.Parse(server_url)
		if err != nil {
			log.Fatal(err)
		}

		backnedServer := backend.NewBackendServer(url)
		backnedServers = append(backnedServers, backnedServer)
	}

	serverpool := serverpool.NewServerPool(backnedServers)

	go healthCheck(serverpool)

	server := http.Server{
		Addr: port,
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				loadBalancer(w, r, serverpool)
			},
		),
	}

	log.Println("Load balancer started")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func loadBalancer(w http.ResponseWriter, r *http.Request, serverpool *serverpool.ServerPool) {
	peer := serverpool.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// * A passive health check where we check the health of the backends every 20 seconds
func healthCheck(serverpool *serverpool.ServerPool) {
	// initial health check
	serverpool.HealthCheck()
	
	t := time.NewTicker(time.Second * 20)
	for range t.C {
		log.Println("Starting health check...")
		serverpool.HealthCheck()
		log.Println("Health check completed")
	}
}

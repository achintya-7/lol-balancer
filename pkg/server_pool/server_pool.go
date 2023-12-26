package serverpool

import (
	"log"
	"lol-balancer/pkg/backend"
	"net"
	"net/url"
	"sync/atomic"
	"time"
)

type ServerPool struct {
	backends []*backend.Backend
	current  uint32
}

func NewServerPool(backends []*backend.Backend) *ServerPool {
	return &ServerPool{
		backends: backends,
		current:  0,
	}
}

// ? Get the next active backend from the server pool
// We use the atomic package to safely handle concurrent access from multiple goroutines
// and then we use the modulo operator to get the next index in the slice without going out of bounds
func (s *ServerPool) NextIndex() uint32 {
	totalLength := uint32(len(s.backends))
	idx := atomic.AddUint32(&s.current, uint32(1)) % totalLength
	return idx
}

// ? Get the next active backend from the server pool
func (s *ServerPool) GetNextPeer() *backend.Backend {
	// we will first get the next peer
	next := s.NextIndex()

	// now we have to traverse the backends array
	// and find the next active peer starting from the next index
	totalLength := len(s.backends)
	for i := 0; i < totalLength; i++ {
		index := int((i + int(next)) % totalLength)

		// if we have an active peer, assign it as the current index
		// and return the backend item from the array
		if s.backends[index].IsAlive() {
			log.Printf("current peer index -> %d\n", index)
			if i != int(next) {
				atomic.StoreUint32(&s.current, uint32(index))
			}
			return s.backends[index]
		}
	}

	return nil
}

// ? Function to health check all the backends in the server pool
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"

		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		if !alive {
			status = "down"
		}

		log.Printf("URL -> %s is %s\n", b.URL, status)
	}
}

// ? A simple utility function to see if a backend is alive or not
func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second

	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		return false
	}

	_ = conn.Close()
	return true
}

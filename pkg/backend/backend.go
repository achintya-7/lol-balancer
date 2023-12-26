package backend

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// Function to create a new backend instance
func NewBackendServer(url *url.URL) *Backend {
	proxy := newProxyServer(url)

	return &Backend{
		URL:          url,
		Alive:        true,
		ReverseProxy: proxy,
	}
}

func newProxyServer(url *url.URL) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(url)
}

// Function to set Alive bool to the backend instance
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// Function to check the current value of Alive bool
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	alive := b.Alive
	b.mux.RUnlock()
	return alive
}

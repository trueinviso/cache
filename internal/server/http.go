package server

import (
	"encoding/json"
	"lru/internal/lru"
	"net/http"
	"strconv"
)

type Cache struct {
	lru *lru.LRU
}

type WriteRequest struct {
	Key  int64  `json:"key"`
	Data []byte `json:"byte"`
}

type ReadResponse struct {
	Data []byte `json:"data"`
}

func NewHTTPServer(address string) *http.Server {
	cache := NewCache()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", cache.Write)
	mux.HandleFunc("GET /", cache.Read)

	return &http.Server{
		Addr:    address,
		Handler: mux,
	}
}

func NewCache() *Cache {
	return &Cache{
		lru: lru.NewLRU(1024),
	}
}

func (c *Cache) Write(w http.ResponseWriter, r *http.Request) {
	var wReq WriteRequest
	err := json.NewDecoder(r.Body).Decode(&wReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.lru.Add(wReq.Key, wReq.Data)
	w.WriteHeader(http.StatusOK)
}

func (c *Cache) Read(w http.ResponseWriter, r *http.Request) {
	keyStr := r.URL.Query().Get("key")
	if keyStr == "" {
		http.Error(w, "Missing 'key' query parameter", http.StatusBadRequest)
		return
	}

	key, err := strconv.ParseInt(keyStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	node := c.lru.Read(key)
	if node == nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(ReadResponse{Data: node.Data.([]byte)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

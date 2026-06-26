package storage

import "sync"

// Engine represent our thread-safe, in-memory key-value database
type Engine struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewEngine is the constructor that init the map
func NewEngine() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

// Set writes a KV pair to the database
func (e *Engine) Set(key, value string) {
	e.mu.Lock()         // write lock
	defer e.mu.Unlock() // guarantee lock release

	e.data[key] = value
}

// Get retrieves a value from the database safely
// returns a boolean 'exists' so we can differentiate between an empty string value and a missing key
func (e *Engine) Get(key string) (string, bool) {
	e.mu.RLock() // shared read lock
	defer e.mu.RUnlock()

	value, exists := e.data[key]
	return value, exists
}

// Delete removes a key from the database
func (e *Engine) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.data, key)
}

package core

import (
	"net/http"
	"sort"
)

var (
	handlers map[int]func(http.Handler) http.Handler
)

// Init creates the handler map
func init() {
	handlers = make(map[int]func(http.Handler) http.Handler, 0)
}

// GetHandler returns the current handler map
func GetHandler() map[int]func(http.Handler) http.Handler {
	return handlers
}

// AddHandler adds an http Handler to the map, requires a priority set
func AddHandler(priority int, handle func(http.Handler) http.Handler) {
	handlers[priority] = handle
}

// DeleteHandler removes an httpHandler based on the priority set
func DeleteHandler(priority int) {
	delete(handlers, priority)
}

// handler returns the final http handler from the map
// Builds the final handler based on the priorities
func handler(a http.Handler) http.Handler {
	out := a
	keys := make([]int, 0)
	for i := range handlers {
		keys = append(keys, i)
	}
	sort.Ints(keys)
	for _, i := range keys {
		out = handlers[i](out)
	}
	return out
}

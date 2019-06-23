package core

import (
	"net/http"
)

type Handler struct {
	handlers []func(http.Handler) http.Handler
}

func newHandler() *Handler {
	h := new(Handler)
	h.handlers = make([]func(http.Handler) http.Handler, 0)
	return h
}

func (h *Handler) Init() {
	h.handlers = make([]func(http.Handler) http.Handler, 0)
}

func (h Handler) Get() []func(http.Handler) http.Handler {
	return h.handlers
}

func (h *Handler) AddHandler(handle func(http.Handler) http.Handler) {
	h.handlers = append(h.handlers, handle)
}

func (h *Handler) AppendHandlers(hands []func(http.Handler) http.Handler) {
	for i := range hands {
		h.handlers = append(h.handlers, hands[i])
	}
}

func (c *Handler) Final(a http.Handler) http.Handler {
	handle := a
	for i := range c.handlers {
		handle = c.handlers[len(c.handlers)-1-i](handle)
	}
	return handle
}

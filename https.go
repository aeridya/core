package core

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func autoSSL() *autocert.Manager {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(AutoSSLDirectory),
	}

	go http.ListenAndServe(Port, certManager.HTTPHandler(nil))

	return &certManager
}

func hostPolicy(ctx context.Context, host string) error {
	// Note: change to your real domain
	allowedHost := Domain
	if host == allowedHost {
		return nil
	}
	return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
}

func http2https(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Use HTTPS", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, (FullDomain + r.URL.RequestURI()), http.StatusFound)
}

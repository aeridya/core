package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/securecookie"

	"github.com/hlfstr/configurit"
	"github.com/hlfstr/quit"
	"github.com/hlfstr/quit/quitters"

	"github.com/aeridya/logit"
)

var (
	server *http.Server

	Handle *Handler
	Config *configurit.Conf

	Port        string
	Domain      string
	FullDomain  string
	Development bool
	HTTPS       bool

	Serve func(*Response)
	Error func(*Response)

	cookieHash  []byte
	cookieBlock []byte

	isInit  bool
	limiter chan struct{}
)

func Create(conf string) error {
	var err error
	Config, err = loadConfig(conf)
	if err != nil {
		return err
	}
	quitters.AddQuit(quit.Quit)
	quitters.AddQuit(logit.Quit)
	Handle = newHandler()
	cookieHandler = securecookie.New(cookieHash, cookieBlock)
	if HTTPS {
		FullDomain = "https://" + Domain
	} else {
		FullDomain = "http://" + Domain
	}
	server = &http.Server{Addr: Port}
	isInit = true
	return nil
}

func loadConfig(conf string) (*configurit.Conf, error) {
	c, err := configurit.Open(conf)
	if err != nil {
		return nil, err
	}

	if Domain, err = c.GetString("", "Domain"); err != nil {
		return nil, err
	}

	if s, err := c.GetString("", "Port"); err != nil {
		return nil, err
	} else {
		Port = ":" + s
	}

	if n, err := c.GetInt("", "Workers"); err != nil {
		return nil, err
	} else {
		limiter = make(chan struct{}, n)
	}

	if h, err := c.GetString("", "CookieHash"); err != nil {
		return nil, err
	} else {
		cookieHash = []byte(h)
	}

	if h, err := c.GetString("", "CookieBlock"); err != nil {
		return nil, err
	} else {
		cookieBlock = []byte(h)
	}

	if log, err := c.GetString("", "Log"); err != nil {
		return nil, err
	} else {
		if log == "stdout" {
			if err = logit.Start(logit.TermLog()); err != nil {
				return nil, err
			}
		} else {
			if file, err := logit.OpenFile(log); err != nil {
				return nil, err
			} else {
				if err = logit.Start(file); err != nil {
					return nil, err
				}
			}
		}
	}

	if Development, err = c.GetBool("", "Development"); err != nil {
		return nil, err
	}

	if HTTPS, err = c.GetBool("", "HTTPS"); err != nil {
		return nil, err
	}

	return c, err
}

func Run() error {
	if !isInit {
		return fmt.Errorf("Aeridya[Error]: Must use Create(\"/path/to/config\") before Run()")
	}
	Config = nil //Delete the Config reference, FREE THE MEMORY!!
	go quit.Run(shutdown, shutdown)
	defer panicAttack()
	defer quitters.Quit()
	logit.Logf(logit.MSG, "Starting %s for %s | Listening on %s", Version(), Domain, Port)
	http.Handle("/", Handle.Final(http.HandlerFunc(serve)))
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		logit.LogError(1, err)
		return nil
	}
	logit.Logf(0, "Shutting Down Aeridya for %s", Domain)
	return nil
}

func serve(w http.ResponseWriter, r *http.Request) {
	resp := &Response{W: w, R: r}
	Serve(resp)
	if resp.Err != nil {
		Error(resp)
		if Development {
			logit.Logf(logit.ERROR, "[Error(%d)] %s", resp.Status, resp.Err.Error())
		}
	}
	return
}

func limit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter <- struct{}{}
		h.ServeHTTP(w, r)
		<-limiter
	})
}

/*
func AddTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.Path
		if u[len(u)-1:] != "/" {
			u = u + "/"
			http.Redirect(w, r, u, 301)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func NoTrailingSlash(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL.Path
		if u[len(u)-1:] == "/" {
			u = u[:len(u)-1]
			http.Redirect(w, r, u, 301)
			return
		}

		h.ServeHTTP(w, r)
	})
}
*/

func shutdown() {
	if err := server.Shutdown(context.Background()); err != nil {
		logit.LogError(1, err)
		server.Close()
	}
}

func panicAttack() {
	err := recover()
	if err != nil {
		logit.Logf(logit.PANIC, "PANIC!\n  %#v\n", err)
		buf := make([]byte, 4096)
		buf = buf[:runtime.Stack(buf, true)]
		logit.Logf(logit.PANIC, "Stack Trace:\n%s\n", buf)
		os.Exit(1)
	}
}

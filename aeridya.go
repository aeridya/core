package core

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/aeridya/core/configurit"
	"github.com/aeridya/core/logit"
	"github.com/aeridya/core/quit"
)

var (
	Port        string
	Domain      string
	FullDomain  string
	Development bool
	HTTPS       bool

	Serve func(*Response)
	Error func(*Response)

	isInit  bool
	limiter chan struct{}

	server *http.Server
)

//Create loads the configuration file and sets Aeridya up
func Create(conf string) error {
	var err error
	err = loadConfig(conf)
	if err != nil {
		return err
	}

	quit.AddQuit(shutdown)
	quit.AddQuit(logit.Quit)

	if Development {
		logit.Log(logit.NOTICE, "Development mode enabled")
	}

	if HTTPS {
		FullDomain = "https://" + Domain
	} else {
		FullDomain = "http://" + Domain
	}

	if Development {
		logit.Log(logit.DEBUG, "Working Domain: ", FullDomain)
	}

	if limiter == nil {
		logit.Log(logit.NOTICE, "Limiter disabled due to configuration")
	} else {
		AddHandler(1000, limit)
	}

	server = &http.Server{Addr: Port}

	isInit = true
	return nil
}

// Load the Configuration options from the Configuration
func loadConfig(conf string) error {
	// Open configuration file
	err := configurit.Open(conf)
	if err != nil {
		return err
	}

	// Load Domain name
	if Domain, err = configurit.Config.GetString("Aeridya", "Domain"); err != nil {
		return err
	}

	// Load Port
	if s, err := configurit.Config.GetString("Aeridya", "Port"); err != nil {
		return err
	} else {
		Port = ":" + s
	}

	// Load Workers amount
	if n, err := configurit.Config.GetInt("Aeridya", "Workers"); err != nil {
		return err
	} else {
		if n > 0 {
			limiter = make(chan struct{}, n)
		}
	}

	// Load log location
	if log, err := configurit.Config.GetString("Aeridya", "Log"); err != nil {
		return err
	} else {
		if log == "stdout" {
			if err = logit.Start(logit.TermLog()); err != nil {
				return err
			}
		} else {
			if file, err := logit.OpenFile(log); err != nil {
				return err
			} else {
				if err = logit.Start(file); err != nil {
					return err
				}
			}
		}
	}

	// Load Development setting
	if Development, err = configurit.Config.GetBool("Aeridya", "Development"); err != nil {
		return err
	}

	// Load HTTPS setting
	if HTTPS, err = configurit.Config.GetBool("Aeridya", "HTTPS"); err != nil {
		return err
	}

	return err
}

//Run starts the server, begins Aeridya
func Run() error {
	if !isInit {
		return fmt.Errorf("Aeridya[Error]: Must use Create(\"/path/to/config\") before Run()")
	}

	//Start listening on the quits
	quit.Run(quit.Quit, shutdown)
	//Defer catching panics
	defer panicAttack()
	//Defer running the quits on exit
	defer quit.Exit()

	logit.Logf(logit.MSG, "Starting %s for %s | Listening on %s", Version(), Domain, Port)

	//Set the handler for Aeridya.  Using "/" does all URLs as well
	http.Handle("/", handler(http.HandlerFunc(serve)))
	if Development {
		logit.Log(logit.DEBUG, "Aeridya server starting...")
	}
	//Run the server
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		logit.LogError(1, err)
		return err
	}
	logit.Log(logit.DEBUG, "Aeridya server shutting down...")
	return nil
}

//serve is the function used on every connection to run the theme
func serve(w http.ResponseWriter, r *http.Request) {
	resp := mkResponse(w, r)
	Serve(resp)
	if resp.Err != nil {
		Error(resp)
		if Development {
			logit.Logf(logit.DEBUG, "[Error(%d)] %s", resp.Status, resp.Err.Error())
		}
	}
	return
}

//limit is the internal handler function to limit the amount of connections
func limit(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter <- struct{}{}
		h.ServeHTTP(w, r)
		<-limiter
	})
}

//shutdown closes the server
func shutdown() {
	if err := server.Shutdown(context.Background()); err != nil {
		logit.LogError(1, err)
		server.Close()
	}
	logit.Logf(0, "Shutting Down Aeridya for %s", Domain)
}

//panicAttack attempts to catch a panic made in Aeridya
//will log the panic to via logit
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

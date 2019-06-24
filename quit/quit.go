package quit

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	//channel to recieve os quit on
	sigchan chan os.Signal
	//channel to receive quit
	done chan bool
	// Array to store quit functions
	quits []func()
	// isDone maintains state of quit function
	// Prevents running after channels closed
	isDone bool
)

func init() {
	quits = make([]func(), 0)
	sigchan = make(chan os.Signal, 1)
	done = make(chan bool, 1)
}

// AddQuit will add the function passed to the
// quits array
func AddQuit(f func()) {
	quits = append(quits, f)
}

// Quit runs all quit functions
func Quit() {
	for i := range quits {
		quits[i]()
	}
}

// Run begins to listen to the signals and runs the functions
// passed to it.  These should be the functions that shuts down your
// application, the other restarts it. Starts the run() goroutine
func Run(quit func(), restart func()) {
	go run(quit, restart)
}

// run listens on a goroutine for quit and if to quit
// Listens to all POSIX signals that would result in a
// termination and for SIGUSR2 to restart.
// On Windows, only INTERRUPT and KILL are available.
func run(quit func(), restart func()) {
	defer close(sigchan)
	defer close(done)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGIOT, syscall.SIGTERM, syscall.SIGUSR2)
loop:
	for {
		select {
		case s := <-sigchan:
			switch s {
			//Quit
			case os.Signal(syscall.SIGHUP): //hangup
				fallthrough
			case os.Signal(syscall.SIGABRT), os.Signal(syscall.SIGIOT): //abort
				fallthrough
			case os.Signal(syscall.SIGTERM): //terminate
				fallthrough
			case os.Interrupt: //interrupt
				fallthrough
			case os.Kill: //kill
				quit()
				isDone = true
				break loop
			//Restart
			case os.Signal(syscall.SIGUSR2):
				restart()
			}
		//Done
		case <-done:
			quit()
			isDone = true
			break loop
		}
	}
}

// Quit sends the signal to end the goroutine
// The quit() function will be run if the state == true
// blocks to completion
func Exit() {
	if !isDone {
		done <- true
		<-done
	}
}

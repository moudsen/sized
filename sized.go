package main

import (
	"fmt"
	"log"
	"net/http"
        "path/filepath"
	"os"
	"os/signal"
	"syscall"
        "time"

	"github.com/takama/daemon"
)

const (
	name        = "sized"
	description = "sized: Retrieve directory size"
)

var stdlog, errlog *log.Logger

type Service struct {
	daemon.Daemon
}

func handleIndex(w http.ResponseWriter, req *http.Request) {
	fmt.Println("use /size?dir=<dirname>")
}

func DirSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return err
    })
    return size, err
}

func DirCount(path string) (int64, error) {
    var count int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            count += 1
        }
        return err
    })
    return count, err
}

func handleSizeRequest(w http.ResponseWriter, req *http.Request) {
        var dirsize int64
        var start time.Time

	// Obtain the ip address from the requestor. As this routine likely sits behind a reverse proxy,
	// first test for an ip in the header. Only if not there, use the RemoteAddr method.

	ip := req.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = req.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = req.RemoteAddr
	}

        // Fetch the request directory name and process

	dirs, ok := req.URL.Query()["dir"]

        dirsize = -1

        start = time.Now()

        if ok {
            dirsize,_ = DirSize(dirs[0])
        }

        used := time.Since(start)

        fmt.Fprint(w,dirsize)

        stdlog.Printf("Directory %s has %d bytes, %f secs", dirs[0], dirsize, used.Seconds())
}

func handleCountRequest(w http.ResponseWriter, req *http.Request) {
        var dircount int64
        var start time.Time

	// Obtain the ip address from the requestor. As this routine likely sits behind a reverse proxy,
	// first test for an ip in the header. Only if not there, use the RemoteAddr method.

	ip := req.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = req.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = req.RemoteAddr
	}

        // Fetch the request directory name and process

	dirs, ok := req.URL.Query()["dir"]

        dircount = -1

        start = time.Now()

        if ok {
            dircount,_ = DirCount(dirs[0])
        }

        used := time.Since(start)

        fmt.Fprint(w,dircount)

        stdlog.Printf("Directory %s has %d files, %f secs", dirs[0], dircount, used.Seconds())
}

func (service *Service) Manage() (string, error) {
	// The deamon control section has been copied from an example of the Takama/daemon library. It's
	// straightforward and install/using/removing the daemon is a breeze ...

	usage := "Usage: sized install | remove | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/size", handleSizeRequest)
	http.HandleFunc("/count", handleCountRequest)

	// Start the daemon in a child process. We can handle multiple requests in parallel from here.

	go func() {
		http.ListenAndServe(":7007", nil)
	}()

	// Log that our service is ready and listening on port 7007.

	stdlog.Println("Service started, listening on port 7007")

	// Wait for signals in an infinite loop. Note that we only accept a kill signal; no other signals are caught.

	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

func init() {
	// As we are a daemon we need to divert our standard and error output.

	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

func main() {
	// Create a new daemon process.

	srv, err := daemon.New(name, description, daemon.SystemDaemon)

	// If we failed, report the error and halt.

	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}

	// Otherwise initialize and run the daemon service.

	service := &Service{srv}
	status, err := service.Manage()

	// If initialization failed (short of memory for examples) report the error and halt.

	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}

	// Report our daemon status.
	fmt.Println(status)
}

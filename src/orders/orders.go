package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/url"
	"orders/db"
	"orders/server"
	"orders/tmpl"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Start an 'orders' server using a random TCP port.
// The specified 'orders' file is loaded if the file exists at startup,
// otherwise the server starts with an empty DB and the file is created at
// the first operation.
//
// ./orders [-d ./data.orders]
func main() {
	// Read CLI flags
	var data string
	flag.StringVar(&data, "d", "data.orders", "App data file")
	data, _ = filepath.Abs(data)

	// Load the DB if the file exist:
	d, err := db.Load(data)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Fail to load the DB at path %q: %v", data, err)
		os.Exit(2)
	}

	// Init template assets:
	t, err := tmpl.LoadFromAssets("templates")
	if err != nil {
		log.Printf("Fail to load the templates: %v", err)
		os.Exit(2)
	}

	// Randomly open a TCP port for the server:
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Println("Fail to listen a tcp port:", err)
		os.Exit(2)
	}

	// Server instance:
	srv := &server.OrderServer{
		DB:        &d,
		Templater: t,
	}

	// A mutex to synchronize the browser cmd and the server startup:
	var waitForServer sync.Mutex
	waitForServer.Lock()

	// Start the server in a separated routine: (async start)
	go func() {
		s := http.Server{
			Handler: srv,
		}
		err = s.Serve(l)
		waitForServer.Unlock()
		if err != nil {
			log.Println("Fail to listen incoming requests:", err)
			os.Exit(2)
		}
	}()

	// Open the default system browser with application root URL:
	u, err := url.Parse(l.Addr().String())
	if err != nil {
		log.Println("Unknown addr value:", err)
		os.Exit(2)
	}
	u.Scheme = "http"
	log.Println("Start server at", u.String())
	exec.Command("cmd", "/C", "start", u.String()).Start()

	waitForServer.Lock()
}

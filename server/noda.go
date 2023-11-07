package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func Run() {
	connectToDatabase()
	db := getDatabase()
	defer db.Close()
	r := startRouter()
	config := getServerConfig()
	server := http.Server{
		WriteTimeout:      config.WriteTimeout,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		IdleTimeout:       config.IdleTimeout,
		Handler:           r,
		ErrorLog:          log.New(os.Stderr, "\033[0;31mfatal: \033[0m", log.LstdFlags),
		ConnState:         tcpConnStatLogger,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Serving HTTP at \033[0;32m`%s'\033[0m ...\n", listener.Addr())
	log.Fatal(server.Serve(listener))
}

func tcpConnStatLogger(c net.Conn, cs http.ConnState) {
	var color string
	switch cs {
	case http.StateNew:
		color = "\033[0;32m" /* Green.  */
	case http.StateActive:
		color = "\033[0;34m" /* Blue.  */
	case http.StateIdle:
		color = "\033[0;33m" /* Yellow.  */
	case http.StateHijacked:
		color = "\033[0;35m" /* Purple.  */
	case http.StateClosed:
		color = "\033[0;31m" /* Red.  */
	}
	fmt.Printf("\033[0;32mTCP connection:\033[0m (me) <--> (%s): %s%s\033[0m\n", c.RemoteAddr(), color, cs)
}

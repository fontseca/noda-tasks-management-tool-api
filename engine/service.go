package engine

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"noda"
	noda_middleware "noda/engine/internal/middleware"
	"noda/engine/internal/routes"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() {
	noda.ConnectToDatabase()
	db := noda.GetDatabase()
	defer db.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Headers", "*"))
	router.Use(middleware.SetHeader("Access-Control-Allow-Credentials", "true"))
	router.Use(noda_middleware.LetOptionsPassThrough)
	router.NotFound(noda_middleware.NotFound)

	routes.InitializeForAuthentication(router)
	routes.InitializeForUsers(router)
	routes.InitializeForTasks(router)
	routes.InitializeForGroups(router)

	config := noda.GetServerConfig()
	server := http.Server{
		WriteTimeout:      config.WriteTimeout,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		IdleTimeout:       config.IdleTimeout,
		Handler:           router,
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

package engine

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"noda/api/handler"
	"noda/api/repository"
	"noda/api/service"
	"noda/config"
	"noda/database"
	"os"

	"github.com/go-chi/chi/v5"
)

func Run() {
	db := database.ConnectAndGet()
	defer db.Close()
	hdls := getServiceHandlers(db)
	r := initializeRouter(hdls)
	startService(r)
}

func getServiceHandlers(db *sql.DB) *serviceHandlers {
	return &serviceHandlers{
		TaskHandler: getTaskHandler(db),
	}
}

func getTaskHandler(db *sql.DB) *handler.TaskHandler {
	taskRepository := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepository)
	return handler.NewTaskHandler(taskService)
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
	log.Printf("\033[0;32mTCP connection:\033[0m (me) <--> (%s): %s%s\033[0m\n", c.RemoteAddr(), color, cs)
}

func startService(r *chi.Mux) {
	c := config.GetServerConfig()
	s := http.Server{
		WriteTimeout:      c.WriteTimeout,
		ReadTimeout:       c.ReadTimeout,
		ReadHeaderTimeout: c.ReadHeaderTimeout,
		IdleTimeout:       c.IdleTimeout,
		Handler:           r,
		ErrorLog:          log.New(os.Stderr, "\033[0;31mfatal: \033[0m", log.LstdFlags),
		ConnState:         tcpConnStatLogger,
	}
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Printf("Serving HTTP at \033[0;32m`%s'\033[0m ...", l.Addr())
	log.Fatal(s.Serve(l))
}

package server

import (
	"context"
	"ex-server/internal/adaptor"
	"ex-server/internal/entity"
	"ex-server/internal/handler"
	"ex-server/pkg/config"
	"ex-server/pkg/db"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

const (
	waitTimeout  = time.Second * 15
	writeTimeout = time.Second * 15
	readTimeout  = time.Second * 15
	idleTimeout  = time.Second * 60
)

type Server struct {
	config  config.Config
	Handler *handler.Handler
}

func Init(configPath string) (*Server, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	taskRepo := adaptor.NewTaskRepository(db)

	handler := handler.Init(*taskRepo)

	return &Server{config: cfg, Handler: handler}, nil
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.config.GetString("App.Host"), s.config.GetString("App.Port")),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      s.initRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
}

func (s *Server) initRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/task/list", s.Handler.GetTasksList).Methods("GET")
	r.HandleFunc("/task/{id}", s.Handler.GetTask).Methods("GET")
	r.HandleFunc("/task", s.Handler.CreateTask).Methods("POST")
	r.HandleFunc("/task/{id}", s.Handler.UpdateTask).Methods("PUT")
	r.HandleFunc("/task/{id}", s.Handler.DeleteTask).Methods("DELETE")

	return r
}

func initDB(cfg config.Config) (*gorm.DB, error) {

	return db.NewConnection(cfg,
		&entity.Task{},
	)
}

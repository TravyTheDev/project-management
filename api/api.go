package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"project-management/service/i18n"
	"project-management/service/jwt"
	"project-management/service/mailer"
	"project-management/service/projects"
	"project-management/service/users"
	"project-management/service/ws"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	secretKey := os.Getenv("JWT_KEY")
	frontUrl := os.Getenv("FRONT_URL")
	router := mux.NewRouter()

	subRouter := router.PathPrefix("/api/v1").Subrouter()
	//if I make a website client
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontUrl},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
	})
	langMap := i18n.LoadLocaleFiles()
	handler := c.Handler(router)
	userStore := users.NewUserStore(s.db)
	sessionStore := jwt.NewSessionStore(s.db)
	userHandler := users.NewHandler(userStore, sessionStore, secretKey, langMap)
	userHandler.RegisterRoutes(subRouter)
	mailHandler := mailer.NewMailHandler(userStore, secretKey)
	mailHandler.RegisterRoutes(subRouter)
	projectsStore := projects.NewProjectsStore(s.db)
	projectsHandler := projects.NewProjectsHandler(projectsStore)
	projectsHandler.RegisterRoutes(subRouter)

	hub := ws.NewHub()
	wsHandler := ws.NewWSHandler(hub, secretKey)
	wsHandler.RegisterRoutes(subRouter)
	go hub.Run()

	inviteBroker := ws.NewBroker(secretKey)
	inviteBroker.RegisterRoutes(subRouter)

	fmt.Printf("Listening on : %s\n", s.addr)
	return http.ListenAndServe(s.addr, handler)
}

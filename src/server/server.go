package server

import (
	"auth-service/src/api"
	"auth-service/src/config"
	"auth-service/src/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	cfg     *config.Config
	service service.UserService
}

func NewServer(cfg *config.Config, userService service.UserService) *Server {
	return &Server{
		cfg:     cfg,
		service: userService,
	}
}

func (s *Server) Run() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	apiHandler := api.NewHandler(s.service, s.cfg)

	// --- Configuração das Rotas ---
	// Rotas Públicas
	router.Post("/register", apiHandler.HandleRegister)
	router.Post("/login", apiHandler.HandleLogin)

	// Rotas Protegidas
	router.Group(func(r chi.Router) {
		r.Use(apiHandler.APIKeyAuthMiddleware)
		r.Post("/auth/validate", apiHandler.HandleAuthValidate)
	})
	router.Group(func(r chi.Router) {
		r.Use(apiHandler.JWTAuthMiddleware)
		r.Get("/profile", apiHandler.HandleGetProfile)
	})

	log.Printf("Servidor de Autenticação iniciado em %s", s.cfg.ListenAddr)
	if err := http.ListenAndServe(s.cfg.ListenAddr, router); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}

// Comando api: lê a configuração e monta o sistema.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"receitas-api/internal/config"
	"receitas-api/internal/handler"
	"receitas-api/internal/middleware"
	"receitas-api/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("erro ao obter conexão: %v", err)
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	handler.NewReceitaHandler(repository.NewReceitaDAO(db)).Registrar(mux)

	// Seguro por padrão: tudo exige chave, exceto o que estiver listado aqui.
	publicas := map[string]bool{"GET /health": true}
	auth := middleware.Auth(repository.NewAPIKeyDAO(db), publicas)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           auth(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("servidor ouvindo em :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

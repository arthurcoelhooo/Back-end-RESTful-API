// Comando genkey: gera uma nova API key, guarda só o hash no banco
// e mostra a chave UMA única vez.
//
// Uso: go run ./cmd/genkey -nome "grupo-dev"
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"receitas-api/internal/config"
	"receitas-api/internal/middleware"
	"receitas-api/internal/repository"
)

func main() {
	nome := flag.String("nome", "", "nome/descrição da chave (obrigatório)")
	flag.Parse()
	if *nome == "" {
		log.Fatal("informe o nome da chave: -nome \"minha-chave\"")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if sqlDB, err := db.DB(); err == nil {
		defer sqlDB.Close()
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	chave := hex.EncodeToString(bytes)

	repo := repository.NewAPIKeyDAO(db)
	if err := repo.Criar(context.Background(), *nome, middleware.HashChave(chave)); err != nil {
		log.Fatalf("erro ao salvar chave: %v", err)
	}

	fmt.Println("Chave criada. Guarde-a agora, ela não será exibida novamente:")
	fmt.Println(chave)
}

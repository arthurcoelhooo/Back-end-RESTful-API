// Package middleware contém os middlewares HTTP.
package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"receitas-api/internal/repository"
)

// HashChave devolve o SHA-256 (hexadecimal) da chave. É o que vai para o banco.
func HashChave(chave string) string {
	soma := sha256.Sum256([]byte(chave))
	return hex.EncodeToString(soma[:])
}

// Auth protege TODAS as rotas, exceto as declaradas em publicas
// (no formato "MÉTODO /caminho", ex.: "GET /health"). Seguro por padrão.
func Auth(repo repository.APIKeyRepository, publicas map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if publicas[r.Method+" "+r.URL.Path] {
				next.ServeHTTP(w, r)
				return
			}

			chave := r.Header.Get("X-API-Key")
			if chave == "" {
				escreverErro(w, http.StatusUnauthorized, "chave de API ausente ou inválida")
				return
			}

			existe, err := repo.ExisteHash(r.Context(), HashChave(chave))
			if err != nil {
				log.Printf("erro ao validar chave: %v", err)
				escreverErro(w, http.StatusInternalServerError, "erro interno do servidor")
				return
			}
			if !existe {
				escreverErro(w, http.StatusUnauthorized, "chave de API ausente ou inválida")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func escreverErro(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"erro": msg})
}

// Package handler contém os handlers HTTP
package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func escreverJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func erroCliente(w http.ResponseWriter, status int, msg string, detalhes ...string) {
	corpo := map[string]any{"erro": msg}
	if len(detalhes) > 0 {
		corpo["detalhes"] = detalhes
	}
	escreverJSON(w, status, corpo)
}

// erroInterno registra o erro real no log e devolve uma mensagem genérica ao cliente.
func erroInterno(w http.ResponseWriter, err error) {
	log.Printf("erro interno: %v", err)
	erroCliente(w, http.StatusInternalServerError, "erro interno do servidor")
}

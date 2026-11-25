// Package model contém as entidades e suas regras de validação.
package model

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

var CategoriasValidas = map[string]bool{
	"doce": true, "salgado": true, "bebida": true, "massa": true,
	"sobremesa": true, "lanche": true, "salada": true,
}

type Receita struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	Titulo          string    `json:"titulo" gorm:"column:titulo"`
	Categoria       string    `json:"categoria" gorm:"column:categoria"`
	Ingredientes    string    `json:"ingredientes" gorm:"column:ingredientes"`
	ModoPreparo     string    `json:"modo_preparo" gorm:"column:modo_preparo"`
	TempoPreparoMin int       `json:"tempo_preparo_min" gorm:"column:tempo_preparo_min"`
	Porcoes         int       `json:"porcoes" gorm:"column:porcoes"`
	CriadaEm        time.Time `json:"criada_em" gorm:"column:criada_em;autoCreateTime"`
}

// TableName diz ao GORM o nome real da tabela (a estrutura vem do schema.sql).
func (Receita) TableName() string { return "receitas" }

// APIKey representa uma linha da tabela api_keys (só o hash, nunca a chave).
type APIKey struct {
	ID       int       `gorm:"primaryKey"`
	Nome     string    `gorm:"column:nome"`
	KeyHash  string    `gorm:"column:key_hash"`
	CriadaEm time.Time `gorm:"column:criada_em;autoCreateTime"`
}

func (APIKey) TableName() string { return "api_keys" }

// Validar normaliza os campos de texto e devolve a lista de problemas
// encontrados. Lista vazia significa que a receita é válida.
func (r *Receita) Validar() []string {
	r.Titulo = strings.TrimSpace(r.Titulo)
	r.Categoria = strings.ToLower(strings.TrimSpace(r.Categoria))
	r.Ingredientes = strings.TrimSpace(r.Ingredientes)
	r.ModoPreparo = strings.TrimSpace(r.ModoPreparo)

	var erros []string

	if n := utf8.RuneCountInString(r.Titulo); n < 3 || n > 120 {
		erros = append(erros, "titulo é obrigatório e deve ter entre 3 e 120 caracteres")
	}
	if !CategoriasValidas[r.Categoria] {
		erros = append(erros, fmt.Sprintf("categoria inválida; use: %s", listaCategorias()))
	}
	if r.Ingredientes == "" {
		erros = append(erros, "ingredientes é obrigatório")
	}
	if r.ModoPreparo == "" {
		erros = append(erros, "modo_preparo é obrigatório")
	}
	if r.TempoPreparoMin <= 0 {
		erros = append(erros, "tempo_preparo_min deve ser maior que 0")
	}
	if r.Porcoes <= 0 {
		erros = append(erros, "porcoes deve ser maior que 0")
	}
	return erros
}

func listaCategorias() string {
	ordem := []string{"doce", "salgado", "bebida", "massa", "sobremesa", "lanche", "salada"}
	return strings.Join(ordem, ", ")
}

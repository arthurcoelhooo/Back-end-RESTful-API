package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"receitas-api/internal/model"
	"receitas-api/internal/repository"
)

type ReceitaHandler struct {
	repo repository.ReceitaRepository // depende só da interface
}

func NewReceitaHandler(repo repository.ReceitaRepository) *ReceitaHandler {
	return &ReceitaHandler{repo: repo}
}

func (h *ReceitaHandler) Registrar(mux *http.ServeMux) {
	mux.HandleFunc("GET /receitas", h.listar)
	mux.HandleFunc("POST /receitas", h.criar)
	mux.HandleFunc("GET /receitas/{id}", h.buscar)
	mux.HandleFunc("PUT /receitas/{id}", h.atualizar)
	mux.HandleFunc("DELETE /receitas/{id}", h.excluir)
}

func (h *ReceitaHandler) listar(w http.ResponseWriter, r *http.Request) {
	receitas, err := h.repo.Listar(r.Context(), r.URL.Query().Get("categoria"))
	if err != nil {
		erroInterno(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, receitas)
}

func (h *ReceitaHandler) criar(w http.ResponseWriter, r *http.Request) {
	rec, ok := lerReceita(w, r)
	if !ok {
		return
	}
	if err := h.repo.Criar(r.Context(), rec); err != nil {
		erroInterno(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/receitas/%d", rec.ID))
	escreverJSON(w, http.StatusCreated, rec)
}

func (h *ReceitaHandler) buscar(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}
	rec, err := h.repo.BuscarPorID(r.Context(), id)
	if errors.Is(err, repository.ErrNaoEncontrado) {
		erroCliente(w, http.StatusNotFound, "receita não encontrada")
		return
	}
	if err != nil {
		erroInterno(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, rec)
}

func (h *ReceitaHandler) atualizar(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}
	rec, ok := lerReceita(w, r)
	if !ok {
		return
	}
	rec.ID = id // o ID vem da URL, nunca do corpo

	err := h.repo.Atualizar(r.Context(), rec)
	if errors.Is(err, repository.ErrNaoEncontrado) {
		erroCliente(w, http.StatusNotFound, "receita não encontrada")
		return
	}
	if err != nil {
		erroInterno(w, err)
		return
	}
	escreverJSON(w, http.StatusOK, rec)
}

func (h *ReceitaHandler) excluir(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}
	err := h.repo.Excluir(r.Context(), id)
	if errors.Is(err, repository.ErrNaoEncontrado) {
		erroCliente(w, http.StatusNotFound, "receita não encontrada")
		return
	}
	if err != nil {
		erroInterno(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// lerID extrai e valida o {id} do caminho. Responde 400 se for inválido.
func lerID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		erroCliente(w, http.StatusBadRequest, "id inválido")
		return 0, false
	}
	return id, true
}

// lerReceita decodifica o JSON do corpo e valida. Responde 400 se falhar.
func lerReceita(w http.ResponseWriter, r *http.Request) (*model.Receita, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var rec model.Receita
	if err := dec.Decode(&rec); err != nil {
		erroCliente(w, http.StatusBadRequest, "JSON inválido ou com campos desconhecidos")
		return nil, false
	}
	if erros := rec.Validar(); len(erros) > 0 {
		erroCliente(w, http.StatusBadRequest, "dados inválidos", erros...)
		return nil, false
	}
	return &rec, true
}

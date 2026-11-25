// Package repository define as interfaces de acesso a dados e os DAOs que as implementam.
package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"receitas-api/internal/model"
)

var ErrNaoEncontrado = errors.New("registro não encontrado")

// ReceitaRepository é a interface da qual os handlers dependem.
type ReceitaRepository interface {
	Listar(ctx context.Context, categoria string) ([]model.Receita, error)
	BuscarPorID(ctx context.Context, id int) (*model.Receita, error)
	Criar(ctx context.Context, r *model.Receita) error
	Atualizar(ctx context.Context, r *model.Receita) error
	Excluir(ctx context.Context, id int) error
}

// ReceitaDAO implementa ReceitaRepository com GORM.
type ReceitaDAO struct {
	db *gorm.DB
}

func NewReceitaDAO(db *gorm.DB) *ReceitaDAO {
	return &ReceitaDAO{db: db}
}

func (d *ReceitaDAO) Listar(ctx context.Context, categoria string) ([]model.Receita, error) {
	q := d.db.WithContext(ctx).Order("id")
	if categoria != "" {
		// Sempre com "?" e o valor como argumento: nunca concatenar.
		q = q.Where("categoria = ?", categoria)
	}

	receitas := []model.Receita{} // vazio (e não nil) para o JSON sair como []
	if err := q.Find(&receitas).Error; err != nil {
		return nil, err
	}
	return receitas, nil
}

func (d *ReceitaDAO) BuscarPorID(ctx context.Context, id int) (*model.Receita, error) {
	var r model.Receita
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNaoEncontrado
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (d *ReceitaDAO) Criar(ctx context.Context, r *model.Receita) error {
	// O GORM preenche r.ID (via RETURNING) e r.CriadaEm.
	return d.db.WithContext(ctx).Create(r).Error
}

func (d *ReceitaDAO) Atualizar(ctx context.Context, r *model.Receita) error {
	// Select lista os campos a gravar: assim o GORM não ignora valores zero
	// e não mexe em criada_em.
	res := d.db.WithContext(ctx).
		Model(&model.Receita{}).
		Where("id = ?", r.ID).
		Select("Titulo", "Categoria", "Ingredientes", "ModoPreparo", "TempoPreparoMin", "Porcoes").
		Updates(r)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNaoEncontrado
	}

	// Recarrega para devolver criada_em correta na resposta.
	atual, err := d.BuscarPorID(ctx, r.ID)
	if err != nil {
		return err
	}
	*r = *atual
	return nil
}

func (d *ReceitaDAO) Excluir(ctx context.Context, id int) error {
	res := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Receita{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNaoEncontrado
	}
	return nil
}

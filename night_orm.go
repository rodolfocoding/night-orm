package night_orm

import (
	"context"
	"github.com/rodolfocoding/night-orm/pkg/core"
	"github.com/rodolfocoding/night-orm/pkg/postgres"
)

// ORM é a interface principal que define as operações básicas do ORM
type ORM = core.ORM

// Model é a interface que todos os modelos devem implementar
type Model = core.Model

// ModelWithPrimaryKey é uma interface para modelos com chave primária
type ModelWithPrimaryKey = core.ModelWithPrimaryKey

// Transaction representa uma transação de banco de dados
type Transaction = core.Transaction

// NewPostgresORM cria uma nova instância do ORM para PostgreSQL
func NewPostgresORM() ORM {
	return postgres.NewPostgresORM()
}

// Connect é uma função auxiliar para conectar ao banco de dados PostgreSQL
func Connect(ctx context.Context, connectionString string) (ORM, error) {
	orm := NewPostgresORM()
	err := orm.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return orm, nil
}

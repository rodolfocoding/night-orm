package core

import (
	"context"
	"database/sql"
)

// ORM é a interface principal que define as operações básicas do ORM
type ORM interface {
	// Connect estabelece uma conexão com o banco de dados
	Connect(ctx context.Context, connectionString string) error
	
	// Close fecha a conexão com o banco de dados
	Close() error
	
	// DB retorna a conexão subjacente com o banco de dados
	DB() *sql.DB
	
	// Create insere um novo registro no banco de dados
	Create(ctx context.Context, model Model) error
	
	// FindByID busca um registro pelo ID
	FindByID(ctx context.Context, model ModelWithPrimaryKey, id interface{}) error
	
	// FindAll busca todos os registros de um modelo
	FindAll(ctx context.Context, model Model, dest interface{}) error
	
	// Update atualiza um registro existente
	Update(ctx context.Context, model ModelWithPrimaryKey) error
	
	// Delete remove um registro do banco de dados
	Delete(ctx context.Context, model ModelWithPrimaryKey) error
	
	// Query executa uma consulta SQL personalizada
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	
	// Exec executa um comando SQL personalizado
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	
	// Transaction inicia uma nova transação
	Transaction(ctx context.Context) (Transaction, error)
}

// Transaction representa uma transação de banco de dados
type Transaction interface {
	// Commit confirma a transação
	Commit() error
	
	// Rollback reverte a transação
	Rollback() error
	
	// Create insere um novo registro dentro da transação
	Create(ctx context.Context, model Model) error
	
	// Update atualiza um registro dentro da transação
	Update(ctx context.Context, model ModelWithPrimaryKey) error
	
	// Delete remove um registro dentro da transação
	Delete(ctx context.Context, model ModelWithPrimaryKey) error
	
	// Query executa uma consulta SQL personalizada dentro da transação
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	
	// Exec executa um comando SQL personalizado dentro da transação
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

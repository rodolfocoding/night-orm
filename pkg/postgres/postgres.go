package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/rodolfocoding/night-orm/pkg/core"
	"github.com/rodolfocoding/night-orm/pkg/utils"

	"github.com/lib/pq"
)

// PostgresORM é a implementação do ORM para PostgreSQL
type PostgresORM struct {
	db *sql.DB
}

// NewPostgresORM cria uma nova instância do ORM para PostgreSQL
func NewPostgresORM() *PostgresORM {
	return &PostgresORM{}
}

// Connect estabelece uma conexão com o banco de dados PostgreSQL
func (p *PostgresORM) Connect(ctx context.Context, connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao PostgreSQL: %w", err)
	}

	// Testa a conexão
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("erro ao testar conexão com PostgreSQL: %w", err)
	}

	p.db = db
	return nil
}

// Close fecha a conexão com o banco de dados
func (p *PostgresORM) Close() error {
	if p.db == nil {
		return errors.New("conexão não estabelecida")
	}
	return p.db.Close()
}

// DB retorna a conexão subjacente com o banco de dados
func (p *PostgresORM) DB() *sql.DB {
	return p.db
}

// Create insere um novo registro no banco de dados
func (p *PostgresORM) Create(ctx context.Context, model core.Model) error {
	if p.db == nil {
		return errors.New("conexão não estabelecida")
	}

	// Obtém os campos da estrutura
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("erro ao obter campos da estrutura: %w", err)
	}

	// Prepara a consulta de inserção
	qb := utils.NewQueryBuilder()
	columns := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields))

	for column, value := range fields {
		columns = append(columns, column)
		values = append(values, value)
	}

	qb.WriteInsert(model.TableName(), columns, values)
	query, args := qb.Build()

	// Executa a consulta
	_, err = p.db.ExecContext(ctx, query, args...)
	if err != nil {
		// Verifica se é um erro de violação de chave única
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("registro já existe: %w", err)
		}
		return fmt.Errorf("erro ao inserir registro: %w", err)
	}

	return nil
}

// FindByID busca um registro pelo ID
func (p *PostgresORM) FindByID(ctx context.Context, model core.ModelWithPrimaryKey, id interface{}) error {
	if p.db == nil {
		return errors.New("conexão não estabelecida")
	}

	// Constrói a consulta
	qb := utils.NewQueryBuilder()
	qb.WriteSelect().
		WriteFrom(model.TableName()).
		WriteWhere(fmt.Sprintf("%s = %s", model.PrimaryKey(), qb.AddParam(id)))

	query, args := qb.Build()

	// Executa a consulta
	row := p.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return fmt.Errorf("erro ao executar consulta: %w", row.Err())
	}

	// Obtém os campos da estrutura
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("modelo deve ser um ponteiro não-nil")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return errors.New("modelo deve ser um ponteiro para uma estrutura")
	}

	// Prepara os destinos para o scan
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("erro ao obter campos da estrutura: %w", err)
	}

	// Cria slices para os nomes das colunas e os destinos dos valores
	columns := make([]string, 0, len(fields))
	destinations := make([]interface{}, 0, len(fields))

	for column := range fields {
		columns = append(columns, column)
		// Cria um destino para cada campo
		dest := reflect.New(reflect.TypeOf(fields[column])).Interface()
		destinations = append(destinations, dest)
	}

	// Faz o scan dos valores
	if err := row.Scan(destinations...); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("registro não encontrado")
		}
		return fmt.Errorf("erro ao fazer scan dos valores: %w", err)
	}

	// Define os valores nos campos da estrutura
	for i, column := range columns {
		// Obtém o valor do destino
		destVal := reflect.ValueOf(destinations[i]).Elem().Interface()
		// Define o valor no campo da estrutura
		if err := utils.SetStructField(model, column, destVal); err != nil {
			return fmt.Errorf("erro ao definir valor no campo %s: %w", column, err)
		}
	}

	return nil
}

// FindAll busca todos os registros de um modelo
func (p *PostgresORM) FindAll(ctx context.Context, model core.Model, dest interface{}) error {
	if p.db == nil {
		return errors.New("conexão não estabelecida")
	}

	// Verifica se o destino é um slice de ponteiros para o tipo do modelo
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return errors.New("destino deve ser um ponteiro não-nil para um slice")
	}
	destVal = destVal.Elem()
	if destVal.Kind() != reflect.Slice {
		return errors.New("destino deve ser um ponteiro para um slice")
	}

	// Constrói a consulta
	qb := utils.NewQueryBuilder()
	qb.WriteSelect().WriteFrom(model.TableName())
	query, args := qb.Build()

	// Executa a consulta
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("erro ao executar consulta: %w", err)
	}
	defer rows.Close()

	// Obtém as colunas da consulta
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("erro ao obter colunas: %w", err)
	}

	// Cria um slice para armazenar os resultados
	sliceType := destVal.Type()
	elemType := sliceType.Elem()
	
	// Itera sobre os resultados
	for rows.Next() {
		// Cria uma nova instância do tipo do elemento
		elemVal := reflect.New(elemType.Elem()).Elem()
		
		// Prepara os destinos para o scan
		destinations := make([]interface{}, len(columns))
		for i, column := range columns {
			// Cria um destino para cada coluna
			field := elemVal.FieldByNameFunc(func(name string) bool {
				return strings.EqualFold(name, column) || strings.EqualFold(utils.GetTagName(elemType.Elem(), name, "db"), column)
			})
			
			if field.IsValid() && field.CanAddr() {
				destinations[i] = field.Addr().Interface()
			} else {
				// Se o campo não for encontrado, usa um destino descartável
				var dest interface{}
				destinations[i] = &dest
			}
		}
		
		// Faz o scan dos valores
		if err := rows.Scan(destinations...); err != nil {
			return fmt.Errorf("erro ao fazer scan dos valores: %w", err)
		}
		
		// Adiciona o elemento ao slice de destino
		destVal.Set(reflect.Append(destVal, reflect.New(elemType.Elem())))
		destVal.Index(destVal.Len() - 1).Set(elemVal.Addr())
	}
	
	if err := rows.Err(); err != nil {
		return fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}
	
	return nil
}

// Update atualiza um registro existente
func (p *PostgresORM) Update(ctx context.Context, model core.ModelWithPrimaryKey) error {
	if p.db == nil {
		return errors.New("conexão não estabelecida")
	}

	// Obtém os campos da estrutura
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("erro ao obter campos da estrutura: %w", err)
	}

	// Remove a chave primária dos campos a serem atualizados
	primaryKey := model.PrimaryKey()
	primaryKeyValue := model.PrimaryKeyValue()
	delete(fields, primaryKey)

	// Prepara a consulta de atualização
	qb := utils.NewQueryBuilder()
	columns := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields))

	for column, value := range fields {
		columns = append(columns, column)
		values = append(values, value)
	}

	qb.WriteUpdate(model.TableName(), columns, values).
		WriteWhere(fmt.Sprintf("%s = %s", primaryKey, qb.AddParam(primaryKeyValue)))

	query, args := qb.Build()

	// Executa a consulta
	result, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("erro ao atualizar registro: %w", err)
	}

	// Verifica se algum registro foi afetado
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter número de linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum registro foi atualizado")
	}

	return nil
}

// Delete remove um registro do banco de dados
func (p *PostgresORM) Delete(ctx context.Context, model core.ModelWithPrimaryKey) error {
	if p.db == nil {
		return errors.New("conexão não estabelecida")
	}

	// Constrói a consulta
	qb := utils.NewQueryBuilder()
	qb.WriteDelete(model.TableName()).
		WriteWhere(fmt.Sprintf("%s = %s", model.PrimaryKey(), qb.AddParam(model.PrimaryKeyValue())))

	query, args := qb.Build()

	// Executa a consulta
	result, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("erro ao excluir registro: %w", err)
	}

	// Verifica se algum registro foi afetado
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter número de linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum registro foi excluído")
	}

	return nil
}

// Query executa uma consulta SQL personalizada
func (p *PostgresORM) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if p.db == nil {
		return nil, errors.New("conexão não estabelecida")
	}
	return p.db.QueryContext(ctx, query, args...)
}

// Exec executa um comando SQL personalizado
func (p *PostgresORM) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if p.db == nil {
		return nil, errors.New("conexão não estabelecida")
	}
	return p.db.ExecContext(ctx, query, args...)
}

// Transaction inicia uma nova transação
func (p *PostgresORM) Transaction(ctx context.Context) (core.Transaction, error) {
	if p.db == nil {
		return nil, errors.New("conexão não estabelecida")
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	return &PostgresTransaction{tx: tx}, nil
}

// PostgresTransaction é a implementação de Transaction para PostgreSQL
type PostgresTransaction struct {
	tx *sql.Tx
}

// Commit confirma a transação
func (t *PostgresTransaction) Commit() error {
	return t.tx.Commit()
}

// Rollback reverte a transação
func (t *PostgresTransaction) Rollback() error {
	return t.tx.Rollback()
}

// Create insere um novo registro dentro da transação
func (t *PostgresTransaction) Create(ctx context.Context, model core.Model) error {
	// Obtém os campos da estrutura
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("erro ao obter campos da estrutura: %w", err)
	}

	// Prepara a consulta de inserção
	qb := utils.NewQueryBuilder()
	columns := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields))

	for column, value := range fields {
		columns = append(columns, column)
		values = append(values, value)
	}

	qb.WriteInsert(model.TableName(), columns, values)
	query, args := qb.Build()

	// Executa a consulta
	_, err = t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		// Verifica se é um erro de violação de chave única
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("registro já existe: %w", err)
		}
		return fmt.Errorf("erro ao inserir registro: %w", err)
	}

	return nil
}

// Update atualiza um registro dentro da transação
func (t *PostgresTransaction) Update(ctx context.Context, model core.ModelWithPrimaryKey) error {
	// Obtém os campos da estrutura
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("erro ao obter campos da estrutura: %w", err)
	}

	// Remove a chave primária dos campos a serem atualizados
	primaryKey := model.PrimaryKey()
	primaryKeyValue := model.PrimaryKeyValue()
	delete(fields, primaryKey)

	// Prepara a consulta de atualização
	qb := utils.NewQueryBuilder()
	columns := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields))

	for column, value := range fields {
		columns = append(columns, column)
		values = append(values, value)
	}

	qb.WriteUpdate(model.TableName(), columns, values).
		WriteWhere(fmt.Sprintf("%s = %s", primaryKey, qb.AddParam(primaryKeyValue)))

	query, args := qb.Build()

	// Executa a consulta
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("erro ao atualizar registro: %w", err)
	}

	// Verifica se algum registro foi afetado
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter número de linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum registro foi atualizado")
	}

	return nil
}

// Delete remove um registro dentro da transação
func (t *PostgresTransaction) Delete(ctx context.Context, model core.ModelWithPrimaryKey) error {
	// Constrói a consulta
	qb := utils.NewQueryBuilder()
	qb.WriteDelete(model.TableName()).
		WriteWhere(fmt.Sprintf("%s = %s", model.PrimaryKey(), qb.AddParam(model.PrimaryKeyValue())))

	query, args := qb.Build()

	// Executa a consulta
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("erro ao excluir registro: %w", err)
	}

	// Verifica se algum registro foi afetado
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter número de linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum registro foi excluído")
	}

	return nil
}

// Query executa uma consulta SQL personalizada dentro da transação
func (t *PostgresTransaction) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

// Exec executa um comando SQL personalizado dentro da transação
func (t *PostgresTransaction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

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

// PostgresORM is the PostgreSQL ORM implementation
type PostgresORM struct {
	db *sql.DB
}

// NewPostgresORM creates a new instance of the PostgreSQL ORM
func NewPostgresORM() *PostgresORM {
	return &PostgresORM{}
}

// Connect establishes a connection to the PostgreSQL database
func (p *PostgresORM) Connect(ctx context.Context, connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("error connecting to PostgreSQL: %w", err)
	}

	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("error pinging PostgreSQL connection: %w", err)
	}

	p.db = db
	return nil
}

// Close closes the database connection
func (p *PostgresORM) Close() error {
	if p.db == nil {
		return errors.New("connection not established")
	}
	return p.db.Close()
}

// DB returns the underlying database connection
func (p *PostgresORM) DB() *sql.DB {
	return p.db
}

// Create inserts a new record into the database
func (p *PostgresORM) Create(ctx context.Context, model core.Model) error {
	if p.db == nil {
		return errors.New("connection not established")
	}

	// Get the struct fields
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("error retrieving struct fields: %w", err)
	}

	// Prepare the insert query
	qb := utils.NewQueryBuilder()
	columns := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields))

	// Check if the model implements ModelWithPrimaryKey to identify the primary key
	var primaryKey string
	var primaryKeyValue interface{}
	if modelWithPK, ok := model.(core.ModelWithPrimaryKey); ok {
		primaryKey = modelWithPK.PrimaryKey()
		primaryKeyValue = modelWithPK.PrimaryKeyValue()
	}

	// Filter fields, omitting the primary key if its value is zero
	for column, value := range fields {
		if column == primaryKey && reflect.ValueOf(primaryKeyValue).IsZero() {
			continue // Omit the primary key if its value is zero
		}
		columns = append(columns, column)
		values = append(values, value)
	}

	qb.WriteInsert(model.TableName(), columns, values)
	// Add RETURNING to retrieve the generated ID
	if primaryKey != "" {
		qb.WriteReturning(primaryKey)
	}
	query, args := qb.Build()

	// Execute the query and capture the returned ID
	var generatedID int
	if primaryKey != "" {
		err = p.db.QueryRowContext(ctx, query, args...).Scan(&generatedID)
	} else {
		_, err = p.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("record already exists: %w", err)
		}
		return fmt.Errorf("error inserting record: %w", err)
	}

	// Update the model with the generated ID, if applicable
	if primaryKey != "" {
		if err := utils.SetStructField(model, primaryKey, generatedID); err != nil {
			return fmt.Errorf("error setting primary key value: %w", err)
		}
	}

	return nil
}

// FindByID retrieves a record by ID
func (p *PostgresORM) FindByID(ctx context.Context, model core.ModelWithPrimaryKey, id interface{}) error {
	if p.db == nil {
		return errors.New("connection not established")
	}

	// Build the query
	qb := utils.NewQueryBuilder()
	qb.WriteSelect().
		WriteFrom(model.TableName()).
		WriteWhere(fmt.Sprintf("%s = %s", model.PrimaryKey(), qb.AddParam(id)))

	query, args := qb.Build()

	// Execute the query
	row := p.db.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return fmt.Errorf("error executing query: %w", row.Err())
	}

	// Get the struct fields
	val := reflect.ValueOf(model)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("model must be a non-nil pointer")
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return errors.New("model must be a pointer to a struct")
	}

	// Prepare destinations for scanning
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("error retrieving struct fields: %w", err)
	}

	// Create slices for column names and value destinations
	columns := make([]string, 0, len(fields))
	destinations := make([]interface{}, 0, len(fields))

	for column := range fields {
		columns = append(columns, column)
		// Create a destination for each field
		dest := reflect.New(reflect.TypeOf(fields[column])).Interface()
		destinations = append(destinations, dest)
	}

	// Scan the values
	if err := row.Scan(destinations...); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("record not found")
		}
		return fmt.Errorf("error scanning values: %w", err)
	}

	// Set the values in the struct fields
	for i, column := range columns {
		// Get the destination value
		destVal := reflect.ValueOf(destinations[i]).Elem().Interface()
		// Set the value in the struct field
		if err := utils.SetStructField(model, column, destVal); err != nil {
			return fmt.Errorf("error setting value for field %s: %w", column, err)
		}
	}

	return nil
}

// FindAll retrieves all records of a model
func (p *PostgresORM) FindAll(ctx context.Context, model core.Model, dest interface{}) error {
	if p.db == nil {
		return errors.New("connection not established")
	}

	// Verify that the destination is a slice pointer
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return errors.New("destination must be a non-nil pointer to a slice")
	}
	destVal = destVal.Elem()
	if destVal.Kind() != reflect.Slice {
		return errors.New("destination must be a pointer to a slice")
	}

	// Build the query
	qb := utils.NewQueryBuilder()
	qb.WriteSelect().WriteFrom(model.TableName())
	query, args := qb.Build()

	// Execute the query
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Get the query columns
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error retrieving columns: %w", err)
	}

	// Create a slice to store the results
	sliceType := destVal.Type()
	elemType := sliceType.Elem()

	// Iterate over the results
	for rows.Next() {
		// Create a new instance of the element type
		elemVal := reflect.New(elemType.Elem()).Elem()

		// Prepare destinations for scanning
		destinations := make([]interface{}, len(columns))
		for i, column := range columns {
			// Create a destination for each column
			field := elemVal.FieldByNameFunc(func(name string) bool {
				return strings.EqualFold(name, column) || strings.EqualFold(utils.GetTagName(elemType.Elem(), name, "db"), column)
			})

			if field.IsValid() && field.CanAddr() {
				destinations[i] = field.Addr().Interface()
			} else {
				// Use a disposable destination if the field is not found
				var dest interface{}
				destinations[i] = &dest
			}
		}

		// Scan the values
		if err := rows.Scan(destinations...); err != nil {
			return fmt.Errorf("error scanning values: %w", err)
		}

		// Add the element to the destination slice
		destVal.Set(reflect.Append(destVal, reflect.New(elemType.Elem())))
		destVal.Index(destVal.Len()-1).Set(elemVal.Addr())
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over results: %w", err)
	}

	return nil
}

// Update updates an existing record
func (p *PostgresORM) Update(ctx context.Context, model core.ModelWithPrimaryKey) error {
	if p.db == nil {
		return errors.New("connection not established")
	}

	// Get the struct fields
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("error retrieving struct fields: %w", err)
	}

	// Remove the primary key from the fields to be updated
	primaryKey := model.PrimaryKey()
	primaryKeyValue := model.PrimaryKeyValue()
	delete(fields, primaryKey)

	// Prepare the update query
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

	// Execute the query
	result, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating record: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows count: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no records were updated")
	}

	return nil
}

// Delete removes a record from the database
func (p *PostgresORM) Delete(ctx context.Context, model core.ModelWithPrimaryKey) error {
	if p.db == nil {
		return errors.New("connection not established")
	}

	// Build the query
	qb := utils.NewQueryBuilder()
	qb.WriteDelete(model.TableName()).
		WriteWhere(fmt.Sprintf("%s = %s", model.PrimaryKey(), qb.AddParam(model.PrimaryKeyValue())))

	query, args := qb.Build()

	// Execute the query
	result, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error deleting record: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows count: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no records were deleted")
	}

	return nil
}

// Query executes a custom SQL query
func (p *PostgresORM) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if p.db == nil {
		return nil, errors.New("connection not established")
	}
	return p.db.QueryContext(ctx, query, args...)
}

// Exec executes a custom SQL command
func (p *PostgresORM) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if p.db == nil {
		return nil, errors.New("connection not established")
	}
	return p.db.ExecContext(ctx, query, args...)
}

// Transaction starts a new transaction
func (p *PostgresORM) Transaction(ctx context.Context) (core.Transaction, error) {
	if p.db == nil {
		return nil, errors.New("connection not established")
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	return &PostgresTransaction{tx: tx}, nil
}

// PostgresTransaction is the PostgreSQL transaction implementation
type PostgresTransaction struct {
	tx *sql.Tx
}

// Commit commits the transaction
func (t *PostgresTransaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *PostgresTransaction) Rollback() error {
	return t.tx.Rollback()
}

// Create inserts a new record within the transaction
func (t *PostgresTransaction) Create(ctx context.Context, model core.Model) error {
	// Get the struct fields
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("error retrieving struct fields: %w", err)
	}

	// Prepare the insert query
	qb := utils.NewQueryBuilder()
	columns := make([]string, 0, len(fields))
	values := make([]interface{}, 0, len(fields))

	for column, value := range fields {
		columns = append(columns, column)
		values = append(values, value)
	}

	qb.WriteInsert(model.TableName(), columns, values)
	query, args := qb.Build()

	// Execute the query
	_, err = t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		// Check for unique constraint violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("record already exists: %w", err)
		}
		return fmt.Errorf("error inserting record: %w", err)
	}

	return nil
}

// Update updates a record within the transaction
func (t *PostgresTransaction) Update(ctx context.Context, model core.ModelWithPrimaryKey) error {
	// Get the struct fields
	fields, err := utils.GetStructFields(model)
	if err != nil {
		return fmt.Errorf("error retrieving struct fields: %w", err)
	}

	// Remove the primary key from the fields to be updated
	primaryKey := model.PrimaryKey()
	primaryKeyValue := model.PrimaryKeyValue()
	delete(fields, primaryKey)

	// Prepare the update query
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

	// Execute the query
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating record: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows count: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no records were updated")
	}

	return nil
}

// Delete removes a record within the transaction
func (t *PostgresTransaction) Delete(ctx context.Context, model core.ModelWithPrimaryKey) error {
	// Build the query
	qb := utils.NewQueryBuilder()
	qb.WriteDelete(model.TableName()).
		WriteWhere(fmt.Sprintf("%s = %s", model.PrimaryKey(), qb.AddParam(model.PrimaryKeyValue())))

	query, args := qb.Build()

	// Execute the query
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error deleting record: %w", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving affected rows count: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no records were deleted")
	}

	return nil
}

// Query executes a custom SQL query within the transaction
func (t *PostgresTransaction) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

// Exec executes a custom SQL command within the transaction
func (t *PostgresTransaction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}
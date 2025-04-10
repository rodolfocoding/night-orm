package utils

import (
	"fmt"
	"strings"
)

// QueryBuilder é um construtor de consultas SQL
type QueryBuilder struct {
	query      strings.Builder
	args       []interface{}
	paramIndex int
}

// NewQueryBuilder cria um novo construtor de consultas
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		query:      strings.Builder{},
		args:       make([]interface{}, 0),
		paramIndex: 1,
	}
}

// Reset limpa o construtor de consultas
func (qb *QueryBuilder) Reset() {
	qb.query.Reset()
	qb.args = make([]interface{}, 0)
	qb.paramIndex = 1
}

// AddParam adiciona um parâmetro à consulta e retorna o placeholder
func (qb *QueryBuilder) AddParam(value interface{}) string {
	qb.args = append(qb.args, value)
	placeholder := fmt.Sprintf("$%d", qb.paramIndex)
	qb.paramIndex++
	return placeholder
}

// Write adiciona texto à consulta
func (qb *QueryBuilder) Write(s string) *QueryBuilder {
	qb.query.WriteString(s)
	return qb
}

// WriteWithParams adiciona texto com parâmetros à consulta
func (qb *QueryBuilder) WriteWithParams(format string, args ...interface{}) *QueryBuilder {
	placeholders := make([]interface{}, len(args))
	for i, arg := range args {
		placeholders[i] = qb.AddParam(arg)
	}
	qb.query.WriteString(fmt.Sprintf(format, placeholders...))
	return qb
}

// WriteSelect adiciona uma cláusula SELECT à consulta
func (qb *QueryBuilder) WriteSelect(columns ...string) *QueryBuilder {
	qb.Write("SELECT ")
	if len(columns) == 0 {
		qb.Write("*")
	} else {
		qb.Write(strings.Join(columns, ", "))
	}
	return qb
}

// WriteFrom adiciona uma cláusula FROM à consulta
func (qb *QueryBuilder) WriteFrom(table string) *QueryBuilder {
	qb.Write(" FROM ")
	qb.Write(table)
	return qb
}

// WriteWhere adiciona uma cláusula WHERE à consulta
func (qb *QueryBuilder) WriteWhere(condition string, args ...interface{}) *QueryBuilder {
	qb.Write(" WHERE ")
	return qb.WriteWithParams(condition, args...)
}

// WriteAnd adiciona uma cláusula AND à consulta
func (qb *QueryBuilder) WriteAnd(condition string, args ...interface{}) *QueryBuilder {
	qb.Write(" AND ")
	return qb.WriteWithParams(condition, args...)
}

// WriteOr adiciona uma cláusula OR à consulta
func (qb *QueryBuilder) WriteOr(condition string, args ...interface{}) *QueryBuilder {
	qb.Write(" OR ")
	return qb.WriteWithParams(condition, args...)
}

// WriteOrderBy adiciona uma cláusula ORDER BY à consulta
func (qb *QueryBuilder) WriteOrderBy(columns ...string) *QueryBuilder {
	if len(columns) > 0 {
		qb.Write(" ORDER BY ")
		qb.Write(strings.Join(columns, ", "))
	}
	return qb
}

// WriteLimit adiciona uma cláusula LIMIT à consulta
func (qb *QueryBuilder) WriteLimit(limit int) *QueryBuilder {
	if limit > 0 {
		qb.Write(fmt.Sprintf(" LIMIT %d", limit))
	}
	return qb
}

// WriteOffset adiciona uma cláusula OFFSET à consulta
func (qb *QueryBuilder) WriteOffset(offset int) *QueryBuilder {
	if offset > 0 {
		qb.Write(fmt.Sprintf(" OFFSET %d", offset))
	}
	return qb
}

// WriteInsert adiciona uma cláusula INSERT à consulta
func (qb *QueryBuilder) WriteInsert(table string, columns []string, values []interface{}) *QueryBuilder {
	qb.Write(fmt.Sprintf("INSERT INTO %s (", table))
	qb.Write(strings.Join(columns, ", "))
	qb.Write(") VALUES (")

	placeholders := make([]string, len(values))
	for i, value := range values {
		placeholders[i] = qb.AddParam(value)
	}
	qb.Write(strings.Join(placeholders, ", "))
	qb.Write(")")
	return qb
}

// WriteUpdate adiciona uma cláusula UPDATE à consulta
func (qb *QueryBuilder) WriteUpdate(table string, columns []string, values []interface{}) *QueryBuilder {
	qb.Write(fmt.Sprintf("UPDATE %s SET ", table))

	for i := 0; i < len(columns); i++ {
		if i > 0 {
			qb.Write(", ")
		}
		qb.Write(fmt.Sprintf("%s = %s", columns[i], qb.AddParam(values[i])))
	}
	return qb
}

// WriteDelete adiciona uma cláusula DELETE à consulta
func (qb *QueryBuilder) WriteDelete(table string) *QueryBuilder {
	qb.Write(fmt.Sprintf("DELETE FROM %s", table))
	return qb
}

// WriteReturning adiciona uma cláusula RETURNING à consulta
func (qb *QueryBuilder) WriteReturning(columns ...string) *QueryBuilder {
	if len(columns) > 0 {
		qb.Write(" RETURNING ")
		qb.Write(strings.Join(columns, ", "))
	}
	return qb
}

// Build retorna a consulta SQL e os argumentos
func (qb *QueryBuilder) Build() (string, []interface{}) {
	return qb.query.String(), qb.args
}

package utils

import (
	"testing"
)

func TestQueryBuilder(t *testing.T) {
	t.Run("WriteSelect", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect("id", "name", "email")
		query, _ := qb.Build()
		expected := "SELECT id, name, email"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
	})

	t.Run("WriteSelectAll", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect()
		query, _ := qb.Build()
		expected := "SELECT *"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
	})

	t.Run("WriteFrom", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect().WriteFrom("users")
		query, _ := qb.Build()
		expected := "SELECT * FROM users"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
	})

	t.Run("WriteWhere", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect().WriteFrom("users").WriteWhere("id = %s", 1)
		query, args := qb.Build()
		expected := "SELECT * FROM users WHERE id = $1"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args to be [1], got %v", args)
		}
	})

	t.Run("WriteAnd", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect().WriteFrom("users").
			WriteWhere("id = %s", 1).
			WriteAnd("name = %s", "John")
		query, args := qb.Build()
		expected := "SELECT * FROM users WHERE id = $1 AND name = $2"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
		if len(args) != 2 || args[0] != 1 || args[1] != "John" {
			t.Errorf("Expected args to be [1, 'John'], got %v", args)
		}
	})

	t.Run("WriteOr", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect().WriteFrom("users").
			WriteWhere("id = %s", 1).
			WriteOr("id = %s", 2)
		query, args := qb.Build()
		expected := "SELECT * FROM users WHERE id = $1 OR id = $2"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
		if len(args) != 2 || args[0] != 1 || args[1] != 2 {
			t.Errorf("Expected args to be [1, 2], got %v", args)
		}
	})

	t.Run("WriteOrderBy", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect().WriteFrom("users").WriteOrderBy("name", "id DESC")
		query, _ := qb.Build()
		expected := "SELECT * FROM users ORDER BY name, id DESC"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
	})

	t.Run("WriteLimit", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect().WriteFrom("users").WriteLimit(10)
		query, _ := qb.Build()
		expected := "SELECT * FROM users LIMIT 10"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
	})

	t.Run("WriteOffset", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect().WriteFrom("users").WriteLimit(10).WriteOffset(5)
		query, _ := qb.Build()
		expected := "SELECT * FROM users LIMIT 10 OFFSET 5"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
	})

	t.Run("WriteInsert", func(t *testing.T) {
		qb := NewQueryBuilder()
		columns := []string{"name", "email"}
		values := []interface{}{"John", "john@example.com"}
		qb.WriteInsert("users", columns, values)
		query, args := qb.Build()
		expected := "INSERT INTO users (name, email) VALUES ($1, $2)"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
		if len(args) != 2 || args[0] != "John" || args[1] != "john@example.com" {
			t.Errorf("Expected args to be ['John', 'john@example.com'], got %v", args)
		}
	})

	t.Run("WriteUpdate", func(t *testing.T) {
		qb := NewQueryBuilder()
		columns := []string{"name", "email"}
		values := []interface{}{"John", "john@example.com"}
		qb.WriteUpdate("users", columns, values).WriteWhere("id = %s", 1)
		query, args := qb.Build()
		expected := "UPDATE users SET name = $1, email = $2 WHERE id = $3"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
		if len(args) != 3 || args[0] != "John" || args[1] != "john@example.com" || args[2] != 1 {
			t.Errorf("Expected args to be ['John', 'john@example.com', 1], got %v", args)
		}
	})

	t.Run("WriteDelete", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteDelete("users").WriteWhere("id = %s", 1)
		query, args := qb.Build()
		expected := "DELETE FROM users WHERE id = $1"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
		if len(args) != 1 || args[0] != 1 {
			t.Errorf("Expected args to be [1], got %v", args)
		}
	})

	t.Run("WriteReturning", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteInsert("users", []string{"name"}, []interface{}{"John"}).
			WriteReturning("id", "created_at")
		query, _ := qb.Build()
		expected := "INSERT INTO users (name) VALUES ($1) RETURNING id, created_at"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect().WriteFrom("users")
		qb.Reset()
		query, args := qb.Build()
		if query != "" {
			t.Errorf("Expected query to be empty after reset, got '%s'", query)
		}
		if len(args) != 0 {
			t.Errorf("Expected args to be empty after reset, got %v", args)
		}
	})

	t.Run("ComplexQuery", func(t *testing.T) {
		qb := NewQueryBuilder()
		qb.WriteSelect("u.id", "u.name", "u.email").
			WriteFrom("users u").
			WriteWhere("u.active = %s", true).
			WriteAnd("u.created_at > %s", "2023-01-01").
			WriteOrderBy("u.name ASC").
			WriteLimit(10).
			WriteOffset(20)
		query, args := qb.Build()
		expected := "SELECT u.id, u.name, u.email FROM users u WHERE u.active = $1 AND u.created_at > $2 ORDER BY u.name ASC LIMIT 10 OFFSET 20"
		if query != expected {
			t.Errorf("Expected query to be '%s', got '%s'", expected, query)
		}
		if len(args) != 2 || args[0] != true || args[1] != "2023-01-01" {
			t.Errorf("Expected args to be [true, '2023-01-01'], got %v", args)
		}
	})
}

# Database Support in NightORM

This document describes the current database support in NightORM and how to add support for new databases.

## Supported Databases

Currently, NightORM supports the following databases:

### PostgreSQL

PostgreSQL is the first database supported by NightORM. The implementation is in the `pkg/postgres` package.

To use NightORM with PostgreSQL:

```go
import (
    "context"
    "log"

    "github.com/your-username/night-orm"
)

func main() {
    // Connection string for PostgreSQL
    connectionString := "postgres://username:password@localhost:5432/database?sslmode=disable"
    ctx := context.Background()

    // Connect to PostgreSQL
    orm, err := night_orm.Connect(ctx, connectionString)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    defer orm.Close()

    // Now you can use the ORM to interact with PostgreSQL
}
```

## Adding Support for New Databases

NightORM was designed to be easily extensible to support different databases. To add support for a new database, follow the steps below:

### 1. Create a New Package

Create a new package inside the `pkg/` directory with the name of the database. For example, to add support for MySQL, create the directory `pkg/mysql/`.

### 2. Implement the ORM Interface

Implement the `ORM` interface defined in `pkg/core/orm.go`. You will need to implement all the methods defined in the interface.

Example structure for MySQL:

```go
package mysql

import (
    "context"
    "database/sql"

    _ "github.com/go-sql-driver/mysql" // MySQL driver
    "github.com/your-username/night-orm/pkg/core"
)

// MySQLORM is the ORM implementation for MySQL
type MySQLORM struct {
    db *sql.DB
}

// NewMySQLORM creates a new instance of the ORM for MySQL
func NewMySQLORM() *MySQLORM {
    return &MySQLORM{}
}

// Connect establishes a connection to the MySQL database
func (m *MySQLORM) Connect(ctx context.Context, connectionString string) error {
    // Implementation of the connection to MySQL
}

// Implement the other methods of the ORM interface...
```

### 3. Add Factory Functions

Add factory functions in the main `night_orm.go` file to facilitate the creation of ORM instances for the new database.

```go
// NewMySQLORM creates a new instance of the ORM for MySQL
func NewMySQLORM() ORM {
    return mysql.NewMySQLORM()
}

// ConnectMySQL is a helper function to connect to the MySQL database
func ConnectMySQL(ctx context.Context, connectionString string) (ORM, error) {
    orm := NewMySQLORM()
    err := orm.Connect(ctx, connectionString)
    if err != nil {
        return nil, err
    }
    return orm, nil
}
```

### 4. Add Tests

Add tests for the new implementation to ensure it works correctly.

### 5. Update Documentation

Update the documentation to include information about the newly supported database.

## Considerations for Different Databases

When implementing support for different databases, consider the following differences:

### SQL Syntax

Different databases may have slightly different SQL syntaxes. For example, PostgreSQL uses `$1`, `$2`, etc. for parameters, while MySQL uses `?`.

### Data Types

Data types may vary between databases. Make sure to correctly map Go data types to database data types.

### Specific Features

Some databases have specific features that may be useful for the ORM. For example, PostgreSQL has the `RETURNING` operator that allows returning values from rows affected by an insert, update, or delete operation.

### Error Handling

Different database drivers may return errors in different ways. Make sure to properly handle errors specific to each database.

## Planned Databases

The following databases are planned for future support:

- MySQL
- SQLite
- Microsoft SQL Server
- Oracle

If you are interested in contributing support for any of these databases, see the [CONTRIBUTING.md](../CONTRIBUTING.md) file for information on how to contribute to the project.

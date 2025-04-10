# NightORM Documentation

Welcome to the NightORM documentation, a simple and flexible ORM (Object-Relational Mapping) for Go.

## Table of Contents

- [Introduction](#introduction)
- [Guides](#guides)
- [Reference](#reference)
- [Examples](#examples)

## Introduction

NightORM is an ORM for Go that facilitates interaction with relational databases. It provides a simple and intuitive interface for CRUD (Create, Read, Update, Delete) operations and supports transactions.

Currently, NightORM supports PostgreSQL, with plans to expand to other databases in the future.

## Guides

- [Struct Tags](struct_tags.en.md) - How to use struct tags to customize the mapping between Go structures and database tables.
- [Transactions](transactions.en.md) - How to use transactions to ensure data integrity in operations that involve multiple changes to the database.
- [Database Support](database_support.en.md) - Information about supported databases and how to add support for new databases.

## Reference

### Main Interfaces

- `ORM` - Main interface that defines the basic operations of the ORM.
- `Model` - Interface that all models must implement.
- `ModelWithPrimaryKey` - Interface for models with a primary key.
- `Transaction` - Interface that represents a database transaction.

### Packages

- `pkg/core` - Main ORM interfaces and types.
- `pkg/postgres` - ORM implementation for PostgreSQL.
- `pkg/utils` - Utilities for reflection and SQL query building.

## Examples

See the [examples](../examples) directory for examples of using NightORM.

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/your-username/night-orm"
)

type User struct {
    ID        int       `db:"id,primary"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    Active    bool      `db:"active"`
}

func (u *User) TableName() string {
    return "users"
}

func (u *User) PrimaryKey() string {
    return "id"
}

func (u *User) PrimaryKeyValue() interface{} {
    return u.ID
}

func main() {
    // Connect to PostgreSQL database
    connectionString := "postgres://username:password@localhost:5432/database?sslmode=disable"
    ctx := context.Background()

    orm, err := night_orm.Connect(ctx, connectionString)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    defer orm.Close()

    // Create a new user
    user := &User{
        Name:      "John Smith",
        Email:     "john@example.com",
        CreatedAt: time.Now(),
        Active:    true,
    }

    // Insert the user into the database
    if err := orm.Create(ctx, user); err != nil {
        log.Fatalf("Error creating user: %v", err)
    }
    fmt.Printf("User created with ID: %d\n", user.ID)

    // Find a user by ID
    fetchedUser := &User{ID: user.ID}
    if err := orm.FindByID(ctx, fetchedUser, user.ID); err != nil {
        log.Fatalf("Error finding user: %v", err)
    }
    fmt.Printf("User found: %s (%s)\n", fetchedUser.Name, fetchedUser.Email)
}
```

## Contributing

See the [CONTRIBUTING.md](../CONTRIBUTING.md) file for information on how to contribute to the project.

## License

This project is licensed under the [MIT License](../LICENSE).

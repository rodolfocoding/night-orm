# Transactions in NightORM

This document describes how to use transactions in NightORM to ensure data integrity in operations that involve multiple changes to the database.

## Introduction

Transactions are used to group multiple database operations into a single logical unit of work. This ensures that all operations are completed successfully or none of them are applied, maintaining data integrity.

NightORM provides support for transactions through the `Transaction` interface, which allows you to execute CRUD operations within a transaction.

## Starting a Transaction

To start a transaction, use the `Transaction()` method of the `ORM` interface:

```go
tx, err := orm.Transaction(ctx)
if err != nil {
    log.Fatalf("Error starting transaction: %v", err)
}
```

## Operations within a Transaction

Once you have a transaction, you can execute CRUD operations within it:

### Creating a Record

```go
user := &User{
    Name:      "John Smith",
    Email:     "john@example.com",
    CreatedAt: time.Now(),
    Active:    true,
}

if err := tx.Create(ctx, user); err != nil {
    tx.Rollback()
    log.Fatalf("Error creating user: %v", err)
}
```

### Updating a Record

```go
user.Name = "John Smith Updated"
if err := tx.Update(ctx, user); err != nil {
    tx.Rollback()
    log.Fatalf("Error updating user: %v", err)
}
```

### Deleting a Record

```go
if err := tx.Delete(ctx, user); err != nil {
    tx.Rollback()
    log.Fatalf("Error deleting user: %v", err)
}
```

### Custom Queries

You can also execute custom SQL queries within a transaction:

```go
rows, err := tx.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
if err != nil {
    tx.Rollback()
    log.Fatalf("Error executing query: %v", err)
}
defer rows.Close()

// Process the results...
```

```go
result, err := tx.Exec(ctx, "UPDATE users SET active = $1 WHERE id = $2", false, 1)
if err != nil {
    tx.Rollback()
    log.Fatalf("Error executing command: %v", err)
}

// Check the result...
```

## Committing or Rolling Back a Transaction

After executing all necessary operations, you need to commit or roll back the transaction:

### Committing a Transaction

To commit a transaction (apply all changes), use the `Commit()` method:

```go
if err := tx.Commit(); err != nil {
    log.Fatalf("Error committing transaction: %v", err)
}
```

### Rolling Back a Transaction

To roll back a transaction (discard all changes), use the `Rollback()` method:

```go
if err := tx.Rollback(); err != nil {
    log.Fatalf("Error rolling back transaction: %v", err)
}
```

## Usage Pattern

A common pattern for using transactions is:

```go
// Start a transaction
tx, err := orm.Transaction(ctx)
if err != nil {
    log.Fatalf("Error starting transaction: %v", err)
}

// Cleanup function to ensure the transaction is rolled back in case of error
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
        panic(r) // Re-panic after rollback
    }
}()

// Execute operations within the transaction
if err := tx.Create(ctx, user1); err != nil {
    tx.Rollback()
    log.Fatalf("Error creating user 1: %v", err)
}

if err := tx.Create(ctx, user2); err != nil {
    tx.Rollback()
    log.Fatalf("Error creating user 2: %v", err)
}

// Commit the transaction
if err := tx.Commit(); err != nil {
    log.Fatalf("Error committing transaction: %v", err)
}
```

## Complete Example

Here's a complete example of how to use transactions in NightORM:

```go
package main

import (
    "context"
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
    // Connect to the database
    connectionString := "postgres://username:password@localhost:5432/database?sslmode=disable"
    ctx := context.Background()

    orm, err := night_orm.Connect(ctx, connectionString)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    defer orm.Close()

    // Start a transaction
    tx, err := orm.Transaction(ctx)
    if err != nil {
        log.Fatalf("Error starting transaction: %v", err)
    }

    // Create two users within the transaction
    user1 := &User{
        Name:      "John Smith",
        Email:     "john@example.com",
        CreatedAt: time.Now(),
        Active:    true,
    }

    user2 := &User{
        Name:      "Mary Johnson",
        Email:     "mary@example.com",
        CreatedAt: time.Now(),
        Active:    true,
    }

    // Try to create the first user
    if err := tx.Create(ctx, user1); err != nil {
        tx.Rollback()
        log.Fatalf("Error creating user 1: %v", err)
    }

    // Try to create the second user
    if err := tx.Create(ctx, user2); err != nil {
        tx.Rollback()
        log.Fatalf("Error creating user 2: %v", err)
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        log.Fatalf("Error committing transaction: %v", err)
    }

    log.Println("Transaction completed successfully!")
    log.Printf("User 1 created with ID: %d\n", user1.ID)
    log.Printf("User 2 created with ID: %d\n", user2.ID)
}
```

## Important Considerations

### Error Handling

It's important to always check the errors returned by operations within a transaction and call `Rollback()` in case of error.

### Resource Cleanup

Make sure to close resources like `*sql.Rows` even in case of error in the transaction.

### Context

Use the same context (`context.Context`) for all operations within a transaction to ensure consistency.

### Isolation

The isolation level of the transaction depends on the database and configuration. Consult your database documentation for more information.

## Conclusion

Transactions are a powerful tool for ensuring data integrity in operations that involve multiple changes to the database. NightORM provides a simple and intuitive interface for working with transactions, allowing you to execute CRUD operations and custom queries within a transaction.

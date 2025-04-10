# NightORM

NightORM is a simple and flexible ORM (Object-Relational Mapping) for Go, designed to facilitate interaction with relational databases. Currently, NightORM supports PostgreSQL, with plans to expand to other databases in the future.

[Versão em português](README.md)

## Features

- Simple and intuitive interface
- Support for basic CRUD operations
- Automatic mapping between Go structures and database tables
- Transaction support
- Custom SQL queries
- Support for struct tags to customize mapping

## Installation

```bash
go get github.com/rodolfocoding/night-orm
```

## Basic Usage

### Defining a Model

```go
type User struct {
    ID        int       `db:"id,primary"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    Active    bool      `db:"active"`
}

// TableName returns the table name in the database
func (u *User) TableName() string {
    return "users"
}

// PrimaryKey returns the primary key column name
func (u *User) PrimaryKey() string {
    return "id"
}

// PrimaryKeyValue returns the primary key value
func (u *User) PrimaryKeyValue() interface{} {
    return u.ID
}
```

### Connecting to the Database

```go
import (
    "context"
    "log"

    "github.com/rodolfocoding/night-orm"
)

func main() {
    // Connect to PostgreSQL database
    connectionString := "postgres://username:password@localhost:5432/database?sslmode=disable"
    ctx := context.Background()

    orm, err := night_orm.Connect(ctx, connectionString)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    defer orm.Close()

    // Now you can use the ORM to interact with the database
}
```

### CRUD Operations

#### Create a Record

```go
user := &User{
    Name:      "John Smith",
    Email:     "john@example.com",
    CreatedAt: time.Now(),
    Active:    true,
}

if err := orm.Create(ctx, user); err != nil {
    log.Fatalf("Error creating user: %v", err)
}
```

#### Find a Record by ID

```go
user := &User{ID: 1}
if err := orm.FindByID(ctx, user, 1); err != nil {
    log.Fatalf("Error finding user: %v", err)
}
```

#### Find All Records

```go
var users []*User
if err := orm.FindAll(ctx, &User{}, &users); err != nil {
    log.Fatalf("Error finding all users: %v", err)
}
```

#### Update a Record

```go
user.Name = "John Smith Updated"
if err := orm.Update(ctx, user); err != nil {
    log.Fatalf("Error updating user: %v", err)
}
```

#### Delete a Record

```go
if err := orm.Delete(ctx, user); err != nil {
    log.Fatalf("Error deleting user: %v", err)
}
```

### Custom Queries

```go
rows, err := orm.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
if err != nil {
    log.Fatalf("Error executing query: %v", err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil {
        log.Fatalf("Error scanning results: %v", err)
    }
    fmt.Printf("ID: %d, Name: %s\n", id, name)
}
```

### Transactions

```go
tx, err := orm.Transaction(ctx)
if err != nil {
    log.Fatalf("Error starting transaction: %v", err)
}

// Create a new user within the transaction
newUser := &User{
    Name:      "Mary Johnson",
    Email:     "mary@example.com",
    CreatedAt: time.Now(),
    Active:    true,
}

if err := tx.Create(ctx, newUser); err != nil {
    tx.Rollback()
    log.Fatalf("Error creating user in transaction: %v", err)
}

// Commit the transaction
if err := tx.Commit(); err != nil {
    log.Fatalf("Error committing transaction: %v", err)
}
```

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

### Adding Support for New Databases

To add support for a new database, you need to implement the `ORM` interface defined in `pkg/core/orm.go`. See the implementation for PostgreSQL in `pkg/postgres/postgres.go` as an example.

## License

This project is licensed under the [MIT License](LICENSE).

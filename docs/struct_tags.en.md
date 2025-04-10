# Struct Tags in NightORM

This document describes how to use struct tags in NightORM to customize the mapping between Go structures and database tables.

## Introduction

NightORM uses struct tags to determine how fields in a Go structure are mapped to columns in a database table. Struct tags are defined using Go's tag syntax:

```go
type User struct {
    ID        int       `db:"id,primary"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    Active    bool      `db:"active"`
}
```

## The `db` Tag

The `db` tag is the main tag used by NightORM. It can contain the column name and additional options, separated by commas.

### Column Name

The first value in the `db` tag is the column name in the database. If not specified, NightORM will use the lowercase field name.

```go
type User struct {
    // Mapped to the "id" column
    ID int `db:"id"`

    // Mapped to the "full_name" column
    Name string `db:"full_name"`

    // Mapped to the "email" column
    Email string `db:"email"`

    // Mapped to the "created_at" column
    CreatedAt time.Time `db:"created_at"`

    // Mapped to the "active" column
    Active bool `db:"active"`

    // Mapped to the "notag" column (lowercase field name)
    NoTag string
}
```

### The `primary` Option

The `primary` option indicates that the field is the primary key of the table. This is used by NightORM for operations like `FindByID`, `Update`, and `Delete`.

```go
type User struct {
    // Primary key
    ID int `db:"id,primary"`

    // Normal fields
    Name  string `db:"name"`
    Email string `db:"email"`
}
```

### Ignoring Fields

To ignore a field (not map it to a column), use `-` as the column name:

```go
type User struct {
    ID    int    `db:"id,primary"`
    Name  string `db:"name"`
    Email string `db:"email"`

    // This field will be ignored by NightORM
    Password string `db:"-"`

    // This field will also be ignored
    TempField string `db:"-"`
}
```

## Unexported Fields

Unexported fields (starting with lowercase letter) are automatically ignored by NightORM:

```go
type User struct {
    ID    int    `db:"id,primary"`
    Name  string `db:"name"`
    Email string `db:"email"`

    // This field will be ignored because it's unexported
    password string

    // This field will also be ignored
    tempField string
}
```

## Implementing the Required Interfaces

To use a structure with NightORM, you need to implement the `Model` or `ModelWithPrimaryKey` interface:

### The `Model` Interface

The `Model` interface requires implementing the `TableName()` method:

```go
// Model is the interface that all models must implement
type Model interface {
    // TableName returns the table name in the database
    TableName() string
}
```

Example implementation:

```go
// TableName returns the table name in the database
func (u *User) TableName() string {
    return "users"
}
```

### The `ModelWithPrimaryKey` Interface

The `ModelWithPrimaryKey` interface extends the `Model` interface and requires implementing the `PrimaryKey()` and `PrimaryKeyValue()` methods:

```go
// ModelWithPrimaryKey is an interface for models with a primary key
type ModelWithPrimaryKey interface {
    Model
    // PrimaryKey returns the primary key column name
    PrimaryKey() string
    // PrimaryKeyValue returns the primary key value
    PrimaryKeyValue() interface{}
}
```

Example implementation:

```go
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

## Complete Example

Here's a complete example of a structure with tags and implementation of the required interfaces:

```go
type User struct {
    ID        int       `db:"id,primary"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    Active    bool      `db:"active"`
    Password  string    `db:"-"` // Ignored
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

## Additional Considerations

### Nested Fields

NightORM does not directly support mapping nested fields. If you need to map nested fields, consider using simple fields or implementing custom methods for serialization/deserialization.

### Custom Types

NightORM supports custom types as long as they implement the necessary methods for conversion between Go and SQL (such as `Scan` and `Value` from the `sql.Scanner` and `driver.Valuer` interfaces).

### Calculated Fields

Calculated fields (which do not directly correspond to a column in the database) should be marked with `db:"-"` to prevent NightORM from trying to map them to columns.

## Conclusion

Struct tags are an important part of NightORM, allowing for flexible mapping between Go structures and database tables. Use them to customize how your models are mapped and to implement the necessary interfaces to work with NightORM.

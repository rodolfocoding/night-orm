package lib

import (
	"context"
	"fmt"
	"log"
	"time"

	night_orm "github.com/rodolfocoding/night-orm"
)

// User is an example model
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

// RunExample demonstrates basic usage of NightORM
func RunExample() {
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

	// Update a user
	fetchedUser.Name = "John Smith Updated"
	if err := orm.Update(ctx, fetchedUser); err != nil {
		log.Fatalf("Error updating user: %v", err)
	}
	fmt.Println("User updated successfully")

	// Find all users
	var users []*User
	if err := orm.FindAll(ctx, &User{}, &users); err != nil {
		log.Fatalf("Error finding all users: %v", err)
	}
	fmt.Printf("Found %d users\n", len(users))
	for _, u := range users {
		fmt.Printf("- %s (%s)\n", u.Name, u.Email)
	}

	// Execute a custom query
	rows, err := orm.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	fmt.Println("Active users:")
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("Error scanning results: %v", err)
		}
		fmt.Printf("- ID: %d, Name: %s\n", id, name)
	}

	// Use a transaction
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
	fmt.Println("Transaction committed successfully")

	// Delete a user
	if err := orm.Delete(ctx, fetchedUser); err != nil {
		log.Fatalf("Error deleting user: %v", err)
	}
	fmt.Println("User deleted successfully")
}

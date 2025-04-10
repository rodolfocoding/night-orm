package main

import (
	"context"
	"fmt"
	"log"
	"time"

	night_orm "night-orm"
)

// User é um exemplo de modelo
type User struct {
	ID        int       `db:"id,primary"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	Active    bool      `db:"active"`
}

// TableName retorna o nome da tabela no banco de dados
func (u *User) TableName() string {
	return "users"
}

// PrimaryKey retorna o nome da coluna de chave primária
func (u *User) PrimaryKey() string {
	return "id"
}

// PrimaryKeyValue retorna o valor da chave primária
func (u *User) PrimaryKeyValue() interface{} {
	return u.ID
}

func main() {
	// Conecta ao banco de dados PostgreSQL
	connectionString := "postgres://username:password@localhost:5432/database?sslmode=disable"
	ctx := context.Background()

	orm, err := night_orm.Connect(ctx, connectionString)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer orm.Close()

	// Cria um novo usuário
	user := &User{
		Name:      "João Silva",
		Email:     "joao@example.com",
		CreatedAt: time.Now(),
		Active:    true,
	}

	// Insere o usuário no banco de dados
	if err := orm.Create(ctx, user); err != nil {
		log.Fatalf("Erro ao criar usuário: %v", err)
	}
	fmt.Printf("Usuário criado com ID: %d\n", user.ID)

	// Busca um usuário pelo ID
	fetchedUser := &User{ID: user.ID}
	if err := orm.FindByID(ctx, fetchedUser, user.ID); err != nil {
		log.Fatalf("Erro ao buscar usuário: %v", err)
	}
	fmt.Printf("Usuário encontrado: %s (%s)\n", fetchedUser.Name, fetchedUser.Email)

	// Atualiza um usuário
	fetchedUser.Name = "João Silva Atualizado"
	if err := orm.Update(ctx, fetchedUser); err != nil {
		log.Fatalf("Erro ao atualizar usuário: %v", err)
	}
	fmt.Println("Usuário atualizado com sucesso")

	// Busca todos os usuários
	var users []*User
	if err := orm.FindAll(ctx, &User{}, &users); err != nil {
		log.Fatalf("Erro ao buscar todos os usuários: %v", err)
	}
	fmt.Printf("Encontrados %d usuários\n", len(users))
	for _, u := range users {
		fmt.Printf("- %s (%s)\n", u.Name, u.Email)
	}

	// Executa uma consulta personalizada
	rows, err := orm.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
	if err != nil {
		log.Fatalf("Erro ao executar consulta: %v", err)
	}
	defer rows.Close()

	fmt.Println("Usuários ativos:")
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("Erro ao fazer scan dos resultados: %v", err)
		}
		fmt.Printf("- ID: %d, Nome: %s\n", id, name)
	}

	// Usa uma transação
	tx, err := orm.Transaction(ctx)
	if err != nil {
		log.Fatalf("Erro ao iniciar transação: %v", err)
	}

	// Cria um novo usuário dentro da transação
	newUser := &User{
		Name:      "Maria Souza",
		Email:     "maria@example.com",
		CreatedAt: time.Now(),
		Active:    true,
	}

	if err := tx.Create(ctx, newUser); err != nil {
		tx.Rollback()
		log.Fatalf("Erro ao criar usuário na transação: %v", err)
	}

	// Confirma a transação
	if err := tx.Commit(); err != nil {
		log.Fatalf("Erro ao confirmar transação: %v", err)
	}
	fmt.Println("Transação confirmada com sucesso")

	// Exclui um usuário
	if err := orm.Delete(ctx, fetchedUser); err != nil {
		log.Fatalf("Erro ao excluir usuário: %v", err)
	}
	fmt.Println("Usuário excluído com sucesso")
}

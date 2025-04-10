# Transações no NightORM

Este documento descreve como usar transações no NightORM para garantir a integridade dos dados em operações que envolvem múltiplas alterações no banco de dados.

[English version](transactions.en.md)

## Introdução

Transações são usadas para agrupar várias operações de banco de dados em uma única unidade lógica de trabalho. Isso garante que todas as operações sejam concluídas com sucesso ou que nenhuma delas seja aplicada, mantendo a integridade dos dados.

O NightORM fornece suporte para transações através da interface `Transaction`, que permite executar operações CRUD dentro de uma transação.

## Iniciando uma Transação

Para iniciar uma transação, use o método `Transaction()` da interface `ORM`:

```go
tx, err := orm.Transaction(ctx)
if err != nil {
    log.Fatalf("Erro ao iniciar transação: %v", err)
}
```

## Operações dentro de uma Transação

Uma vez que você tenha uma transação, você pode executar operações CRUD dentro dela:

### Criar um Registro

```go
user := &User{
    Name:      "João Silva",
    Email:     "joao@example.com",
    CreatedAt: time.Now(),
    Active:    true,
}

if err := tx.Create(ctx, user); err != nil {
    tx.Rollback()
    log.Fatalf("Erro ao criar usuário: %v", err)
}
```

### Atualizar um Registro

```go
user.Name = "João Silva Atualizado"
if err := tx.Update(ctx, user); err != nil {
    tx.Rollback()
    log.Fatalf("Erro ao atualizar usuário: %v", err)
}
```

### Excluir um Registro

```go
if err := tx.Delete(ctx, user); err != nil {
    tx.Rollback()
    log.Fatalf("Erro ao excluir usuário: %v", err)
}
```

### Consultas Personalizadas

Você também pode executar consultas SQL personalizadas dentro de uma transação:

```go
rows, err := tx.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
if err != nil {
    tx.Rollback()
    log.Fatalf("Erro ao executar consulta: %v", err)
}
defer rows.Close()

// Processa os resultados...
```

```go
result, err := tx.Exec(ctx, "UPDATE users SET active = $1 WHERE id = $2", false, 1)
if err != nil {
    tx.Rollback()
    log.Fatalf("Erro ao executar comando: %v", err)
}

// Verifica o resultado...
```

## Confirmando ou Revertendo uma Transação

Após executar todas as operações necessárias, você precisa confirmar ou reverter a transação:

### Confirmando uma Transação

Para confirmar uma transação (aplicar todas as alterações), use o método `Commit()`:

```go
if err := tx.Commit(); err != nil {
    log.Fatalf("Erro ao confirmar transação: %v", err)
}
```

### Revertendo uma Transação

Para reverter uma transação (descartar todas as alterações), use o método `Rollback()`:

```go
if err := tx.Rollback(); err != nil {
    log.Fatalf("Erro ao reverter transação: %v", err)
}
```

## Padrão de Uso

Um padrão comum para usar transações é:

```go
// Inicia uma transação
tx, err := orm.Transaction(ctx)
if err != nil {
    log.Fatalf("Erro ao iniciar transação: %v", err)
}

// Função de limpeza para garantir que a transação seja revertida em caso de erro
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
        panic(r) // Re-panic após o rollback
    }
}()

// Executa operações dentro da transação
if err := tx.Create(ctx, user1); err != nil {
    tx.Rollback()
    log.Fatalf("Erro ao criar usuário 1: %v", err)
}

if err := tx.Create(ctx, user2); err != nil {
    tx.Rollback()
    log.Fatalf("Erro ao criar usuário 2: %v", err)
}

// Confirma a transação
if err := tx.Commit(); err != nil {
    log.Fatalf("Erro ao confirmar transação: %v", err)
}
```

## Exemplo Completo

Aqui está um exemplo completo de como usar transações no NightORM:

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/seu-usuario/night-orm"
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
    // Conecta ao banco de dados
    connectionString := "postgres://username:password@localhost:5432/database?sslmode=disable"
    ctx := context.Background()

    orm, err := night_orm.Connect(ctx, connectionString)
    if err != nil {
        log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
    }
    defer orm.Close()

    // Inicia uma transação
    tx, err := orm.Transaction(ctx)
    if err != nil {
        log.Fatalf("Erro ao iniciar transação: %v", err)
    }

    // Cria dois usuários dentro da transação
    user1 := &User{
        Name:      "João Silva",
        Email:     "joao@example.com",
        CreatedAt: time.Now(),
        Active:    true,
    }

    user2 := &User{
        Name:      "Maria Souza",
        Email:     "maria@example.com",
        CreatedAt: time.Now(),
        Active:    true,
    }

    // Tenta criar o primeiro usuário
    if err := tx.Create(ctx, user1); err != nil {
        tx.Rollback()
        log.Fatalf("Erro ao criar usuário 1: %v", err)
    }

    // Tenta criar o segundo usuário
    if err := tx.Create(ctx, user2); err != nil {
        tx.Rollback()
        log.Fatalf("Erro ao criar usuário 2: %v", err)
    }

    // Confirma a transação
    if err := tx.Commit(); err != nil {
        log.Fatalf("Erro ao confirmar transação: %v", err)
    }

    log.Println("Transação concluída com sucesso!")
    log.Printf("Usuário 1 criado com ID: %d\n", user1.ID)
    log.Printf("Usuário 2 criado com ID: %d\n", user2.ID)
}
```

## Considerações Importantes

### Tratamento de Erros

É importante sempre verificar os erros retornados pelas operações dentro de uma transação e chamar `Rollback()` em caso de erro.

### Fechamento de Recursos

Certifique-se de fechar recursos como `*sql.Rows` mesmo em caso de erro na transação.

### Contexto

Use o mesmo contexto (`context.Context`) para todas as operações dentro de uma transação para garantir consistência.

### Isolamento

O nível de isolamento da transação depende do banco de dados e da configuração. Consulte a documentação do seu banco de dados para mais informações.

## Conclusão

Transações são uma ferramenta poderosa para garantir a integridade dos dados em operações que envolvem múltiplas alterações no banco de dados. O NightORM fornece uma interface simples e intuitiva para trabalhar com transações, permitindo que você execute operações CRUD e consultas personalizadas dentro de uma transação.

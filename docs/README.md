<div align="center">
  <img src="../pkg/assets/Brand.svg" alt="NightORM Logo" width="500">
</div>

# Documentação do NightORM

Bem-vindo à documentação do NightORM, um ORM (Object-Relational Mapping) simples e flexível para Go.

[English version](README.en.md)

## Índice

- [Introdução](#introdução)
- [Guias](#guias)
- [Referência](#referência)
- [Exemplos](#exemplos)

## Introdução

O NightORM é um ORM para Go que facilita a interação com bancos de dados relacionais. Ele fornece uma interface simples e intuitiva para operações CRUD (Create, Read, Update, Delete) e suporta transações.

Atualmente, o NightORM suporta PostgreSQL, com planos para expandir para outros bancos de dados no futuro.

## Guias

- [Tags de Estrutura](struct_tags.md) - Como usar tags de estrutura para personalizar o mapeamento entre estruturas Go e tabelas de banco de dados.
- [Transações](transactions.md) - Como usar transações para garantir a integridade dos dados em operações que envolvem múltiplas alterações no banco de dados.
- [Suporte a Bancos de Dados](database_support.md) - Informações sobre os bancos de dados suportados e como adicionar suporte para novos bancos de dados.

## Referência

### Interfaces Principais

- `ORM` - Interface principal que define as operações básicas do ORM.
- `Model` - Interface que todos os modelos devem implementar.
- `ModelWithPrimaryKey` - Interface para modelos com chave primária.
- `Transaction` - Interface que representa uma transação de banco de dados.

### Pacotes

- `pkg/core` - Interfaces e tipos principais do ORM.
- `pkg/postgres` - Implementação do ORM para PostgreSQL.
- `pkg/utils` - Utilitários para reflexão e construção de consultas SQL.

## Exemplos

Veja o diretório [examples](../examples) para exemplos de uso do NightORM.

### Exemplo Básico

```go
package main

import (
    "context"
    "fmt"
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
}
```

## Contribuindo

Veja o arquivo [CONTRIBUTING.md](../CONTRIBUTING.md) para informações sobre como contribuir para o projeto.

## Licença

Este projeto está licenciado sob a [Licença MIT](../LICENSE).

# Suporte a Bancos de Dados no NightORM

Este documento descreve o suporte atual a bancos de dados no NightORM e como adicionar suporte para novos bancos de dados.

[English version](database_support.en.md)

## Bancos de Dados Suportados

Atualmente, o NightORM suporta os seguintes bancos de dados:

### PostgreSQL

O PostgreSQL é o primeiro banco de dados suportado pelo NightORM. A implementação está no pacote `pkg/postgres`.

Para usar o NightORM com PostgreSQL:

```go
import (
    "context"
    "log"

    "github.com/seu-usuario/night-orm"
)

func main() {
    // String de conexão para PostgreSQL
    connectionString := "postgres://username:password@localhost:5432/database?sslmode=disable"
    ctx := context.Background()

    // Conecta ao PostgreSQL
    orm, err := night_orm.Connect(ctx, connectionString)
    if err != nil {
        log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
    }
    defer orm.Close()

    // Agora você pode usar o ORM para interagir com o PostgreSQL
}
```

## Adicionando Suporte para Novos Bancos de Dados

O NightORM foi projetado para ser facilmente extensível para suportar diferentes bancos de dados. Para adicionar suporte para um novo banco de dados, siga os passos abaixo:

### 1. Crie um Novo Pacote

Crie um novo pacote dentro do diretório `pkg/` com o nome do banco de dados. Por exemplo, para adicionar suporte para MySQL, crie o diretório `pkg/mysql/`.

### 2. Implemente a Interface ORM

Implemente a interface `ORM` definida em `pkg/core/orm.go`. Você precisará implementar todos os métodos definidos na interface.

Exemplo de estrutura para MySQL:

```go
package mysql

import (
    "context"
    "database/sql"

    _ "github.com/go-sql-driver/mysql" // Driver MySQL
    "github.com/seu-usuario/night-orm/pkg/core"
)

// MySQLORM é a implementação do ORM para MySQL
type MySQLORM struct {
    db *sql.DB
}

// NewMySQLORM cria uma nova instância do ORM para MySQL
func NewMySQLORM() *MySQLORM {
    return &MySQLORM{}
}

// Connect estabelece uma conexão com o banco de dados MySQL
func (m *MySQLORM) Connect(ctx context.Context, connectionString string) error {
    // Implementação da conexão com MySQL
}

// Implemente os outros métodos da interface ORM...
```

### 3. Adicione Funções de Fábrica

Adicione funções de fábrica no arquivo principal `night_orm.go` para facilitar a criação de instâncias do ORM para o novo banco de dados.

```go
// NewMySQLORM cria uma nova instância do ORM para MySQL
func NewMySQLORM() ORM {
    return mysql.NewMySQLORM()
}

// ConnectMySQL é uma função auxiliar para conectar ao banco de dados MySQL
func ConnectMySQL(ctx context.Context, connectionString string) (ORM, error) {
    orm := NewMySQLORM()
    err := orm.Connect(ctx, connectionString)
    if err != nil {
        return nil, err
    }
    return orm, nil
}
```

### 4. Adicione Testes

Adicione testes para a nova implementação para garantir que ela funciona corretamente.

### 5. Atualize a Documentação

Atualize a documentação para incluir informações sobre o novo banco de dados suportado.

## Considerações para Diferentes Bancos de Dados

Ao implementar suporte para diferentes bancos de dados, considere as seguintes diferenças:

### Sintaxe SQL

Diferentes bancos de dados podem ter sintaxes SQL ligeiramente diferentes. Por exemplo, o PostgreSQL usa `$1`, `$2`, etc. para parâmetros, enquanto o MySQL usa `?`.

### Tipos de Dados

Os tipos de dados podem variar entre bancos de dados. Certifique-se de mapear corretamente os tipos de dados do Go para os tipos de dados do banco de dados.

### Funcionalidades Específicas

Alguns bancos de dados têm funcionalidades específicas que podem ser úteis para o ORM. Por exemplo, o PostgreSQL tem o operador `RETURNING` que permite retornar valores de linhas afetadas por uma operação de inserção, atualização ou exclusão.

### Tratamento de Erros

Diferentes drivers de banco de dados podem retornar erros de maneiras diferentes. Certifique-se de tratar corretamente os erros específicos de cada banco de dados.

## Bancos de Dados Planejados

Os seguintes bancos de dados estão planejados para suporte futuro:

- MySQL
- SQLite
- Microsoft SQL Server
- Oracle

Se você estiver interessado em contribuir com suporte para algum desses bancos de dados, consulte o arquivo [CONTRIBUTING.md](../CONTRIBUTING.md) para obter informações sobre como contribuir para o projeto.

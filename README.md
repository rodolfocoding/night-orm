# NightORM

NightORM é um ORM (Object-Relational Mapping) simples e flexível para Go, projetado para facilitar a interação com bancos de dados relacionais. Atualmente, o NightORM suporta PostgreSQL, com planos para expandir para outros bancos de dados no futuro.

## Características

- Interface simples e intuitiva
- Suporte para operações CRUD básicas
- Mapeamento automático entre estruturas Go e tabelas de banco de dados
- Suporte para transações
- Consultas SQL personalizadas
- Suporte para tags de estrutura para personalizar o mapeamento

## Instalação

```bash
go get github.com/seu-usuario/night-orm
```

## Uso Básico

### Definindo um Modelo

```go
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
```

### Conectando ao Banco de Dados

```go
import (
    "context"
    "log"

    "github.com/seu-usuario/night-orm"
)

func main() {
    // Conecta ao banco de dados PostgreSQL
    connectionString := "postgres://username:password@localhost:5432/database?sslmode=disable"
    ctx := context.Background()

    orm, err := night_orm.Connect(ctx, connectionString)
    if err != nil {
        log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
    }
    defer orm.Close()

    // Agora você pode usar o ORM para interagir com o banco de dados
}
```

### Operações CRUD

#### Criar um Registro

```go
user := &User{
    Name:      "João Silva",
    Email:     "joao@example.com",
    CreatedAt: time.Now(),
    Active:    true,
}

if err := orm.Create(ctx, user); err != nil {
    log.Fatalf("Erro ao criar usuário: %v", err)
}
```

#### Buscar um Registro pelo ID

```go
user := &User{ID: 1}
if err := orm.FindByID(ctx, user, 1); err != nil {
    log.Fatalf("Erro ao buscar usuário: %v", err)
}
```

#### Buscar Todos os Registros

```go
var users []*User
if err := orm.FindAll(ctx, &User{}, &users); err != nil {
    log.Fatalf("Erro ao buscar todos os usuários: %v", err)
}
```

#### Atualizar um Registro

```go
user.Name = "João Silva Atualizado"
if err := orm.Update(ctx, user); err != nil {
    log.Fatalf("Erro ao atualizar usuário: %v", err)
}
```

#### Excluir um Registro

```go
if err := orm.Delete(ctx, user); err != nil {
    log.Fatalf("Erro ao excluir usuário: %v", err)
}
```

### Consultas Personalizadas

```go
rows, err := orm.Query(ctx, "SELECT id, name FROM users WHERE active = $1", true)
if err != nil {
    log.Fatalf("Erro ao executar consulta: %v", err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil {
        log.Fatalf("Erro ao fazer scan dos resultados: %v", err)
    }
    fmt.Printf("ID: %d, Nome: %s\n", id, name)
}
```

### Transações

```go
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
```

## Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

### Adicionando Suporte para Novos Bancos de Dados

Para adicionar suporte para um novo banco de dados, você precisa implementar a interface `ORM` definida em `pkg/core/orm.go`. Veja a implementação para PostgreSQL em `pkg/postgres/postgres.go` como exemplo.

## Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).

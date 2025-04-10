# Tags de Estrutura no NightORM

Este documento descreve como usar as tags de estrutura no NightORM para personalizar o mapeamento entre estruturas Go e tabelas de banco de dados.

## Introdução

O NightORM usa tags de estrutura para determinar como os campos de uma estrutura Go são mapeados para colunas em uma tabela de banco de dados. As tags de estrutura são definidas usando a sintaxe de tags do Go:

```go
type User struct {
    ID        int       `db:"id,primary"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    Active    bool      `db:"active"`
}
```

## Tag `db`

A tag `db` é a principal tag usada pelo NightORM. Ela pode conter o nome da coluna e opções adicionais, separados por vírgulas.

### Nome da Coluna

O primeiro valor na tag `db` é o nome da coluna no banco de dados. Se não for especificado, o NightORM usará o nome do campo em minúsculas.

```go
type User struct {
    // Mapeado para a coluna "id"
    ID int `db:"id"`

    // Mapeado para a coluna "full_name"
    Name string `db:"full_name"`

    // Mapeado para a coluna "email"
    Email string `db:"email"`

    // Mapeado para a coluna "created_at"
    CreatedAt time.Time `db:"created_at"`

    // Mapeado para a coluna "active"
    Active bool `db:"active"`

    // Mapeado para a coluna "notag" (nome do campo em minúsculas)
    NoTag string
}
```

### Opção `primary`

A opção `primary` indica que o campo é a chave primária da tabela. Isso é usado pelo NightORM para operações como `FindByID`, `Update` e `Delete`.

```go
type User struct {
    // Chave primária
    ID int `db:"id,primary"`

    // Campos normais
    Name  string `db:"name"`
    Email string `db:"email"`
}
```

### Ignorando Campos

Para ignorar um campo (não mapeá-lo para uma coluna), use `-` como nome da coluna:

```go
type User struct {
    ID    int    `db:"id,primary"`
    Name  string `db:"name"`
    Email string `db:"email"`

    // Este campo será ignorado pelo NightORM
    Password string `db:"-"`

    // Este campo também será ignorado
    TempField string `db:"-"`
}
```

## Campos Não Exportados

Campos não exportados (começando com letra minúscula) são automaticamente ignorados pelo NightORM:

```go
type User struct {
    ID    int    `db:"id,primary"`
    Name  string `db:"name"`
    Email string `db:"email"`

    // Este campo será ignorado por ser não exportado
    password string

    // Este campo também será ignorado
    tempField string
}
```

## Implementando as Interfaces Necessárias

Para usar uma estrutura com o NightORM, você precisa implementar a interface `Model` ou `ModelWithPrimaryKey`:

### Interface `Model`

A interface `Model` requer a implementação do método `TableName()`:

```go
// Model é a interface que todos os modelos devem implementar
type Model interface {
    // TableName retorna o nome da tabela no banco de dados
    TableName() string
}
```

Exemplo de implementação:

```go
// TableName retorna o nome da tabela no banco de dados
func (u *User) TableName() string {
    return "users"
}
```

### Interface `ModelWithPrimaryKey`

A interface `ModelWithPrimaryKey` estende a interface `Model` e requer a implementação dos métodos `PrimaryKey()` e `PrimaryKeyValue()`:

```go
// ModelWithPrimaryKey é uma interface para modelos com chave primária
type ModelWithPrimaryKey interface {
    Model
    // PrimaryKey retorna o nome da coluna de chave primária
    PrimaryKey() string
    // PrimaryKeyValue retorna o valor da chave primária
    PrimaryKeyValue() interface{}
}
```

Exemplo de implementação:

```go
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

## Exemplo Completo

Aqui está um exemplo completo de uma estrutura com tags e implementação das interfaces necessárias:

```go
type User struct {
    ID        int       `db:"id,primary"`
    Name      string    `db:"name"`
    Email     string    `db:"email"`
    CreatedAt time.Time `db:"created_at"`
    Active    bool      `db:"active"`
    Password  string    `db:"-"` // Ignorado
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

## Considerações Adicionais

### Campos Aninhados

O NightORM não suporta diretamente o mapeamento de campos aninhados. Se você precisar mapear campos aninhados, considere usar campos simples ou implementar métodos personalizados para serialização/desserialização.

### Tipos Personalizados

O NightORM suporta tipos personalizados, desde que eles implementem os métodos necessários para conversão entre Go e SQL (como `Scan` e `Value` da interface `sql.Scanner` e `driver.Valuer`).

### Campos Calculados

Campos calculados (que não correspondem diretamente a uma coluna no banco de dados) devem ser marcados com `db:"-"` para evitar que o NightORM tente mapeá-los para colunas.

## Conclusão

As tags de estrutura são uma parte importante do NightORM, permitindo um mapeamento flexível entre estruturas Go e tabelas de banco de dados. Use-as para personalizar como seus modelos são mapeados e para implementar as interfaces necessárias para trabalhar com o NightORM.

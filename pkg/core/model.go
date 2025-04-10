package core

// Model é a interface que todos os modelos devem implementar
type Model interface {
	// TableName retorna o nome da tabela no banco de dados
	TableName() string
}

// ModelWithPrimaryKey é uma interface para modelos com chave primária
type ModelWithPrimaryKey interface {
	Model
	// PrimaryKey retorna o nome da coluna de chave primária
	PrimaryKey() string
	// PrimaryKeyValue retorna o valor da chave primária
	PrimaryKeyValue() interface{}
}

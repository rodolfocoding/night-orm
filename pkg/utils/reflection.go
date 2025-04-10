package utils

import (
	"errors"
	"reflect"
	"strings"
)

// GetStructFields retorna um mapa de campos da estrutura com seus nomes e valores
func GetStructFields(obj interface{}) (map[string]interface{}, error) {
	if obj == nil {
		return nil, errors.New("objeto não pode ser nil")
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, errors.New("objeto deve ser uma estrutura ou um ponteiro para uma estrutura")
	}

	fields := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Ignora campos não exportados
		if !fieldType.IsExported() {
			continue
		}

		// Verifica a tag "db" para o nome da coluna
		tag := fieldType.Tag.Get("db")
		if tag == "-" {
			continue // Ignora campos marcados com db:"-"
		}

		// Se não houver tag, usa o nome do campo em minúsculas
		columnName := tag
		if columnName == "" {
			columnName = strings.ToLower(fieldType.Name)
		}

		// Extrai o valor do campo
		var fieldValue interface{}
		if field.CanInterface() {
			fieldValue = field.Interface()
		}

		fields[columnName] = fieldValue
	}

	return fields, nil
}

// SetStructField define o valor de um campo em uma estrutura
func SetStructField(obj interface{}, fieldName string, value interface{}) error {
	if obj == nil {
		return errors.New("objeto não pode ser nil")
	}

	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("objeto deve ser um ponteiro não-nil")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return errors.New("objeto deve ser um ponteiro para uma estrutura")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Verifica se o campo corresponde ao nome fornecido
		tag := fieldType.Tag.Get("db")
		if tag == fieldName || (tag == "" && strings.ToLower(fieldType.Name) == fieldName) {
			if !field.CanSet() {
				return errors.New("campo não pode ser definido")
			}

			fieldVal := reflect.ValueOf(value)
			if field.Type() != fieldVal.Type() {
				// Tenta converter o valor para o tipo do campo
				if fieldVal.Type().ConvertibleTo(field.Type()) {
					fieldVal = fieldVal.Convert(field.Type())
				} else {
					return errors.New("tipo de valor incompatível com o tipo do campo")
				}
			}

			field.Set(fieldVal)
			return nil
		}
	}

	return errors.New("campo não encontrado")
}

// GetTagName obtém o nome da tag de um campo
func GetTagName(structType reflect.Type, fieldName, tagName string) string {
	field, ok := structType.FieldByName(fieldName)
	if !ok {
		return ""
	}
	
	tag := field.Tag.Get(tagName)
	if tag == "" {
		return ""
	}
	
	parts := strings.Split(tag, ",")
	return parts[0]
}

// GetPrimaryKeyField retorna o nome e o valor do campo marcado como chave primária
func GetPrimaryKeyField(obj interface{}) (string, interface{}, error) {
	if obj == nil {
		return "", nil, errors.New("objeto não pode ser nil")
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", nil, errors.New("objeto deve ser uma estrutura ou um ponteiro para uma estrutura")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Verifica a tag "db" para identificar a chave primária
		tag := fieldType.Tag.Get("db")
		if strings.Contains(tag, ",primary") {
			// Extrai o nome da coluna da tag
			columnName := strings.Split(tag, ",")[0]
			if columnName == "" {
				columnName = strings.ToLower(fieldType.Name)
			}

			// Extrai o valor do campo
			var fieldValue interface{}
			if field.CanInterface() {
				fieldValue = field.Interface()
			}

			return columnName, fieldValue, nil
		}
	}

	// Se não encontrar uma tag de chave primária, procura por um campo chamado "ID" ou "Id"
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if strings.ToLower(fieldType.Name) == "id" {
			tag := fieldType.Tag.Get("db")
			columnName := tag
			if columnName == "" || columnName == "-" {
				columnName = "id"
			}

			var fieldValue interface{}
			if field.CanInterface() {
				fieldValue = field.Interface()
			}

			return columnName, fieldValue, nil
		}
	}

	return "", nil, errors.New("chave primária não encontrada")
}

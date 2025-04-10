package utils

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	ID        int    `db:"id,primary"`
	Name      string `db:"name"`
	Email     string `db:"email"`
	Ignored   string `db:"-"`
	NoTag     string
	unexported string
}

func TestGetStructFields(t *testing.T) {
	testStruct := TestStruct{
		ID:         1,
		Name:       "Test Name",
		Email:      "test@example.com",
		Ignored:    "This should be ignored",
		NoTag:      "No tag field",
		unexported: "Unexported field",
	}

	fields, err := GetStructFields(testStruct)
	if err != nil {
		t.Fatalf("GetStructFields returned error: %v", err)
	}

	// Verifica se o número de campos está correto (ID, Name, Email, NoTag)
	if len(fields) != 4 {
		t.Errorf("Expected 4 fields, got %d", len(fields))
	}

	// Verifica se os campos estão corretos
	if id, ok := fields["id"]; !ok || id != 1 {
		t.Errorf("Expected field 'id' with value 1, got %v", id)
	}

	if name, ok := fields["name"]; !ok || name != "Test Name" {
		t.Errorf("Expected field 'name' with value 'Test Name', got %v", name)
	}

	if email, ok := fields["email"]; !ok || email != "test@example.com" {
		t.Errorf("Expected field 'email' with value 'test@example.com', got %v", email)
	}

	if notag, ok := fields["notag"]; !ok || notag != "No tag field" {
		t.Errorf("Expected field 'notag' with value 'No tag field', got %v", notag)
	}

	// Verifica se os campos ignorados não estão presentes
	if _, ok := fields["ignored"]; ok {
		t.Errorf("Field 'ignored' should not be present")
	}

	if _, ok := fields["unexported"]; ok {
		t.Errorf("Field 'unexported' should not be present")
	}
}

func TestSetStructField(t *testing.T) {
	testStruct := &TestStruct{
		ID:   1,
		Name: "Original Name",
	}

	// Testa definir um campo com tag
	err := SetStructField(testStruct, "name", "Updated Name")
	if err != nil {
		t.Fatalf("SetStructField returned error: %v", err)
	}

	if testStruct.Name != "Updated Name" {
		t.Errorf("Expected Name to be 'Updated Name', got '%s'", testStruct.Name)
	}

	// Testa definir um campo sem tag
	err = SetStructField(testStruct, "notag", "Updated NoTag")
	if err != nil {
		t.Fatalf("SetStructField returned error: %v", err)
	}

	if testStruct.NoTag != "Updated NoTag" {
		t.Errorf("Expected NoTag to be 'Updated NoTag', got '%s'", testStruct.NoTag)
	}

	// Testa definir um campo que não existe
	err = SetStructField(testStruct, "nonexistent", "value")
	if err == nil {
		t.Errorf("Expected error when setting nonexistent field, got nil")
	}
}

func TestGetPrimaryKeyField(t *testing.T) {
	testStruct := TestStruct{
		ID:   42,
		Name: "Test Name",
	}

	// Testa obter a chave primária
	name, value, err := GetPrimaryKeyField(testStruct)
	if err != nil {
		t.Fatalf("GetPrimaryKeyField returned error: %v", err)
	}

	if name != "id" {
		t.Errorf("Expected primary key name to be 'id', got '%s'", name)
	}

	if value != 42 {
		t.Errorf("Expected primary key value to be 42, got %v", value)
	}

	// Testa com uma estrutura sem chave primária
	type NoKeyStruct struct {
		Name string `db:"name"`
	}

	noKeyStruct := NoKeyStruct{
		Name: "No Key",
	}

	_, _, err = GetPrimaryKeyField(noKeyStruct)
	if err == nil {
		t.Errorf("Expected error when getting primary key from struct without primary key, got nil")
	}
}

func TestGetTagName(t *testing.T) {
	typ := reflect.TypeOf(TestStruct{})

	// Testa obter o nome da tag para um campo com tag
	tagName := GetTagName(typ, "Name", "db")
	if tagName != "name" {
		t.Errorf("Expected tag name to be 'name', got '%s'", tagName)
	}

	// Testa obter o nome da tag para um campo com tag e opções
	tagName = GetTagName(typ, "ID", "db")
	if tagName != "id" {
		t.Errorf("Expected tag name to be 'id', got '%s'", tagName)
	}

	// Testa obter o nome da tag para um campo sem tag
	tagName = GetTagName(typ, "NoTag", "db")
	if tagName != "" {
		t.Errorf("Expected tag name to be empty, got '%s'", tagName)
	}

	// Testa obter o nome da tag para um campo que não existe
	tagName = GetTagName(typ, "NonExistent", "db")
	if tagName != "" {
		t.Errorf("Expected tag name to be empty, got '%s'", tagName)
	}
}

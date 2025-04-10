# Makefile para o NightORM

# Variáveis
GO = go
GOFLAGS = -v
GOFMT = gofmt
GOTEST = $(GO) test
GOVET = $(GO) vet
GOLINT = golint
GOCOVER = $(GO) tool cover
GOMOD = $(GO) mod
GOBUILD = $(GO) build
GORUN = $(GO) run

# Diretórios
PKG_LIST = ./...
EXAMPLES_DIR = ./examples

# Targets
.PHONY: all build test fmt lint vet clean cover tidy examples help

all: fmt lint vet test build

# Compila o projeto
build:
	@echo "Compilando o projeto..."
	$(GOBUILD) $(GOFLAGS) $(PKG_LIST)

# Executa os testes
test:
	@echo "Executando testes..."
	$(GOTEST) $(GOFLAGS) $(PKG_LIST)

# Formata o código
fmt:
	@echo "Formatando o código..."
	$(GOFMT) -s -w .

# Executa o linter
lint:
	@echo "Executando linter..."
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) $(PKG_LIST); \
	else \
		echo "golint não está instalado. Execute: go install golang.org/x/lint/golint@latest"; \
	fi

# Executa o vet
vet:
	@echo "Executando vet..."
	$(GOVET) $(PKG_LIST)

# Limpa os arquivos gerados
clean:
	@echo "Limpando arquivos gerados..."
	@rm -rf ./bin
	@rm -rf ./vendor
	@rm -f ./coverage.out

# Executa os testes com cobertura
cover:
	@echo "Executando testes com cobertura..."
	$(GOTEST) -coverprofile=coverage.out $(PKG_LIST)
	$(GOCOVER) -html=coverage.out

# Atualiza as dependências
tidy:
	@echo "Atualizando dependências..."
	$(GOMOD) tidy

# Executa os exemplos
examples:
	@echo "Executando exemplos..."
	@for example in $(wildcard $(EXAMPLES_DIR)/*.go); do \
		echo "Executando $$example..."; \
		$(GORUN) $$example; \
	done

# Exibe a ajuda
help:
	@echo "Makefile para o NightORM"
	@echo ""
	@echo "Uso:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all      Executa fmt, lint, vet, test e build"
	@echo "  build    Compila o projeto"
	@echo "  test     Executa os testes"
	@echo "  fmt      Formata o código"
	@echo "  lint     Executa o linter"
	@echo "  vet      Executa o vet"
	@echo "  clean    Limpa os arquivos gerados"
	@echo "  cover    Executa os testes com cobertura"
	@echo "  tidy     Atualiza as dependências"
	@echo "  examples Executa os exemplos"
	@echo "  help     Exibe esta ajuda"

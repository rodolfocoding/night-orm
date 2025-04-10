# Guia de Publicação do NightORM

Este documento descreve como publicar o NightORM como um módulo Go.

## Preparação

1. Certifique-se de que o código está pronto para ser publicado:

   - Todos os testes passam
   - A documentação está atualizada
   - O código segue as convenções de estilo do Go

2. Atualize o arquivo `go.mod` para usar o caminho correto do repositório:

```bash
go mod edit -module github.com/seu-usuario/night-orm
```

3. Atualize todas as importações no código para usar o novo caminho do módulo:

```bash
find . -type f -name "*.go" -exec sed -i 's|night-orm|github.com/seu-usuario/night-orm|g' {} \;
```

## Publicação

1. Crie um repositório no GitHub com o nome `night-orm`.

2. Inicialize o repositório Git localmente e faça o primeiro commit:

```bash
git init
git add .
git commit -m "Versão inicial do NightORM"
```

3. Adicione o repositório remoto e envie o código:

```bash
git remote add origin https://github.com/seu-usuario/night-orm.git
git push -u origin main
```

4. Crie uma tag para a versão inicial:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Uso do Módulo

Após a publicação, os usuários podem instalar o NightORM usando:

```bash
go get github.com/seu-usuario/night-orm@v0.1.0
```

## Atualizações

Para publicar novas versões:

1. Faça as alterações necessárias no código.
2. Atualize a versão no arquivo `go.mod` se necessário.
3. Commit e push das alterações.
4. Crie uma nova tag para a versão:

```bash
git tag v0.2.0
git push origin v0.2.0
```

## Versionamento Semântico

O NightORM segue o [Versionamento Semântico](https://semver.org/):

- **MAJOR** (x.0.0): Alterações incompatíveis com versões anteriores
- **MINOR** (0.x.0): Adição de funcionalidades de forma compatível
- **PATCH** (0.0.x): Correções de bugs de forma compatível

## Notas Adicionais

- Certifique-se de que o repositório é público para que os usuários possam acessá-lo.
- Mantenha o CHANGELOG.md atualizado com as alterações em cada versão.
- Considere configurar integração contínua (CI) para executar testes automaticamente.

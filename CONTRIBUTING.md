# Contribuindo para o NightORM

Obrigado pelo seu interesse em contribuir para o NightORM! Este documento fornece diretrizes para contribuir com o projeto.

## Código de Conduta

Por favor, seja respeitoso e construtivo em todas as interações com o projeto. Valorizamos a diversidade e inclusão na nossa comunidade.

## Como Contribuir

### Reportando Bugs

Se você encontrar um bug, por favor, crie uma issue no GitHub com as seguintes informações:

- Título claro e descritivo
- Passos detalhados para reproduzir o bug
- Comportamento esperado e comportamento atual
- Versão do Go e do NightORM que você está usando
- Qualquer informação adicional relevante

### Sugerindo Melhorias

Para sugerir melhorias, crie uma issue no GitHub descrevendo:

- O que você gostaria de ver melhorado
- Por que essa melhoria seria útil
- Como você imagina que essa melhoria poderia ser implementada

### Enviando Pull Requests

1. Faça um fork do repositório
2. Clone o seu fork: `git clone https://github.com/seu-usuario/night-orm.git`
3. Crie uma branch para sua contribuição: `git checkout -b minha-contribuicao`
4. Faça suas alterações
5. Execute os testes: `go test ./...`
6. Commit suas alterações: `git commit -m "Descrição da alteração"`
7. Push para o seu fork: `git push origin minha-contribuicao`
8. Crie um Pull Request para o repositório original

### Diretrizes para Código

- Siga as convenções de estilo do Go
- Escreva testes para novas funcionalidades
- Mantenha a documentação atualizada
- Use comentários claros e úteis
- Mantenha o código simples e legível

## Adicionando Suporte para Novos Bancos de Dados

Para adicionar suporte para um novo banco de dados:

1. Crie um novo pacote em `pkg/` com o nome do banco de dados (ex: `pkg/mysql/`)
2. Implemente a interface `ORM` definida em `pkg/core/orm.go`
3. Adicione funções de fábrica no arquivo principal `night_orm.go`
4. Adicione testes para a nova implementação
5. Atualize a documentação para incluir o novo banco de dados

## Processo de Desenvolvimento

1. As issues são triadas e priorizadas pelos mantenedores
2. As contribuições são revisadas pelos mantenedores
3. Após a aprovação, as contribuições são mescladas na branch principal
4. Novas versões são lançadas periodicamente seguindo o versionamento semântico

## Comunicação

- Use as issues do GitHub para discussões relacionadas ao projeto
- Para questões mais complexas, considere abrir uma discussão no GitHub Discussions

## Agradecimentos

Agradecemos a todos os contribuidores que ajudam a melhorar o NightORM!

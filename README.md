# FINANCE MANAGER

## Escopo

### MVP

- Autenticação com google
- Conectar com open finance do seu banco
- Tela Home:
  - Balanço, receitas e despesas totais
  - Listar bancos conectados, com botão + no final para adicionar um novo
  - Listar cartões de créditos conectados, com botão + no final para adicionar um novo
  - Gráfico de pizza para despesas por categoria
  - Orçamento (progress bar com total e restante (verde até 80%, amarela até 100% e vermelha ao exceder))
- Tela transações:
  - Listar transações
  - Pesquisar
  - Filtrar por categoria, mês, banco, método de pagamento, despesas e receitas
  - Header com overview (balanço e soma de receitas e despesas levando em consideração o filtro/pesquisa)
  - Detalhe da transação (Banco, Nome, Descrição, Valor, Data, Categoria)
  - Alterar categoria da transação
  - Ignorar transação (não incluí-la nas somas de receita e despesas)
- Tela de orçamento (limite de gastos referente ao mês atual):
  - Definir um gasto limite
  - Atualizar o orçamento
  - Deletar o orçamento
  - Definir orçamento por categoria (Opcional para o usuário, pode criar apenas o orçamento geral)
  - Listar informações gerais do orçamento (orçamento total, valor gasto, e recomendado por dia)
  - Orçamento por (progress bar de orçamento total e valor gasto)

### Backlog

- Enviar notificações de orçamento quando os gastos estiverem em 80% e >=100% do orçamento definido
- Transações customizadas: Criar, Atualizar e Deletar transações (criadas pelo usuário)
- Carteira (representa a carteira física, quantia que pode ser alterada pelo usuário e é adicionada no balanço total)
- Transações customizadas podem ser vinculadas a carteira
- Criar categorias customizadas
- CRUD Objetivos para acumular dinheiro, referente ao balanço total da sua conta (sugestão de viagem, comprar carro, reserva de emergência ou outro customizável)
- Tela de analytics com gráfico e insights
- Tela de investimentos

### Pluggy

- account.item_id -> consent.item_id -> item = institution

## Referências

### Mobills

- [Play Store](https://play.google.com/store/apps/details?id=br.com.gerenciadorfinanceiro.controller)
- [Apple Store](https://apps.apple.com/us/app/mobills-budget-planner/id921838244)

### Minhas Economias

- [Play Store](https://play.google.com/store/apps/details?id=com.minhaseconomias&hl=pt)
- [Apple Store](https://apps.apple.com/br/app/minhas-economias/id520869873)

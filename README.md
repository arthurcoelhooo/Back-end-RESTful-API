# API de Receitas

API em Go para cadastrar, listar, atualizar e excluir receitas em um banco PostgreSQL.

## Integrantes

- Arthur Coelho

## O que você precisa

- Go 1.22 ou mais novo
- PostgreSQL
- Terminal

Confira o Go:

```bash
go version
```

## 1. Preparar o banco

No PostgreSQL, crie um usuário e o banco:

```sql
CREATE USER receitas_app WITH PASSWORD 'SUA_SENHA';
CREATE DATABASE receitas OWNER receitas_app;
```

Conecte-se ao banco `receitas` e execute todo o conteúdo de [`schema.sql`](schema.sql). Isso cria as tabelas `receitas` e `api_keys`.

Se o PostgreSQL mostrar erro de permissão, execute:

```sql
GRANT USAGE ON SCHEMA public TO receitas_app;
GRANT ALL PRIVILEGES ON TABLE public.api_keys, public.receitas TO receitas_app;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO receitas_app;
```

## 2. Configurar a API

Abra a pasta do projeto no terminal.

**macOS ou Linux:**

```bash
export DATABASE_URL='postgres://receitas_app:SUA_SENHA@localhost:5432/receitas?sslmode=disable'
export PORT=8080
```

**Windows PowerShell:**

```powershell
$env:DATABASE_URL = "postgres://receitas_app:SUA_SENHA@localhost:5432/receitas?sslmode=disable"
$env:PORT = "8080"
```

Troque `SUA_SENHA` pela senha definida no PostgreSQL.

## 3. Iniciar o servidor

```bash
go mod tidy
go run ./cmd/api
```

Abra outro terminal e confirme que a API está funcionando:

```bash
curl http://localhost:8080/health
```

Resposta esperada:

```json
{"status":"ok"}
```

No Windows, use `curl.exe` caso `curl` não funcione.

Mantenha o servidor aberto enquanto usa a API. Para parar, pressione `Ctrl+C`.

## 4. Criar uma chave de API

Em outro terminal, configure novamente `DATABASE_URL` como no passo 2 e execute:

```bash
go run ./cmd/genkey -nome "minha-chave"
```

Copie a chave exibida. Ela não será mostrada novamente e deve ser enviada no cabeçalho `X-API-Key`.

Teste:

```bash
curl -H "X-API-Key: SUA_CHAVE" http://localhost:8080/receitas
```

## Usar a API

Nas URLs abaixo, use `SUA_CHAVE` no lugar da chave gerada.

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/health` | Verifica se a API está ativa |
| `GET` | `/receitas` | Lista todas as receitas |
| `GET` | `/receitas?categoria=doce` | Filtra por categoria |
| `POST` | `/receitas` | Cadastra uma receita |
| `GET` | `/receitas/{id}` | Busca uma receita |
| `PUT` | `/receitas/{id}` | Atualiza uma receita |
| `DELETE` | `/receitas/{id}` | Exclui uma receita |

O `/health` não usa chave. Todas as rotas `/receitas` usam `X-API-Key`.

### Criar uma receita

```bash
curl -X POST http://localhost:8080/receitas \
  -H "X-API-Key: SUA_CHAVE" \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Bolo de cenoura",
    "categoria": "doce",
    "ingredientes": "cenoura, ovos, farinha e açúcar",
    "modo_preparo": "Misture os ingredientes e asse.",
    "tempo_preparo_min": 60,
    "porcoes": 8
  }'
```

### Outros exemplos

Listar:

```bash
curl -H "X-API-Key: SUA_CHAVE" http://localhost:8080/receitas
```

Buscar pelo ID:

```bash
curl -H "X-API-Key: SUA_CHAVE" http://localhost:8080/receitas/1
```

Atualizar pelo ID:

```bash
curl -X PUT http://localhost:8080/receitas/1 \
  -H "X-API-Key: SUA_CHAVE" \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Bolo de cenoura com cobertura",
    "categoria": "doce",
    "ingredientes": "cenoura, ovos, farinha, açúcar e chocolate",
    "modo_preparo": "Misture, asse e cubra com chocolate.",
    "tempo_preparo_min": 70,
    "porcoes": 8
  }'
```

Excluir pelo ID:

```bash
curl -X DELETE \
  -H "X-API-Key: SUA_CHAVE" \
  http://localhost:8080/receitas/1
```

Categorias disponíveis: `doce`, `salgado`, `bebida`, `massa`, `sobremesa`, `lanche` e `salada`.

## Testar no VS Code

1. Instale a extensão **REST Client**.
2. Copie `.env.example` para `.env` e coloque a chave gerada em `API_KEY`.
3. Abra [`requests.http`](requests.http).
4. Clique em **Send Request** acima da requisição desejada.

As requisições identificadas como “sem chave” ou “chave inválida” são testes: elas devem retornar `401`.

## Problemas comuns

- **`DATABASE_URL não definida`:** configure a variável no mesmo terminal que executa o servidor ou o gerador de chave.
- **Erro ao conectar ao PostgreSQL:** confira usuário, senha, porta `5432` e se o banco está em execução.
- **Tabela não existe:** execute o conteúdo de `schema.sql` no banco `receitas`.
- **Permissão negada:** execute os comandos `GRANT` mostrados na preparação do banco.
- **API Key inválida:** gere uma nova chave e use exatamente o valor exibido.

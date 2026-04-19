# PsicoSistem Backend

Base de backend em Go para o SaaS de psicopedagogia, organizada por domínio e preparada para:

- autenticação JWT com `tenant_id`, `role` e `type`
- multi-tenant
- painel admin com controle por papel
- assinatura por organização
- respostas padronizadas para o frontend
- documentação inicial em OpenAPI

## Estrutura

```txt
/cmd
  /api
    main.go
/internal
  /auth
  /subscription
  /tenant
  /user
  /shared
    /config
    /database
    /errors
    /middleware
    /response
    /security
    /validator
/pkg
  /hash
  /jwt
  /logger
/docs
  openapi.yaml
```

## Rotas implementadas nesta fundação

- `POST /v1/auth/register`
- `POST /v1/auth/login`
- `POST /v1/auth/refresh`
- `GET /v1/tenant/me`
- `GET /v1/tenant/subscription`
- `GET /v1/users/me`
- `GET /v1/users`
- `POST /v1/users`
- `PATCH /v1/users/:id`
- `DELETE /v1/users/:id`
- `GET /v1/health`

## Variáveis de ambiente

Crie um arquivo `.env` na raiz do backend:

```env
JWT_SECRET=troque-este-segredo
# ou, se preferir:
SECRET_KEY=troque-este-segredo
PORT=8080
JWT_ISSUER=psicosistem-backend
DATA_FILE=data/app_state.json
APP_ENV=development
```

## Como subir o projeto

```bash
go run ./cmd/api
```

## Como rodar testes

```bash
go test ./...
```

## Persistência atual

A fundação já não usa mais repositório volátil em memória. O adapter atual persiste o estado em `data/app_state.json`, o que resolve perda de dados em reinício local e mantém os módulos desacoplados.

Para produção, o próximo passo natural é substituir apenas a implementação de `repository` por um adapter PostgreSQL/pgx, preservando handlers, usecases e contratos de resposta.

## Convenções da API

### Sucesso

```json
{
  "data": {},
  "meta": {},
  "error": null
}
```

### Erro

```json
{
  "data": null,
  "meta": null,
  "error": {
    "code": "INVALID_EMAIL",
    "message": "invalid email format"
  }
}
```

## Segurança aplicada

- JWT obrigatório, sem fallback inseguro
- token com `user_id`, `tenant_id`, `role`, `email` e `type`
- isolamento por `tenant_id`
- gestão de usuários restrita a `owner` e `admin`
- soft delete por status `inactive`

## Próximos módulos sugeridos

- `patient`
- `session`
- `dashboard`
- `finance`
- `guardian`
- `ai`
- adapter PostgreSQL com migrations

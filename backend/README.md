# PsicoSistem Backend

Base de backend em Go para o SaaS de psicopedagogia, organizada por domínio e preparada para:

- autenticação JWT com `tenant_id`, `role` e `type`
- multi-tenant
- painel admin com controle por papel
- permissões segmentadas por conta e por escopo de dados
- assinatura por organização
- persistência com `json` ou `postgres`
- perfil de infraestrutura para operação local ou AWS
- respostas padronizadas para o frontend
- documentação inicial em OpenAPI

## Estrutura

```txt
/cmd
  /api
    main.go
/internal
  /auth
  /permission
  /subscription
  /tenant
  /user
  /shared
    /bootstrap
    /config
    /database
    /errors
    /infra
    /middleware
    /permissions
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
- `POST /v1/auth/forgot-password`
- `POST /v1/auth/reset-password`
- `POST /v1/auth/refresh`
- `GET /v1/tenant/me`
- `GET /v1/tenant/subscription`
- `GET /v1/permissions/me`
- `GET /v1/account/me`
- `GET /v1/users/me`
- `GET /v1/users`
- `POST /v1/users`
- `GET /v1/users/:id/permissions`
- `PATCH /v1/users/:id/permissions`
- `PATCH /v1/users/:id`
- `DELETE /v1/users/:id`
- `GET /v1/health`

## Onboarding e antifraude

O fluxo inicial de autenticação/cadastro agora considera:

- confirmação obrigatória da sessão de checkout antes do `register`
- campo `cpf_cnpj` para validar CPF do responsável ou CNPJ da clínica
- unicidade de telefone da clínica para reduzir reutilização indevida
- trial de 1 mês no plano `intermediario`, com assinatura criada em `trialing`

## Variáveis de ambiente

Crie um arquivo `.env` na raiz do backend:

```env
JWT_SECRET=troque-este-segredo
PORT=8080
JWT_ISSUER=psicosistem-backend
APP_ENV=development
FRONTEND_URL=http://localhost:3000
PASSWORD_RESET_TTL_MINUTES=30

# SMTP para recuperacao de senha
SMTP_HOST=
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=

# Persistencia
STORAGE_DRIVER=json
DATA_FILE=data/app_state.json
DATABASE_URL=postgres://postgres:postgres@localhost:5432/psicosistem?sslmode=disable
DATABASE_AUTO_MIGRATE=true
DATABASE_MAX_CONNS=10

# Infraestrutura
CLOUD_PROVIDER=local
AWS_REGION=sa-east-1
AWS_S3_BUCKET=
AWS_SECRETS_PREFIX=
AWS_USE_IAM_ROLE=true
AWS_CLOUDWATCH_NAMESPACE=PsicoSistem
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

A fundação já não usa mais repositório volátil em memória. O backend suporta duas estratégias:

- `STORAGE_DRIVER=json`: persiste estado em `data/app_state.json`
- `STORAGE_DRIVER=postgres`: conecta em `DATABASE_URL`, aplica schema inicial e usa repositórios SQL

O bootstrap que escolhe isso fica em `internal/shared/bootstrap`.

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
- identidade protegida por recarga do usuário salvo, não só pelas claims do token
- isolamento por `tenant_id`
- gestão de usuários e permissões restrita a `owner` e `admin`
- soft delete por status `inactive`

## Modelo de permissões

Além do `role`, cada conta possui um mapa de escopos:

- `user_directory`
- `patients`
- `services`
- `calendar`
- `finance`
- `ai_history`
- `plans`

Cada escopo aceita:

- `all`: visão completa do tenant
- `own`: visão apenas dos recursos ligados ao próprio usuário
- `none`: sem acesso

Defaults relevantes:

- `owner` e `admin`: `all` em tudo
- `professional`: `own` em pacientes, serviços, agenda, financeiro, IA e planos
- `financial`: `all` apenas em financeiro

## AWS

O backend já aceita perfil de infraestrutura `local` ou `aws` via `CLOUD_PROVIDER`.
O runtime AWS foi preparado para:

- RDS PostgreSQL
- S3 para documentos/artefatos
- Secrets Manager ou Parameter Store
- CloudWatch

O desenho recomendado está em `docs/infrastructure.md`.

## Próximos módulos sugeridos

- `patient`
- `session`
- `dashboard`
- `finance`
- `guardian`
- `ai`
- adapter PostgreSQL com migrations

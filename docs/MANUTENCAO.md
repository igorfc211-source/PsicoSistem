# PsicoSistem - Guia de alteracao e manutencao

Este documento e o mapa principal para qualquer pessoa alterar, corrigir, evoluir ou revisar o PsicoSistem. Ele explica como o frontend SvelteKit e o backend Go se conectam, onde cada responsabilidade vive, quais comandos rodar e quais cuidados tomar antes de entregar uma mudanca.

## Objetivo

Manter o sistema facil de evoluir, seguro e previsivel para um time trabalhar. Sempre que uma nova funcionalidade for criada, ela deve respeitar a separacao por modulos, ter regra de negocio em lugar correto, validar dados, manter isolamento por clinica e passar nos testes/checks antes de ser considerada pronta.

## Visao geral do repositorio

```txt
PsicoSistem/
  backend/
    cmd/api/main.go
    internal/
      auth/
      permission/
      shared/
      subscription/
      tenant/
      user/
    pkg/
      hash/
      jwt/
      logger/
    docs/
      infrastructure.md
      openapi.yaml
  frontend/
    src/
      routes/
      lib/
        auth/
        modules/
          clinic-shell/
          learners/
          scheduling/
        shared/
  docs/
    MANUTENCAO.md
```

## Comandos essenciais

### Backend

```powershell
cd backend
go run ./cmd/api
go test ./...
```

Se o cache do Go precisar ficar dentro do projeto:

```powershell
cd backend
$env:GOCACHE = Join-Path (Get-Location) '.gocache'
go test ./...
```

### Frontend

```powershell
cd frontend
npm install
npm run dev
npm run check
npm run build
```

### Validacao antes de entregar

```powershell
cd backend
go test ./...

cd ../frontend
npm run check
npm run build
```

## Arquitetura do backend

O backend e uma API REST em Go com Gin. A entrada da aplicacao fica em:

```txt
backend/cmd/api/main.go
```

Esse arquivo faz quatro coisas principais:

1. Carrega configuracao do `.env`.
2. Decide qual persistencia usar: `json` ou `postgres`.
3. Monta repositories, usecases e handlers.
4. Registra rotas HTTP e inicia o servidor.

### Camadas do backend

Cada modulo de dominio segue este desenho:

```txt
internal/nome-do-modulo/
  dto.go
  handler.go
  usecase.go
  model.go
  repository.go
  postgres_repository.go
  usecase_test.go
```

Nem todo modulo precisa ter todos esses arquivos, mas este e o padrao preferido.

### Responsabilidade de cada arquivo

`dto.go`

Define os dados que entram e saem da API. Exemplo: `RegisterInput`, `LoginInput`, `AuthPayload`.

`handler.go`

Recebe a requisicao HTTP. Ele deve apenas ler parametros/body, chamar o usecase e devolver uma resposta padronizada. Nao coloque regra de negocio pesada no handler.

`usecase.go`

Contem regra de negocio. Aqui ficam validacoes, permissoes, regras de plano, criacao de entidades, bloqueios e decisoes importantes.

`model.go`

Define a entidade principal do dominio. Exemplo: `User`, `Tenant`, `Subscription`.

`repository.go`

Define a interface do repositorio e a implementacao JSON local. A interface protege o usecase de depender diretamente de banco.

`postgres_repository.go`

Implementa a mesma interface usando PostgreSQL.

`usecase_test.go`

Testa as regras de negocio sem precisar subir servidor.

## Fluxo de uma requisicao no backend

Exemplo: `POST /v1/auth/login`.

```txt
Frontend
  -> POST /v1/auth/login
    -> handler Login
      -> usecase Login
        -> repository FindByEmail
          -> json ou postgres
        -> hash.Compare
        -> jwt.GenerateToken
      -> response.Success
  <- JSON padronizado
```

## Padrao de resposta da API

Toda resposta deve seguir este formato.

Sucesso:

```json
{
  "data": {},
  "meta": {},
  "error": null
}
```

Erro:

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

Use sempre:

```go
response.Success(c, status, data, meta)
response.Fail(c, err)
```

## Autenticacao e seguranca

O backend usa JWT com dados essenciais:

```txt
user_id
tenant_id
role
email
type
```

Rotas privadas usam middleware:

```go
middleware.AuthRequiredWithResolver(jwtSvc, identityResolver, security.UserTypeInternal)
```

O resolver recarrega o usuario do armazenamento antes de aceitar a identidade. Isso evita confiar apenas no token antigo quando usuario, status ou permissoes mudaram.

### Regra de ouro multi-tenant

Nunca aceite `tenant_id` vindo do frontend para operacoes internas. O `tenant_id` deve vir sempre do JWT/identity.

Exemplo correto:

```go
actor, ok := middleware.GetIdentity(c)
tenantID := actor.TenantID
```

Exemplo perigoso:

```go
tenantID := input.TenantID
```

## Persistencia do backend

O backend suporta dois modos:

### JSON local

Usado para desenvolvimento simples.

```env
STORAGE_DRIVER=json
DATA_FILE=data/app_state.json
```

Os dados ficam em um arquivo JSON local. Bom para aprender, testar e demonstrar.

### PostgreSQL

Usado para ambiente mais proximo de producao.

```env
STORAGE_DRIVER=postgres
DATABASE_URL=postgres://postgres:postgres@localhost:5432/psicosistem?sslmode=disable
DATABASE_AUTO_MIGRATE=true
```

Quando `postgres` esta ativo, o backend abre um pool com `pgx`, cria o schema basico e usa os repositories SQL.

## Arquitetura do frontend

O frontend usa SvelteKit com TypeScript. A tela principal autenticada fica em:

```txt
frontend/src/routes/app/+page.svelte
```

Essa pagina e a orquestradora. Ela guarda o estado principal, chama operacoes de dominio e passa dados/callbacks para componentes menores.

### Modulos principais

```txt
frontend/src/lib/modules/
  auth/
  clinic-shell/
  learners/
  scheduling/
```

`auth`

Guarda funcoes de sessao e comunicacao de login/cadastro.

`clinic-shell`

Layout geral: sidebar, topbar, banner, estilos globais do painel, navegacao e tema.

`learners`

Modulo dos aprendentes: dominio, calendario, lista, detalhe, abas, documentos, anamnese, plano de acao e relatorios.

`scheduling`

Modulo de agenda: eventos livres, sessoes, timeline diaria e regras de horarios.

`shared`

Componentes e utilitarios reutilizaveis, como formatadores e editor rico.

## Organizacao preferida de componentes Svelte

Use este padrao:

```txt
components/
  area/
    NomeDoComponente.svelte
    index.ts
```

Exemplo:

```txt
learners/components/tabs/ReportsTab.svelte
learners/components/tabs/ActionPlanTab.svelte
learners/components/calendar/CalendarPanel.svelte
```

Cada componente deve:

1. Receber dados por props.
2. Emitir acoes por callbacks.
3. Evitar acessar armazenamento diretamente.
4. Ter comentarios de bloco nas divisoes principais.
5. Manter CSS usando classes existentes quando possivel.

## Estado e dados no frontend

Atualmente, grande parte dos dados clinicos do painel fica local no navegador:

```txt
localStorage
IndexedDB
```

`localStorage`

Guarda metadados leves, como aprendentes, visitas, relatorios e configuracoes simples.

`IndexedDB`

Guarda blobs de documentos, porque arquivos grandes nao devem ir para localStorage.

Ao migrar esses dados para backend, mantenha o frontend consumindo uma camada de infraestrutura, para nao espalhar `fetch` por todos os componentes.

## Como adicionar uma nova funcionalidade no backend

Exemplo: criar modulo `patient`.

1. Criar pasta:

```txt
backend/internal/patient/
```

2. Criar modelo:

```go
type Patient struct {
    ID uuid.UUID
    TenantID uuid.UUID
    Name string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

3. Criar DTOs:

```go
type CreateInput struct {
    Name string `json:"name"`
}
```

4. Criar repository interface:

```go
type Repository interface {
    Create(ctx context.Context, item *Patient) error
    ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Patient, error)
}
```

5. Criar usecase com regras:

```go
func (u *Usecase) Create(ctx context.Context, actor security.Identity, input CreateInput) (*Response, error) {
    if !actor.IsInternal() {
        return nil, sharederrors.Forbidden("invalid actor")
    }
    // validar, montar entidade, salvar
}
```

6. Criar handler HTTP.

7. Registrar no `cmd/api/main.go`.

8. Adicionar implementacao JSON e Postgres.

9. Adicionar testes em `usecase_test.go`.

10. Atualizar `docs/openapi.yaml`.

## Como adicionar uma nova tela no frontend

1. Identifique o modulo correto.

Exemplo:

```txt
learners
scheduling
clinic-shell
```

2. Crie o componente em uma pasta especifica:

```txt
frontend/src/lib/modules/learners/components/example/ExamplePanel.svelte
```

3. Exporte no `index.ts` da pasta.

4. Passe dados por props.

5. Envie acoes por callbacks.

6. Conecte no workspace ou rota orquestradora.

7. Rode:

```powershell
npm run check
npm run build
```

## Como alterar uma regra existente

Antes de alterar regra, responda:

1. Essa regra pertence ao backend ou ao frontend?
2. A regra afeta seguranca, permissao, pagamento ou dados clinicos?
3. Existe teste cobrindo essa regra?
4. Precisa atualizar DTO/API?
5. Precisa atualizar documentacao?

### Exemplo

Regra: "profissional nao pode listar todos os usuarios".

Lugar correto:

```txt
backend/internal/user/usecase.go
```

Teste correto:

```txt
backend/internal/user/usecase_test.go
```

Nao coloque essa regra apenas no frontend. O frontend pode esconder botao, mas o backend precisa bloquear.

## Como escrever testes no backend

O Go reconhece arquivos terminados em:

```txt
*_test.go
```

E funcoes com nome:

```go
func TestNomeDaRegra(t *testing.T)
```

Padrao do projeto:

```go
func TestMinhaRegra(t *testing.T) {
    t.Parallel()

    store := database.NewStore(filepath.Join(t.TempDir(), "state.json"))
    if err := store.Initialize(); err != nil {
        t.Fatalf("initialize store: %v", err)
    }

    repo := modulo.NewRepository(store)
    usecase := modulo.NewUsecase(repo)

    result, err := usecase.Acao(context.Background(), input)

    if err != nil {
        t.Fatalf("acao: %v", err)
    }

    if result.Campo != "esperado" {
        t.Fatalf("expected esperado, got %s", result.Campo)
    }
}
```

### Quando criar teste

Crie teste quando a mudanca envolver:

1. Autenticacao.
2. Permissao.
3. Plano ou assinatura.
4. Limite do plano.
5. Validacao de CPF/CNPJ, email, telefone ou senha.
6. Criacao/alteracao de entidade importante.
7. Bug corrigido que nao pode voltar.

## Testes do frontend

O frontend ainda nao possui testes automatizados de componente, entao a validacao obrigatoria atual e:

```powershell
npm run check
npm run build
```

`npm run check`

Valida TypeScript, Svelte e contratos de props.

`npm run build`

Garante que a aplicacao compila para producao.

Quando o projeto crescer, recomenda-se adicionar testes com:

```txt
Vitest
Testing Library
Playwright
```

## Checklist antes de entregar uma mudanca

### Codigo

- A mudanca esta no modulo correto.
- Componentes Svelte continuam pequenos e separados.
- Regras de negocio criticas estao no backend.
- O backend nao confia em `tenant_id` vindo do frontend.
- Erros usam `sharederrors`.
- Respostas usam `response.Success` ou `response.Fail`.
- Novas rotas privadas usam middleware de auth.
- Novos dados sensiveis nao ficam expostos no JSON.
- Arquivos grandes nao vao para localStorage.

### Testes

- `go test ./...` passou.
- `npm run check` passou.
- `npm run build` passou.
- Foi criado teste para regra nova ou bug corrigido.

### Documentacao

- README/docs foram atualizados se houve rota nova.
- OpenAPI foi atualizado se houve contrato novo.
- Este guia foi atualizado se houve mudanca estrutural.

## Padroes de codigo

### Go

Use nomes claros:

```go
userRepo
subscriptionRepo
authUsecase
tenantHandler
```

Retorne erro cedo:

```go
if err != nil {
    return nil, err
}
```

Nao misture HTTP com regra de negocio. Handler chama usecase; usecase decide.

Use `context.Context` em operacoes de repositorio e usecase.

### Svelte

Use props tipadas:

```svelte
let { learner, onUpdateLearner } = $props<{
    learner: Learner;
    onUpdateLearner: (patch: Partial<Learner>) => void;
}>();
```

Comentarios devem explicar regioes importantes do template:

```svelte
<!-- Lista de documentos do prontuario. -->
<div class="document-list">
```

Evite comentarios obvios:

```svelte
<!-- Botao que salva -->
```

## Padrao para permissao

O sistema usa escopos:

```txt
all
own
none
```

Regra esperada:

```txt
all  -> ve tudo dentro do tenant
own  -> ve apenas dados ligados ao usuario
none -> nao acessa
```

Ao criar modulo novo, pense:

1. Qual permissao controla esse modulo?
2. Owner/admin acessam tudo?
3. Professional acessa apenas proprios dados?
4. Financial acessa apenas financeiro?

## Padrao para banco PostgreSQL

Ao adicionar tabela:

1. Adicione `CREATE TABLE IF NOT EXISTS` em `EnsureSchema`.
2. Use `tenant_id` quando o dado pertencer a uma clinica.
3. Crie indice por `tenant_id`.
4. Use `created_at` e `updated_at`.
5. Use soft delete quando for dado de usuario ou dado clinico sensivel.

Exemplo:

```sql
CREATE TABLE IF NOT EXISTS patients (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_patients_tenant_id ON patients (tenant_id);
```

## Integracao frontend-backend

Rotas principais atuais:

```txt
POST /v1/auth/register
POST /v1/auth/login
POST /v1/auth/refresh
GET  /v1/tenant/me
GET  /v1/tenant/subscription
GET  /v1/permissions/me
GET  /v1/users/me
GET  /v1/users
POST /v1/users
PATCH /v1/users/:id
DELETE /v1/users/:id
```

Ao consumir uma rota no frontend:

1. Centralize o `fetch` em uma camada de API.
2. Leia `data`, `meta` e `error`.
3. Nunca trate erro apenas por texto solto.
4. Guarde o token com cuidado.
5. Envie `Authorization: Bearer <token>` nas rotas privadas.

## Dados sensiveis e LGPD

Este sistema lida com dados clinicos e de criancas/adolescentes. Portanto:

1. Colete apenas dados necessarios.
2. Evite logs com CPF, CNPJ, telefone, relatorios ou documentos.
3. Proteja documentos fora de localStorage.
4. Use tenant isolation em todas as queries.
5. Planeje auditoria para alteracoes importantes.
6. Evite excluir dados clinicos fisicamente sem politica clara.

## Troubleshooting

### Backend nao sobe

Verifique:

```txt
backend/.env
JWT_SECRET
PORT
STORAGE_DRIVER
DATABASE_URL se STORAGE_DRIVER=postgres
```

Rode:

```powershell
cd backend
go run ./cmd/api
```

### Login falha no frontend

Verifique:

1. Backend esta rodando na porta `8080`.
2. Frontend esta apontando para endpoint correto.
3. Usuario existe em `data/app_state.json` ou no Postgres.
4. Senha esta correta.
5. Resposta da API contem `error.code`.

### Frontend nao compila

Rode:

```powershell
cd frontend
npm run check
```

O erro normalmente indica componente, prop, import ou tipo TypeScript quebrado.

### Build passa, mas tela quebra

Verifique:

1. Dados antigos no localStorage.
2. Normalizadores em `learner-storage.ts`.
3. Props obrigatorias adicionadas sem fallback.
4. Eventos/callbacks nao passados pelo workspace.

## Como manter este guia

Atualize este arquivo quando:

1. Um novo modulo for criado.
2. Uma rota nova entrar.
3. O formato de dados mudar.
4. O fluxo de deploy mudar.
5. A arquitetura do frontend ou backend mudar.
6. Um padrao novo for combinado pelo time.

## Definition of Done

Uma alteracao esta pronta quando:

1. A regra esta implementada no lugar correto.
2. A UI esta responsiva quando necessario.
3. Erros sao tratados de forma clara.
4. Testes/checks passaram.
5. Documentacao relevante foi atualizada.
6. Nao ha regressao conhecida.
7. O codigo esta simples o suficiente para outro dev manter.


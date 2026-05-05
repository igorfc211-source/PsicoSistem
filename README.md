# PsicoSistem

Sistema SaaS para clinicas de psicopedagogia, com backend em Go e frontend em SvelteKit.

## Estrutura

```txt
backend/   API em Go, autenticacao, tenants, usuarios, permissoes e assinaturas
frontend/  Interface SvelteKit, painel clinico, agenda e aprendentes
docs/      Documentacao geral de manutencao do produto
```

## Documentacao principal

Leia primeiro:

[Guia de alteracao e manutencao](./docs/MANUTENCAO.md)

## Comandos rapidos

Backend:

```powershell
cd backend
go run ./cmd/api
go test ./...
```

Frontend:

```powershell
cd frontend
npm install
npm run dev
npm run check
npm run build
```

## Portas usadas

```txt
Backend:  localhost:8080
Frontend: localhost:5173
```


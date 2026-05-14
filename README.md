# PsiSistem

Sistema SaaS para clinicas de psicopedagogia, com backend em Go e frontend em SvelteKit.
  com o objetivo de organizar os pacientes(aprendentes) e seus resposáveis, o sistema contém uma UI mais objetiva para quem não tem afinidade com sistemas tipo CRM ou ERP
tendo seções que se organizam por prioridade

## A tela principal "aprendentes"
 aonde pode-se criar um card do paciente, aonde recebe as informações,foto do aprendente, nome, idade, e-mail, nome e número do responsável(guardian), gênero, Status(ativo ou inativo), ínicio e fim das sessões, número de visitas, valor por sessão e/ou valor geral (para planos e coisas do tipo)

  **tabs da seção de aprendentes**
- Plano de ação do paciente
- Agenda, aonde você pode ver os dias agendados, e após clicar no dia, consegue clicar para agendar, aonde é mandado para seção de agenda daquele dia selecionado
- Anamnese, aonde pode-se escrever o relátorio de ananmese do aprendente, e pode também anexar um arquivo de até 20mb para manter a anamnese
- Documento, aonde pode-se anexar os arquivos necessários
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


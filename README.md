# PsiSistem

Sistema SaaS para clinicas de psicopedagogia, com backend em Go e frontend em SvelteKit.
  com o objetivo de organizar os pacientes(aprendentes) e seus resposáveis, o sistema contém uma UI mais objetiva para quem não tem afinidade com sistemas tipo CRM ou ERP
tendo seções que se organizam por prioridade

### *Sistema em andamento, sujeito a correção de bugs e adições(já planejadas)*

## Tela de criação da conta
<img width="1876" height="1004" alt="Screenshot_1" src="https://github.com/user-attachments/assets/807862cd-9fb4-459f-8fb7-1286d2d8b845" />
>>
<img width="1859" height="957" alt="Screenshot_2" src="https://github.com/user-attachments/assets/d0451dc7-9cdb-4a4f-9f28-4caa2863283f" />

## tela de login
<img width="1738" height="512" alt="Screenshot_16" src="https://github.com/user-attachments/assets/be3edadc-db85-4ae7-b88a-3f83086984cd" />



## A tela principal "aprendentes"
 aonde pode-se criar um card do paciente, aonde recebe as informações,foto do aprendente, nome, idade, e-mail, nome e número do responsável(guardian), gênero, Status(ativo ou inativo), ínicio e fim das sessões, número de visitas, valor por sessão e/ou valor geral (para planos e coisas do tipo)
 
<img width="938" height="482" alt="Screenshot_3" src="https://github.com/user-attachments/assets/3a66e3c2-755c-4910-9650-8b02a3c09221" />
<img width="784" height="951" alt="Screenshot_5" src="https://github.com/user-attachments/assets/9721a62b-e80e-4df8-9554-57a9f925a58a" />




  **tabs da seção de aprendentes**
- Resumo
  <img width="1877" height="965" alt="Screenshot_6" src="https://github.com/user-attachments/assets/2935a551-615e-435b-ad91-c6d23a065c5f" />
  
- Agenda, aonde você pode ver os dias agendados, e após clicar no dia, consegue clicar para agendar, aonde é mandado para seção de agenda daquele dia selecionado
  <img width="889" height="862" alt="Screenshot_7" src="https://github.com/user-attachments/assets/ce184351-fa49-4b5e-b527-2e464ee24031" />

## (Caso tenha uma seção marcada)
  <img width="871" height="943" alt="Screenshot_14" src="https://github.com/user-attachments/assets/a6a89a77-0055-4a8c-a3c6-eeb8948aaab4" />


- Relatório
  
<img width="884" height="508" alt="Screenshot_10" src="https://github.com/user-attachments/assets/ae28383c-2254-4d68-ada7-19a238816a98" />

- Anamnese, aonde pode-se escrever o relátorio de ananmese do aprendente, e pode também anexar um arquivo de até 20mb para manter a anamnese
  
  <img width="877" height="549" alt="Screenshot_8" src="https://github.com/user-attachments/assets/eed1bc1e-f3da-45c9-8ded-77c646cfdd6b" />

- Plano de ação do paciente
  
  <img width="879" height="829" alt="Screenshot_9" src="https://github.com/user-attachments/assets/750d8b3a-b630-4fc1-890c-1f34cb08a0d0" />
  
- Documento, aonde pode-se anexar os arquivos necessários
  <img width="879" height="282" alt="Screenshot_11" src="https://github.com/user-attachments/assets/aec2d40d-69c3-429a-9cf1-5f246d9d9d59" />

## Seção Calendário
Tela primária
<img width="1876" height="958" alt="Screenshot_12" src="https://github.com/user-attachments/assets/8b9f0522-5eb5-4c94-ab6e-e7f7d326cdfe" />
 Sessão agendada, ou evento
 <img width="1874" height="964" alt="Screenshot_13" src="https://github.com/user-attachments/assets/fa4413ca-d230-4da7-adac-cc8b6d07a4cb" />

## Seção Financeiro
<img width="1870" height="917" alt="financeiro" src="https://github.com/user-attachments/assets/d8eb20d3-b547-4c02-aa8a-347df04f831c" />

## Seção comunicação/contatos de todos paciente e responsáveis
<img width="1873" height="964" alt="comunicação" src="https://github.com/user-attachments/assets/f78eed9a-38b9-41bc-8007-c5f9feaea159" />

Card selecionado
<img width="1875" height="957" alt="comunicação 2" src="https://github.com/user-attachments/assets/7670a0a0-9b07-41aa-b44e-189cbf326184" />

## Requisitos 
- Postgres
- Go
- node.js

## Instalação
````txt
Git Clone https://github.com/igorfc211-source/PsicoSistem
````

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


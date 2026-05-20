# PsiSistem

Sistema SaaS para clínicas de psicopedagogia desenvolvido com foco em organização clínica, acompanhamento de aprendentes e simplificação operacional.

O projeto foi criado para oferecer uma interface objetiva e intuitiva para profissionais que não possuem familiaridade com sistemas complexos como CRMs ou ERPs tradicionais.

# Visão Geral

O PsiSistem centraliza:

Gestão de aprendentes
Controle de responsáveis
Agenda clínica
Relatórios
Financeiro
Comunicação
Documentação clínica

A plataforma foi estruturada utilizando arquitetura SaaS multi-tenant, permitindo expansão futura para múltiplas clínicas e diferentes níveis de acesso.

Stack Utilizada
- Frontend
- SvelteKit
- TypeScript
- TailwindCSS
- Backend
- Go
- PostgreSQL
- JWT Authentication
- REST API
- Infraestrutura
- Estrutura preparada para AWS
- Upload de arquivos
- Organização modular
- Separação por domínio
- Funcionalidades
- Gestão de Aprendentes

Cadastro completo de aprendentes contendo:

- Foto
- Nome
- Idade
- E-mail
- Responsável
- Contato do responsável
- Gênero
- Status
- Controle de sessões
- Valores por sessão/plano
  
## Tela principal

<img width="938" height="482" alt="Screenshot_3" src="https://github.com/user-attachments/assets/3a66e3c2-755c-4910-9650-8b02a3c09221" />

### Cadastro do aprendente

<img width="784" height="951" alt="Screenshot_5" src="https://github.com/user-attachments/assets/9721a62b-e80e-4df8-9554-57a9f925a58a" />

### Área do Aprendente


Cada aprendente possui uma área individual organizada em seções clínicas.

### Resumo Clínico

Visualização rápida das principais informações do aprendente.

<img width="1877" height="965" alt="Screenshot_6" src="https://github.com/user-attachments/assets/2935a551-615e-435b-ad91-c6d23a065c5f" />

### Agenda

Permite visualizar dias agendados e criar sessões diretamente pelo calendário.

<img width="889" height="862" alt="Screenshot_7" src="https://github.com/user-attachments/assets/ce184351-fa49-4b5e-b527-2e464ee24031" />

### Sessão agendada

<img width="871" height="943" alt="Screenshot_14" src="https://github.com/user-attachments/assets/a6a89a77-0055-4a8c-a3c6-eeb8948aaab4" />

### Relatórios

Área destinada à criação e acompanhamento de relatórios clínicos.

<img width="884" height="508" alt="Screenshot_10" src="https://github.com/user-attachments/assets/ae28383c-2254-4d68-ada7-19a238816a98" />

### Anamnese

<img width="877" height="549" alt="Screenshot_8" src="https://github.com/user-attachments/assets/eed1bc1e-f3da-45c9-8ded-77c646cfdd6b" />

Registro completo da anamnese do aprendente, incluindo upload de arquivos.

Upload de arquivos até 20MB

### Plano de Ação

Organização de estratégias, metas e acompanhamento clínico.

<img width="879" height="829" alt="Screenshot_9" src="https://github.com/user-attachments/assets/750d8b3a-b630-4fc1-890c-1f34cb08a0d0" />

### Documentos

Upload e gerenciamento de documentos importantes.

<img width="879" height="282" alt="Screenshot_11" src="https://github.com/user-attachments/assets/aec2d40d-69c3-429a-9cf1-5f246d9d9d59" />

## Calendário

Sistema de calendário para gerenciamento de sessões e eventos clínicos.

## Tela principal

<img width="1876" height="958" alt="Screenshot_12" src="https://github.com/user-attachments/assets/8b9f0522-5eb5-4c94-ab6e-e7f7d326cdfe" />

### Evento agendado

<img width="1874" height="964" alt="Screenshot_13" src="https://github.com/user-attachments/assets/fa4413ca-d230-4da7-adac-cc8b6d07a4cb" />

## Financeiro

 Controle financeiro integrado para:

- Sessões
- Planos
- Pagamentos
- Valores pendentes

<img width="1870" height="917" alt="financeiro" src="https://github.com/user-attachments/assets/d8eb20d3-b547-4c02-aa8a-347df04f831c" />

## Comunicação

 Centralização de contatos de aprendentes e responsáveis.

<img width="1873" height="964" alt="comunicação" src="https://github.com/user-attachments/assets/f78eed9a-38b9-41bc-8007-c5f9feaea159" />

Contato selecionado

<img width="1875" height="957" alt="comunicação 2" src="https://github.com/user-attachments/assets/7670a0a0-9b07-41aa-b44e-189cbf326184" />

## Autenticação

Criação de conta

<img width="1876" height="1004" alt="Screenshot_1" src="https://github.com/user-attachments/assets/807862cd-9fb4-459f-8fb7-1286d2d8b845" /> <img width="1859" height="957" alt="Screenshot_2" src="https://github.com/user-attachments/assets/d0451dc7-9cdb-4a4f-9f28-4caa2863283f" />

Login

<img width="1738" height="512" alt="Screenshot_16" src="https://github.com/user-attachments/assets/be3edadc-db85-4ae7-b88a-3f83086984cd" />

## Arquitetura do Projeto

backend/

 ├── cmd/
 │    └── api/
 │
 ├── internal/
 │    ├── auth/
 │    ├── users/
 │    ├── tenants/
 │    ├── learners/
 │    ├── scheduling/
 │    ├── reports/
 │    ├── finance/
 │    ├── communication/
 │    └── uploads/
 │
 ├── middleware/
 ├── database/
 ├── storage/
 ├── config/
 └── utils/


frontend/
 ├── src/
 │    ├── routes/
 │    ├── lib/
 │    │    ├── components/
 │    │    ├── services/
 │    │    ├── stores/
 │    │    ├── hooks/
 │    │    └── utils/
 │
 └── static/
 
## Recursos Técnicos

Arquitetura SaaS
Estrutura multi-tenant
Autenticação JWT
Upload de arquivos
Organização modular
API REST
Controle de permissões
Estrutura preparada para cloud storage
Sistema escalável para múltiplas clínicas
Requisitos
Go
PostgreSQL
Node.js

## Instalação

Clonar o projeto
git clone https://github.com/igorfc211-source/PsicoSistem

Executando o Backend
`cd backend`

`go run ./cmd/api
Testes
go test ./... `

Executando o Frontend

cd frontend

`npm install`
`
npm run dev
Build de produção
npm run build
`
Portas utilizadas
`
Backend:  localhost:8080
Frontend: localhost:5173
`
Documentação

Leia primeiro:

/docs/MANUTENCAO.md
Status do Projeto

Projeto em desenvolvimento contínuo.

Possíveis novas funcionalidades e melhorias podem ser encontradas na branch:

*functionalities*

Objetivo do Projeto

O PsiSistem foi desenvolvido com o objetivo de simplificar o fluxo operacional de clínicas psicopedagógicas através de uma interface acessível, moderna e organizada por prioridade de uso.

O foco principal do sistema é reduzir complexidade operacional e melhorar a experiência dos profissionais durante o acompanhamento clínico dos aprendentes.

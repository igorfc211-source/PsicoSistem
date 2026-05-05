# Guia dos modulos do frontend

Esta pasta segue uma organizacao por feature, parecida com projetos React/Next.js escalaveis.

## Camadas

- `components`: componentes Svelte agrupados por feature. Devem renderizar UI e chamar callbacks, sem controlar persistencia.
- `domain`: tipos, fabricas e helpers puros de dominio. Esta camada nao deve conhecer localStorage ou IndexedDB.
- `application`: funcoes de caso de uso que coordenam objetos do dominio. Use para regras reutilizaveis.
- `infrastructure`: storage do navegador, IndexedDB e detalhes de integracao.
- `presentation`: tipos e constantes usados apenas pela UI.
- `styles`: CSS global de um shell visual ou feature.

## Estilo de importacao

Prefira barrels por feature:

```ts
import { LearnersWorkspace } from '$lib/modules/learners/components';
import { filterLearners } from '$lib/modules/learners';
```

Mantenha arquivos de rota enxutos. Uma rota deve carregar estado, chamar funcoes de `application` e compor componentes.

## Organizacao dos componentes Svelte

Dentro de `components`, os arquivos ficam em subpastas por responsabilidade:

- `workspaces`: telas/areas grandes que juntam varios componentes.
- `tabs`: secoes internas de um prontuario ou painel.
- `forms`: formularios de criacao/edicao.
- `list`: listas e cards de selecao.
- `detail`: paineis de detalhe de uma entidade.
- `calendar`: componentes reutilizaveis de calendario.
- `navigation`, `topbar`, `feedback`: componentes do shell da aplicacao.

Cada subpasta tem seu proprio `index.ts`, e o `components/index.ts` reexporta tudo. Assim o time pode reorganizar internamente sem quebrar imports das rotas.

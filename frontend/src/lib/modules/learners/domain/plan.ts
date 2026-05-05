import type { CoreActionPlanKey } from './types';

export type PlanCategory = {
	key: CoreActionPlanKey;
	label: string;
	description: string;
};

export const PLAN_CATEGORIES: PlanCategory[] = [
	{
		key: 'educationProcess',
		label: 'Processo de educacao',
		description: 'Objetivos escolares, metodos e adaptacoes.'
	},
	{
		key: 'familyGuidance',
		label: 'Orientacao familiar',
		description: 'Combinados, rotina e comunicacao com responsaveis.'
	},
	{
		key: 'cognitiveSkills',
		label: 'Habilidades cognitivas',
		description: 'Atencao, memoria, linguagem, leitura e escrita.'
	},
	{
		key: 'behavior',
		label: 'Comportamento',
		description: 'Regulacao, autonomia e engajamento nas sessoes.'
	},
	{
		key: 'clinicGoals',
		label: 'Metas clinicas',
		description: 'Prioridades, indicadores e proximos passos.'
	}
];

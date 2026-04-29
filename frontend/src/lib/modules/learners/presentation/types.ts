export type LearnerFilter = 'active' | 'inactive' | 'all';
export type DetailTab = 'resumo' | 'agenda' | 'anamnese' | 'documentos' | 'plano' | 'relatorios';

export type DetailTabItem = {
	value: DetailTab;
	label: string;
};

export const DETAIL_TABS: DetailTabItem[] = [
	{ value: 'resumo', label: 'Resumo' },
	{ value: 'agenda', label: 'Agenda' },
	{ value: 'anamnese', label: 'Anamnese' },
	{ value: 'plano', label: 'Plano' },
	{ value: 'relatorios', label: 'Relatorios' },
	{ value: 'documentos', label: 'Documentos' }
];

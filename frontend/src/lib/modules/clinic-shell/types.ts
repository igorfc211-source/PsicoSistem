export type NavSection =
	| 'aprendentes'
	| 'agenda'
	| 'comunicacoes'
	| 'financeiro'
	| 'configuracoes';

export type BannerTone = 'info' | 'success' | 'error';

export type Banner = {
	tone: BannerTone;
	text: string;
};

export type NavItem = {
	value: NavSection;
	label: string;
	icon: string;
};

export const NAV_ITEMS: NavItem[] = [
	{ value: 'aprendentes', label: 'Aprendentes', icon: 'P' },
	{ value: 'agenda', label: 'Agenda', icon: 'A' },
	{ value: 'financeiro', label: 'Financeiro', icon: 'F' },
	{ value: 'comunicacoes', label: 'Comunicacoes', icon: 'C' },
	{ value: 'configuracoes', label: 'Configurações', icon: '>>' }
];

// Traduz a secao ativa para o titulo exibido no workspace.
export function getSectionTitle(section: NavSection) {
	const titles: Record<NavSection, string> = {
		aprendentes: 'Aprendentes',
		agenda: 'Agenda',
		financeiro: 'Financeiro',
		comunicacoes: 'Comunicacoes',
		configuracoes: 'Configurações'
	};

	return titles[section];
}

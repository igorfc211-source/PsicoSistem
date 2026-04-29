export type NavSection =
	| 'aprendentes'
	| 'agenda'
	| 'avaliacoes'
	| 'relatorios'
	| 'comunicacoes'
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
	{ value: 'avaliacoes', label: 'Avaliacoes', icon: 'V' },
	{ value: 'relatorios', label: 'Relatorios', icon: 'R' },
	{ value: 'comunicacoes', label: 'Comunicacoes', icon: 'C' },
	{ value: 'configuracoes', label: 'Configuracoes', icon: 'S' }
];

// Traduz a secao ativa para o titulo exibido no workspace.
export function getSectionTitle(section: NavSection) {
	const titles: Record<NavSection, string> = {
		aprendentes: 'Aprendentes',
		agenda: 'Agenda',
		avaliacoes: 'Avaliacoes',
		relatorios: 'Relatorios',
		comunicacoes: 'Comunicacoes',
		configuracoes: 'Configuracoes'
	};

	return titles[section];
}

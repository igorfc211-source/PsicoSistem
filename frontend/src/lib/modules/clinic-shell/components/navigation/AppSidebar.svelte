<script lang="ts">
	import { NAV_ITEMS, type NavItem, type NavSection } from '../../types';

	let {
		tenantName,
		activeSection,
		navItems = NAV_ITEMS,
		isOpen,
		onSelectSection,
		onClose
	} = $props<{
		tenantName: string;
		activeSection: NavSection;
		navItems?: NavItem[];
		isOpen: boolean;
		onSelectSection: (section: NavSection) => void;
		onClose: () => void;
	}>();
</script>

<aside class="sidebar" class:open={isOpen}>

	<!-- Identidade da clinica: logo simples, nome do tenant e futuro seletor para celular -->
	<div class="clinic-brand">
		<div class="brand-mark">PC</div>
		<strong>{tenantName}</strong>
	
		<button type="button" class="sidebar-close" aria-label="Fechar menu" onclick={onClose}>x</button>
	</div>

	<!-- Navegacao principal: controla qual workspace aparece no painel. -->
	<nav class="side-nav" aria-label="Navegacao principal">
		{#each navItems as item}
			<button
				type="button"
				class:active={activeSection === item.value}
				onclick={() => onSelectSection(item.value)}
			>
				<span>{item.icon}</span>
				{item.label}
			</button>
		{/each}
	</nav>



	<!-- Card comercial: isolado para evoluir depois para consumo real do plano. -->
	<div class="plan-card">
	
		<strong>Plano Premium</strong>
		
	</div>
</aside>

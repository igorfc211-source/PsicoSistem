<script lang="ts">
	import { NAV_ITEMS, type NavSection } from '../../types';

	let {
		tenantName,
		activeSection,
		onSelectSection
	} = $props<{
		tenantName: string;
		activeSection: NavSection;
		onSelectSection: (section: NavSection) => void;
	}>();
</script>

<aside class="sidebar">
	<div class="clinic-brand">
		<div class="brand-mark">PC</div>
		<strong>{tenantName}</strong>
		<span>v</span>
	</div>

	<nav class="side-nav" aria-label="Navegacao principal">
		{#each NAV_ITEMS as item}
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

	<div class="plan-card">
		<div>*</div>
		<strong>Plano Premium</strong>
		<p>Sua clinica esta com todos os recursos ativados.</p>
		<button type="button">Saiba mais</button>
	</div>
</aside>

<script lang="ts">
	import { getInitials } from '$lib/shared/formatters';

	let {
		name,
		photoUrl = '',
		size = 'regular'
	} = $props<{
		name: string;
		photoUrl?: string;
		size?: 'small' | 'regular' | 'large';
	}>();
</script>

<!-- Avatar reutilizavel: usa foto quando existe e cai para iniciais quando nao existe imagem. -->
<div class={`avatar ${size === 'regular' ? '' : size}`}>
	{#if photoUrl}
		<img src={photoUrl} alt={`Foto de ${name}`} />
	{:else}
		{getInitials(name)}
	{/if}
</div>

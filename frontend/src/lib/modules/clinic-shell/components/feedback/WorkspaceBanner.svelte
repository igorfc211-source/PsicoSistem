<script lang="ts">
	import type { Banner } from '../../types';

	let { banner } = $props<{ banner: Banner | null }>();
</script>

{#if banner}
	<!-- Feedback global da pagina: sucesso, erro ou informacao depois de uma acao. -->
	
	<div class={`workspace-banner ${banner.tone}`}>{banner.text}</div>
{/if}

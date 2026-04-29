<script lang="ts">
	import type { Banner } from '../../types';

	let { banner } = $props<{ banner: Banner | null }>();
</script>

{#if banner}
	<div class={`workspace-banner ${banner.tone}`}>{banner.text}</div>
{/if}

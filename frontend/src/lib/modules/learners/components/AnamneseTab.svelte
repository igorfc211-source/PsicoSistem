<script lang="ts">
	let {
		value,
		onChange
	} = $props<{
		value: string;
		onChange: (value: string) => void;
	}>();
</script>

<section class="tab-panel">
	<textarea
		class="large-textarea"
		value={value}
		oninput={(event) => onChange((event.currentTarget as HTMLTextAreaElement).value)}
		placeholder="Anamnese do aprendente."
	></textarea>
</section>

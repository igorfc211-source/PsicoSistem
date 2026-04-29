<script lang="ts">
	import type { ActionPlan, PlanCategory } from '$lib/modules/learners';

	let {
		actionPlan,
		categories,
		onChange
	} = $props<{
		actionPlan: ActionPlan;
		categories: PlanCategory[];
		onChange: (key: keyof ActionPlan, value: string) => void;
	}>();
</script>

<section class="tab-panel plan-grid">
	{#each categories as category}
		<label class="plan-category card">
			<span>{category.label}</span>
			<small>{category.description}</small>
			<textarea
				value={actionPlan[category.key]}
				oninput={(event) =>
					onChange(category.key, (event.currentTarget as HTMLTextAreaElement).value)}
			></textarea>
		</label>
	{/each}
</section>

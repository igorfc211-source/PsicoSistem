<script lang="ts">
	import type { ActionPlan, CoreActionPlanKey, PlanCategory } from '$lib/modules/learners';

	let {
		actionPlan,
		categories,
		onChange,
		onAddCustomField,
		onChangeCustomField,
		onRemoveCustomField
	} = $props<{
		actionPlan: ActionPlan;
		categories: PlanCategory[];
		onChange: (key: CoreActionPlanKey, value: string) => void;
		onAddCustomField: (label: string, description: string) => boolean;
		onChangeCustomField: (fieldId: string, value: string) => void;
		onRemoveCustomField: (fieldId: string) => void;
	}>();

	let customLabel = $state('');
	let customDescription = $state('');

	// Permite criar novas secoes do plano conforme a necessidade clinica do aprendente.
	function handleAddCustomField(event: SubmitEvent) {
		event.preventDefault();
		if (!customLabel.trim()) return;

		if (onAddCustomField(customLabel, customDescription)) {
			customLabel = '';
			customDescription = '';
		}
	}
</script>

<section class="tab-panel">
	<!-- Grade principal do plano: mistura secoes padrao e secoes personalizadas. -->
	<div class="plan-grid">
		<!-- Secoes padrao do protocolo clinico, definidas no dominio do modulo. -->
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

		<!-- Secoes personalizadas criadas pela clinica para este aprendente especifico. -->
		{#each actionPlan.customFields as field}
			<label class="plan-category card custom-plan-field">
				<div class="custom-plan-head">
					<div>
						<span>{field.label}</span>
						<small>{field.description || 'Campo personalizado do plano.'}</small>
					</div>
					<button type="button" class="danger-button" onclick={() => onRemoveCustomField(field.id)}>
						Remover
					</button>
				</div>
				<textarea
					value={field.value}
					oninput={(event) =>
						onChangeCustomField(field.id, (event.currentTarget as HTMLTextAreaElement).value)}
				></textarea>
			</label>
		{/each}
	</div>

	<!-- Criador de novas categorias: deixa o plano extensivel sem alterar o codigo. -->
	<form class="custom-plan-form card" onsubmit={handleAddCustomField}>
		<div>
			<strong>Adicionar campo ao plano</strong>
			<p>Crie uma nova categoria para objetivos especificos deste aprendente.</p>
		</div>

		<div class="form-grid">
			<label>
				<span>Nome do campo</span>
				<input bind:value={customLabel} placeholder="Ex: Regulacao emocional" required />
			</label>

			<label>
				<span>Descricao curta</span>
				<input bind:value={customDescription} placeholder="Ex: estrategias, indicadores e rotina" />
			</label>
		</div>

		<button type="submit" class="primary-button">Adicionar campo</button>
	</form>
</section>

<script lang="ts">
	import type { Learner } from '../domain/types';
	import { getInitials } from '$lib/shared/formatters';

	let {
		learners,
		selectedLearnerId,
		onSelectLearner
	} = $props<{
		learners: Learner[];
		selectedLearnerId: string | null;
		onSelectLearner: (id: string) => void;
	}>();
</script>

<div class="learner-list">
	{#each learners as learner}
		<button
			type="button"
			class:selected={learner.id === selectedLearnerId}
			onclick={() => onSelectLearner(learner.id)}
		>
			<div class="avatar">{getInitials(learner.name)}</div>
			<div>
				<strong>{learner.name}</strong>
				<span>{learner.age || 'Idade nao informada'} - {learner.gender || 'Genero nao informado'}</span>
			</div>
			<small class={learner.status}>{learner.status === 'active' ? 'Ativo' : 'Inativo'}</small>
			<b>&gt;</b>
		</button>
	{:else}
		<p class="empty-state">Nenhum aprendente encontrado.</p>
	{/each}
</div>

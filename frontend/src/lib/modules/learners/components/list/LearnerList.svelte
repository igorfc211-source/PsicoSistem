<script lang="ts">
	import type { Learner } from '../../domain/types';
	import LearnerAvatar from '../avatar/LearnerAvatar.svelte';

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
	<!-- Lista de cards compactos: cada botao abre o prontuario na coluna de detalhe. -->
	{#each learners as learner}
		<button
			type="button"
			class:selected={learner.id === selectedLearnerId}
			onclick={() => onSelectLearner(learner.id)}
		>
			<!-- Identificacao visual do aprendente. -->
			<LearnerAvatar name={learner.name} photoUrl={learner.photoUrl} />

			<!-- Dados resumidos exibidos antes de abrir o prontuario completo. -->
			<div>
				<strong>{learner.name}</strong>
				<span>{learner.age || 'Idade nao informada'} - {learner.gender || 'Genero nao informado'}</span>
			</div>

			<small class={learner.status}>{learner.status === 'active' ? 'Ativo' : 'Inativo'}</small>
			<b>&gt;</b>
		</button>
	{:else}
		<!-- Estado vazio quando busca/filtro nao encontra aprendentes. -->
		<p class="empty-state">Nenhum aprendente encontrado.</p>
	{/each}
</div>

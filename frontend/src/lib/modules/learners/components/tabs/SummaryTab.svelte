<script lang="ts">
	import { getAttendanceRate, getNextVisit, type Learner } from '$lib/modules/learners';
	import { formatDateTime, formatLongDate } from '$lib/shared/formatters';

	let {
		learner,
		onOpenAgenda
	} = $props<{
		learner: Learner;
		onOpenAgenda: () => void;
	}>();

	const nextVisit = $derived(getNextVisit(learner));
</script>

<section class="summary-grid">
	<!-- Card de proxima sessao: atalho direto para a agenda do aprendente. -->
	<div class="next-session card">
		<div>
			<strong>Proxima sessao</strong>
			{#if nextVisit}
				<p>{formatLongDate(nextVisit.date)}</p>
				<span>{nextVisit.status}</span>
			{:else}
				<p>Nenhuma sessao criada</p>
			{/if}
		</div>
		<button type="button" onclick={onOpenAgenda}>Ver agenda</button>
	</div>

	<!-- Indicador de frequencia calculado a partir das visitas planejadas/realizadas. -->
	<div class="card stat-card">
		<strong>Frequencia</strong>
		<div>{getAttendanceRate(learner)}%</div>
		<p>{learner.visits.length} sessoes planejadas</p>
	</div>

	<!-- Indicador de documentos gerais anexados ao prontuario. -->
	<div class="card stat-card">
		<strong>Documentos</strong>
		<div>{learner.documents.length}</div>
		<p>Arquivos armazenados</p>
	</div>

	<!-- Ultimo relatorio como observacao rapida no resumo do aprendente. -->
	<div class="card observations">
		<strong>Ultimas observacoes</strong>
		<p>{learner.reports[0]?.text ?? 'Nenhum relatorio salvo para este aprendente.'}</p>
		{#if learner.reports[0]}
			<span>{formatDateTime(learner.reports[0].createdAt)}</span>
		{/if}
	</div>
</section>

<script lang="ts">
	import { getAttendanceRate, getNextVisit, type Learner } from '$lib/modules/learners';
	import { formatDateTime, formatLongDate, formatTimeForNBR } from '$lib/shared/formatters';

	let {
		learner,
		onOpenAgenda
	} = $props<{
		learner: Learner;
		onOpenAgenda: () => void;
	}>();

	const nextVisit = $derived(getNextVisit(learner));
	const learnerMeetings = $derived(
		[...learner.visits]
			.sort((left, right) =>
				`${left.date} ${left.startTime}`.localeCompare(`${right.date} ${right.startTime}`)
			)
			.slice(0, 6)
	);

	function getVisitStatusLabel(status: Learner['visits'][number]['status']) {
		if (status === 'completed') return 'Realizada';
		if (status === 'missed') return 'Falta';
		return 'Agendada';
	}
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

	<!-- Reunioes/atendimentos do aprendente para enxergar a agenda direto no resumo. -->
	<div class="card learner-meetings">
		<div class="summary-card-head">
			<strong>Reunioes do aprendente</strong>
			<span>{learner.visits.length} registro{learner.visits.length === 1 ? '' : 's'}</span>
		</div>

		<div class="summary-meeting-list">
			{#each learnerMeetings as visit}
				<div class="summary-meeting-row">
					<div>
						<strong>{visit.title || 'Sessao individual'}</strong>
						<span>
							{formatLongDate(visit.date)} - {formatTimeForNBR(visit.startTime)} as
							{formatTimeForNBR(visit.endTime)}
						</span>
					</div>
					<small>{getVisitStatusLabel(visit.status)}</small>
				</div>
			{:else}
				<p>Nenhuma reuniao registrada para este aprendente.</p>
			{/each}
		</div>
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

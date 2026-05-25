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
	const totalSessions = $derived(learner.visits.length);

const completedSessions = $derived(
	learner.visits.filter((visit) => visit.status === 'completed').length
);

const missedSessions = $derived(
	learner.visits.filter((visit) => visit.status === 'missed').length
);

const scheduledSessions = $derived(
	learner.visits.filter((visit) => visit.status === 'scheduled').length
);


	function getVisitStatusLabel(status: Learner['visits'][number]['status']) {
		if (status === 'completed') return 'Realizada';
		if (status === 'missed') return 'Falta';
		if (status === 'scheduled') return 'Agendado';
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
			{:else}
				<p>Nenhuma sessao criada</p>
			{/if}
		</div>
		<button type="button" onclick={onOpenAgenda}>Ver agenda</button>
	</div>

	<!-- Indicador de frequencia calculado a partir das visitas planejadas/realizadas. -->

<div class="card stat-card flex flex-col gap-4 overflow-hidden rounded-2xl p-4">
	<div class="flex items-start justify-between gap-4">
		<div class="flex flex-col">
			<strong class="text-xs font-medium text-zinc-500">
				Total de sessões
			</strong>

			<p class="text-3xl font-bold leading-none text-zinc-900">
				{totalSessions}
			</p>
		</div>

		<div class="flex flex-col items-end">
			<h2 class="text-xs font-medium text-zinc-500">
				Frequência
			</h2>

			<div class="text-2xl font-bold leading-none text-zinc-900">
				{getAttendanceRate(learner)}%
			</div>
		</div>
	</div>

	<div class="grid grid-cols-3 gap-2">
		<span class="flex flex-col rounded-xl bg-zinc-100 px-3 py-2">
			<label class="text-[11px] text-zinc-500">
				Realizadas
			</label>

			<strong class="text-base font-semibold text-zinc-900">
				{completedSessions}
			</strong>
		</span>

		<span class="flex flex-col rounded-xl bg-zinc-100 px-3 py-2">
			<label class="text-[11px] text-zinc-500">
				Agendadas
			</label>

			<strong class="text-base font-semibold text-zinc-900">
				{scheduledSessions}
			</strong>
		</span>

		<span class="flex flex-col rounded-xl bg-zinc-100 px-3 py-2">
			<label class="text-[11px] text-zinc-500">
				Faltas
			</label>

			<strong class="text-base font-semibold text-zinc-900">
				{missedSessions}
			</strong>
		</span>
	</div>
</div>

<!-- Indicador de documentos gerais anexados ao prontuario. -->
<div class="card stat-card flex flex-col rounded-2xl p-4">
	<strong class="text-xs font-medium text-zinc-500">
		Documentos
	</strong>

	<div class="text-3xl font-bold leading-none text-zinc-900">
		{learner.documents.length}
	</div>

	<p class="text-sm text-zinc-500">
		Arquivos armazenados
	</p>
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

<style>
	.session-stats {
	display: flex;
	gap: 1rem;
	margin-top: 1rem;
	flex-wrap: wrap;
}

.session-stats span {
	display: flex;
	flex-direction: column;
	padding: 0.65rem 0.9rem;
	border-radius: 14px;
	background: #f8f8fa;
	min-width: 90px;
}

.session-stats label {
	font-size: 0.72rem;
	font-weight: 500;
	color: #8a8a96;
	margin-bottom: 0.2rem;
	letter-spacing: 0.02em;
}

.session-stats strong {
	font-size: 1rem;
	font-weight: 600;
	color: #1c1c1e;
}

</style>
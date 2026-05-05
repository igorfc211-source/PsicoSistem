<script lang="ts">
	import type { CalendarDay, Learner, LearnerVisitSchedule } from '$lib/modules/learners';
	import type {
		AgendaEvent,
		NewAgendaEventInput,
		NewSessionAppointmentInput,
		ScheduleItem
	} from '$lib/modules/scheduling';
	import { formatDateForNBR } from '$lib/shared/formatters';
	import CalendarPanel from '../calendar/CalendarPanel.svelte';
	import { DaySchedulerPanel } from '$lib/modules/scheduling/components';

	let {
		calendarDays,
		monthLabel,
		selectedDate,
		currentDateLabel,
		learners,
		selectedLearnerId,
		userName,
		dayItems,
		pendingVisits,
		onShiftMonth,
		onSelectCalendarDate,
		onSelectLearnerId,
		onOpenLearner,
		onCreateSession,
		onCreateEvent,
		onRemoveSession,
		onRemoveEvent
	} = $props<{
		calendarDays: CalendarDay[];
		monthLabel: string;
		selectedDate: string;
		currentDateLabel: string;
		learners: Learner[];
		selectedLearnerId: string | null;
		userName: string;
		dayItems: ScheduleItem[];
		pendingVisits: LearnerVisitSchedule[];
		onShiftMonth: (delta: number) => void;
		onSelectCalendarDate: (date: string) => void;
		onSelectLearnerId: (id: string) => void;
		onOpenLearner: (id: string) => void;
		onCreateSession: (input: NewSessionAppointmentInput) => boolean;
		onCreateEvent: (input: NewAgendaEventInput) => boolean;
		onRemoveSession: (learnerId: string, visitId: string) => void;
		onRemoveEvent: (event: AgendaEvent) => void;
	}>();
</script>

<section class="agenda-workspace">
	<!-- Lateral da agenda: calendario compacto e filtros do contexto atual. -->
	<div class="agenda-sidebar">
		<h2>Agenda</h2>

		<!-- Mini calendario: define o dia que alimenta a linha do tempo principal. -->
		<CalendarPanel
			days={calendarDays}
			{monthLabel}
			{selectedDate}
			variant="mini"
			onShiftMonth={onShiftMonth}
			onSelectDate={onSelectCalendarDate}
		/>

		<!-- Filtros de agenda: preparados para crescer com equipe, sala ou aprendente. -->
		<div class="agenda-filters">
			<label>
				<span>Profissional</span>
				<select>
					<option>{userName}</option>
				</select>
			</label>

			<label>
				<span>Aprendente</span>
				<select
					value={selectedLearnerId ?? ''}
					onchange={(event) => {
						const id = (event.currentTarget as HTMLSelectElement).value;
						if (id) onSelectLearnerId(id);
					}}
				>
					<option value="" disabled>Selecione</option>
					{#each learners as learner}
						<option value={learner.id}>{learner.name}</option>
					{/each}
				</select>
			</label>
		</div>

		<!-- Visitas pendentes: fila rapida para localizar sessoes ainda nao realizadas. -->
		<div class="pending-visits card">
			<div class="pending-head">
				<strong>Visitas pendentes</strong>
				<span>{pendingVisits.length}</span>
			</div>

			<div class="pending-list">
				{#each pendingVisits.slice(0, 6) as item}
					<button
						type="button"
						onclick={() => {
							onSelectLearnerId(item.learner.id);
							onSelectCalendarDate(item.visit.date);
						}}
					>
						<strong>{item.learner.name}</strong>
						<span>{formatDateForNBR(item.visit.date)} - {item.visit.startTime}</span>
					</button>
				{:else}
					<p class="empty-state">Nenhuma visita pendente.</p>
				{/each}
			</div>
		</div>
	</div>

	<!-- Area principal da agenda: navegacao do dia e formulario de novos compromissos. -->
	<div class="agenda-main">
		<!-- Toolbar superior: troca de visualizacao e atalho para agendar no dia selecionado. -->
		<div class="agenda-toolbar">
			<div class="view-switch">
				<button type="button" class="active">Dia</button>
				<button type="button">Semana</button>
				<button type="button">Mes</button>
				<button type="button">Ano</button>
			</div>
			<strong>{currentDateLabel}</strong>
			<button type="button" class="primary-button" onclick={() => onSelectCalendarDate(selectedDate)}>
				+ Agendar neste dia
			</button>
		</div>

		<!-- Linha do tempo diaria: lista, adiciona e remove sessoes/eventos do dia. -->
		<DaySchedulerPanel
			{selectedDate}
			selectedDateLabel={currentDateLabel}
			{learners}
			{selectedLearnerId}
			{dayItems}
			onCreateSession={onCreateSession}
			onCreateEvent={onCreateEvent}
			onOpenLearner={onOpenLearner}
			onRemoveSession={onRemoveSession}
			onRemoveEvent={onRemoveEvent}
		/>
	</div>
</section>

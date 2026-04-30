<script lang="ts">
	import type { CalendarDay, Learner } from '$lib/modules/learners';
	import type {
		AgendaEvent,
		NewAgendaEventInput,
		NewSessionAppointmentInput,
		ScheduleItem
	} from '$lib/modules/scheduling';
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
	<div class="agenda-sidebar">
		<h2>Agenda</h2>
		<CalendarPanel
			days={calendarDays}
			{monthLabel}
			{selectedDate}
			variant="mini"
			onShiftMonth={onShiftMonth}
			onSelectDate={onSelectCalendarDate}
		/>

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
	</div>

	<div class="agenda-main">
		<div class="agenda-toolbar">
			<div class="view-switch">
				<button type="button" class="active">Dia</button>
				<button type="button">Semana</button>
				<button type="button">Mes</button>
			</div>
			<strong>{currentDateLabel}</strong>
			<button type="button" class="primary-button" onclick={() => onSelectCalendarDate(selectedDate)}>
				+ Agendar neste dia
			</button>
		</div>

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

<script lang="ts">
	import type { CalendarDay, Visit } from '../../domain/types';
	import CalendarPanel from '../calendar/CalendarPanel.svelte';
	import VisitEditor from './VisitEditor.svelte';

	let {
		calendarDays,
		monthLabel,
		selectedDate,
		selectedVisit,
		onShiftMonth,
		onSelectCalendarDate,
		onUpdateVisit,
		onRemoveVisit
	} = $props<{
		calendarDays: CalendarDay[];
		monthLabel: string;
		selectedDate: string;
		selectedVisit: Visit | null;
		onShiftMonth: (delta: number) => void;
		onSelectCalendarDate: (date: string) => void;
		onUpdateVisit: (visitId: string, patch: Partial<Visit>) => void;
		onRemoveVisit: (visitId: string) => void;
	}>();
</script>

<section class="tab-panel">
	<!-- Calendario do aprendente: seleciona uma visita especifica para edicao. -->
	<CalendarPanel
		days={calendarDays}
		{monthLabel}
		{selectedDate}
		onShiftMonth={onShiftMonth}
		onSelectDate={onSelectCalendarDate}
	/>

	<!-- Editor da visita selecionada: altera horario, status, local e observacoes. -->
	<VisitEditor
		visit={selectedVisit}
		onUpdateVisit={onUpdateVisit}
		onRemoveVisit={onRemoveVisit}
	/>
</section>

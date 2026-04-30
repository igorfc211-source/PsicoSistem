<script lang="ts">
	import type { CalendarDay } from '../../domain/types';

	let {
		days,
		monthLabel,
		selectedDate = '',
		variant = 'full',
		onShiftMonth,
		onSelectDate
	} = $props<{
		days: CalendarDay[];
		monthLabel: string;
		selectedDate?: string;
		variant?: 'full' | 'mini';
		onShiftMonth: (delta: number) => void;
		onSelectDate: (date: string) => void;
	}>();

	const isMini = $derived(variant === 'mini');
</script>

<div class={isMini ? 'mini-calendar' : 'calendar-panel'}>
	<!-- Cabecalho do mes: navegacao simples usada nos calendarios full e mini. -->
	<div class={isMini ? 'mini-head' : 'calendar-head'}>
		<button type="button" onclick={() => onShiftMonth(-1)}>&lt;</button>
		<strong>{monthLabel}</strong>
		<button type="button" onclick={() => onShiftMonth(1)}>&gt;</button>
	</div>

	<!-- Linha fixa de dias da semana, compartilhada entre agenda e prontuario. -->
	<div class={isMini ? 'mini-weekdays' : 'weekdays'}>
		<span>Seg</span>
		<span>Ter</span>
		<span>Qua</span>
		<span>Qui</span>
		<span>Sex</span>
		<span>Sab</span>
		<span>Dom</span>
	</div>

	<!-- Grade de datas: estados visuais indicam dia fora do mes, hoje, selecionado e com eventos. -->
	<div class={isMini ? 'mini-grid' : 'calendar-grid'}>
		{#each days as day}
			<button
				type="button"
				class:muted={!day.inMonth}
				class:today={day.isToday}
				class:selected={day.isSelected || selectedDate === day.date}
				class:hasVisit={day.visits.length + day.eventCount > 0}
				onclick={() => onSelectDate(day.date)}
				aria-label={`${day.date}, ${day.pendingVisitCount} visita pendente`}
			>
				{day.day}
				{#if day.pendingVisitCount > 0}
					<small class="calendar-badge">{day.pendingVisitCount}</small>
				{/if}
			</button>
		{/each}
	</div>
</div>

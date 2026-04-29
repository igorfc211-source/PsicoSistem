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
	<div class={isMini ? 'mini-head' : 'calendar-head'}>
		<button type="button" onclick={() => onShiftMonth(-1)}>&lt;</button>
		<strong>{monthLabel}</strong>
		<button type="button" onclick={() => onShiftMonth(1)}>&gt;</button>
	</div>

	<div class={isMini ? 'mini-weekdays' : 'weekdays'}>
		<span>D</span>
		<span>S</span>
		<span>T</span>
		<span>Q</span>
		<span>Q</span>
		<span>S</span>
		<span>S</span>
	</div>

	<div class={isMini ? 'mini-grid' : 'calendar-grid'}>
		{#each days as day}
			<button
				type="button"
				class:muted={!day.inMonth}
				class:today={day.isToday}
				class:selected={day.isSelected || selectedDate === day.date}
				class:hasVisit={day.visits.length + day.eventCount > 0}
				onclick={() => onSelectDate(day.date)}
			>
				{day.day}
			</button>
		{/each}
	</div>
</div>

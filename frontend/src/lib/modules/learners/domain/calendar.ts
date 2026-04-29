import { toDateInputValue } from './factories';
import type { CalendarDay, Learner, Visit } from './types';

// Gera sempre seis semanas para que o calendario nao mude de altura entre meses.
export function buildCalendarDays(
	monthDate: Date,
	visits: Visit[],
	eventDates: string[] = [],
	selectedDate = ''
): CalendarDay[] {
	const year = monthDate.getFullYear();
	const month = monthDate.getMonth();
	const firstDay = new Date(year, month, 1);
	const start = new Date(firstDay);
	start.setDate(firstDay.getDate() - firstDay.getDay());

	return Array.from({ length: 42 }, (_, index) => {
		const date = new Date(start);
		date.setDate(start.getDate() + index);
		const value = toDateInputValue(date);

			return {
				date: value,
				day: date.getDate(),
				inMonth: date.getMonth() === month,
				isToday: value === toDateInputValue(new Date()),
				isSelected: value === selectedDate,
				eventCount: eventDates.filter((eventDate) => eventDate === value).length,
				visits: visits.filter((visit) => visit.date === value)
			};
		});
}

// Calcula a frequencia considerando apenas sessoes marcadas como realizadas.
export function getAttendanceRate(learner: Learner) {
	if (learner.visits.length === 0) return 0;

	const completed = learner.visits.filter((visit) => visit.status === 'completed').length;
	return Math.round((completed / learner.visits.length) * 100);
}

// Encontra a proxima sessao futura, usando a primeira sessao como fallback quando necessario.
export function getNextVisit(learner: Learner) {
	const todayValue = toDateInputValue(new Date());
	return learner.visits.find((visit) => visit.date >= todayValue) ?? learner.visits[0] ?? null;
}

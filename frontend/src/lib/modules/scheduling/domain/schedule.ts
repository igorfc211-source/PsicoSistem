import type { AgendaEvent, ScheduleItem } from './types';
import type { LearnerVisitSchedule } from '$lib/modules/learners';

// Ordena compromissos como uma agenda real: primeiro horario, depois titulo.
export function sortScheduleItems(items: ScheduleItem[]) {
	return [...items].sort((left, right) => {
		const timeCompare = `${left.startTime}-${left.endTime}`.localeCompare(
			`${right.startTime}-${right.endTime}`
		);

		if (timeCompare !== 0) return timeCompare;
		return left.title.localeCompare(right.title);
	});
}

// Combina sessoes de aprendentes e eventos livres em uma unica linha do tempo diaria.
export function buildDayScheduleItems(
	date: string,
	visits: LearnerVisitSchedule[],
	events: AgendaEvent[]
): ScheduleItem[] {
	const sessionItems: ScheduleItem[] = visits
		.filter(({ visit }) => visit.date === date)
		.map(({ learner, visit }) => ({
			id: visit.id,
			date: visit.date,
			startTime: visit.startTime,
			endTime: visit.endTime,
			title: visit.title || 'Sessao individual',
			subtitle: `${learner.name} - ${visit.location || 'Consultorio'}`,
			kind: 'session',
			tone: 'purple',
			learner,
			visit
		}));

	const eventItems: ScheduleItem[] = events
		.filter((event) => event.date === date)
		.map((event) => ({
			id: event.id,
			date: event.date,
			startTime: event.startTime,
			endTime: event.endTime,
			title: event.title,
			subtitle: event.description || 'Evento da agenda',
			kind: 'event',
			tone: event.kind === 'block' ? 'gray' : event.kind === 'supervision' ? 'green' : 'amber',
			event
		}));

	return sortScheduleItems([...sessionItems, ...eventItems]);
}

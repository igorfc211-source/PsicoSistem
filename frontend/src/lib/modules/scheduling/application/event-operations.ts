import { createAgendaEvent } from '../domain/factories';
import type { AgendaEvent, NewAgendaEventInput } from '../domain/types';

// Valida o intervalo usado por sessoes e eventos livres.
export function isInvalidScheduleTime(startTime: string, endTime: string) {
	return !startTime || !endTime || endTime <= startTime;
}

// Ordena eventos livres pelo mesmo criterio visual da agenda profissional.
export function sortAgendaEvents(events: AgendaEvent[]) {
	return [...events].sort((left, right) =>
		`${left.date}-${left.startTime}`.localeCompare(`${right.date}-${right.startTime}`)
	);
}

// Cria um evento e devolve a lista pronta para persistencia.
export function addAgendaEvent(events: AgendaEvent[], input: NewAgendaEventInput) {
	const event = createAgendaEvent(input);

	return {
		event,
		events: sortAgendaEvents([...events, event])
	};
}

// Remove um evento livre sem tocar nas sessoes vinculadas a aprendentes.
export function removeAgendaEventById(events: AgendaEvent[], eventId: string) {
	return events.filter((event) => event.id !== eventId);
}

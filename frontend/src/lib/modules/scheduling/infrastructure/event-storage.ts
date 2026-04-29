import type { AgendaEvent } from '../domain/types';

export const AGENDA_EVENTS_STORAGE_KEY = 'psicosistem.agenda-events';

// Normaliza eventos antigos para manter compatibilidade com evolucoes do formulario.
function normalizeAgendaEvent(event: Partial<AgendaEvent>): AgendaEvent {
	const now = new Date().toISOString();

	return {
		id: event.id ?? '',
		title: event.title ?? 'Evento sem titulo',
		date: event.date ?? '',
		startTime: event.startTime ?? '09:00',
		endTime: event.endTime ?? '10:00',
		kind: event.kind ?? 'other',
		description: event.description ?? '',
		createdAt: event.createdAt ?? now,
		updatedAt: event.updatedAt ?? now
	};
}

// Carrega eventos livres da agenda profissional armazenados no navegador.
export function loadAgendaEvents() {
	const rawEvents = localStorage.getItem(AGENDA_EVENTS_STORAGE_KEY);
	if (!rawEvents) return [];

	try {
		const parsedEvents = JSON.parse(rawEvents) as Array<Partial<AgendaEvent>>;
		return parsedEvents.map(normalizeAgendaEvent).filter((event) => event.id && event.date);
	} catch {
		localStorage.removeItem(AGENDA_EVENTS_STORAGE_KEY);
		return [];
	}
}

// Persiste reunioes, bloqueios e eventos que nao pertencem a um aprendente especifico.
export function saveAgendaEvents(events: AgendaEvent[]) {
	localStorage.setItem(AGENDA_EVENTS_STORAGE_KEY, JSON.stringify(events));
}

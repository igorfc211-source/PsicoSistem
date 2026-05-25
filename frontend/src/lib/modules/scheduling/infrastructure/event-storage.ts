import type { AgendaEvent } from '../domain/types';

export const AGENDA_EVENTS_STORAGE_KEY = 'psicosistem.agenda-events';

function getScopedStorageKey(scope?: string | null) {
	return scope ? `${AGENDA_EVENTS_STORAGE_KEY}.${scope}` : AGENDA_EVENTS_STORAGE_KEY;
}

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
export function loadAgendaEvents(scope?: string | null) {
	const storageKey = getScopedStorageKey(scope);
	const rawEvents = localStorage.getItem(storageKey);
	if (!rawEvents) return [];

	try {
		const parsedEvents = JSON.parse(rawEvents) as Array<Partial<AgendaEvent>>;
		return parsedEvents.map(normalizeAgendaEvent).filter((event) => event.id && event.date);
	} catch {
		localStorage.removeItem(storageKey);
		return [];
	}
}

// Persiste reunioes, bloqueios e eventos que nao pertencem a um aprendente especifico.
export function saveAgendaEvents(events: AgendaEvent[], scope?: string | null) {
	localStorage.setItem(getScopedStorageKey(scope), JSON.stringify(events));
}

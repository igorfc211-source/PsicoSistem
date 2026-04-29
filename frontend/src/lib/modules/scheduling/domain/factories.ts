import type { AgendaEvent, NewAgendaEventInput } from './types';

// Cria ids proprios para eventos da agenda, separados dos ids de aprendentes e visitas.
export function createScheduleId(prefix: string) {
	const value = globalThis.crypto?.randomUUID
		? globalThis.crypto.randomUUID()
		: `${Date.now()}-${Math.random()}`;

	return `${prefix}-${value}`;
}

// Fabrica eventos livres, como reunioes, supervisoes ou bloqueios de horario.
export function createAgendaEvent(input: NewAgendaEventInput): AgendaEvent {
	const now = new Date().toISOString();

	return {
		id: createScheduleId('event'),
		title: input.title.trim(),
		date: input.date,
		startTime: input.startTime,
		endTime: input.endTime,
		kind: input.kind,
		description: input.description.trim(),
		createdAt: now,
		updatedAt: now
	};
}

// Retorna um nome amigavel para cada tipo de evento no painel diario.
export function getAgendaEventKindLabel(kind: AgendaEvent['kind']) {
	const labels: Record<AgendaEvent['kind'], string> = {
		meeting: 'Reuniao',
		supervision: 'Supervisao',
		block: 'Bloqueio',
		other: 'Evento'
	};

	return labels[kind];
}

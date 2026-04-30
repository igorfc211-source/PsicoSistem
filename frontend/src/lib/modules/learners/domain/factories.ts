import type { ActionPlan, Learner, NewLearnerInput, Visit } from './types';

// Gera ids previsiveis por dominio, mantendo unicidade para aprendentes, visitas e documentos.
export function createId(prefix: string) {
	const value = globalThis.crypto?.randomUUID
		? globalThis.crypto.randomUUID()
		: `${Date.now()}-${Math.random()}`;

	return `${prefix}-${value}`;
}

// Cria a estrutura padrao do plano de acao, pronta para receber secoes customizadas.
export function createEmptyActionPlan(): ActionPlan {
	return {
		educationProcess: '',
		familyGuidance: '',
		cognitiveSkills: '',
		behavior: '',
		clinicGoals: '',
		customFields: []
	};
}

// Monta o estado inicial do formulario de cadastro com um ciclo de um mes.
export function createDefaultLearnerInput(): NewLearnerInput {
	const today = new Date();
	const nextMonth = new Date();
	nextMonth.setMonth(nextMonth.getMonth() + 1);

	return {
		name: '',
		photoUrl: '',
		gender: '',
		guardian: '',
		age: '',
		status: 'active',
		startDate: toDateInputValue(today),
		endDate: toDateInputValue(nextMonth),
		visitCount: 8
	};
}

// Fabrica a entidade completa do aprendente a partir do formulario da interface.
export function createLearner(input: NewLearnerInput): Learner {
	const now = new Date().toISOString();

	return {
		id: createId('learner'),
		name: input.name.trim(),
		photoUrl: input.photoUrl,
		gender: input.gender.trim(),
		guardian: input.guardian.trim(),
		age: input.age.trim(),
		status: input.status,
		startDate: input.startDate,
		endDate: input.endDate,
		visitCount: input.visitCount,
		anamnese: '',
		anamneseDocuments: [],
		actionPlan: createEmptyActionPlan(),
		visits: buildVisits(input.startDate, input.endDate, input.visitCount),
		documents: [],
		reports: [],
		createdAt: now,
		updatedAt: now
	};
}

// Distribui as visitas entre a data inicial e final para criar uma agenda base editavel.
export function buildVisits(startDate: string, endDate: string, visitCount: number): Visit[] {
	if (!startDate || !endDate || visitCount <= 0) return [];

	const start = new Date(`${startDate}T12:00:00`);
	const end = new Date(`${endDate}T12:00:00`);
	if (Number.isNaN(start.getTime()) || Number.isNaN(end.getTime()) || end < start) {
		return [];
	}

	const totalDays = Math.max(1, Math.round((end.getTime() - start.getTime()) / 86400000));
	const spacing = visitCount === 1 ? 0 : totalDays / (visitCount - 1);

	return Array.from({ length: visitCount }, (_, index) => {
		const date = new Date(start);
		date.setDate(start.getDate() + Math.round(spacing * index));
		const startTime = buildDefaultVisitStartTime(index);

		return {
			id: createId('visit'),
			date: toDateInputValue(date),
			title: 'Sessao individual',
			startTime,
			endTime: addMinutesToTime(startTime, 50),
			kind: 'session',
			location: 'Consultorio',
			status: 'scheduled',
			notes: ''
		};
	});
}

// Normaliza datas para o formato esperado por inputs nativos e chaves do calendario.
export function toDateInputValue(date: Date) {
	return date.toISOString().slice(0, 10);
}

// Sugere horarios comerciais variados para a agenda inicial nao ficar toda no mesmo horario.
function buildDefaultVisitStartTime(index: number) {
	const baseMinutes = 8 * 60;
	const slotMinutes = baseMinutes + (index % 8) * 60;
	return minutesToTime(slotMinutes);
}

// Soma minutos a um horario HH:mm, usado para gerar duracoes padrao de sessao.
export function addMinutesToTime(time: string, minutes: number) {
	const [hour = '0', minute = '0'] = time.split(':');
	const totalMinutes = Number(hour) * 60 + Number(minute) + minutes;

	return minutesToTime(totalMinutes);
}

function minutesToTime(totalMinutes: number) {
	const normalizedMinutes = ((totalMinutes % 1440) + 1440) % 1440;
	const hour = Math.floor(normalizedMinutes / 60);
	const minute = normalizedMinutes % 60;

	return `${String(hour).padStart(2, '0')}:${String(minute).padStart(2, '0')}`;
}

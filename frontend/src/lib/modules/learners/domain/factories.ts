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
		visitCount: 8,
		sessionPriceCents: 0
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
		sessionPriceCents: normalizeAmountCents(input.sessionPriceCents),
		anamnese: '',
		anamneseDocuments: [],
		actionPlan: createEmptyActionPlan(),
		visits: buildVisits(),
		documents: [],
		reports: [],
		createdAt: now,
		updatedAt: now
	};
}

// Distribui as visitas entre a data inicial e final para criar uma agenda base editavel.
export function buildVisits() {
	return [];
}
// Normaliza datas para o formato esperado por inputs nativos e chaves do calendario.
export function toDateInputValue(date: Date) {
	return date.toISOString().slice(0, 10);
}

function normalizeAmountCents(value: number) {
	if (!Number.isFinite(value)) return 0;
	return Math.max(0, Math.round(value));
}

// Sugere horarios comerciais variados para a agenda inicial nao ficar toda no mesmo horario.
function buildDefaultVisitStartTime(index: number) {
	const baseMinutes = 8 * 60;
	const slotMinutes = baseMinutes + (index % 8) * 60;
	return minutesToTime(slotMinutes);
}


function minutesToTime(totalMinutes: number) {
	const normalizedMinutes = ((totalMinutes % 1440) + 1440) % 1440;
	const hour = Math.floor(normalizedMinutes / 60);
	const minute = normalizedMinutes % 60;

	return `${String(hour).padStart(2, '0')}:${String(minute).padStart(2, '0')}`;
}

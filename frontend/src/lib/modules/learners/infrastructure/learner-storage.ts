import { createEmptyActionPlan } from '../domain/factories';
import type { Learner, Visit } from '../domain/types';

export const LEARNERS_STORAGE_KEY = 'psicosistem.learners';

// Corrige dados antigos do localStorage para o contrato atual do dominio.
function normalizeLearner(learner: Partial<Learner>): Learner {
	const now = new Date().toISOString();

	return {
		id: learner.id ?? '',
		name: learner.name ?? '',
		gender: learner.gender ?? '',
		guardian: learner.guardian ?? '',
		age: learner.age ?? '',
		status: learner.status ?? 'active',
		startDate: learner.startDate ?? '',
		endDate: learner.endDate ?? '',
		visitCount: learner.visitCount ?? learner.visits?.length ?? 0,
		anamnese: learner.anamnese ?? '',
		actionPlan: {
			...createEmptyActionPlan(),
			...(learner.actionPlan ?? {})
		},
		visits: (learner.visits ?? []).map(normalizeVisit),
		documents: learner.documents ?? [],
		reports: learner.reports ?? [],
		createdAt: learner.createdAt ?? now,
		updatedAt: learner.updatedAt ?? now
	};
}

// Completa visitas antigas com campos profissionais de horario, titulo e local.
function normalizeVisit(visit: Partial<Visit>): Visit {
	return {
		id: visit.id ?? '',
		date: visit.date ?? '',
		title: visit.title ?? 'Sessao individual',
		startTime: visit.startTime ?? '09:00',
		endTime: visit.endTime ?? '09:50',
		kind: visit.kind ?? 'session',
		location: visit.location ?? 'Consultorio',
		status: visit.status ?? 'scheduled',
		notes: visit.notes ?? ''
	};
}

// Carrega aprendentes persistidos localmente e limpa o cache se o JSON estiver invalido.
export function loadLearners() {
	const rawLearners = localStorage.getItem(LEARNERS_STORAGE_KEY);
	if (!rawLearners) return [];

	try {
		const parsedLearners = JSON.parse(rawLearners) as Array<Partial<Learner>>;
		return parsedLearners.map(normalizeLearner).filter((learner) => learner.id);
	} catch {
		localStorage.removeItem(LEARNERS_STORAGE_KEY);
		return [];
	}
}

// Persiste o snapshot completo para manter a interface responsiva sem backend de prontuario ainda.
export function saveLearners(learners: Learner[]) {
	localStorage.setItem(LEARNERS_STORAGE_KEY, JSON.stringify(learners));
}

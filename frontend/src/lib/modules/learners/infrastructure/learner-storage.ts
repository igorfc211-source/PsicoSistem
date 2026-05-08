import { createEmptyActionPlan } from '../domain/factories';
import type { Learner, LearnerGuardian, Visit } from '../domain/types';

export const LEARNERS_STORAGE_KEY = 'psicosistem.learners';

// Corrige dados antigos do localStorage para o contrato atual do dominio.
function normalizeLearner(learner: Partial<Learner>): Learner {
	const now = new Date().toISOString();
	const actionPlan = {
		...createEmptyActionPlan(),
		...(learner.actionPlan ?? {})
	};
	actionPlan.customFields = actionPlan.customFields ?? [];
	const guardians = normalizeLearnerGuardians(learner);
	const primaryGuardian = guardians[0] ?? null;

	return {
		id: learner.id ?? '',
		name: learner.name ?? '',
		photoUrl: learner.photoUrl ?? '',
		gender: learner.gender ?? '',
		guardian: primaryGuardian?.name ?? learner.guardian ?? '',
		guardianRelationship: primaryGuardian?.relationship ?? learner.guardianRelationship ?? '',
		guardians,
		age: learner.age ?? '',
		status: learner.status ?? 'active',
		startDate: learner.startDate ?? '',
		endDate: learner.endDate ?? '',
		visitCount: learner.visitCount ?? learner.visits?.length ?? 0,
		sessionPriceCents: normalizeAmountCents(learner.sessionPriceCents),
		generalValueCents: normalizeAmountCents(learner.generalValueCents),
		anamnese: learner.anamnese ?? '',
		anamneseDocuments: learner.anamneseDocuments ?? [],
		actionPlan,
		visits: (learner.visits ?? []).map(normalizeVisit),
		documents: learner.documents ?? [],
		reports: learner.reports ?? [],
		createdAt: learner.createdAt ?? now,
		updatedAt: learner.updatedAt ?? now
	};
}

function normalizeLearnerGuardians(learner: Partial<Learner>): LearnerGuardian[] {
	const rawGuardians =
		learner.guardians && learner.guardians.length > 0
			? learner.guardians
			: [
					{
						id: '',
						sourceKey: '',
						name: learner.guardian ?? '',
						relationship: learner.guardianRelationship ?? '',
						phone: ''
					}
				];
	const seenKeys = new Set<string>();

	return rawGuardians
		.map((guardian) => {
			const name = guardian.name ?? '';
			const sourceKey = guardian.sourceKey || normalizeGuardianKey(name);

			return {
				id: guardian.id || `learner-guardian-${sourceKey}`,
				sourceKey,
				name: name.trim(),
				relationship: guardian.relationship ?? '',
				phone: normalizePhone(guardian.phone ?? '')
			};
		})
		.filter((guardian) => {
			if (!guardian.name || !guardian.sourceKey || seenKeys.has(guardian.sourceKey)) return false;
			seenKeys.add(guardian.sourceKey);
			return true;
		})
		.slice(0, 2);
}

function normalizeAmountCents(value: unknown) {
	const numericValue = typeof value === 'number' ? value : Number(value ?? 0);
	if (!Number.isFinite(numericValue)) return 0;

	return Math.max(0, Math.round(numericValue));
}

function normalizePhone(value: string) {
	return value.replace(/\D/g, '').slice(0, 11);
}

function normalizeGuardianKey(value: string) {
	return value.trim().toLowerCase().replace(/\s+/g, ' ');
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

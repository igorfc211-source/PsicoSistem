import { createEmptyActionPlan } from '../domain/factories';
import type { Learner, LearnerGuardian, LearnerStatus } from '../domain/types';

export type BackendLearner = {
	id?: string;
	tenant_id?: string;
	name?: string;
	photo_url?: string;
	gender?: string;
	guardian?: string;
	age?: string;
	status?: string;
	start_date?: string;
	end_date?: string;
	visit_count?: number;
	session_price_cents?: number;
	general_value_cents?: number;
	guardian_ids?: string[];
	created_at?: string;
	updated_at?: string;
};

type ApiResponse<T> = {
	data: T | null;
	meta: unknown;
	error: {
		code: string;
		message: string;
	} | null;
};

export async function fetchBackendLearners(token: string) {
	if (!token) return [];

	const response = await fetch('/api/learners?per_page=100', {
		headers: {
			authorization: `Bearer ${token}`
		}
	});
	const payload = (await response.json()) as ApiResponse<BackendLearner[]>;

	if (!response.ok || payload.error) {
		throw new Error(payload.error?.message ?? 'Nao foi possivel carregar aprendentes do backend.');
	}

	return Array.isArray(payload.data) ? payload.data.filter((learner) => learner.id) : [];
}

export async function updateBackendLearner(token: string, learner: Learner) {
	if (!token || !learner.backendId) return null;

	const response = await fetch(`/api/learners/${learner.backendId}`, {
		method: 'PATCH',
		headers: {
			authorization: `Bearer ${token}`,
			'content-type': 'application/json'
		},
		body: JSON.stringify(toBackendLearnerInput(learner))
	});
	const payload = (await response.json()) as ApiResponse<BackendLearner>;

	if (!response.ok || payload.error) {
		throw new Error(payload.error?.message ?? 'Nao foi possivel atualizar o aprendente no backend.');
	}

	return payload.data;
}

export function mergeBackendLearners(localLearners: Learner[], backendLearners: BackendLearner[]) {
	const mergedLearners = [...localLearners];

	for (const backendLearner of backendLearners) {
		if (!backendLearner.id) continue;

		const localIndex = findLocalLearnerIndex(mergedLearners, backendLearner);
		const currentLearner = localIndex >= 0 ? mergedLearners[localIndex] : null;
		const nextLearner = normalizeBackendLearner(backendLearner, currentLearner);

		if (localIndex >= 0) {
			mergedLearners[localIndex] = nextLearner;
		} else {
			mergedLearners.push(nextLearner);
		}
	}

	return mergedLearners.filter((learner) => learner.id && learner.name);
}

function findLocalLearnerIndex(localLearners: Learner[], backendLearner: BackendLearner) {
	const backendId = backendLearner.id ?? '';
	const idIndex = localLearners.findIndex(
		(learner) => learner.id === backendId || learner.backendId === backendId
	);
	if (idIndex >= 0) return idIndex;

	const backendName = normalizeName(backendLearner.name ?? '');
	if (!backendName) return -1;

	const matchingIndexes = localLearners
		.map((learner, index) => ({ learner, index }))
		.filter(({ learner }) => normalizeName(learner.name) === backendName);

	return matchingIndexes.length === 1 ? matchingIndexes[0].index : -1;
}

function normalizeBackendLearner(backendLearner: BackendLearner, currentLearner: Learner | null) {
	const now = new Date().toISOString();
	const backendId = backendLearner.id ?? currentLearner?.backendId ?? currentLearner?.id ?? '';
	const guardianName = backendLearner.guardian?.trim() || currentLearner?.guardian || '';
	const guardians =
		currentLearner?.guardians && currentLearner.guardians.length > 0
			? currentLearner.guardians
			: buildBackendGuardianFallbacks(backendLearner.guardian_ids ?? [], guardianName);

	return {
		id: currentLearner?.id ?? backendId,
		backendId,
		backendGuardianIds: backendLearner.guardian_ids ?? currentLearner?.backendGuardianIds ?? [],
		name: backendLearner.name?.trim() || currentLearner?.name || '',
		photoUrl: backendLearner.photo_url ?? currentLearner?.photoUrl ?? '',
		gender: backendLearner.gender?.trim() || currentLearner?.gender || '',
		guardian: guardianName,
		guardianRelationship: currentLearner?.guardianRelationship ?? '',
		guardians,
		age: backendLearner.age?.trim() || currentLearner?.age || '',
		status: normalizeStatus(backendLearner.status, currentLearner?.status),
		startDate: backendLearner.start_date ?? currentLearner?.startDate ?? '',
		endDate: backendLearner.end_date ?? currentLearner?.endDate ?? '',
		visitCount: normalizeNonNegativeInteger(
			backendLearner.visit_count ?? currentLearner?.visitCount ?? currentLearner?.visits.length ?? 0
		),
		sessionPriceCents: normalizeAmountCents(
			backendLearner.session_price_cents ?? currentLearner?.sessionPriceCents ?? 0
		),
		generalValueCents: normalizeAmountCents(
			backendLearner.general_value_cents ?? currentLearner?.generalValueCents ?? 0
		),
		anamnese: currentLearner?.anamnese ?? '',
		anamneseDocuments: currentLearner?.anamneseDocuments ?? [],
		actionPlan: currentLearner?.actionPlan ?? createEmptyActionPlan(),
		visits: currentLearner?.visits ?? [],
		documents: currentLearner?.documents ?? [],
		reports: currentLearner?.reports ?? [],
		createdAt: backendLearner.created_at ?? currentLearner?.createdAt ?? now,
		updatedAt: backendLearner.updated_at ?? currentLearner?.updatedAt ?? now
	} satisfies Learner;
}

function toBackendLearnerInput(learner: Learner) {
	const guardianIds = learner.backendGuardianIds ?? [];
	if (guardianIds.length < 1) {
		throw new Error('Este aprendente nao tem responsaveis do backend vinculados para atualizar.');
	}

	return {
		name: learner.name,
		photo_url: learner.photoUrl,
		gender: learner.gender,
		guardian: learner.guardian,
		age: learner.age,
		status: learner.status,
		start_date: learner.startDate,
		end_date: learner.endDate,
		visit_count: learner.visitCount,
		session_price_cents: learner.sessionPriceCents,
		general_value_cents: learner.generalValueCents,
		guardian_ids: guardianIds.slice(0, 2)
	};
}

function buildBackendGuardianFallbacks(guardianIds: string[], guardianName: string): LearnerGuardian[] {
	return guardianIds.slice(0, 2).map((guardianId, index) => ({
		id: `backend-guardian-${guardianId}`,
		sourceKey: guardianId,
		name: guardianName || `Responsavel ${index + 1}`,
		relationship: '',
		phone: ''
	}));
}

function normalizeStatus(status: string | undefined, fallback: LearnerStatus | undefined): LearnerStatus {
	return status === 'inactive' || status === 'active' ? status : fallback ?? 'active';
}

function normalizeAmountCents(value: unknown) {
	const numericValue = typeof value === 'number' ? value : Number(value ?? 0);
	if (!Number.isFinite(numericValue)) return 0;

	return Math.max(0, Math.round(numericValue));
}

function normalizeNonNegativeInteger(value: unknown) {
	const numericValue = typeof value === 'number' ? value : Number(value ?? 0);
	if (!Number.isFinite(numericValue)) return 0;

	return Math.max(0, Math.round(numericValue));
}

function normalizeName(value: string) {
	return value.trim().toLowerCase().replace(/\s+/g, ' ');
}

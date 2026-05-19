import {
	COMMUNICATION_STAGES,
	type CommunicationContact,
	type CommunicationFamily,
	type CommunicationResponsible,
	type CommunicationStage,
	type ContactChannelType
} from '../domain/communications';
import { createId } from '../domain/factories';

export const COMMUNICATION_FAMILIES_STORAGE_KEY = 'psicosistem.communication-families';
export const HIDDEN_COMMUNICATION_SOURCE_KEYS_STORAGE_KEY =
	'psicosistem.hidden-communication-source-keys';

const CONTACT_TYPES = new Set<ContactChannelType>([
	'phone',
	'whatsapp',
	'instagram',
	'email',
	'other'
]);
const STAGES = new Set<CommunicationStage>(COMMUNICATION_STAGES.map((stage) => stage.value));

function getScopedStorageKey(key: string, scope?: string | null) {
	return scope ? `${key}.${scope}` : key;
}

export function loadCommunicationFamilies(scope?: string | null) {
	const storageKey = getScopedStorageKey(COMMUNICATION_FAMILIES_STORAGE_KEY, scope);
	const rawFamilies = localStorage.getItem(storageKey);
	if (!rawFamilies) return [];

	try {
		const parsedFamilies = JSON.parse(rawFamilies) as Array<Partial<CommunicationFamily>>;
		return parsedFamilies.map(normalizeFamily).filter((family) => family.id && family.familyName);
	} catch {
		localStorage.removeItem(storageKey);
		return [];
	}
}

export function saveCommunicationFamilies(families: CommunicationFamily[], scope?: string | null) {
	localStorage.setItem(
		getScopedStorageKey(COMMUNICATION_FAMILIES_STORAGE_KEY, scope),
		JSON.stringify(families)
	);
}

export function loadHiddenCommunicationSourceKeys(scope?: string | null) {
	const storageKey = getScopedStorageKey(HIDDEN_COMMUNICATION_SOURCE_KEYS_STORAGE_KEY, scope);
	const rawKeys = localStorage.getItem(storageKey);
	if (!rawKeys) return [];

	try {
		const parsedKeys = JSON.parse(rawKeys) as string[];
		return normalizeIds(parsedKeys);
	} catch {
		localStorage.removeItem(storageKey);
		return [];
	}
}

export function saveHiddenCommunicationSourceKeys(keys: string[], scope?: string | null) {
	localStorage.setItem(
		getScopedStorageKey(HIDDEN_COMMUNICATION_SOURCE_KEYS_STORAGE_KEY, scope),
		JSON.stringify(normalizeIds(keys))
	);
}

function normalizeFamily(family: Partial<CommunicationFamily>): CommunicationFamily {
	const now = new Date().toISOString();
	const stage = family.stage && STAGES.has(family.stage) ? family.stage : 'novo';

	return {
		id: family.id ?? createId('family'),
		familyName: family.familyName ?? '',
		sourceGuardianKey: family.sourceGuardianKey,
		stage,
		learnerIds: normalizeIds(family.learnerIds),
		responsibles: normalizeResponsibles(family.responsibles),
		contacts: normalizeContacts(family.contacts),
		nextStep: family.nextStep ?? '',
		notes: family.notes ?? '',
		lastContactAt: family.lastContactAt ?? '',
		createdAt: family.createdAt ?? now,
		updatedAt: family.updatedAt ?? now
	};
}

function normalizeResponsibles(input?: CommunicationResponsible[]) {
	return (input ?? []).map((item) => ({
		id: item.id ?? createId('responsible'),
		name: item.name ?? '',
		relationship: item.relationship ?? '',
		phone: item.phone ?? ''
	}));
}

function normalizeContacts(input?: CommunicationContact[]) {
	return (input ?? []).map((item) => ({
		id: item.id ?? createId('contact'),
		type: item.type && CONTACT_TYPES.has(item.type) ? item.type : 'other',
		label: item.label ?? '',
		value: item.value ?? '',
		notes: item.notes ?? ''
	}));
}

function normalizeIds(ids?: string[]) {
	return Array.from(new Set((ids ?? []).filter(Boolean)));
}

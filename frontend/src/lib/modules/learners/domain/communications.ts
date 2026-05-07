import { createId } from './factories';
import type { Learner } from './types';

export type CommunicationStage = 'novo' | 'em_contato' | 'aguardando' | 'acompanhamento';
export type ContactChannelType = 'phone' | 'whatsapp' | 'instagram' | 'email' | 'other';

export type CommunicationResponsible = {
	id: string;
	name: string;
	relationship: string;
	phone: string;
};

export type CommunicationContact = {
	id: string;
	type: ContactChannelType;
	label: string;
	value: string;
	notes: string;
};

export type CommunicationFamily = {
	id: string;
	familyName: string;
	sourceGuardianKey?: string;
	stage: CommunicationStage;
	learnerIds: string[];
	responsibles: CommunicationResponsible[];
	contacts: CommunicationContact[];
	nextStep: string;
	notes: string;
	lastContactAt: string;
	createdAt: string;
	updatedAt: string;
};

export type NewCommunicationFamilyInput = {
	familyName: string;
	responsibleName: string;
	responsiblePhone: string;
	relationship: string;
	learnerIds: string[];
};

export type NewCommunicationResponsibleInput = {
	name: string;
	relationship: string;
	phone: string;
};

export type NewCommunicationContactInput = {
	type: ContactChannelType;
	label: string;
	value: string;
	notes: string;
};

export const COMMUNICATION_STAGES: Array<{ value: CommunicationStage; label: string }> = [
	{ value: 'novo', label: 'Novo' },
	{ value: 'em_contato', label: 'Em contato' },
	{ value: 'aguardando', label: 'Aguardando' },
	{ value: 'acompanhamento', label: 'Acompanhamento' }
];

export const CONTACT_CHANNEL_TYPES: Array<{ value: ContactChannelType; label: string }> = [
	{ value: 'phone', label: 'Telefone' },
	{ value: 'whatsapp', label: 'WhatsApp' },
	{ value: 'instagram', label: 'Instagram' },
	{ value: 'email', label: 'E-mail' },
	{ value: 'other', label: 'Outro' }
];

export function createCommunicationFamily(input: NewCommunicationFamilyInput): CommunicationFamily {
	const now = new Date().toISOString();
	const responsibleName = input.responsibleName.trim();

	return {
		id: createId('family'),
		familyName: input.familyName.trim(),
		stage: 'novo',
		learnerIds: normalizeIds(input.learnerIds),
		responsibles: responsibleName
			? [
					{
						id: createId('responsible'),
						name: responsibleName,
						relationship: input.relationship.trim(),
						phone: normalizePhone(input.responsiblePhone)
					}
				]
			: [],
		contacts: [],
		nextStep: '',
		notes: '',
		lastContactAt: '',
		createdAt: now,
		updatedAt: now
	};
}

export function addResponsibleToFamily(
	family: CommunicationFamily,
	input: NewCommunicationResponsibleInput
): CommunicationFamily {
	const now = new Date().toISOString();
	const responsible: CommunicationResponsible = {
		id: createId('responsible'),
		name: input.name.trim(),
		relationship: input.relationship.trim(),
		phone: normalizePhone(input.phone)
	};

	return {
		...family,
		responsibles: [...family.responsibles, responsible],
		updatedAt: now
	};
}

export function addContactToFamily(
	family: CommunicationFamily,
	input: NewCommunicationContactInput
): CommunicationFamily {
	const now = new Date().toISOString();
	const contact: CommunicationContact = {
		id: createId('contact'),
		type: input.type,
		label: input.label.trim(),
		value: normalizeContactValue(input.type, input.value),
		notes: input.notes.trim()
	};

	return {
		...family,
		contacts: [...family.contacts, contact],
		lastContactAt: now,
		updatedAt: now
	};
}

export function removeResponsibleFromFamily(
	family: CommunicationFamily,
	responsibleId: string
): CommunicationFamily {
	return {
		...family,
		responsibles: family.responsibles.filter((item) => item.id !== responsibleId),
		updatedAt: new Date().toISOString()
	};
}

export function removeContactFromFamily(
	family: CommunicationFamily,
	contactId: string
): CommunicationFamily {
	return {
		...family,
		contacts: family.contacts.filter((item) => item.id !== contactId),
		updatedAt: new Date().toISOString()
	};
}

export function syncCommunicationFamiliesWithLearners(
	families: CommunicationFamily[],
	learners: Learner[],
	hiddenSourceGuardianKeys: string[] = []
): CommunicationFamily[] {
	const hiddenGuardianKeys = new Set(hiddenSourceGuardianKeys);
	const validLearnerIds = new Set(learners.map((learner) => learner.id));
	const nextFamilies = families
		.filter((family) => !family.sourceGuardianKey || !hiddenGuardianKeys.has(family.sourceGuardianKey))
		.map((family) => ({
			...family,
			learnerIds: normalizeIds(family.learnerIds).filter((id) => validLearnerIds.has(id))
		}));
	const familiesByGuardianKey = new Map<string, CommunicationFamily>();

	for (const family of nextFamilies) {
		if (family.sourceGuardianKey) {
			familiesByGuardianKey.set(family.sourceGuardianKey, family);
		}
	}

	for (const learner of learners) {
		const guardianName = learner.guardian.trim();
		const sourceGuardianKey = normalizeKey(guardianName);
		if (!sourceGuardianKey) continue;
		if (hiddenGuardianKeys.has(sourceGuardianKey)) continue;

		const existingFamily = familiesByGuardianKey.get(sourceGuardianKey);
		if (existingFamily) {
			if (!existingFamily.learnerIds.includes(learner.id)) {
				existingFamily.learnerIds = [...existingFamily.learnerIds, learner.id];
				existingFamily.updatedAt = new Date().toISOString();
			}
			continue;
		}

		const createdFamily = createFamilyFromLearner(learner, sourceGuardianKey);
		nextFamilies.push(createdFamily);
		familiesByGuardianKey.set(sourceGuardianKey, createdFamily);
	}

	return nextFamilies;
}

export function getFamilyLearners(family: CommunicationFamily, learners: Learner[]) {
	const learnerMap = new Map(learners.map((learner) => [learner.id, learner]));
	return family.learnerIds.map((id) => learnerMap.get(id)).filter((item): item is Learner => Boolean(item));
}

export function getContactTypeLabel(type: ContactChannelType) {
	return CONTACT_CHANNEL_TYPES.find((item) => item.value === type)?.label ?? 'Contato';
}

export function normalizePhone(value: string) {
	return value.replace(/\D/g, '');
}

export function isValidPhoneNumber(value: string) {
	return /^\d{10,11}$/.test(normalizePhone(value));
}

export function normalizeInstagramHandle(value: string) {
	const handle = value
		.trim()
		.replace(/^https?:\/\/(www\.)?instagram\.com\//, '')
		.replace(/\/+$/, '')
		.replace(/^@/, '');

	return handle ? `@${handle}` : '';
}

export function isValidInstagramHandle(value: string) {
	return /^@[a-zA-Z0-9._]{2,30}$/.test(normalizeInstagramHandle(value));
}

export function isValidEmailAddress(value: string) {
	return /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/.test(value.trim());
}

function normalizeContactValue(type: ContactChannelType, value: string) {
	if (type === 'phone' || type === 'whatsapp') return normalizePhone(value);
	if (type === 'instagram') return normalizeInstagramHandle(value);
	if (type === 'email') return value.trim().toLowerCase();
	return value.trim();
}

function createFamilyFromLearner(learner: Learner, sourceGuardianKey: string): CommunicationFamily {
	const now = new Date().toISOString();
	const guardianName = learner.guardian.trim();

	return {
		id: createId('family'),
		familyName: inferFamilyName(guardianName, learner.name),
		sourceGuardianKey,
		stage: 'novo',
		learnerIds: [learner.id],
		responsibles: [
			{
				id: createId('responsible'),
				name: guardianName,
				relationship: '',
				phone: ''
			}
		],
		contacts: [],
		nextStep: '',
		notes: '',
		lastContactAt: '',
		createdAt: now,
		updatedAt: now
	};
}

function inferFamilyName(guardianName: string, learnerName: string) {
	const referenceName = guardianName || learnerName;
	const mainName = referenceName.split(/\s+e\s+|,|\/|&/i)[0]?.trim() ?? referenceName;
	const tokens = mainName.split(/\s+/).filter(Boolean);
	const surname = tokens.at(-1) ?? 'Contato';

	return `Familia ${surname}`;
}

function normalizeIds(ids: string[]) {
	return Array.from(new Set(ids.filter(Boolean)));
}

function normalizeKey(value: string) {
	return value.trim().toLowerCase().replace(/\s+/g, ' ');
}

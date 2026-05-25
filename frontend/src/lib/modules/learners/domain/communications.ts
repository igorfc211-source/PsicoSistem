import { createId } from './factories';
import type { Learner, LearnerGuardian } from './types';

export type CommunicationStage =
	| 'novo'
	| 'em_contato'
	| 'aguardando'
	| 'acompanhamento'
	| 'inativo';
export type ContactChannelType = 'phone' | 'whatsapp' | 'instagram' | 'email' | 'other';

export type GuardianOption = {
	key: string;
	name: string;
	relationship: string;
	phone: string;
	learnerIds: string[];
};

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
	learnerIds: string[];
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
	{ value: 'acompanhamento', label: 'Acompanhamento' },
	{ value: 'inativo', label: 'Inativo' }
];

export const CONTACT_CHANNEL_TYPES: Array<{ value: ContactChannelType; label: string }> = [
	{ value: 'phone', label: 'Telefone' },
	{ value: 'whatsapp', label: 'WhatsApp' },
	{ value: 'instagram', label: 'Instagram' },
	{ value: 'email', label: 'E-mail' },
	{ value: 'other', label: 'Outro' }
];

export const RELATIONSHIP_OPTIONS = [
	'Mae',
	'Pai',
	'Responsavel legal',
	'Avo',
	'Avo materna',
	'Avo paterna',
	'Tia',
	'Tio',
	'Madrasta',
	'Padrasto',
	'Irma',
	'Irmao',
	'Outro'
];

export function createCommunicationFamily(input: NewCommunicationFamilyInput): CommunicationFamily {
	const now = new Date().toISOString();
	const responsibleName = input.responsibleName.trim();
	const responsibleKey = normalizeResponsibleKey(responsibleName);

	return {
		id: createId('family'),
		familyName: input.familyName.trim(),
		sourceGuardianKey: responsibleKey || undefined,
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
		learnerIds: normalizeIds([...family.learnerIds, ...input.learnerIds]),
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
	const nextFamilies = mergeFamiliesByResponsible(
		families
		.filter((family) => !family.sourceGuardianKey || !hiddenGuardianKeys.has(family.sourceGuardianKey))
		.map((family) => ({
			...family,
			learnerIds: normalizeIds(family.learnerIds).filter((id) => validLearnerIds.has(id))
		}))
	);
	const familiesByGuardianKey = new Map<string, CommunicationFamily>();
	const familiesByResponsibleKey = new Map<string, CommunicationFamily>();

	for (const family of nextFamilies) {
		if (family.sourceGuardianKey) {
			familiesByGuardianKey.set(family.sourceGuardianKey, family);
		}
		for (const key of getCommunicationFamilyResponsibleKeys(family)) {
			familiesByResponsibleKey.set(key, family);
		}
	}

	for (const learner of learners) {
		const learnerGuardians = getLearnerGuardianEntries(learner).filter(
			(guardian) => !hiddenGuardianKeys.has(guardian.sourceKey)
		);
		if (!learnerGuardians.length) continue;

		const existingFamily = learnerGuardians
			.map((guardian) => familiesByGuardianKey.get(guardian.sourceKey) ?? familiesByResponsibleKey.get(guardian.sourceKey))
			.find(Boolean);
		if (existingFamily) {
			if (!existingFamily.learnerIds.includes(learner.id)) {
				existingFamily.learnerIds = [...existingFamily.learnerIds, learner.id];
				existingFamily.updatedAt = new Date().toISOString();
			}
			if (!existingFamily.sourceGuardianKey) {
				existingFamily.sourceGuardianKey = learnerGuardians[0].sourceKey;
			}
			existingFamily.responsibles = upsertFamilyResponsibles(existingFamily, learnerGuardians);
			for (const guardian of learnerGuardians) {
				familiesByGuardianKey.set(guardian.sourceKey, existingFamily);
				familiesByResponsibleKey.set(guardian.sourceKey, existingFamily);
			}
			continue;
		}

		const createdFamily = createFamilyFromLearner(learner, learnerGuardians);
		nextFamilies.push(createdFamily);
		for (const guardian of learnerGuardians) {
			familiesByGuardianKey.set(guardian.sourceKey, createdFamily);
			familiesByResponsibleKey.set(guardian.sourceKey, createdFamily);
		}
	}

	return mergeFamiliesByResponsible(nextFamilies);
}

export function buildGuardianOptionsFromLearners(
	learners: Learner[],
	families: CommunicationFamily[] = []
): GuardianOption[] {
	const optionsByKey = new Map<string, GuardianOption>();

	for (const learner of learners) {
		for (const guardian of getLearnerGuardianEntries(learner)) {
			upsertGuardianOption(optionsByKey, {
				key: guardian.sourceKey,
				name: guardian.name,
				relationship: guardian.relationship,
				phone: guardian.phone,
				learnerIds: [learner.id]
			});
		}
	}

	for (const family of families) {
		for (const responsible of family.responsibles) {
			const key = normalizeResponsibleKey(responsible.name);
			if (!key) continue;

			upsertGuardianOption(optionsByKey, {
				key,
				name: responsible.name,
				relationship: responsible.relationship,
				phone: responsible.phone,
				learnerIds: family.learnerIds
			});
		}
	}

	return Array.from(optionsByKey.values()).sort((left, right) =>
		left.name.localeCompare(right.name)
	);
}

export function getCommunicationFamilyResponsibleKeys(family: CommunicationFamily) {
	const keys = new Set<string>();
	if (family.sourceGuardianKey) keys.add(family.sourceGuardianKey);

	for (const responsible of family.responsibles) {
		const key = normalizeResponsibleKey(responsible.name);
		if (key) keys.add(key);
	}

	return Array.from(keys);
}

export function normalizeResponsibleKey(value: string) {
	return normalizeKey(value);
}

export function getLearnerGuardianEntries(learner: Learner): LearnerGuardian[] {
	const rawGuardians =
		learner.guardians && learner.guardians.length > 0
			? learner.guardians
			: [
					{
						id: '',
						sourceKey: '',
						name: learner.guardian,
						relationship: learner.guardianRelationship,
						phone: ''
					}
				];
	const seenKeys = new Set<string>();

	return rawGuardians
		.map((guardian) => {
			const name = guardian.name.trim();
			const sourceKey = guardian.sourceKey || normalizeResponsibleKey(name);

			return {
				id: guardian.id || `learner-guardian-${sourceKey}`,
				sourceKey,
				name,
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

function createFamilyFromLearner(
	learner: Learner,
	learnerGuardians: LearnerGuardian[]
): CommunicationFamily {
	const now = new Date().toISOString();
	const primaryGuardian = learnerGuardians[0];

	return {
		id: createId('family'),
		familyName: inferFamilyName(primaryGuardian?.name ?? '', learner.name),
		sourceGuardianKey: primaryGuardian?.sourceKey,
		stage: 'novo',
		learnerIds: [learner.id],
		responsibles: learnerGuardians.slice(0, 2).map((guardian) => ({
			id: createId('responsible'),
			name: guardian.name,
			relationship: guardian.relationship,
			phone: normalizePhone(guardian.phone)
		})),
		contacts: [],
		nextStep: '',
		notes: '',
		lastContactAt: '',
		createdAt: now,
		updatedAt: now
	};
}

function inferFamilyName(guardianName: string, learnerName: string) {
	const learnerReference = learnerName.trim();
	if (learnerReference) return `Contatos de ${learnerReference}`;

	const guardianReference = guardianName.split(/\s+e\s+|,|\/|&/i)[0]?.trim() ?? guardianName;
	return guardianReference ? `Contatos de ${guardianReference}` : 'Contatos do aprendente';
}

function normalizeIds(ids: string[]) {
	return Array.from(new Set(ids.filter(Boolean)));
}

function normalizeKey(value: string) {
	return value.trim().toLowerCase().replace(/\s+/g, ' ');
}

function upsertFamilyResponsibles(
	family: CommunicationFamily,
	learnerGuardians: LearnerGuardian[]
) {
	let responsibles = family.responsibles;

	for (const guardian of learnerGuardians) {
		const guardianKey = normalizeResponsibleKey(guardian.name);
		if (!guardianKey) continue;

		const existingResponsible = responsibles.find(
			(responsible) => normalizeResponsibleKey(responsible.name) === guardianKey
		);
		if (existingResponsible) {
			responsibles = responsibles.map((responsible) =>
				responsible.id === existingResponsible.id
					? {
							...responsible,
							relationship: responsible.relationship || guardian.relationship,
							phone: responsible.phone || normalizePhone(guardian.phone)
						}
					: responsible
			);
			continue;
		}

		if (responsibles.length >= 2) continue;
		responsibles = [
			...responsibles,
			{
				id: createId('responsible'),
				name: guardian.name,
				relationship: guardian.relationship,
				phone: normalizePhone(guardian.phone)
			}
		];
	}

	return responsibles;
}

function mergeFamiliesByResponsible(families: CommunicationFamily[]) {
	const mergedFamilies: CommunicationFamily[] = [];
	const familyByResponsibleKey = new Map<string, CommunicationFamily>();

	for (const family of families) {
		const keys = getCommunicationFamilyResponsibleKeys(family);
		const existingFamily = keys.map((key) => familyByResponsibleKey.get(key)).find(Boolean);

		if (!existingFamily) {
			mergedFamilies.push(family);
			for (const key of keys) familyByResponsibleKey.set(key, family);
			continue;
		}

		existingFamily.learnerIds = normalizeIds([...existingFamily.learnerIds, ...family.learnerIds]);
		existingFamily.responsibles = mergeResponsibles(
			existingFamily.responsibles,
			family.responsibles
		);
		existingFamily.contacts = mergeContacts(existingFamily.contacts, family.contacts);
		existingFamily.sourceGuardianKey ??= family.sourceGuardianKey;
		existingFamily.updatedAt = new Date().toISOString();

		for (const key of getCommunicationFamilyResponsibleKeys(existingFamily)) {
			familyByResponsibleKey.set(key, existingFamily);
		}
	}

	return mergedFamilies;
}

function mergeResponsibles(
	left: CommunicationResponsible[],
	right: CommunicationResponsible[]
) {
	const responsibles = [...left];
	const keys = new Set(left.map((responsible) => normalizeResponsibleKey(responsible.name)));

	for (const responsible of right) {
		const key = normalizeResponsibleKey(responsible.name);
		if (key && keys.has(key)) continue;
		if (responsibles.length >= 2) break;
		responsibles.push(responsible);
		if (key) keys.add(key);
	}

	return responsibles;
}

function mergeContacts(left: CommunicationContact[], right: CommunicationContact[]) {
	const contactIds = new Set(left.map((contact) => contact.id));
	return [...left, ...right.filter((contact) => !contactIds.has(contact.id))];
}

function upsertGuardianOption(options: Map<string, GuardianOption>, option: GuardianOption) {
	const existing = options.get(option.key);
	if (existing) {
		existing.learnerIds = normalizeIds([...existing.learnerIds, ...option.learnerIds]);
		existing.relationship ||= option.relationship;
		existing.phone ||= normalizePhone(option.phone);
		return;
	}

	options.set(option.key, {
		...option,
		phone: normalizePhone(option.phone),
		learnerIds: normalizeIds(option.learnerIds)
	});
}

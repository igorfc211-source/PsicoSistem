<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { AppSidebar, AppTopbar, WorkspaceBanner } from '$lib/modules/clinic-shell/components';
	import '$lib/modules/clinic-shell/styles/clinic-app.css';
	import type { Banner, NavSection } from '$lib/modules/clinic-shell/types';
	import {
		AgendaWorkspace,
		CommunicationsWorkspace,
		LearnersWorkspace
	} from '$lib/modules/learners/components';
	import { clearStoredSession, getStoredSession, type StoredSession } from '$lib/auth';
	import {
		PLAN_CATEGORIES,
		addContactToFamily,
		addCustomActionPlanField,
		addResponsibleToFamily,
		appendAnamneseDocumentsToLearner,
		appendDocumentsToLearner,
		buildCalendarDays,
		createCommunicationFamily,
		createId,
		createLearner,
		createReportEntry,
		createSessionVisit,
		deleteDocumentBlob,
		filterLearners,
		getDocumentBlob,
		getDocumentStorageKey,
		getLearnerGuardianEntries,
		isValidEmailAddress,
		isValidInstagramHandle,
		isValidPhoneNumber,
		buildGuardianOptionsFromLearners,
		loadCommunicationFamilies,
		loadHiddenCommunicationSourceKeys,
		loadLearners,
		getCommunicationFamilyResponsibleKeys,
		normalizeResponsibleKey,
		patchLearnerInList,
		prepareDocumentBlob,
		prependReportToLearner,
		putDocumentBlob,
		removeContactFromFamily,
		removeCustomActionPlanField,
		removeReportFromLearner,
		removeResponsibleFromFamily,
		removeVisitFromList,
		saveCommunicationFamilies,
		saveHiddenCommunicationSourceKeys,
		saveLearners,
		sortVisitsBySchedule,
		syncCommunicationFamiliesWithLearners,
		toDateInputValue,
		updateActionPlanValue,
		updateCustomActionPlanField,
		updateVisitInList,
		type CommunicationFamily,
		type CoreActionPlanKey,
		type DetailTab,
		type Learner,
		type LearnerDocument,
		type LearnerFilter,
		type LearnerGuardianInput,
		type NewCommunicationContactInput,
		type NewCommunicationFamilyInput,
		type NewCommunicationResponsibleInput,
		type NewLearnerInput,
		type Visit
	} from '$lib/learners';
	import {
		addAgendaEvent,
		buildDayScheduleItems,
		isInvalidScheduleTime,
		loadAgendaEvents,
		removeAgendaEventById,
		saveAgendaEvents,
		type AgendaEvent,
		type NewAgendaEventInput,
		type NewSessionAppointmentInput
	} from '$lib/modules/scheduling';
	import { formatLongDate, formatMonth } from '$lib/shared/formatters';

	type DeletionConfirmation = {
		title: string;
		message: string;
		cancelLabel: string;
		confirmLabel: string;
		resolve: (confirmed: boolean) => void;
	};

	let session = $state<StoredSession | null>(null);
	let learners = $state<Learner[]>([]);
	let communicationFamilies = $state<CommunicationFamily[]>([]);
	let hiddenCommunicationSourceKeys = $state<string[]>([]);
	let agendaEvents = $state<AgendaEvent[]>([]);
	let selectedLearnerId = $state<string | null>(null);
	let selectedFamilyId = $state<string | null>(null);
	let activeSection = $state<NavSection>('aprendentes');
	let learnerFilter = $state<LearnerFilter>('active');
	let detailTab = $state<DetailTab>('resumo');
	let currentMonth = $state(new Date());
	let selectedAgendaDate = $state(toDateInputValue(new Date()));
	let selectedVisitId = $state<string | null>(null);
	let searchTerm = $state('');
	let isUploading = $state(false);
	let showAddForm = $state(false);
	let isSidebarOpen = $state(false);
	let theme = $state<'light' | 'dark'>('light');
	let banner = $state<Banner | null>(null);
	let deletionConfirmation = $state<DeletionConfirmation | null>(null);

	const filteredLearners = $derived(filterLearners(learners, searchTerm, learnerFilter));
	const selectedLearner = $derived(
		learners.find((learner) => learner.id === selectedLearnerId) ?? null
	);
	const selectedVisit = $derived(
		selectedLearner?.visits.find((visit) => visit.id === selectedVisitId) ?? null
	);
	const selectedLearnerVisits = $derived(selectedLearner?.visits ?? []);
	const allVisits = $derived(
		learners
			.flatMap((learner) => learner.visits.map((visit) => ({ learner, visit })))
			.sort((left, right) =>
				`${left.visit.date}-${left.visit.startTime}`.localeCompare(
					`${right.visit.date}-${right.visit.startTime}`
				)
			)
	);
	const selectedLearnerCalendarDays = $derived(
		buildCalendarDays(currentMonth, selectedLearnerVisits, [], selectedAgendaDate)
	);
	const agendaCalendarDays = $derived(
		buildCalendarDays(
			currentMonth,
			allVisits.map(({ visit }) => visit),
			agendaEvents.map((event) => event.date),
			selectedAgendaDate
		)
	);
	const selectedDayItems = $derived(
		buildDayScheduleItems(selectedAgendaDate, allVisits, agendaEvents)
	);
	const pendingVisits = $derived(
		allVisits.filter(({ visit }) => visit.status === 'scheduled').slice(0, 12)
	);
	const monthLabel = $derived(formatMonth(currentMonth));
	const currentDateLabel = $derived(formatLongDate(selectedAgendaDate));
	const tenantName = $derived(session?.payload.tenant?.name ?? 'PsicoClinica');
	const userName = $derived(session?.payload.user?.name ?? 'Usuario');
	const guardianOptions = $derived(buildGuardianOptionsFromLearners(learners, communicationFamilies));

	onMount(() => {
		if (!browser) return;

		const storedSession = getStoredSession();
		if (!storedSession) {
			void goto('/');
			return;
		}

		session = storedSession;
		learners = loadLearners();
		hiddenCommunicationSourceKeys = loadHiddenCommunicationSourceKeys();
		communicationFamilies = syncCommunicationFamiliesWithLearners(
			loadCommunicationFamilies(),
			learners,
			hiddenCommunicationSourceKeys
		);
		agendaEvents = loadAgendaEvents();
		selectedLearnerId = learners[0]?.id ?? null;
		selectedFamilyId = null;
		saveCommunicationFamilies(communicationFamilies);
		theme = localStorage.getItem('psicosistem.theme') === 'dark' ? 'dark' : 'light';
	});

	// Abre/fecha o menu lateral mobile sem afetar a navegacao desktop.
	function toggleSidebar() {
		isSidebarOpen = !isSidebarOpen;
	}

	// Fecha o menu lateral depois de selecionar uma secao no celular.
	function closeSidebar() {
		isSidebarOpen = false;
	}

	// Alterna tema claro/escuro e salva a preferencia local do usuario.
	function toggleTheme() {
		theme = theme === 'dark' ? 'light' : 'dark';
		if (browser) {
			localStorage.setItem('psicosistem.theme', theme);
		}
	}

	// Persiste o snapshot de aprendentes sempre que uma operacao clinica muda o estado.
	function persistLearners(nextLearners = learners) {
		if (!browser) return;
		saveLearners(nextLearners);
	}

	// Persiste eventos livres, que vivem fora do prontuario de um aprendente especifico.
	function persistAgendaEvents(nextEvents = agendaEvents) {
		if (!browser) return;
		saveAgendaEvents(nextEvents);
	}

	// Persiste os cards de comunicacao, que funcionam como um mini-CRM de familias.
	function persistCommunicationFamilies(nextFamilies = communicationFamilies) {
		if (!browser) return;
		saveCommunicationFamilies(nextFamilies);
	}

	function setCommunicationFamilies(
		nextFamilies: CommunicationFamily[],
		nextLearners = learners,
		hiddenSourceKeys = hiddenCommunicationSourceKeys
	) {
		const syncedFamilies = syncCommunicationFamiliesWithLearners(
			nextFamilies,
			nextLearners,
			hiddenSourceKeys
		);
		communicationFamilies = syncedFamilies;
		persistCommunicationFamilies(syncedFamilies);

		if (!selectedFamilyId || !syncedFamilies.some((family) => family.id === selectedFamilyId)) {
			selectedFamilyId = null;
		}
	}

	function syncFamiliesForLearners(nextLearners: Learner[]) {
		if (!browser) return;
		setCommunicationFamilies(communicationFamilies, nextLearners);
	}

	function confirmDeletion(message: string, title = 'Confirmar remocao') {
		if (!browser) return Promise.resolve(false);

		if (deletionConfirmation) {
			deletionConfirmation.resolve(false);
		}

		return new Promise<boolean>((resolve) => {
			deletionConfirmation = {
				title,
				message,
				cancelLabel: 'Cancelar',
				confirmLabel: 'Confirmar exclusao',
				resolve
			};
		});
	}

	function closeDeletionConfirmation(confirmed: boolean) {
		const currentConfirmation = deletionConfirmation;
		if (!currentConfirmation) return;

		deletionConfirmation = null;
		currentConfirmation.resolve(confirmed);
	}

	function handleDeletionConfirmationKeydown(event: KeyboardEvent) {
		if (!deletionConfirmation) return;

		if (event.key === 'Escape') {
			event.preventDefault();
			closeDeletionConfirmation(false);
		}

		if (event.key === 'Enter') {
			event.preventDefault();
			closeDeletionConfirmation(true);
		}
	}

	function getLearnerDocumentStorageKeys(learner: Learner) {
		return [...learner.documents, ...learner.anamneseDocuments].map((document) =>
			getDocumentStorageKey(learner.id, document.id)
		);
	}

	// Atualiza qualquer aprendente por id, mantendo a data de edicao em um unico lugar.
	function updateLearnerById(learnerId: string, patch: Partial<Learner>) {
		const nextLearners = patchLearnerInList(learners, learnerId, patch);
		learners = nextLearners;
		persistLearners(nextLearners);
		syncFamiliesForLearners(nextLearners);
	}

	// Aplica mudancas no aprendente selecionado, usado por anamnese, documentos, plano e agenda.
	function updateSelectedLearner(patch: Partial<Learner>) {
		if (!selectedLearner) return;
		updateLearnerById(selectedLearner.id, patch);
	}

	// Mantem as secoes do plano de acao independentes para facilitar novas categorias.
	function updateActionPlan(key: CoreActionPlanKey, value: string) {
		if (!selectedLearner) return;

		updateSelectedLearner({
			actionPlan: updateActionPlanValue(selectedLearner.actionPlan, key, value)
		});
	}

	// Adiciona campos personalizados no plano para clinicas que precisam de protocolos proprios.
	function addCustomActionPlanFieldToSelected(label: string, description: string) {
		if (!selectedLearner) return false;

		if (!label.trim()) {
			banner = {
				tone: 'error',
				text: 'Informe o nome do novo campo do plano.'
			};
			return false;
		}

		updateSelectedLearner({
			actionPlan: addCustomActionPlanField(selectedLearner.actionPlan, label, description)
		});
		banner = {
			tone: 'success',
			text: 'Campo adicionado ao plano de acao.'
		};
		return true;
	}

	// Atualiza apenas o campo personalizado editado, sem tocar nas secoes padrao.
	function updateCustomActionPlanFieldForSelected(fieldId: string, value: string) {
		if (!selectedLearner) return;

		updateSelectedLearner({
			actionPlan: updateCustomActionPlanField(selectedLearner.actionPlan, fieldId, value)
		});
	}

	// Remove campos personalizados do plano quando a clinica nao precisa mais daquela secao.
	async function removeCustomActionPlanFieldFromSelected(fieldId: string) {
		if (!selectedLearner) return;

		const field = selectedLearner.actionPlan.customFields.find((item) => item.id === fieldId);
		const fieldName = field?.label ? `"${field.label}"` : 'este campo';
		if (!(await confirmDeletion(`Remover ${fieldName} do plano de acao?`))) return;

		updateSelectedLearner({
			actionPlan: removeCustomActionPlanField(selectedLearner.actionPlan, fieldId)
		});
		banner = {
			tone: 'success',
			text: 'Campo removido do plano de acao.'
		};
	}

	// Edita uma visita existente e reordena a agenda pelo campo de data.
	function updateVisit(visitId: string, patch: Partial<Visit>) {
		if (!selectedLearner) return;

		updateSelectedLearner({
			visits: updateVisitInList(selectedLearner.visits, visitId, patch)
		});
	}

	// Cria um aprendente completo, incluindo agenda inicial calculada pelo dominio.
	function handleCreateLearner(input: NewLearnerInput) {
		if (!input.name.trim() || !input.startDate || !input.endDate) {
			banner = {
				tone: 'error',
				text: 'Informe nome, data inicial e data final do aprendente.'
			};
			return false;
		}

		const filledGuardians = getFilledGuardianInputs(input.guardians);
		if (filledGuardians.length === 0) {
			banner = {
				tone: 'error',
				text: 'Informe pelo menos um responsavel do aprendente.'
			};
			return false;
		}

		if (filledGuardians.some((guardian) => !guardian.relationship.trim())) {
			banner = {
				tone: 'error',
				text: 'Selecione o parentesco de cada responsavel informado.'
			};
			return false;
		}

		if (hasDuplicateGuardian(filledGuardians)) {
			banner = {
				tone: 'error',
				text: 'O mesmo responsavel nao pode ser informado duas vezes no aprendente.'
			};
			return false;
		}

		const createdLearner = createLearner(input);
		const nextLearners = [createdLearner, ...learners];
		learners = nextLearners;
		selectedLearnerId = createdLearner.id;
		activeSection = 'aprendentes';
		detailTab = 'resumo';
		showAddForm = false;
		persistLearners(nextLearners);
		syncFamiliesForLearners(nextLearners);

		banner = {
			tone: 'success',
			text: 'Aprendente adicionado ao painel.'
		};
		return true;
	}

	async function deleteLearner(learnerId: string) {
		const learner = learners.find((item) => item.id === learnerId);
		if (!learner) return;

		if (
			!(await confirmDeletion(
				`Excluir o aprendente "${learner.name}"? Isso remove o prontuario, agenda, documentos e relatorios desse aprendente.`
			))
		) {
			return;
		}

		try {
			await Promise.all(getLearnerDocumentStorageKeys(learner).map((key) => deleteDocumentBlob(key)));
		} catch (error) {
			banner = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel excluir os arquivos do aprendente.'
			};
			return;
		}

		const nextLearners = learners.filter((item) => item.id !== learnerId);
		learners = nextLearners;
		persistLearners(nextLearners);
		syncFamiliesForLearners(nextLearners);

		if (selectedLearnerId === learnerId) {
			selectedLearnerId = nextLearners[0]?.id ?? null;
			selectedVisitId = null;
			detailTab = 'resumo';
		}

		banner = {
			tone: 'success',
			text: 'Aprendente excluido com sucesso.'
		};
	}

	function getFilledGuardianInputs(guardians: LearnerGuardianInput[]) {
		return guardians.filter(
			(guardian) =>
				guardian.name.trim() ||
				guardian.relationship.trim() ||
				guardian.phone.trim() ||
				guardian.sourceKey.trim()
		);
	}

	function hasDuplicateGuardian(guardians: LearnerGuardianInput[]) {
		const keys = new Set<string>();
		for (const guardian of guardians) {
			const key = guardian.sourceKey || normalizeResponsibleKey(guardian.name);
			if (!key) continue;
			if (keys.has(key)) return true;
			keys.add(key);
		}

		return false;
	}

	// Abre o prontuario do aprendente e volta para a secao principal quando necessario.
	function selectLearner(id: string) {
		selectedLearnerId = id;
		activeSection = 'aprendentes';
		detailTab = 'resumo';
		selectedVisitId = null;
		banner = null;
	}

	function selectFamily(id: string) {
		selectedFamilyId = id;
		banner = null;
	}

	function closeFamily() {
		selectedFamilyId = null;
		banner = null;
	}

	function openLearnerResponsible(learner: Learner) {
		let family = communicationFamilies.find((item) => item.learnerIds.includes(learner.id)) ?? null;
		const learnerGuardians = getLearnerGuardianEntries(learner);
		const primaryGuardian = learnerGuardians[0] ?? null;

		if (!family && primaryGuardian) {
			family = createCommunicationFamily({
				familyName: buildFamilyNameFromLearner(learner),
				responsibleName: primaryGuardian.name,
				responsiblePhone: primaryGuardian.phone,
				relationship: primaryGuardian.relationship,
				learnerIds: [learner.id]
			});
			const sourceKey = family.sourceGuardianKey;
			let nextHiddenSourceKeys = hiddenCommunicationSourceKeys;
			if (sourceKey && hiddenCommunicationSourceKeys.includes(sourceKey)) {
				nextHiddenSourceKeys = hiddenCommunicationSourceKeys.filter((key) => key !== sourceKey);
				hiddenCommunicationSourceKeys = nextHiddenSourceKeys;
				saveHiddenCommunicationSourceKeys(nextHiddenSourceKeys);
			}
			setCommunicationFamilies([family, ...communicationFamilies], learners, nextHiddenSourceKeys);
		}

		if (family) {
			selectedFamilyId = family.id;
		}

		activeSection = 'comunicacoes';
		showAddForm = false;
		banner = family
			? null
			: {
					tone: 'info',
					text: 'Informe um responsavel no cadastro do aprendente para abrir a comunicacao.'
				};
	}

	function buildFamilyNameFromLearner(learner: Learner) {
		return `Contatos de ${learner.name}`;
	}

	function createCommunicationFamilyCard(input: NewCommunicationFamilyInput) {
		const normalizedInput = {
			...input,
			familyName: input.familyName.trim() || buildCommunicationFamilyName(input)
		};

		if (!normalizedInput.responsibleName.trim()) {
			banner = {
				tone: 'error',
				text: 'Informe o responsavel principal.'
			};
			return false;
		}

		if (!isValidPhoneNumber(normalizedInput.responsiblePhone)) {
			banner = {
				tone: 'error',
				text: 'Informe um numero com 10 ou 11 digitos.'
			};
			return false;
		}

		if (isResponsibleUsedInAnotherFamily(normalizedInput.responsibleName)) {
			banner = {
				tone: 'error',
				text: 'Esse responsavel ja possui um card de comunicacao.'
			};
			return false;
		}

		const createdFamily = createCommunicationFamily(normalizedInput);
		const createdSourceKey = createdFamily.sourceGuardianKey;
		let nextHiddenSourceKeys = hiddenCommunicationSourceKeys;
		if (createdSourceKey && hiddenCommunicationSourceKeys.includes(createdSourceKey)) {
			nextHiddenSourceKeys = hiddenCommunicationSourceKeys.filter((key) => key !== createdSourceKey);
			hiddenCommunicationSourceKeys = nextHiddenSourceKeys;
			saveHiddenCommunicationSourceKeys(nextHiddenSourceKeys);
		}

		selectedFamilyId = createdFamily.id;
		setCommunicationFamilies([createdFamily, ...communicationFamilies], learners, nextHiddenSourceKeys);
		banner = {
			tone: 'success',
			text: 'Card de comunicacao criado.'
		};
		return true;
	}

	function buildCommunicationFamilyName(input: NewCommunicationFamilyInput) {
		const linkedLearners = input.learnerIds
			.map((learnerId) => learners.find((learner) => learner.id === learnerId))
			.filter((learner): learner is Learner => Boolean(learner));

		if (linkedLearners.length) {
			const [firstLearner] = linkedLearners;
			return linkedLearners.length === 1
				? `Contatos de ${firstLearner.name}`
				: `Contatos de ${firstLearner.name} +${linkedLearners.length - 1}`;
		}

		const responsibleName = input.responsibleName.trim();
		return responsibleName ? `Contatos de ${responsibleName}` : 'Contatos do aprendente';
	}

	function updateCommunicationFamily(familyId: string, patch: Partial<CommunicationFamily>) {
		const nextFamilies = communicationFamilies.map((family) =>
			family.id === familyId
				? {
						...family,
						...patch,
						updatedAt: new Date().toISOString()
					}
				: family
		);
		setCommunicationFamilies(nextFamilies);
	}

	async function deleteCommunicationFamily(familyId: string) {
		const removedFamily = communicationFamilies.find((family) => family.id === familyId);
		if (!removedFamily) return;

		if (!(await confirmDeletion(`Excluir o card de comunicacao "${removedFamily.familyName}"?`))) {
			return;
		}

		const nextFamilies = communicationFamilies.filter((family) => family.id !== familyId);
		let nextHiddenSourceKeys = hiddenCommunicationSourceKeys;
		const removedSourceKeys = getCommunicationFamilyResponsibleKeys(removedFamily);

		if (removedSourceKeys.length) {
			nextHiddenSourceKeys = Array.from(
				new Set([...hiddenCommunicationSourceKeys, ...removedSourceKeys])
			);
			hiddenCommunicationSourceKeys = nextHiddenSourceKeys;
			saveHiddenCommunicationSourceKeys(nextHiddenSourceKeys);
		}

		setCommunicationFamilies(nextFamilies, learners, nextHiddenSourceKeys);
		banner = {
			tone: 'success',
			text: 'Card de comunicacao removido.'
		};
	}

	function addResponsibleToCommunicationFamily(
		familyId: string,
		input: NewCommunicationResponsibleInput
	) {
		if (!input.name.trim() || !input.phone.trim()) {
			banner = {
				tone: 'error',
				text: 'Informe nome e numero do responsavel.'
			};
			return false;
		}

		if (!isValidPhoneNumber(input.phone)) {
			banner = {
				tone: 'error',
				text: 'Informe um numero com 10 ou 11 digitos.'
			};
			return false;
		}

		const family = communicationFamilies.find((item) => item.id === familyId);
		if (!family) return false;

		if (family.responsibles.length >= 2) {
			banner = {
				tone: 'error',
				text: 'Cada card pode ter no maximo dois responsaveis.'
			};
			return false;
		}

		if (isResponsibleInFamily(family, input.name)) {
			banner = {
				tone: 'error',
				text: 'Esse responsavel ja esta neste card.'
			};
			return false;
		}

		if (isResponsibleUsedInAnotherFamily(input.name, familyId)) {
			banner = {
				tone: 'error',
				text: 'Esse responsavel ja possui outro card de comunicacao.'
			};
			return false;
		}

		const nextFamilies = communicationFamilies.map((family) =>
			family.id === familyId ? addResponsibleToFamily(family, input) : family
		);
		setCommunicationFamilies(nextFamilies);
		banner = {
			tone: 'success',
			text: 'Responsavel adicionado ao card.'
		};
		return true;
	}

	function isResponsibleUsedInAnotherFamily(responsibleName: string, ignoredFamilyId?: string) {
		const key = normalizeResponsibleKey(responsibleName);
		if (!key) return false;

		return communicationFamilies.some(
			(family) =>
				family.id !== ignoredFamilyId && getCommunicationFamilyResponsibleKeys(family).includes(key)
		);
	}

	function isResponsibleInFamily(family: CommunicationFamily, responsibleName: string) {
		const key = normalizeResponsibleKey(responsibleName);
		if (!key) return false;

		return getCommunicationFamilyResponsibleKeys(family).includes(key);
	}

	async function removeResponsibleFromCommunicationFamily(familyId: string, responsibleId: string) {
		const family = communicationFamilies.find((item) => item.id === familyId);
		const responsible = family?.responsibles.find((item) => item.id === responsibleId);
		if (!family || !responsible) return;

		if (!(await confirmDeletion(`Remover "${responsible.name}" do card "${family.familyName}"?`))) {
			return;
		}

		const nextFamilies = communicationFamilies.map((family) =>
			family.id === familyId ? removeResponsibleFromFamily(family, responsibleId) : family
		);
		setCommunicationFamilies(nextFamilies);
		banner = {
			tone: 'success',
			text: 'Responsavel removido do card.'
		};
	}

	function addContactToCommunicationFamily(familyId: string, input: NewCommunicationContactInput) {
		if (!input.value.trim()) {
			banner = {
				tone: 'error',
				text: 'Informe o contato que sera adicionado.'
			};
			return false;
		}

		if ((input.type === 'phone' || input.type === 'whatsapp') && !isValidPhoneNumber(input.value)) {
			banner = {
				tone: 'error',
				text: 'Informe um numero com 10 ou 11 digitos.'
			};
			return false;
		}

		if (input.type === 'instagram' && !isValidInstagramHandle(input.value)) {
			banner = {
				tone: 'error',
				text: 'Informe um @ de Instagram valido.'
			};
			return false;
		}

		if (input.type === 'email' && !isValidEmailAddress(input.value)) {
			banner = {
				tone: 'error',
				text: 'Informe um e-mail valido.'
			};
			return false;
		}

		if (input.type !== 'instagram' && !input.label.trim()) {
			banner = {
				tone: 'error',
				text: 'Informe o nome do contato.'
			};
			return false;
		}

		const nextFamilies = communicationFamilies.map((family) =>
			family.id === familyId ? addContactToFamily(family, input) : family
		);
		setCommunicationFamilies(nextFamilies);
		banner = {
			tone: 'success',
			text: 'Contato adicionado ao card.'
		};
		return true;
	}

	async function removeContactFromCommunicationFamily(familyId: string, contactId: string) {
		const family = communicationFamilies.find((item) => item.id === familyId);
		const contact = family?.contacts.find((item) => item.id === contactId);
		if (!family || !contact) return;

		if (!(await confirmDeletion(`Remover o contato "${contact.label || contact.value}" deste card?`))) {
			return;
		}

		const nextFamilies = communicationFamilies.map((family) =>
			family.id === familyId ? removeContactFromFamily(family, contactId) : family
		);
		setCommunicationFamilies(nextFamilies);
		banner = {
			tone: 'success',
			text: 'Contato removido do card.'
		};
	}

	// Troca o aprendente dentro da agenda sem tirar o usuario da visualizacao de calendario.
	function selectLearnerInsideAgenda(id: string) {
		selectedLearnerId = id;
		selectedVisitId = null;
		banner = null;
	}

	// Mapeia secoes do menu lateral para a aba interna equivalente do prontuario.
	function selectSection(section: NavSection) {
		activeSection = section;
		showAddForm = false;
		isSidebarOpen = false;

		if (section === 'agenda') {
			detailTab = 'agenda';
		} else if (section === 'financeiro') {
			detailTab = 'anamnese';
		} else if (section === 'comunicacoes') {
			detailTab = 'relatorios';
		}
	}

	// Entrada do dropdown de perfil para conduzir o usuario ate configuracoes futuras.
	function editProfile() {
		activeSection = 'configuracoes';
		showAddForm = false;
		banner = {
			tone: 'info',
			text: 'Edicao de perfil pronta para conectar aos dados do usuario.'
		};
	}

	// Abre a area de configuracoes a partir do menu de perfil.
	function openSettings() {
		activeSection = 'configuracoes';
		showAddForm = false;
		banner = {
			tone: 'info',
			text: 'Configuracoes da clinica e preferencias abertas.'
		};
	}

	// Navega os calendarios mantendo a mesma data base para todas as secoes.
	function shiftMonth(delta: number) {
		const next = new Date(currentMonth);
		next.setMonth(next.getMonth() + delta);
		currentMonth = next;
	}

	// Seleciona o dia na agenda e abre a linha do tempo profissional daquele dia.
	function handleCalendarDate(date: string) {
		selectedAgendaDate = date;
		currentMonth = new Date(`${date}T12:00:00`);

		const existingVisit = selectedLearner?.visits.find((visit) => visit.date === date) ?? null;
		selectedVisitId = existingVisit?.id ?? null;
		detailTab = 'agenda';
		showAddForm = false;
		banner = null;
	}

	// Remove uma visita e sincroniza a contagem exibida no cadastro do aprendente.
	async function removeVisit(id: string) {
		if (!selectedLearner) return;
		const visit = selectedLearner.visits.find((item) => item.id === id);
		if (!visit) return;

		if (!(await confirmDeletion(`Remover a sessao de ${visit.date} as ${visit.startTime}?`))) return;

		const nextVisits = removeVisitFromList(selectedLearner.visits, id);

		updateSelectedLearner({
			visits: nextVisits,
			visitCount: nextVisits.length
		});
		selectedVisitId = null;
		banner = {
			tone: 'success',
			text: 'Sessao removida da agenda.'
		};
	}

	// Remove uma sessao diretamente da agenda diaria, mesmo quando outro aprendente esta aberto.
	async function removeSessionAppointment(learnerId: string, visitId: string) {
		const learner = learners.find((item) => item.id === learnerId);
		if (!learner) return;
		const visit = learner.visits.find((item) => item.id === visitId);
		if (!visit) return;

		if (
			!(await confirmDeletion(
				`Remover a sessao de ${learner.name} em ${visit.date} as ${visit.startTime}?`
			))
		) {
			return;
		}

		const nextVisits = removeVisitFromList(learner.visits, visitId);
		updateLearnerById(learner.id, {
			visits: nextVisits,
			visitCount: nextVisits.length
		});

		if (selectedVisitId === visitId) {
			selectedVisitId = null;
		}

		banner = {
			tone: 'success',
			text: 'Sessao removida da agenda.'
		};
	}

	// Cria uma sessao com aprendente usando horarios precisos e dados de local/modalidade.
	function createSessionAppointment(input: NewSessionAppointmentInput) {
		const learner = learners.find((item) => item.id === input.learnerId);
		if (!learner) {
			banner = {
				tone: 'error',
				text: 'Escolha um aprendente para confirmar a sessao.'
			};
			return false;
		}

		if (isInvalidScheduleTime(input.startTime, input.endTime)) {
			banner = {
				tone: 'error',
				text: 'Confira os horarios: o fim precisa ser depois do inicio.'
			};
			return false;
		}

		const visit = createSessionVisit(input);
		const nextVisits = [...learner.visits, visit].sort(sortVisitsBySchedule);

		updateLearnerById(learner.id, {
			visits: nextVisits,
			visitCount: nextVisits.length
		});
		selectedLearnerId = learner.id;
		selectedVisitId = visit.id;
		selectedAgendaDate = input.date;
		banner = {
			tone: 'success',
			text: `Sessao agendada com ${learner.name}.`
		};
		return true;
	}

	// Cria um evento livre, como reuniao, supervisao ou bloqueio de agenda.
	function createEventAppointment(input: NewAgendaEventInput) {
		if (!input.title.trim()) {
			banner = {
				tone: 'error',
				text: 'Informe um titulo para o evento.'
			};
			return false;
		}

		if (isInvalidScheduleTime(input.startTime, input.endTime)) {
			banner = {
				tone: 'error',
				text: 'Confira os horarios: o fim precisa ser depois do inicio.'
			};
			return false;
		}

		const { events: nextEvents } = addAgendaEvent(agendaEvents, input);

		agendaEvents = nextEvents;
		selectedAgendaDate = input.date;
		persistAgendaEvents(nextEvents);
		banner = {
			tone: 'success',
			text: 'Evento adicionado a agenda.'
		};
		return true;
	}

	// Remove eventos livres sem mexer nas sessoes registradas nos aprendentes.
	async function removeAgendaEvent(event: AgendaEvent) {
		if (!(await confirmDeletion(`Remover o evento "${event.title}" da agenda?`))) return;

		const nextEvents = removeAgendaEventById(agendaEvents, event.id);
		agendaEvents = nextEvents;
		persistAgendaEvents(nextEvents);
		banner = {
			tone: 'success',
			text: 'Evento removido da agenda.'
		};
	}

	// Processa uploads no IndexedDB e grava apenas metadados leves no prontuario.
	async function handleDocumentUpload(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		const files = Array.from(input.files ?? []);
		const targetLearner = selectedLearner;
		if (!targetLearner || files.length === 0) return;

		isUploading = true;
		banner = {
			tone: 'info',
			text: 'Processando documentos.'
		};

		try {
			const uploadedDocuments: LearnerDocument[] = [];

			for (const file of files) {
				const prepared = await prepareDocumentBlob(file);
				const documentId = createId('doc');
				await putDocumentBlob(getDocumentStorageKey(targetLearner.id, documentId), prepared.blob);

				uploadedDocuments.push({
					id: documentId,
					name: file.name,
					type: prepared.compressed ? 'application/gzip' : file.type || 'application/octet-stream',
					size: file.size,
					storedSize: prepared.blob.size,
					compressed: prepared.compressed,
					createdAt: new Date().toISOString()
				});
			}

			const currentLearner = learners.find((learner) => learner.id === targetLearner.id);
			if (currentLearner) {
				const nextLearner = appendDocumentsToLearner(currentLearner, uploadedDocuments);
				updateLearnerById(targetLearner.id, {
					documents: nextLearner.documents
				});
			}

			banner = {
				tone: 'success',
				text: 'Documento armazenado no aprendente.'
			};
			input.value = '';
		} catch (error) {
			banner = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel armazenar o documento.'
			};
		} finally {
			isUploading = false;
		}
	}

	// Processa anexos especificos da anamnese mantendo-os separados dos documentos gerais.
	async function handleAnamneseDocumentUpload(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		const files = Array.from(input.files ?? []);
		const targetLearner = selectedLearner;
		if (!targetLearner || files.length === 0) return;

		isUploading = true;
		banner = {
			tone: 'info',
			text: 'Processando anexos da anamnese.'
		};

		try {
			const uploadedDocuments: LearnerDocument[] = [];

			for (const file of files) {
				const prepared = await prepareDocumentBlob(file);
				const documentId = createId('anam-doc');
				await putDocumentBlob(getDocumentStorageKey(targetLearner.id, documentId), prepared.blob);

				uploadedDocuments.push({
					id: documentId,
					name: file.name,
					type: prepared.compressed ? 'application/gzip' : file.type || 'application/octet-stream',
					size: file.size,
					storedSize: prepared.blob.size,
					compressed: prepared.compressed,
					createdAt: new Date().toISOString()
				});
			}

			const currentLearner = learners.find((learner) => learner.id === targetLearner.id);
			if (currentLearner) {
				const nextLearner = appendAnamneseDocumentsToLearner(currentLearner, uploadedDocuments);
				updateLearnerById(targetLearner.id, {
					anamneseDocuments: nextLearner.anamneseDocuments
				});
			}

			banner = {
				tone: 'success',
				text: 'Anexo adicionado a anamnese.'
			};
			input.value = '';
		} catch (error) {
			banner = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel anexar o arquivo.'
			};
		} finally {
			isUploading = false;
		}
	}

	// Recupera o blob do IndexedDB e dispara o download no navegador.
	async function downloadDocument(document: LearnerDocument) {
		if (!selectedLearner) return;

		const blob = await getDocumentBlob(getDocumentStorageKey(selectedLearner.id, document.id));
		if (!blob) {
			banner = {
				tone: 'error',
				text: 'Arquivo nao encontrado no armazenamento local.'
			};
			return;
		}

		const url = URL.createObjectURL(blob);
		const anchor = window.document.createElement('a');
		anchor.href = url;
		anchor.download = document.compressed ? `${document.name}.gz` : document.name;
		anchor.click();
		URL.revokeObjectURL(url);
	}

	// Baixa anexos da anamnese usando o mesmo cofre local de arquivos do prontuario.
	async function downloadAnamneseDocument(document: LearnerDocument) {
		if (!selectedLearner) return;

		const blob = await getDocumentBlob(getDocumentStorageKey(selectedLearner.id, document.id));
		if (!blob) {
			banner = {
				tone: 'error',
				text: 'Arquivo da anamnese nao encontrado.'
			};
			return;
		}

		const url = URL.createObjectURL(blob);
		const anchor = window.document.createElement('a');
		anchor.href = url;
		anchor.download = document.compressed ? `${document.name}.gz` : document.name;
		anchor.click();
		URL.revokeObjectURL(url);
	}

	// Exclui o arquivo fisico local e remove seus metadados do aprendente.
	async function removeDocument(document: LearnerDocument) {
		if (!selectedLearner) return;

		if (!(await confirmDeletion(`Excluir o documento "${document.name}"?`))) return;

		try {
			await deleteDocumentBlob(getDocumentStorageKey(selectedLearner.id, document.id));
			updateSelectedLearner({
				documents: selectedLearner.documents.filter((item) => item.id !== document.id)
			});
			banner = {
				tone: 'success',
				text: 'Documento excluido.'
			};
		} catch (error) {
			banner = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel excluir o documento.'
			};
		}
	}

	// Remove anexos da anamnese sem afetar os documentos gerais do aprendente.
	async function removeAnamneseDocument(document: LearnerDocument) {
		if (!selectedLearner) return;

		if (!(await confirmDeletion(`Excluir o anexo "${document.name}" da anamnese?`))) return;

		try {
			await deleteDocumentBlob(getDocumentStorageKey(selectedLearner.id, document.id));
			updateSelectedLearner({
				anamneseDocuments: selectedLearner.anamneseDocuments.filter((item) => item.id !== document.id)
			});
			banner = {
				tone: 'success',
				text: 'Anexo da anamnese excluido.'
			};
		} catch (error) {
			banner = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel excluir o anexo.'
			};
		}
	}

	// Registra relatorios com carimbo de criacao e atualizacao.
	function addReport(text: string) {
		if (!selectedLearner || !text.trim()) return;

		const nextLearner = prependReportToLearner(selectedLearner, createReportEntry(text));
		updateSelectedLearner({
			reports: nextLearner.reports
		});
	}

	// Remove um relatorio mantendo os demais registros historicos intactos.
	async function removeReport(id: string) {
		if (!selectedLearner) return;
		const report = selectedLearner.reports.find((item) => item.id === id);
		if (!report) return;

		if (!(await confirmDeletion('Excluir este relatorio do aprendente?'))) return;

		const nextLearner = removeReportFromLearner(selectedLearner, id);
		updateSelectedLearner({
			reports: nextLearner.reports
		});
		banner = {
			tone: 'success',
			text: 'Relatorio excluido.'
		};
	}

	// Encerra a sessao local e retorna para login.
	function logout() {
		clearStoredSession();
		void goto('/');
	}
</script>

<svelte:head>
	<title>PsicoSistem | Painel</title>
</svelte:head>

<svelte:window onkeydown={handleDeletionConfirmationKeydown} />

{#if !session}
	<!-- Estado de carregamento antes de validar a sessao local. -->
	<main class="loading-screen">
		<p>Carregando sessao...</p>
	</main>
{:else}
	<!-- Shell autenticado: concentra a estrutura visual usada em todo o painel. -->
	<div class="clinic-shell" class:dark-mode={theme === 'dark'}>
		<!-- Barra decorativa superior no estilo janela desktop. -->
		<div class="window-bar">
			<p>PsicoSistem</p>
		</div>

		<!-- Frame principal: menu lateral fixo + area de conteudo dinamica. -->
		<div class="app-frame">
			<AppSidebar
				{tenantName}
				{activeSection}
				isOpen={isSidebarOpen}
				onSelectSection={selectSection}
				onClose={closeSidebar}
			/>


				<!-- Componente teste -->
			<button
				type="button"
				class="sidebar-backdrop"
				class:visible={isSidebarOpen}
				aria-label="Fechar menu"
				onclick={closeSidebar}
			></button>

			<!-- Conteudo da rota: topbar, feedback global e workspace ativo. -->
			<section class="content">
				<AppTopbar
					{searchTerm}
					{userName}
					{theme}
					searchPlaceholder={
						activeSection === 'comunicacoes'
							? 'Buscar familia, responsavel ou aprendente...'
							: 'Buscar aprendente ou agendamento...'
					}
					onToggleSidebar={toggleSidebar}
					onSearchTermChange={(value) => (searchTerm = value)}
					onEditProfile={editProfile}
					onOpenSettings={openSettings}
					onToggleTheme={toggleTheme}
					onLogout={logout}
				/>

				<WorkspaceBanner {banner} />

				<!-- Workspace de agenda global: mostra compromissos de todos os aprendentes. -->
				{#if activeSection === 'agenda'}
					<AgendaWorkspace
						calendarDays={agendaCalendarDays}
						{monthLabel}
						selectedDate={selectedAgendaDate}
						{currentDateLabel}
						{learners}
						{selectedLearnerId}
						{userName}
						dayItems={selectedDayItems}
						{pendingVisits}
						onShiftMonth={shiftMonth}
						onSelectCalendarDate={handleCalendarDate}
						onSelectLearnerId={selectLearnerInsideAgenda}
						onOpenLearner={selectLearner}
						onCreateSession={createSessionAppointment}
						onCreateEvent={createEventAppointment}
						onRemoveSession={removeSessionAppointment}
						onRemoveEvent={removeAgendaEvent}
					/>
				{:else if activeSection === 'comunicacoes'}
					<CommunicationsWorkspace
						{learners}
						families={communicationFamilies}
						{selectedFamilyId}
						{searchTerm}
						onCreateFamily={createCommunicationFamilyCard}
						onUpdateFamily={updateCommunicationFamily}
						onDeleteFamily={deleteCommunicationFamily}
						onAddResponsible={addResponsibleToCommunicationFamily}
						onRemoveResponsible={removeResponsibleFromCommunicationFamily}
						onAddContact={addContactToCommunicationFamily}
						onRemoveContact={removeContactFromCommunicationFamily}
						onSelectFamily={selectFamily}
						onCloseFamily={closeFamily}
					/>
				{:else}

					<!-- Workspace de prontuario: lista, cadastro e detalhe do aprendente selecionado. -->
					<LearnersWorkspace
						{activeSection}
						{filteredLearners}
						{selectedLearnerId}
						{selectedLearner}
						{learnerFilter}
						{guardianOptions}
						{showAddForm}
						{detailTab}
						calendarDays={selectedLearnerCalendarDays}
						{monthLabel}
						selectedDate={selectedAgendaDate}
						{selectedVisit}
						{isUploading}
						planCategories={PLAN_CATEGORIES}
						onOpenAddForm={() => (showAddForm = true)}
						onCloseAddForm={() => (showAddForm = false)}
						onCreateLearner={handleCreateLearner}
						onDeleteLearner={deleteLearner}
						onSelectLearner={selectLearner}
						onSetLearnerFilter={(filter) => (learnerFilter = filter)}
						onSelectTab={(tab) => (detailTab = tab)}
						onShiftMonth={shiftMonth}
						onSelectCalendarDate={handleCalendarDate}
						onUpdateLearner={updateSelectedLearner}
						onUpdateActionPlan={updateActionPlan}
						onAddCustomActionPlanField={addCustomActionPlanFieldToSelected}
						onUpdateCustomActionPlanField={updateCustomActionPlanFieldForSelected}
						onRemoveCustomActionPlanField={removeCustomActionPlanFieldFromSelected}
						onUpdateVisit={updateVisit}
						onRemoveVisit={removeVisit}
						onUploadDocuments={handleDocumentUpload}
						onDownloadDocument={downloadDocument}
						onRemoveDocument={removeDocument}
						onUploadAnamneseDocuments={handleAnamneseDocumentUpload}
						onDownloadAnamneseDocument={downloadAnamneseDocument}
						onRemoveAnamneseDocument={removeAnamneseDocument}
						onAddReport={addReport}
						onRemoveReport={removeReport}
						onOpenResponsible={openLearnerResponsible}
					/>
				{/if}
			</section>
		</div>

		{#if deletionConfirmation}
			<div class="confirmation-overlay" role="presentation">
				<button
					type="button"
					class="confirmation-backdrop"
					aria-label="Cancelar exclusao"
					onclick={() => closeDeletionConfirmation(false)}
				></button>

				<div
					class="confirmation-dialog"
					role="dialog"
					aria-modal="true"
					aria-labelledby="deletion-confirmation-title"
				>
					<div class="confirmation-icon" aria-hidden="true">!</div>
					<div class="confirmation-copy">
						<span>Acao irreversivel</span>
						<h2 id="deletion-confirmation-title">{deletionConfirmation.title}</h2>
						<p>{deletionConfirmation.message}</p>
					</div>
					<div class="confirmation-actions">
						<button
							type="button"
							class="secondary-button"
							onclick={() => closeDeletionConfirmation(false)}
						>
							{deletionConfirmation.cancelLabel}
						</button>
						<button
							type="button"
							class="danger-button confirmation-danger"
							onclick={() => closeDeletionConfirmation(true)}
						>
							{deletionConfirmation.confirmLabel}
						</button>
					</div>
				</div>
			</div>
		{/if}
	</div>
{/if}

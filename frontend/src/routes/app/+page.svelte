<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { AppSidebar, AppTopbar, WorkspaceBanner } from '$lib/modules/clinic-shell/components';
	import '$lib/modules/clinic-shell/styles/clinic-app.css';
	import type { Banner, NavSection } from '$lib/modules/clinic-shell/types';
	import { AgendaWorkspace, LearnersWorkspace } from '$lib/modules/learners/components';
	import { clearStoredSession, getStoredSession, type StoredSession } from '$lib/auth';
	import {
		PLAN_CATEGORIES,
		addCustomActionPlanField,
		appendAnamneseDocumentsToLearner,
		appendDocumentsToLearner,
		buildCalendarDays,
		createId,
		createLearner,
		createReportEntry,
		createSessionVisit,
		deleteDocumentBlob,
		filterLearners,
		getDocumentBlob,
		getDocumentStorageKey,
		loadLearners,
		patchLearnerInList,
		prepareDocumentBlob,
		prependReportToLearner,
		putDocumentBlob,
		removeCustomActionPlanField,
		removeReportFromLearner,
		removeVisitFromList,
		saveLearners,
		sortVisitsBySchedule,
		toDateInputValue,
		updateActionPlanValue,
		updateCustomActionPlanField,
		updateVisitInList,
		type CoreActionPlanKey,
		type DetailTab,
		type Learner,
		type LearnerDocument,
		type LearnerFilter,
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

	let session = $state<StoredSession | null>(null);
	let learners = $state<Learner[]>([]);
	let agendaEvents = $state<AgendaEvent[]>([]);
	let selectedLearnerId = $state<string | null>(null);
	let activeSection = $state<NavSection>('aprendentes');
	let learnerFilter = $state<LearnerFilter>('active');
	let detailTab = $state<DetailTab>('resumo');
	let currentMonth = $state(new Date());
	let selectedAgendaDate = $state(toDateInputValue(new Date()));
	let selectedVisitId = $state<string | null>(null);
	let searchTerm = $state('');
	let isUploading = $state(false);
	let showAddForm = $state(false);
	let banner = $state<Banner | null>(null);

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
	const monthLabel = $derived(formatMonth(currentMonth));
	const currentDateLabel = $derived(formatLongDate(selectedAgendaDate));
	const tenantName = $derived(session?.payload.tenant?.name ?? 'PsicoClinica');
	const userName = $derived(session?.payload.user?.name ?? 'Usuario');

	onMount(() => {
		if (!browser) return;

		const storedSession = getStoredSession();
		if (!storedSession) {
			void goto('/');
			return;
		}

		session = storedSession;
		learners = loadLearners();
		agendaEvents = loadAgendaEvents();
		selectedLearnerId = learners[0]?.id ?? null;
	});

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

	// Atualiza qualquer aprendente por id, mantendo a data de edicao em um unico lugar.
	function updateLearnerById(learnerId: string, patch: Partial<Learner>) {
		const nextLearners = patchLearnerInList(learners, learnerId, patch);
		learners = nextLearners;
		persistLearners(nextLearners);
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
	function removeCustomActionPlanFieldFromSelected(fieldId: string) {
		if (!selectedLearner) return;

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

		const createdLearner = createLearner(input);
		const nextLearners = [createdLearner, ...learners];
		learners = nextLearners;
		selectedLearnerId = createdLearner.id;
		activeSection = 'aprendentes';
		detailTab = 'resumo';
		showAddForm = false;
		persistLearners(nextLearners);

		banner = {
			tone: 'success',
			text: 'Aprendente adicionado ao painel.'
		};
		return true;
	}

	// Abre o prontuario do aprendente e volta para a secao principal quando necessario.
	function selectLearner(id: string) {
		selectedLearnerId = id;
		activeSection = 'aprendentes';
		detailTab = 'resumo';
		selectedVisitId = null;
		banner = null;
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

		if (section === 'agenda') {
			detailTab = 'agenda';
		} else if (section === 'avaliacoes') {
			detailTab = 'anamnese';
		} else if (section === 'relatorios') {
			detailTab = 'relatorios';
		}
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
		activeSection = 'agenda';
		detailTab = 'agenda';
		showAddForm = false;
		banner = null;
	}

	// Remove uma visita e sincroniza a contagem exibida no cadastro do aprendente.
	function removeVisit(id: string) {
		if (!selectedLearner) return;
		const nextVisits = removeVisitFromList(selectedLearner.visits, id);

		updateSelectedLearner({
			visits: nextVisits,
			visitCount: nextVisits.length
		});
		selectedVisitId = null;
	}

	// Remove uma sessao diretamente da agenda diaria, mesmo quando outro aprendente esta aberto.
	function removeSessionAppointment(learnerId: string, visitId: string) {
		const learner = learners.find((item) => item.id === learnerId);
		if (!learner) return;

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
	function removeAgendaEvent(event: AgendaEvent) {
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

		await deleteDocumentBlob(getDocumentStorageKey(selectedLearner.id, document.id));
		updateSelectedLearner({
			documents: selectedLearner.documents.filter((item) => item.id !== document.id)
		});
	}

	// Remove anexos da anamnese sem afetar os documentos gerais do aprendente.
	async function removeAnamneseDocument(document: LearnerDocument) {
		if (!selectedLearner) return;

		await deleteDocumentBlob(getDocumentStorageKey(selectedLearner.id, document.id));
		updateSelectedLearner({
			anamneseDocuments: selectedLearner.anamneseDocuments.filter((item) => item.id !== document.id)
		});
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
	function removeReport(id: string) {
		if (!selectedLearner) return;

		const nextLearner = removeReportFromLearner(selectedLearner, id);
		updateSelectedLearner({
			reports: nextLearner.reports
		});
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

{#if !session}
	<main class="loading-screen">
		<p>Carregando sessao...</p>
	</main>
{:else}
	<div class="clinic-shell">
		<div class="window-bar">
			<span></span>
			<span></span>
			<span></span>
		</div>

		<div class="app-frame">
			<AppSidebar
				{tenantName}
				{activeSection}
				onSelectSection={selectSection}
			/>

			<section class="content">
				<AppTopbar
					{searchTerm}
					{userName}
					onSearchTermChange={(value) => (searchTerm = value)}
					onLogout={logout}
				/>

				<WorkspaceBanner {banner} />

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
						onShiftMonth={shiftMonth}
						onSelectCalendarDate={handleCalendarDate}
						onSelectLearnerId={selectLearnerInsideAgenda}
						onOpenLearner={selectLearner}
						onCreateSession={createSessionAppointment}
						onCreateEvent={createEventAppointment}
						onRemoveSession={removeSessionAppointment}
						onRemoveEvent={removeAgendaEvent}
					/>
				{:else}
					<LearnersWorkspace
						{activeSection}
						{filteredLearners}
						{selectedLearnerId}
						{selectedLearner}
						{learnerFilter}
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
					/>
				{/if}
			</section>
		</div>
	</div>
{/if}

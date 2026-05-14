<script lang="ts">
	import { fly } from 'svelte/transition';
	import {
		getLearnerGuardianEntries,
		type ActionPlan,
		type CalendarDay,
		type CoreActionPlanKey,
		type Learner,
		type LearnerDocument,
		type LearnerStatus,
		type PlanCategory,
		type Visit
	} from '$lib/modules/learners';
	import { formatCurrencyFromCents } from '$lib/shared/formatters';
	import { DETAIL_TABS, type DetailTab } from '../../presentation/types';
	import LearnerAvatar from '../avatar/LearnerAvatar.svelte';
	import ActionPlanTab from '../tabs/ActionPlanTab.svelte';
	import AgendaTab from '../tabs/AgendaTab.svelte';
	import AnamneseTab from '../tabs/AnamneseTab.svelte';
	import DocumentsTab from '../tabs/DocumentsTab.svelte';
	import ReportsTab from '../tabs/ReportsTab.svelte';
	import SummaryTab from '../tabs/SummaryTab.svelte';

	let {
		learner,
		detailTab,
		calendarDays,
		monthLabel,
		selectedDate,
		selectedVisit,
		isUploading,
		planCategories,
		onSelectTab,
		onShiftMonth,
		onSelectCalendarDate,
		onUpdateLearner,
		onDeleteLearner,
		onOpenAgendaWorkspace,
		onUpdateActionPlan,
		onAddCustomActionPlanField,
		onUpdateCustomActionPlanField,
		onRemoveCustomActionPlanField,
		onUpdateVisit,
		onRemoveVisit,
		onUploadDocuments,
		onDownloadDocument,
		onRemoveDocument,
		onUploadAnamneseDocuments,
		onDownloadAnamneseDocument,
		onRemoveAnamneseDocument,
		onAddReport,
		onRemoveReport,
		onOpenResponsible
	} = $props<{
		learner: Learner | null;
		detailTab: DetailTab;
		calendarDays: CalendarDay[];
		monthLabel: string;
		selectedDate: string;
		selectedVisit: Visit | null;
		isUploading: boolean;
		planCategories: PlanCategory[];
		onSelectTab: (tab: DetailTab) => void;
		onShiftMonth: (delta: number) => void;
		onSelectCalendarDate: (date: string) => void;
		onUpdateLearner: (patch: Partial<Learner>) => void;
		onDeleteLearner: (learnerId: string) => void | Promise<void>;
		onOpenAgendaWorkspace: () => void;
		onUpdateActionPlan: (key: CoreActionPlanKey, value: string) => void;
		onAddCustomActionPlanField: (label: string, description: string) => boolean;
		onUpdateCustomActionPlanField: (fieldId: string, value: string) => void;
		onRemoveCustomActionPlanField: (fieldId: string) => void;
		onUpdateVisit: (visitId: string, patch: Partial<Visit>) => void;
		onRemoveVisit: (visitId: string) => void;
		onUploadDocuments: (event: Event) => void | Promise<void>;
		onDownloadDocument: (document: LearnerDocument) => void | Promise<void>;
		onRemoveDocument: (document: LearnerDocument) => void | Promise<void>;
		onUploadAnamneseDocuments: (event: Event) => void | Promise<void>;
		onDownloadAnamneseDocument: (document: LearnerDocument) => void | Promise<void>;
		onRemoveAnamneseDocument: (document: LearnerDocument) => void | Promise<void>;
		onAddReport: (text: string) => void;
		onRemoveReport: (id: string) => void;
		onOpenResponsible: (learner: Learner) => void;
	}>();

	const learnerGuardians = $derived(learner ? getLearnerGuardianEntries(learner) : []);

	function handleStatusChange(event: Event) {
		if (!learner) return;

		const status = (event.currentTarget as HTMLSelectElement).value as LearnerStatus;
		if (status === learner.status) return;

		onUpdateLearner({ status });
	}
</script>

<aside class="detail-column">
	{#if learner}
		{#key learner.id}
			<div
				class="detail-motion"
				in:fly={{ x: 18, duration: 180 }}
				out:fly={{ x: -12, duration: 120 }}
			>
				<!-- Cabecalho do prontuario: identifica o aprendente e seu status clinico. -->
				<div class="profile-head">
					<LearnerAvatar name={learner.name} photoUrl={learner.photoUrl} size="large" />
					<div>
						<div class="profile-title">
							<h2>{learner.name}</h2>
							<label class="profile-status-field">
							
								<select
									class={learner.status}
									value={learner.status}
									onchange={handleStatusChange}
									aria-label="Alterar status do aprendente"
								>
									<option value="active">Ativo</option>
									<option value="inactive">Inativo</option>
								</select>
							</label>
						</div>
						<p>{learner.age || 'Idade nao informada'} - {learner.gender || 'Genero nao informado'}</p>
						<div class="profile-responsible-row">
							<div class="profile-responsible-list">
								{#if learnerGuardians.length}
									{#each learnerGuardians as responsible, index}
										<p>
											Responsavel {index + 1}: {responsible.name}
											{responsible.relationship ? ` (${responsible.relationship})` : ''}
										</p>
									{/each}
								{:else}
									<p>Responsavel: Nao informado</p>
								{/if}
							</div>
							{#if learnerGuardians.length}
								<button type="button" onclick={() => onOpenResponsible(learner)}>Abrir</button>
							{/if}
						</div>
						<p>
							Valor por sessao:
							{learner.sessionPriceCents > 0
								? formatCurrencyFromCents(learner.sessionPriceCents)
								: 'Nao informado'}
						</p>
						<p>
							Valor geral:
							{learner.generalValueCents > 0
								? formatCurrencyFromCents(learner.generalValueCents)
								: 'Nao informado'}
						</p>
					</div>
					<div class="profile-head-actions">
						<button
							type="button"
							class="danger-button"
							onclick={() => onDeleteLearner(learner.id)}
						>
							Excluir aprendente
						</button>
					</div>
				</div>

				<!-- Navegacao interna: cada aba isola uma area funcional do prontuario. -->
				<nav class="detail-tabs" aria-label="Secoes do aprendente">
					{#each DETAIL_TABS as tab}
						<button
							type="button"
							class:active={detailTab === tab.value}
							onclick={() => onSelectTab(tab.value)}
						>
							{tab.label}
						</button>
					{/each}
				</nav>

				<!-- Conteudo da aba ativa: mantem cada secao separada em seu proprio componente. -->
				{#key `${learner.id}-${detailTab}`}
					<div
						class="tab-motion"
						in:fly={{ x: 16, duration: 170 }}
						out:fly={{ x: -10, duration: 110 }}
					>
						{#if detailTab === 'resumo'}
							<SummaryTab learner={learner} onOpenAgenda={() => onSelectTab('agenda')} />
						{:else if detailTab === 'agenda'}
							<AgendaTab
								{calendarDays}
								{monthLabel}
								{selectedDate}
								{selectedVisit}
								onShiftMonth={onShiftMonth}
								onSelectCalendarDate={onSelectCalendarDate}
								onUpdateVisit={onUpdateVisit}
								onRemoveVisit={onRemoveVisit}
								onOpenAgendaWorkspace={onOpenAgendaWorkspace}
							/>
						{:else if detailTab === 'anamnese'}
							<AnamneseTab
								value={learner.anamnese}
								documents={learner.anamneseDocuments}
								{isUploading}
								onChange={(value) => onUpdateLearner({ anamnese: value })}
								onUpload={onUploadAnamneseDocuments}
								onDownload={onDownloadAnamneseDocument}
								onRemove={onRemoveAnamneseDocument}
							/>
						{:else if detailTab === 'documentos'}
							<DocumentsTab
								documents={learner.documents}
								{isUploading}
								onUpload={onUploadDocuments}
								onDownload={onDownloadDocument}
								onRemove={onRemoveDocument}
							/>
						{:else if detailTab === 'plano'}
							<ActionPlanTab
								actionPlan={learner.actionPlan}
								categories={planCategories}
								onChange={onUpdateActionPlan}
								onAddCustomField={onAddCustomActionPlanField}
								onChangeCustomField={onUpdateCustomActionPlanField}
								onRemoveCustomField={onRemoveCustomActionPlanField}
							/>
						{:else}
							<ReportsTab
								reports={learner.reports}
								onAddReport={onAddReport}
								onRemoveReport={onRemoveReport}
							/>
						{/if}
					</div>
				{/key}
			</div>
		{/key}
	{:else}
		<!-- Estado vazio da coluna de detalhe quando nenhum aprendente foi selecionado. -->
		<div class="empty-state">Selecione ou adicione um aprendente.</div>
	{/if}
</aside>

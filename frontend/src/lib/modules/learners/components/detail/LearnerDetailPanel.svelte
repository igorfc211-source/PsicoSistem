<script lang="ts">
	import type {
		ActionPlan,
		CalendarDay,
		CoreActionPlanKey,
		Learner,
		LearnerDocument,
		PlanCategory,
		Visit
	} from '$lib/modules/learners';
	import { getInitials } from '$lib/shared/formatters';
	import { DETAIL_TABS, type DetailTab } from '../../presentation/types';
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
		onRemoveReport
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
	}>();
</script>

<aside class="detail-column">
	{#if learner}
		<div class="profile-head">
			<div class="avatar large">{getInitials(learner.name)}</div>
			<div>
				<div class="profile-title">
					<h2>{learner.name}</h2>
					<span class={learner.status}>{learner.status === 'active' ? 'Ativo' : 'Inativo'}</span>
				</div>
				<p>{learner.age || 'Idade nao informada'} - {learner.gender || 'Genero nao informado'}</p>
				<p>Responsavel: {learner.guardian || 'Nao informado'}</p>
			</div>
		</div>

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
	{:else}
		<div class="empty-state">Selecione ou adicione um aprendente.</div>
	{/if}
</aside>

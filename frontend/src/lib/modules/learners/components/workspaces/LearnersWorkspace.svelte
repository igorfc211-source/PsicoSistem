<script lang="ts">
	import { fly } from 'svelte/transition';
	import type { NavSection } from '$lib/modules/clinic-shell/types';
	import { getSectionTitle } from '$lib/modules/clinic-shell/types';
	import type {
		ActionPlan,
		CalendarDay,
		CoreActionPlanKey,
		GuardianOption,
		Learner,
		LearnerDocument,
		NewLearnerInput,
		PlanCategory,
		Visit
	} from '$lib/modules/learners';
	import type { DetailTab, LearnerFilter } from '../../presentation/types';
	import LearnerAddForm from '../forms/LearnerAddForm.svelte';
	import LearnerDetailPanel from '../detail/LearnerDetailPanel.svelte';
	import LearnerList from '../list/LearnerList.svelte';

	let {
		activeSection,
		filteredLearners,
		selectedLearnerId,
		selectedLearner,
		learnerFilter,
		guardianOptions,
		showAddForm,
		detailTab,
		calendarDays,
		monthLabel,
		selectedDate,
		selectedVisit,
		isUploading,
		planCategories,
		onOpenAddForm,
		onCloseAddForm,
		onCreateLearner,
		onDeleteLearner,
		onSelectLearner,
		onSetLearnerFilter,
		onSelectTab,
		onShiftMonth,
		onSelectCalendarDate,
		onUpdateLearner,
		onUpdateActionPlan,
		onOpenAgendaWorkspace,
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
		activeSection: NavSection;
		filteredLearners: Learner[];
		selectedLearnerId: string | null;
		selectedLearner: Learner | null;
		learnerFilter: LearnerFilter;
		guardianOptions: GuardianOption[];
		showAddForm: boolean;
		detailTab: DetailTab;
		calendarDays: CalendarDay[];
		monthLabel: string;
		selectedDate: string;
		selectedVisit: Visit | null;
		isUploading: boolean;
		planCategories: PlanCategory[];
		onOpenAddForm: () => void;
		onCloseAddForm: () => void;
		onCreateLearner: (input: NewLearnerInput) => boolean;
		onDeleteLearner: (learnerId: string) => void | Promise<void>;
		onSelectLearner: (id: string) => void;
		onSetLearnerFilter: (filter: LearnerFilter) => void;
		onSelectTab: (tab: DetailTab) => void;
		onShiftMonth: (delta: number) => void;
		onSelectCalendarDate: (date: string) => void;
		onUpdateLearner: (patch: Partial<Learner>) => void;
		onUpdateActionPlan: (key: CoreActionPlanKey, value: string) => void;
		onOpenAgendaWorkspace: () => void;
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
</script>

<section class="learners-workspace">
	<!-- Coluna esquerda: controles de lista, filtros, cadastro e selecao de aprendentes. -->
	<div class="learners-column">
		<!-- Cabecalho da lista com total filtrado e acao de adicionar. -->
		<div class="section-title">
			<div>
				<h1>{getSectionTitle(activeSection)}</h1>
				<p>{filteredLearners.length} aprendentes</p>
			</div>

			<div class="title-actions">
				<button type="button" class="icon-button" aria-label="Filtros">=</button>
				{#if !showAddForm}
				<button
					type="button"
					class="primary-button"
					onclick={onOpenAddForm}
				>
					+ Adicionar
				</button>
			{:else}
				<button
					type="button"
					class="primary-button close-button"
					onclick={onCloseAddForm}
				>
					✕ Fechar
				</button>
			{/if}
			</div>
		</div>

		<!-- Filtros de status para alternar rapidamente entre ativos, inativos e todos. -->
		<div class="filter-pills" aria-label="Filtrar aprendentes">
			<button
				type="button"
				class:active={learnerFilter === 'active'}
				onclick={() => onSetLearnerFilter('active')}
			>
				Ativos
			</button>
			<button
				type="button"
				class:active={learnerFilter === 'inactive'}
				onclick={() => onSetLearnerFilter('inactive')}
			>
				Inativos
			</button>
			<button
				type="button"
				class:active={learnerFilter === 'all'}
				onclick={() => onSetLearnerFilter('all')}
			>
				Todos
			</button>
		</div>

		<!-- Formulario expansivel: fica perto da lista para manter o fluxo de cadastro curto. -->
		{#if showAddForm}
			<div class="motion-panel" transition:fly={{ y: -14, duration: 180 }}>
				<LearnerAddForm
					{guardianOptions}
					onCreate={onCreateLearner}
					onCancel={onCloseAddForm}
				/>
			</div>
		{/if}

		<!-- Lista navegavel de aprendentes; ao selecionar, o detalhe abre na coluna direita. -->
		<LearnerList
			learners={filteredLearners}
			{selectedLearnerId}
			onSelectLearner={onSelectLearner}
		/>
	</div>

	<!-- Coluna direita: prontuario completo do aprendente selecionado. -->
	<LearnerDetailPanel
		learner={selectedLearner}
		{detailTab}
		{calendarDays}
		{monthLabel}
		{selectedDate}
		{selectedVisit}
		{isUploading}
		{planCategories}
		onSelectTab={onSelectTab}
		onShiftMonth={onShiftMonth}
		onSelectCalendarDate={onSelectCalendarDate}
		onUpdateLearner={onUpdateLearner}
		onDeleteLearner={onDeleteLearner}
		onOpenAgendaWorkspace={onOpenAgendaWorkspace}
		onUpdateActionPlan={onUpdateActionPlan}
		onAddCustomActionPlanField={onAddCustomActionPlanField}
		onUpdateCustomActionPlanField={onUpdateCustomActionPlanField}
		onRemoveCustomActionPlanField={onRemoveCustomActionPlanField}
		onUpdateVisit={onUpdateVisit}
		onRemoveVisit={onRemoveVisit}
		onUploadDocuments={onUploadDocuments}
		onDownloadDocument={onDownloadDocument}
		onRemoveDocument={onRemoveDocument}
		onUploadAnamneseDocuments={onUploadAnamneseDocuments}
		onDownloadAnamneseDocument={onDownloadAnamneseDocument}
		onRemoveAnamneseDocument={onRemoveAnamneseDocument}
		onAddReport={onAddReport}
		onRemoveReport={onRemoveReport}
		onOpenResponsible={onOpenResponsible}
	/>
</section>

<style>
	.close-button {
	background: #e04545;
	color: #1f1f1f;
	border: 1px solid #f07f7f;
}

.close-button:hover {
	background: #e26d6d;
}
</style>
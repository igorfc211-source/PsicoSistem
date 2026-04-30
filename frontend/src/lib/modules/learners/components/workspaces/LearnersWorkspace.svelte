<script lang="ts">
	import type { NavSection } from '$lib/modules/clinic-shell/types';
	import { getSectionTitle } from '$lib/modules/clinic-shell/types';
	import type {
		ActionPlan,
		CalendarDay,
		CoreActionPlanKey,
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
		onSelectLearner,
		onSetLearnerFilter,
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
		activeSection: NavSection;
		filteredLearners: Learner[];
		selectedLearnerId: string | null;
		selectedLearner: Learner | null;
		learnerFilter: LearnerFilter;
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
		onSelectLearner: (id: string) => void;
		onSetLearnerFilter: (filter: LearnerFilter) => void;
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

<section class="learners-workspace">
	<div class="learners-column">
		<div class="section-title">
			<div>
				<h1>{getSectionTitle(activeSection)}</h1>
				<p>{filteredLearners.length} aprendentes</p>
			</div>

			<div class="title-actions">
				<button type="button" class="icon-button" aria-label="Filtros">=</button>
				<button type="button" class="primary-button" onclick={onOpenAddForm}>+ Adicionar</button>
			</div>
		</div>

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

		{#if showAddForm}
			<LearnerAddForm onCreate={onCreateLearner} onCancel={onCloseAddForm} />
		{/if}

		<LearnerList
			learners={filteredLearners}
			{selectedLearnerId}
			onSelectLearner={onSelectLearner}
		/>
	</div>

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
	/>
</section>

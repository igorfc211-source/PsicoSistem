import type { NewSessionAppointmentInput } from '$lib/modules/scheduling';
import { createId } from '../domain/factories';
import type {
	ActionPlan,
	Learner,
	LearnerDocument,
	LearnerReport,
	Visit
} from '../domain/types';
import type { LearnerFilter } from '../presentation/types';

// Filtra aprendentes por busca textual e status, mantendo a tela livre de regra de lista.
export function filterLearners(
	learners: Learner[],
	searchTerm: string,
	learnerFilter: LearnerFilter
) {
	const query = searchTerm.trim().toLowerCase();

	return learners.filter((learner) => {
		const matchesSearch =
			!query ||
			[learner.name, learner.guardian, learner.gender, learner.age]
				.join(' ')
				.toLowerCase()
				.includes(query);
		const matchesFilter = learnerFilter === 'all' || learner.status === learnerFilter;

		return matchesSearch && matchesFilter;
	});
}

// Aplica patches em um aprendente especifico e atualiza o carimbo de edicao.
export function patchLearnerInList(
	learners: Learner[],
	learnerId: string,
	patch: Partial<Learner>
) {
	return learners.map((learner) =>
		learner.id === learnerId
			? {
					...learner,
					...patch,
					updatedAt: new Date().toISOString()
				}
			: learner
	);
}

// Altera uma secao do plano de acao sem recriar manualmente o objeto na rota.
export function updateActionPlanValue(
	actionPlan: ActionPlan,
	key: keyof ActionPlan,
	value: string
) {
	return {
		...actionPlan,
		[key]: value
	};
}

// Ordena visitas pelo dia e horario para qualquer lista clinica usar a mesma regra.
export function sortVisitsBySchedule(left: Visit, right: Visit) {
	return `${left.date}-${left.startTime}`.localeCompare(`${right.date}-${right.startTime}`);
}

// Atualiza uma visita existente e devolve a agenda ja ordenada.
export function updateVisitInList(visits: Visit[], visitId: string, patch: Partial<Visit>) {
	return visits
		.map((visit) => (visit.id === visitId ? { ...visit, ...patch } : visit))
		.sort(sortVisitsBySchedule);
}

// Remove uma visita e centraliza a regra para que contagem e persistencia fiquem previsiveis.
export function removeVisitFromList(visits: Visit[], visitId: string) {
	return visits.filter((visit) => visit.id !== visitId);
}

// Converte o formulario profissional de agenda em uma visita do aprendente.
export function createSessionVisit(input: NewSessionAppointmentInput): Visit {
	return {
		id: createId('visit'),
		date: input.date,
		title: input.title.trim() || 'Sessao individual',
		startTime: input.startTime,
		endTime: input.endTime,
		kind: 'session',
		location: input.location.trim() || 'Consultorio',
		status: 'scheduled',
		notes: input.notes.trim()
	};
}

// Anexa documentos preparados ao aprendente, sem misturar IndexedDB com UI.
export function appendDocumentsToLearner(
	learner: Learner,
	documents: LearnerDocument[]
): Learner {
	return {
		...learner,
		documents: [...learner.documents, ...documents],
		updatedAt: new Date().toISOString()
	};
}

// Cria o registro de relatorio com data de criacao e atualizacao sincronizadas.
export function createReportEntry(text: string): LearnerReport {
	const now = new Date().toISOString();

	return {
		id: createId('report'),
		text: text.trim(),
		createdAt: now,
		updatedAt: now
	};
}

// Coloca o relatorio mais recente no topo, como historico clinico de leitura rapida.
export function prependReportToLearner(learner: Learner, report: LearnerReport): Learner {
	return {
		...learner,
		reports: [report, ...learner.reports],
		updatedAt: new Date().toISOString()
	};
}

// Exclui apenas o relatorio escolhido, preservando o restante do historico.
export function removeReportFromLearner(learner: Learner, reportId: string): Learner {
	return {
		...learner,
		reports: learner.reports.filter((report) => report.id !== reportId),
		updatedAt: new Date().toISOString()
	};
}

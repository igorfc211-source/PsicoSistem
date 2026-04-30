import type { Learner, Visit, VisitKind } from '$lib/modules/learners';

export type AgendaEventKind = 'meeting' | 'supervision' | 'block' | 'other';

export type AgendaEvent = {
	id: string;
	title: string;
	date: string;
	startTime: string;
	endTime: string;
	kind: AgendaEventKind;
	description: string;
	createdAt: string;
	updatedAt: string;
};

export type NewAgendaEventInput = {
	title: string;
	date: string;
	startTime: string;
	endTime: string;
	kind: AgendaEventKind;
	description: string;
};

export type NewSessionAppointmentInput = {
	learnerId: string;
	date: string;
	startTime: string;
	endTime: string;
	kind: VisitKind;
	title: string;
	notes: string;
	location: string;
};

export type ScheduleItem = {
	id: string;
	date: string;
	startTime: string;
	endTime: string;
	title: string;
	subtitle: string;
	kind: 'session' | 'event';
	tone: 'purple' | 'blue' | 'green' | 'amber' | 'gray';
	learner?: Learner;
	visit?: Visit;
	event?: AgendaEvent;
};

export type LearnerStatus = 'active' | 'inactive';
export type VisitStatus = 'scheduled' | 'completed' | 'missed';
export type VisitKind = 'session' | 'assessment' | 'return';

export type Visit = {
	id: string;
	date: string;
	title: string;
	startTime: string;
	endTime: string;
	kind: VisitKind;
	location: string;
	status: VisitStatus;
	notes: string;
};

export type LearnerDocument = {
	id: string;
	name: string;
	type: string;
	size: number;
	storedSize: number;
	compressed: boolean;
	createdAt: string;
};

export type LearnerReport = {
	id: string;
	text: string;
	createdAt: string;
	updatedAt: string;
};

export type ActionPlanCustomField = {
	id: string;
	label: string;
	description: string;
	value: string;
};

export type ActionPlan = {
	educationProcess: string;
	familyGuidance: string;
	cognitiveSkills: string;
	behavior: string;
	clinicGoals: string;
	customFields: ActionPlanCustomField[];
};

export type CoreActionPlanKey = Exclude<keyof ActionPlan, 'customFields'>;

export type Learner = {
	id: string;
	name: string;
	photoUrl: string;
	gender: string;
	guardian: string;
	age: string;
	status: LearnerStatus;
	startDate: string;
	endDate: string;
	visitCount: number;
	anamnese: string;
	anamneseDocuments: LearnerDocument[];
	actionPlan: ActionPlan;
	visits: Visit[];
	documents: LearnerDocument[];
	reports: LearnerReport[];
	createdAt: string;
	updatedAt: string;
};

export type NewLearnerInput = {
	name: string;
	photoUrl: string;
	gender: string;
	guardian: string;
	age: string;
	status: LearnerStatus;
	startDate: string;
	endDate: string;
	visitCount: number;
};

export type CalendarDay = {
	date: string;
	day: number;
	inMonth: boolean;
	isToday: boolean;
	isSelected: boolean;
	eventCount: number;
	pendingVisitCount: number;
	visits: Visit[];
};

export type LearnerVisitSchedule = {
	learner: Learner;
	visit: Visit;
};

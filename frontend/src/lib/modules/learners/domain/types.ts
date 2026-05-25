export type LearnerStatus = 'active' | 'inactive';
export type VisitStatus = 'scheduled' | 'completed' | 'missed';
export type VisitKind = 'session' | 'assessment' | 'return';

export type LearnerGuardian = {
	id: string;
	sourceKey: string;
	name: string;
	relationship: string;
	phone: string;
};

export type LearnerGuardianInput = {
	sourceKey: string;
	name: string;
	relationship: string;
	phone: string;
};

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
	backendId?: string;
	backendGuardianIds?: string[];
	name: string;
	photoUrl: string;
	gender: string;
	guardian: string;
	guardianRelationship: string;
	guardians: LearnerGuardian[];
	age: string;
	status: LearnerStatus;
	startDate: string;
	endDate: string;
	visitCount: number;
	sessionPriceCents: number;
	generalValueCents: number;
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
	guardianRelationship: string;
	guardians: LearnerGuardianInput[];
	age: string;
	status: LearnerStatus;
	startDate: string;
	endDate: string;
	visitCount: number;
	sessionPriceCents: number;
	generalValueCents: number;
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

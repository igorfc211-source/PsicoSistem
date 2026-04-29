export const AUTH_STORAGE_KEY = 'psicosistem.auth';

export type ApiError = {
	code: string;
	message: string;
};

export type ApiResponse<T> = {
	data: T | null;
	meta: unknown;
	error: ApiError | null;
};

export type AccountPermissions = Record<string, 'none' | 'own' | 'all'>;

export type Tenant = {
	id: string;
	name: string;
	slug: string;
	email: string;
	phone: string;
	status: string;
	created_at: string;
	updated_at: string;
};

export type User = {
	user_id: string;
	tenant_id: string;
	name: string;
	email: string;
	role: string;
	status: string;
	permissions: AccountPermissions;
	last_login_at?: string | null;
	created_at: string;
	updated_at: string;
};

export type SubscriptionSummary = {
	plan: string;
	status: string;
	amount_monthly: number;
	next_amount_monthly: number;
	renewal_at?: string | null;
	trial_ends_at?: string | null;
	has_tests_library: boolean;
	has_ai: boolean;
	has_guardian_portal: boolean;
	max_professionals: number;
	max_patients: number;
};

export type AuthPayload = {
	tenant: Tenant | null;
	user: User | null;
	subscription?: SubscriptionSummary | null;
	token: string;
};

export type StoredSession = {
	mode: 'login' | 'register';
	timestamp: string;
	payload: AuthPayload;
};

// Centraliza a leitura da sessao para que qualquer rota trate dados corrompidos do mesmo jeito.
export function getStoredSession() {
	const rawSession = localStorage.getItem(AUTH_STORAGE_KEY);
	if (!rawSession) return null;

	try {
		return JSON.parse(rawSession) as StoredSession;
	} catch {
		localStorage.removeItem(AUTH_STORAGE_KEY);
		return null;
	}
}

// Mantem a gravacao da sessao em uma unica fronteira de infraestrutura.
export function setStoredSession(session: StoredSession) {
	localStorage.setItem(AUTH_STORAGE_KEY, JSON.stringify(session));
}

// Remove credenciais locais antes de redirecionar para a tela de entrada.
export function clearStoredSession() {
	localStorage.removeItem(AUTH_STORAGE_KEY);
}

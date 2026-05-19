import type { AccountPermissions, ApiResponse, SubscriptionSummary, User } from '$lib/auth';

export type UserRole = 'owner' | 'admin' | 'coordinator' | 'professional' | 'financial';
export type UserStatus = 'active' | 'inactive';
export type PermissionScope = 'none' | 'own' | 'all';

export type TeamUser = User & {
	role: UserRole;
	status: UserStatus;
};

export type CreateTeamUserInput = {
	name: string;
	email: string;
	password: string;
	role: UserRole;
	status: UserStatus;
};

export type UpdateTeamUserInput = {
	name: string;
	email: string;
	role: UserRole;
	status: UserStatus;
};

export type TeamListMeta = {
	page: number;
	per_page: number;
	total: number;
	total_pages: number;
};

export const TEAM_ROLES: Array<{ value: UserRole; label: string }> = [
	{ value: 'admin', label: 'Administrador' },
	{ value: 'coordinator', label: 'Coordenador' },
	{ value: 'professional', label: 'Profissional' },
	{ value: 'financial', label: 'Financeiro' },
	{ value: 'owner', label: 'Responsavel da clinica' }
];

export const PERMISSION_FIELDS: Array<{ key: keyof AccountPermissions; label: string }> = [
	{ key: 'user_directory', label: 'Contas e equipe' },
	{ key: 'patients', label: 'Aprendentes' },
	{ key: 'services', label: 'Servicos' },
	{ key: 'calendar', label: 'Agenda' },
	{ key: 'finance', label: 'Financeiro' },
	{ key: 'ai_history', label: 'Historico IA' },
	{ key: 'plans', label: 'Planos e relatorios' }
];

export const PERMISSION_SCOPES: Array<{ value: PermissionScope; label: string }> = [
	{ value: 'none', label: 'Sem acesso' },
	{ value: 'own', label: 'Somente proprio' },
	{ value: 'all', label: 'Tudo da clinica' }
];

export async function fetchTeamUsers(token: string) {
	const payload = await requestTeamApi<TeamUser[]>('/api/users?per_page=100', token);
	return {
		users: payload.data ?? [],
		meta: payload.meta as TeamListMeta | null
	};
}

export async function createTeamUser(token: string, input: CreateTeamUserInput) {
	const payload = await requestTeamApi<TeamUser>('/api/users', token, {
		method: 'POST',
		body: input
	});
	return payload.data;
}

export async function updateTeamUser(token: string, userId: string, input: UpdateTeamUserInput) {
	const payload = await requestTeamApi<TeamUser>(`/api/users/${userId}`, token, {
		method: 'PATCH',
		body: input
	});
	return payload.data;
}

export async function deactivateTeamUser(token: string, userId: string) {
	await requestTeamApi<{ message: string }>(`/api/users/${userId}`, token, {
		method: 'DELETE'
	});
}

export async function updateTeamUserPermissions(
	token: string,
	userId: string,
	permissions: AccountPermissions
) {
	const payload = await requestTeamApi<{ permissions: AccountPermissions }>(
		`/api/users/${userId}/permissions`,
		token,
		{
			method: 'PATCH',
			body: { permissions }
		}
	);

	return payload.data?.permissions ?? permissions;
}

export function canManageTeam(user: User | null | undefined) {
	return Boolean(
		user &&
			(user.role === 'owner' || user.role === 'admin') &&
			user.permissions.user_directory === 'all'
	);
}

export function getPlanSeatSummary(subscription: SubscriptionSummary | null | undefined, users: TeamUser[]) {
	const activeUsers = users.filter((user) => user.status === 'active').length;
	const maxProfessionals = subscription?.max_professionals ?? 0;

	return {
		activeUsers,
		maxProfessionals,
		availableSeats: Math.max(0, maxProfessionals - activeUsers),
		isAtLimit: maxProfessionals > 0 && activeUsers >= maxProfessionals
	};
}

async function requestTeamApi<T>(
	path: string,
	token: string,
	options: { method?: string; body?: unknown } = {}
) {
	const response = await fetch(path, {
		method: options.method ?? 'GET',
		headers: {
			authorization: `Bearer ${token}`,
			...(options.body ? { 'content-type': 'application/json' } : {})
		},
		...(options.body ? { body: JSON.stringify(options.body) } : {})
	});
	const payload = (await response.json()) as ApiResponse<T>;

	if (!response.ok || payload.error) {
		throw new Error(payload.error?.message ?? 'Nao foi possivel concluir a acao de governanca.');
	}

	return payload;
}

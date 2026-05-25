import type { ApiResponse, StoredSession, User } from './session';

export class AccountApiError extends Error {
	constructor(
		message: string,
		public readonly status: number,
		public readonly code: string
	) {
		super(message);
		this.name = 'AccountApiError';
	}
}

export async function fetchCurrentAccount(token: string) {
	const response = await fetch('/api/account/me', {
		headers: {
			authorization: `Bearer ${token}`
		},
		cache: 'no-store'
	});
	const payload = (await response.json()) as ApiResponse<User>;

	if (!response.ok || payload.error || !payload.data) {
		throw new AccountApiError(
			payload.error?.message ?? 'Nao foi possivel buscar os dados da conta.',
			response.status,
			payload.error?.code ?? 'ACCOUNT_REQUEST_FAILED'
		);
	}

	return payload.data;
}

export function mergeCurrentAccountIntoSession(session: StoredSession, user: User): StoredSession {
	const tenant = session.payload.tenant?.id === user.tenant_id ? session.payload.tenant : null;

	return {
		...session,
		timestamp: new Date().toISOString(),
		payload: {
			...session.payload,
			tenant,
			user
		}
	};
}

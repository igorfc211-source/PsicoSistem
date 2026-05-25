import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const BACKEND_BASE_URL = 'http://localhost:8080';
const BACKEND_PREFIXES = ['/v1', ''] as const;

function buildQuery(url: URL) {
	const query = new URLSearchParams();
	for (const key of ['page', 'per_page', 'role', 'status', 'search']) {
		const value = url.searchParams.get(key);
		if (value) query.set(key, value);
	}

	const serializedQuery = query.toString();
	return serializedQuery ? `?${serializedQuery}` : '';
}

async function proxyUsersRequest(
	fetch: typeof globalThis.fetch,
	authorization: string,
	method: 'GET' | 'POST',
	query = '',
	payload?: unknown
) {
	for (const prefix of BACKEND_PREFIXES) {
		try {
			const upstream = await fetch(`${BACKEND_BASE_URL}${prefix}/users${query}`, {
				method,
				headers: {
					authorization,
					...(payload ? { 'content-type': 'application/json' } : {})
				},
				...(payload ? { body: JSON.stringify(payload) } : {})
			});

			let body: unknown = {
				data: null,
				meta: null,
				error: {
					code: 'INVALID_BACKEND_RESPONSE',
					message: 'A API respondeu em um formato inesperado.'
				}
			};

			try {
				body = await upstream.json();
			} catch {
				// Mantem o payload padrao quando o backend nao responde JSON.
			}

			if (upstream.status === 404 && prefix !== BACKEND_PREFIXES[BACKEND_PREFIXES.length - 1]) {
				continue;
			}

			return json(body, { status: upstream.status });
		} catch {
			if (prefix === BACKEND_PREFIXES[BACKEND_PREFIXES.length - 1]) {
				return json(
					{
						data: null,
						meta: null,
						error: {
							code: 'USERS_SERVICE_UNAVAILABLE',
							message: 'Nao foi possivel conectar ao backend em http://localhost:8080.'
						}
					},
					{ status: 503 }
				);
			}
		}
	}

	return json(
		{
			data: null,
			meta: null,
			error: {
				code: 'USERS_PROXY_ERROR',
				message: 'Nao foi possivel encaminhar a requisicao de usuarios para a API.'
			}
		},
		{ status: 500 }
	);
}

function missingAuthorization() {
	return json(
		{
			data: null,
			meta: null,
			error: {
				code: 'MISSING_AUTHORIZATION',
				message: 'Sessao autenticada nao encontrada para governanca.'
			}
		},
		{ status: 401 }
	);
}

export const GET: RequestHandler = async ({ request, url, fetch }) => {
	const authorization = request.headers.get('authorization');
	if (!authorization) return missingAuthorization();

	return proxyUsersRequest(fetch, authorization, 'GET', buildQuery(url));
};

export const POST: RequestHandler = async ({ request, fetch }) => {
	const authorization = request.headers.get('authorization');
	if (!authorization) return missingAuthorization();

	try {
		return proxyUsersRequest(fetch, authorization, 'POST', '', await request.json());
	} catch {
		return json(
			{
				data: null,
				meta: null,
				error: {
					code: 'INVALID_BODY',
					message: 'Nao foi possivel ler o corpo da requisicao.'
				}
			},
			{ status: 400 }
		);
	}
};

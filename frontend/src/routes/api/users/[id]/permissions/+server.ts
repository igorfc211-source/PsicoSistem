import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const BACKEND_BASE_URL = 'http://localhost:8080';
const BACKEND_PREFIXES = ['/v1', ''] as const;

type ProxyMethod = 'GET' | 'PATCH';

async function proxyPermissionsRequest(
	fetch: typeof globalThis.fetch,
	authorization: string,
	id: string,
	method: ProxyMethod,
	payload?: unknown
) {
	for (const prefix of BACKEND_PREFIXES) {
		try {
			const upstream = await fetch(`${BACKEND_BASE_URL}${prefix}/users/${id}/permissions`, {
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
				message: 'Nao foi possivel encaminhar a permissao do usuario para a API.'
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

export const GET: RequestHandler = async ({ params, request, fetch }) => {
	const authorization = request.headers.get('authorization');
	if (!authorization) return missingAuthorization();

	return proxyPermissionsRequest(fetch, authorization, params.id, 'GET');
};

export const PATCH: RequestHandler = async ({ params, request, fetch }) => {
	const authorization = request.headers.get('authorization');
	if (!authorization) return missingAuthorization();

	try {
		return proxyPermissionsRequest(fetch, authorization, params.id, 'PATCH', await request.json());
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

import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const BACKEND_BASE_URL = 'http://localhost:8080';
const BACKEND_PREFIXES = ['/v1', ''] as const;

type AuthMode = 'login' | 'register';

const isAllowedMode = (mode: string): mode is AuthMode => mode === 'login' || mode === 'register';

export const POST: RequestHandler = async ({ params, request, fetch }) => {
	if (!isAllowedMode(params.mode)) {
		return json(
			{
				data: null,
				meta: null,
				error: {
					code: 'AUTH_ROUTE_NOT_FOUND',
					message: 'Acao de autenticacao nao suportada.'
				}
			},
			{ status: 404 }
		);
	}

	let payload: unknown;

	try {
		payload = await request.json();
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

	for (const prefix of BACKEND_PREFIXES) {
		try {
			const upstream = await fetch(`${BACKEND_BASE_URL}${prefix}/auth/${params.mode}`, {
				method: 'POST',
				headers: {
					'content-type': 'application/json'
				},
				body: JSON.stringify(payload)
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
							code: 'AUTH_SERVICE_UNAVAILABLE',
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
				code: 'AUTH_PROXY_ERROR',
				message: 'Nao foi possivel encaminhar a requisicao para a API.'
			}
		},
		{ status: 500 }
	);
};

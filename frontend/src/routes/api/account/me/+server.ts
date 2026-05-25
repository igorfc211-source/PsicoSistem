import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const BACKEND_BASE_URL = 'http://localhost:8080';
const BACKEND_PREFIXES = ['/v1', ''] as const;
const BACKEND_ACCOUNT_PATHS = ['/account/me', '/users/me'] as const;
const NO_STORE_HEADERS = { 'cache-control': 'no-store' };

function apiResponse(body: unknown, status: number) {
	return json(body, { status, headers: NO_STORE_HEADERS });
}

function missingAuthorization() {
	return apiResponse(
		{
			data: null,
			meta: null,
			error: {
				code: 'MISSING_AUTHORIZATION',
				message: 'Sessao autenticada nao encontrada para buscar a conta.'
			}
		},
		401
	);
}

export const GET: RequestHandler = async ({ request, fetch }) => {
	const authorization = request.headers.get('authorization');
	if (!authorization?.trim()) return missingAuthorization();

	for (const [prefixIndex, prefix] of BACKEND_PREFIXES.entries()) {
		for (const [pathIndex, accountPath] of BACKEND_ACCOUNT_PATHS.entries()) {
			const isLastAttempt =
				prefixIndex === BACKEND_PREFIXES.length - 1 &&
				pathIndex === BACKEND_ACCOUNT_PATHS.length - 1;

			try {
				const upstream = await fetch(`${BACKEND_BASE_URL}${prefix}${accountPath}`, {
					headers: { authorization }
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

				if (upstream.status === 404 && !isLastAttempt) {
					continue;
				}

				return apiResponse(body, upstream.status);
			} catch {
				if (isLastAttempt) {
					return apiResponse(
						{
							data: null,
							meta: null,
							error: {
								code: 'ACCOUNT_SERVICE_UNAVAILABLE',
								message: 'Nao foi possivel conectar ao backend em http://localhost:8080.'
							}
						},
						503
					);
				}
			}
		}
	}

	return apiResponse(
		{
			data: null,
			meta: null,
			error: {
				code: 'ACCOUNT_PROXY_ERROR',
				message: 'Nao foi possivel encaminhar a requisicao da conta para a API.'
			}
		},
		500
	);
};

import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const BACKEND_BASE_URL = 'http://localhost:8080';
const BACKEND_PREFIXES = ['/v1', ''] as const;

function buildLearnerQuery(url: URL) {
	const query = new URLSearchParams();
	for (const key of ['page', 'per_page', 'status', 'search']) {
		const value = url.searchParams.get(key);
		if (value) query.set(key, value);
	}

	const serializedQuery = query.toString();
	return serializedQuery ? `?${serializedQuery}` : '';
}

export const GET: RequestHandler = async ({ request, url, fetch }) => {
	const authorization = request.headers.get('authorization');
	if (!authorization) {
		return json(
			{
				data: null,
				meta: null,
				error: {
					code: 'MISSING_AUTHORIZATION',
					message: 'Sessao autenticada nao encontrada para carregar aprendentes.'
				}
			},
			{ status: 401 }
		);
	}

	const query = buildLearnerQuery(url);

	for (const prefix of BACKEND_PREFIXES) {
		try {
			const upstream = await fetch(`${BACKEND_BASE_URL}${prefix}/learners${query}`, {
				headers: {
					authorization
				}
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
							code: 'LEARNERS_SERVICE_UNAVAILABLE',
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
				code: 'LEARNERS_PROXY_ERROR',
				message: 'Nao foi possivel encaminhar a requisicao de aprendentes para a API.'
			}
		},
		{ status: 500 }
	);
};

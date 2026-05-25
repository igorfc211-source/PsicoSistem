<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { ApiResponse } from '$lib/auth';

	type BannerTone = 'info' | 'success' | 'error';

	type Banner = {
		tone: BannerTone;
		title: string;
		message: string;
	};

	let isSubmitting = $state(false);
	let banner = $state<Banner | null>(null);
	let form = $state({
		password: '',
		confirmPassword: ''
	});

	const token = $derived(page.url.searchParams.get('token') ?? '');

	async function submitResetPassword() {
		if (!token) {
			banner = {
				tone: 'error',
				title: 'Link invalido',
				message: 'Solicite um novo link de recuperacao na tela de login.'
			};
			return;
		}

		if (form.password !== form.confirmPassword) {
			banner = {
				tone: 'error',
				title: 'Senhas diferentes',
				message: 'A confirmacao precisa ser igual a nova senha.'
			};
			return;
		}

		isSubmitting = true;
		banner = {
			tone: 'info',
			title: 'Atualizando senha',
			message: 'Validando o link de recuperacao.'
		};

		try {
			const response = await fetch('/api/auth/reset-password', {
				method: 'POST',
				headers: {
					'content-type': 'application/json'
				},
				body: JSON.stringify({
					token,
					password: form.password
				})
			});
			const result = (await response.json()) as ApiResponse<{ message: string }>;

			if (!response.ok || result.error) {
				throw new Error(result.error?.message ?? 'Nao foi possivel atualizar a senha.');
			}

			banner = {
				tone: 'success',
				title: 'Senha atualizada',
				message: result.data?.message ?? 'Agora voce ja pode entrar com a nova senha.'
			};
			form = { password: '', confirmPassword: '' };
		} catch (error) {
			banner = {
				tone: 'error',
				title: 'Falha na recuperacao',
				message: error instanceof Error ? error.message : 'Erro inesperado ao falar com a API.'
			};
		} finally {
			isSubmitting = false;
		}
	}

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		void submitResetPassword();
	}
</script>

<svelte:head>
	<title>PsicoSistem | Redefinir senha</title>
</svelte:head>

<main class="reset-shell">
	<section class="reset-panel">
		<div class="panel-header">
			<p class="eyebrow">Recuperacao de acesso</p>
			<h1>Criar nova senha</h1>
			<p>Digite uma senha nova para sua conta. O link do e-mail so pode ser usado uma vez.</p>
		</div>

		{#if banner}
			<div class={`banner ${banner.tone}`}>
				<strong>{banner.title}</strong>
				<p>{banner.message}</p>
			</div>
		{/if}

		<form class="reset-form" onsubmit={handleSubmit}>
			<label>
				<span>Nova senha</span>
				<input
					type="password"
					bind:value={form.password}
					autocomplete="new-password"
					placeholder="Digite a nova senha"
					required
					disabled={isSubmitting}
				/>
			</label>

			<label>
				<span>Confirmar senha</span>
				<input
					type="password"
					bind:value={form.confirmPassword}
					autocomplete="new-password"
					placeholder="Repita a nova senha"
					required
					disabled={isSubmitting}
				/>
			</label>

			<div class="action-row">
				<button type="button" class="ghost-button" onclick={() => goto('/')} disabled={isSubmitting}>
					Voltar ao login
				</button>
				<button class="primary-button" type="submit" disabled={isSubmitting || !token}>
					{isSubmitting ? 'Salvando...' : 'Salvar senha'}
				</button>
			</div>
		</form>
	</section>
</main>

<style>
	:global(body) {
		margin: 0;
		min-height: 100vh;
		background:
			linear-gradient(135deg, rgba(248, 244, 237, 0.96), rgba(235, 240, 238, 0.94)),
			repeating-linear-gradient(135deg, rgba(29, 73, 79, 0.04) 0 1px, transparent 1px 18px);
		color: #1d2628;
		font-family: 'Trebuchet MS', 'Segoe UI', sans-serif;
	}

	.reset-shell {
		box-sizing: border-box;
		display: grid;
		min-height: 100vh;
		place-items: center;
		padding: 1rem;
	}

	.reset-panel {
		display: grid;
		width: min(100%, 520px);
		gap: 1.25rem;
		border: 1px solid rgba(29, 73, 79, 0.14);
		border-radius: 18px;
		background: rgba(255, 253, 249, 0.9);
		box-shadow: 0 22px 60px rgba(35, 48, 52, 0.11);
		padding: 2rem;
	}

	.panel-header,
	.reset-form {
		display: grid;
		gap: 1rem;
	}

	.eyebrow {
		width: fit-content;
		border-radius: 999px;
		background: rgba(29, 103, 111, 0.1);
		color: #1d6770;
		font-size: 0.76rem;
		font-weight: 800;
		letter-spacing: 0.08em;
		padding: 0.35rem 0.7rem;
		text-transform: uppercase;
	}

	h1,
	p {
		margin: 0;
	}

	h1 {
		font-family: Georgia, 'Times New Roman', serif;
		font-size: 2rem;
		line-height: 1.05;
	}

	p {
		color: #52666a;
		line-height: 1.55;
	}

	.banner {
		display: grid;
		gap: 0.25rem;
		border: 1px solid rgba(29, 73, 79, 0.12);
		border-radius: 14px;
		padding: 0.95rem 1rem;
	}

	.banner.info {
		background: rgba(29, 103, 112, 0.08);
		color: #1d5960;
	}

	.banner.success {
		background: rgba(57, 126, 73, 0.1);
		color: #244e2d;
	}

	.banner.error {
		background: rgba(153, 64, 51, 0.1);
		color: #7a2b22;
	}

	label {
		display: grid;
		gap: 0.4rem;
	}

	label span {
		color: #42565a;
		font-size: 0.9rem;
		font-weight: 800;
	}

	input {
		box-sizing: border-box;
		width: 100%;
		border: 1px solid rgba(29, 73, 79, 0.16);
		border-radius: 12px;
		background: rgba(255, 255, 255, 0.9);
		color: #1d2628;
		font: inherit;
		outline: none;
		padding: 0.9rem 1rem;
	}

	input:focus {
		border-color: #1d6770;
		box-shadow: 0 0 0 4px rgba(29, 103, 112, 0.12);
	}

	.action-row {
		display: flex;
		justify-content: space-between;
		gap: 0.75rem;
	}

	button {
		border: 0;
		cursor: pointer;
		font: inherit;
	}

	.primary-button,
	.ghost-button {
		border-radius: 12px;
		font-weight: 800;
		padding: 0.9rem 1rem;
	}

	.primary-button {
		background: rgba(85, 37, 243, 0.795);
		color: white;
		box-shadow: 0 14px 30px rgba(155, 81, 56, 0.18);
	}

	.ghost-button {
		background: rgba(29, 103, 112, 0.08);
		color: #1d6770;
	}

	button:disabled,
	input:disabled {
		cursor: wait;
		opacity: 0.72;
	}

	@media (max-width: 560px) {
		.reset-panel {
			padding: 1.25rem;
		}

		.action-row {
			flex-direction: column;
		}
		
	}
</style>

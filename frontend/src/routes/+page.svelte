<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import {
		getStoredSession,
		setStoredSession,
		type ApiResponse,
		type AuthPayload,
		type StoredSession
	} from '$lib/auth';

	type AuthMode = 'login' | 'register';
	type AuthView = AuthMode | 'forgot';
	type RegisterStep = 'details' | 'payment';
	type PlanSlug = 'basico' | 'intermediario' | 'premium';
	type BannerTone = 'info' | 'success' | 'error';

	type Banner = {
		tone: BannerTone;
		title: string;
		message: string;
	};

	type PlanCard = {
		value: PlanSlug;
		label: string;
		description: string;
		price: number;
		eyebrow: string;
		highlight: string;
		freeTrial: boolean;
	};

	const plans: PlanCard[] = [
		{
			value: 'basico',
			label: 'Basico',
			description: 'Entrada organizada para clinicas menores.',
			price: 9700,
			eyebrow: 'Essencial',
			highlight: 'Cobranca mensal direta apos a confirmacao.',
			freeTrial: false
		},
		{
			value: 'intermediario',
			label: 'Intermediario',
			description: 'Mais agenda, equipe e portal familiar.',
			price: 14700,
			eyebrow: 'Trial',
			highlight: '1 mes gratis antes da primeira cobranca.',
			freeTrial: true
		},
		{
			value: 'premium',
			label: 'Premium',
			description: 'Capacidade ampliada e recursos completos.',
			price: 19700,
			eyebrow: 'Escala',
			highlight: 'Plano completo para operacoes maiores.',
			freeTrial: false
		}
	];

	let activeMode = $state<AuthView>('login');
	let registerStep = $state<RegisterStep>('details');
	let isSubmitting = $state(false);
	let banner = $state<Banner | null>(null);

	let loginForm = $state({
		email: '',
		password: ''
	});

	let forgotPasswordForm = $state({
		email: ''
	});

	let registerForm = $state({
		clinic_name: '',
		name: '',
		email: '',
		password: '',
		phone: '',
		cpf_cnpj: '',
		plan_slug: 'intermediario' as PlanSlug,
		payment_session_confirmed: false,
		accept_trial_terms: true
	});

	onMount(() => {
		if (!browser) return;
		if (getStoredSession()) {
			void goto('/app');
		}
	});

	function setMode(mode: AuthView) {
		activeMode = mode;
		registerStep = 'details';
		banner = null;
	}

	function getSelectedPlan() {
		return plans.find((plan) => plan.value === registerForm.plan_slug) ?? plans[1];
	}

	function formatCurrency(value: number) {
		return new Intl.NumberFormat('pt-BR', {
			style: 'currency',
			currency: 'BRL'
		}).format(value / 100);
	}

	function onlyDigits(value: string) {
		return value.replace(/\D/g, '');
	}

	function formatPhone(value: string) {
		const digits = onlyDigits(value).slice(0, 11);
		if (digits.length === 0) return '';
		if (digits.length <= 2) return `(${digits}`;
		if (digits.length <= 6) return `(${digits.slice(0, 2)}) ${digits.slice(2)}`;
		if (digits.length <= 10) {
			return `(${digits.slice(0, 2)}) ${digits.slice(2, 6)}-${digits.slice(6)}`;
		}
		return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7)}`;
	}

	function formatCpfCnpj(value: string) {
		const digits = onlyDigits(value).slice(0, 14);

		if (digits.length <= 11) {
			if (digits.length <= 3) return digits;
			if (digits.length <= 6) return `${digits.slice(0, 3)}.${digits.slice(3)}`;
			if (digits.length <= 9) {
				return `${digits.slice(0, 3)}.${digits.slice(3, 6)}.${digits.slice(6)}`;
			}
			return `${digits.slice(0, 3)}.${digits.slice(3, 6)}.${digits.slice(6, 9)}-${digits.slice(9)}`;
		}

		if (digits.length <= 2) return digits;
		if (digits.length <= 5) return `${digits.slice(0, 2)}.${digits.slice(2)}`;
		if (digits.length <= 8) return `${digits.slice(0, 2)}.${digits.slice(2, 5)}.${digits.slice(5)}`;
		if (digits.length <= 12) {
			return `${digits.slice(0, 2)}.${digits.slice(2, 5)}.${digits.slice(5, 8)}/${digits.slice(8)}`;
		}
		return `${digits.slice(0, 2)}.${digits.slice(2, 5)}.${digits.slice(5, 8)}/${digits.slice(8, 12)}-${digits.slice(12)}`;
	}

	function handlePhoneInput(event: Event) {
		registerForm.phone = formatPhone((event.currentTarget as HTMLInputElement).value);
	}

	function handleCpfCnpjInput(event: Event) {
		registerForm.cpf_cnpj = formatCpfCnpj((event.currentTarget as HTMLInputElement).value);
	}

	function getErrorMessage(result: ApiResponse<unknown> | null) {
		return result?.error?.message ?? 'Nao foi possivel concluir a autenticacao.';
	}

	function persistSession(mode: AuthMode, payload: AuthPayload) {
		const session: StoredSession = {
			mode,
			timestamp: new Date().toISOString(),
			payload
		};

		setStoredSession(session);
	}

	function openPaymentStep() {
		const requiredFields = [
			registerForm.clinic_name,
			registerForm.name,
			registerForm.email,
			registerForm.password,
			registerForm.phone,
			registerForm.cpf_cnpj
		];

		if (requiredFields.some((value) => !value.trim())) {
			banner = {
				tone: 'error',
				title: 'Complete os dados',
				message: 'Preencha os dados da clinica e do responsavel antes do checkout.'
			};
			return;
		}

		registerStep = 'payment';
		banner = {
			tone: 'info',
			title: 'Checkout iniciado',
			message: getSelectedPlan().freeTrial
				? 'O plano intermediario inicia com 1 mes gratis.'
				: 'Revise o plano e confirme a sessao antes de criar a conta.'
		};
	}

	async function submitAuth(mode: AuthMode) {
		isSubmitting = true;
		banner = {
			tone: 'info',
			title: mode === 'login' ? 'Validando acesso' : 'Criando conta',
			message: 'Estamos conversando com o backend em localhost:8080.'
		};

		const payload = mode === 'login' ? loginForm : registerForm;

		try {
			const response = await fetch(`/api/auth/${mode}`, {
				method: 'POST',
				headers: {
					'content-type': 'application/json'
				},
				body: JSON.stringify(payload)
			});

			const result = (await response.json()) as ApiResponse<AuthPayload>;

			if (!response.ok || !result.data || result.error) {
				throw new Error(getErrorMessage(result));
			}

			persistSession(mode, result.data);
			await goto('/app');
		} catch (error) {
			banner = {
				tone: 'error',
				title: mode === 'login' ? 'Falha no login' : 'Falha no cadastro',
				message: error instanceof Error ? error.message : 'Erro inesperado ao falar com a API.'
			};
		} finally {
			isSubmitting = false;
		}
	}

	async function submitForgotPassword() {
		isSubmitting = true;
		banner = {
			tone: 'info',
			title: 'Enviando recuperacao',
			message: 'Se a conta existir, enviaremos um link seguro para o e-mail informado.'
		};

		try {
			const response = await fetch('/api/auth/forgot-password', {
				method: 'POST',
				headers: {
					'content-type': 'application/json'
				},
				body: JSON.stringify(forgotPasswordForm)
			});

			const result = (await response.json()) as ApiResponse<{ message: string }>;
			if (!response.ok || result.error) {
				throw new Error(getErrorMessage(result));
			}

			banner = {
				tone: 'success',
				title: 'Confira seu e-mail',
				message:
					result.data?.message ??
					'Se o e-mail informado estiver cadastrado, enviaremos um link de recuperacao.'
			};
		} catch (error) {
			banner = {
				tone: 'error',
				title: 'Nao foi possivel enviar',
				message: error instanceof Error ? error.message : 'Erro inesperado ao falar com a API.'
			};
		} finally {
			isSubmitting = false;
		}
	}

	function handleLoginSubmit(event: SubmitEvent) {
		event.preventDefault();
		void submitAuth('login');
	}

	function handleRegisterSubmit(event: SubmitEvent) {
		event.preventDefault();
		void submitAuth('register');
	}

	function handleForgotPasswordSubmit(event: SubmitEvent) {
		event.preventDefault();
		void submitForgotPassword();
	}
</script>

<svelte:head>
	<title>PsicoSistem | Login</title>
</svelte:head>

<main class="auth-shell">
	<section class="auth-panel">
		<div class="panel-header">
			<div>
				<p class="eyebrow">Acesso da clinica</p>
				<h2>
					{activeMode === 'forgot'
						? 'Recuperar senha'
						: activeMode === 'login'
							? 'Entrar na plataforma'
							: 'Criar primeira conta'}
				</h2>
			</div>

			<div class="mode-switch" aria-label="Alternar formulario">
				<button
					type="button"
					class:active={activeMode === 'login'}
					onclick={() => setMode('login')}
				>
					Login
				</button>
				<button
					type="button"
					class:active={activeMode === 'register'}
					onclick={() => setMode('register')}
				>
					Cadastro
				</button>
			</div>
		</div>

		{#if banner}
			<div class={`banner ${banner.tone}`}>
				<strong>{banner.title}</strong>
				<p>{banner.message}</p>
			</div>
		{/if}

		{#if activeMode === 'login'}
			<form class="auth-form" onsubmit={handleLoginSubmit}>
				<label>
					<span>Email</span>
					<input
						type="email"
						bind:value={loginForm.email}
						placeholder="responsavel@clinica.com"
						autocomplete="email"
						required
						disabled={isSubmitting}
					/>
				</label>

				

				<label>
					<span>Senha</span>
					<input
						type="password"
						bind:value={loginForm.password}
						placeholder="Digite sua senha"
						autocomplete="current-password"
						required
						disabled={isSubmitting}
					/>
				</label>

				<button class="primary-button mt-12" type="submit" disabled={isSubmitting}>
					{isSubmitting ? 'Entrando...' : 'Entrar'}
				</button>

				<button
					type="button"
					class="text-button"
					onclick={() => {
						forgotPasswordForm.email = loginForm.email;
						setMode('forgot');
					}}
					disabled={isSubmitting}
				>
					Esqueci minha senha
				</button>
			</form>
		{:else if activeMode === 'forgot'}
			<form class="auth-form" onsubmit={handleForgotPasswordSubmit}>
				<p>
					Informe o e-mail da sua conta. Se ele estiver cadastrado, enviaremos um link para criar
					uma nova senha.
				</p>

				<label>
					<span>Email</span>
					<input
						type="email"
						bind:value={forgotPasswordForm.email}
						placeholder="responsavel@clinica.com"
						autocomplete="email"
						required
						disabled={isSubmitting}
					/>
				</label>

				<div class="action-row">
					<button type="button" class="ghost-button" onclick={() => setMode('login')} disabled={isSubmitting}>
						Voltar
					</button>
					<button class="primary-button" type="submit" disabled={isSubmitting}>
						{isSubmitting ? 'Enviando...' : 'Enviar link'}
					</button>
				</div>
			</form>
		{:else if registerStep === 'details'}
			<div class="auth-form">
				<div class="grid two-columns">
					<label>
						<span>Nome da clinica</span>
						<input
							type="text"
							bind:value={registerForm.clinic_name}
							placeholder="Clinica Horizonte"
							required
							disabled={isSubmitting}
						/>
					</label>

					<label>
						<span>Telefone unico</span>
						<input
							type="tel"
							value={registerForm.phone}
							oninput={handlePhoneInput}
							placeholder="(11) 99999-0000"
							inputmode="numeric"
							maxlength="15"
							required
							disabled={isSubmitting}
						/>
					</label>
				</div>

				<div class="grid two-columns">
					<label>
						<span>Responsavel</span>
						<input
							type="text"
							bind:value={registerForm.name}
							placeholder="Nome do owner"
							required
							disabled={isSubmitting}
						/>
					</label>

					<label>
						<span>CPF ou CNPJ</span>
						<input
							type="text"
							value={registerForm.cpf_cnpj}
							oninput={handleCpfCnpjInput}
							placeholder="Documento do responsavel ou clinica"
							inputmode="numeric"
							maxlength="18"
							required
							disabled={isSubmitting}
						/>
					</label>
				</div>

				<div class="grid two-columns">
					<label>
						<span>Email</span>
						<input
							type="email"
							bind:value={registerForm.email}
							placeholder="responsavel@clinica.com"
							autocomplete="email"
							required
							disabled={isSubmitting}
						/>
					</label>

					<label>
						<span>Senha</span>
						<input
							type="password"
							bind:value={registerForm.password}
							placeholder="Crie uma senha"
							autocomplete="new-password"
							required
							disabled={isSubmitting}
						/>
					</label>
				</div>

				<div class="action-row">
					<button type="button" class="ghost-button" onclick={() => setMode('login')}>Ja tenho conta</button>
					<button type="button" class="primary-button" onclick={openPaymentStep}>
						Continuar para checkout
					</button>
				</div>
			</div>
		{:else}
			<form class="auth-form" onsubmit={handleRegisterSubmit}>
				<div class="plan-options">
					{#each plans as plan}
						<button
							type="button"
							class:selected={registerForm.plan_slug === plan.value}
							onclick={() => {
								registerForm.plan_slug = plan.value;
								registerForm.accept_trial_terms = plan.freeTrial;
							}}
							disabled={isSubmitting}
						>
							<span>{plan.eyebrow}</span>
							<strong>{plan.label}</strong>
							<small>{plan.description}</small>
							<em>{plan.freeTrial ? '1 mes gratis' : `${formatCurrency(plan.price)}/mes`}</em>
						</button>
					{/each}
				</div>

				<div class="checkout-card">
					<div>
						<p class="eyebrow">Checkout</p>
						<h3>{getSelectedPlan().label}</h3>
					</div>
					<strong>{getSelectedPlan().freeTrial ? formatCurrency(0) : formatCurrency(getSelectedPlan().price)}</strong>
					<p>{getSelectedPlan().highlight}</p>
				</div>

				<label class="checkbox-field">
					<input
						type="checkbox"
						bind:checked={registerForm.payment_session_confirmed}
						required
						disabled={isSubmitting}
					/>
					<span>Confirmo a sessao de checkout e autorizo a criacao da conta.</span>
				</label>

				{#if getSelectedPlan().freeTrial}
					<label class="checkbox-field">
						<input
							type="checkbox"
							bind:checked={registerForm.accept_trial_terms}
							required
							disabled={isSubmitting}
						/>
						<span>
							Aceito o trial de 1 mes e a cobranca de {formatCurrency(getSelectedPlan().price)}
							apos o periodo gratuito.
						</span>
					</label>
				{/if}

				<div class="action-row">
					<button type="button" class="ghost-button" onclick={() => (registerStep = 'details')}>
						Voltar
					</button>
					<button class="primary-button" type="submit" disabled={isSubmitting}>
						{isSubmitting ? 'Criando...' : 'Criar conta'}
					</button>
				</div>
			</form>
		{/if}
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

	:global(code) {
		font-family: 'Consolas', 'Courier New', monospace;
	}

	.auth-shell {
		box-sizing: border-box;
		flex: auto;
		grid-template-columns: minmax(280px, 0.9fr) minmax(340px, 1.1fr);
		gap: 1rem;
		min-height: 100vh;
		padding: 9rem;
	}

	.brand-panel,
	.auth-panel {
		border: 1px solid rgba(29, 73, 79, 0.14);
		border-radius: 18px;
		background: rgba(255, 253, 249, 0.86);
		box-shadow: 0 22px 60px rgba(35, 48, 52, 0.11);
	}

	.brand-panel {
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
		padding: 2rem;
		background:
			linear-gradient(160deg, rgba(255, 252, 245, 0.9), rgba(229, 239, 237, 0.74)),
			repeating-linear-gradient(135deg, rgba(89, 114, 95, 0.06) 0 8px, transparent 8px 16px);
	}

	.auth-panel {
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
		padding: 2rem;
	}

	.brand-chip,
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
	h2,
	h3,
	p {
		margin: 0;
	}

	h1,
	h2,
	h3 {
		font-family: Georgia, 'Times New Roman', serif;
		letter-spacing: 0;
	}

	h1 {
		max-width: 12ch;
		font-size: clamp(2.4rem, 5vw, 4.4rem);
		line-height: 0.98;
	}

	h2 {
		font-size: 2rem;
		line-height: 1.05;
	}

	p {
		color: #52666a;
		line-height: 1.55;
	}

	.route-card,
	.highlights div,
	.banner,
	.checkout-card {
		border: 1px solid rgba(29, 73, 79, 0.12);
		border-radius: 14px;
	}

	.route-card {
		display: grid;
		gap: 0.45rem;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.6);
	}

	.highlights {
		display: grid;
		gap: 0.75rem;
		margin-top: auto;
	}

	.highlights div {
		display: grid;
		gap: 0.3rem;
		padding: 1rem;
		background: rgba(255, 255, 255, 0.58);
	}

	.highlights span {
		color: #7b4d36;
		font-size: 0.78rem;
		font-weight: 800;
	}

	.panel-header,
	.action-row {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: 0.75rem;
	}

	.mode-switch {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		border-radius: 999px;
		background: #e7ece9;
		padding: 0.25rem;
	}

	button {
		border: 0;
		cursor: pointer;
		font: inherit;
	}

	.mode-switch button {
		border-radius: 999px;
		background: transparent;
		color: #52666a;
		font-weight: 800;
		padding: 0.7rem 1rem;
	}

	.mode-switch button.active {
		background: #1d6770;
		color: white;
	}

	.banner {
		display: grid;
		gap: 0.25rem;
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

	.auth-form,
	.grid {
		display: grid;
		gap: 1rem;
	}

	.two-columns {
		grid-template-columns: repeat(2, minmax(0, 1fr));
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

	.text-button {
		width: fit-content;
		justify-self: center;
		background: transparent;
		color: #1d6770;
		font-weight: 800;
		padding: 0.35rem 0.5rem;
	}

	.text-button:hover {
		text-decoration: underline;
	}

	.plan-options {
		display: grid;
		grid-template-columns: repeat(3, minmax(0, 1fr));
		gap: 0.75rem;
	}

	.plan-options button {
		display: grid;
		gap: 0.3rem;
		border: 1px solid rgba(29, 73, 79, 0.12);
		border-radius: 14px;
		background: rgba(255, 255, 255, 0.72);
		padding: 1rem;
		text-align: left;
	}

	.plan-options button.selected {
		border-color: #1d6770;
		background: rgba(29, 103, 112, 0.1);
	}

	.plan-options span,
	.plan-options em {
		color: #1d6770;
		font-size: 0.78rem;
		font-style: normal;
		font-weight: 800;
		text-transform: uppercase;
	}

	.plan-options small {
		color: #52666a;
		line-height: 1.4;
	}

	.checkout-card {
		display: grid;
		gap: 0.5rem;
		padding: 1rem;
		background: rgba(237, 245, 242, 0.8);
	}

	.checkout-card > strong {
		font-family: Georgia, 'Times New Roman', serif;
		font-size: 1.6rem;
	}

	.checkbox-field {
		grid-template-columns: auto 1fr;
		align-items: start;
		border: 1px solid rgba(29, 73, 79, 0.12);
		border-radius: 12px;
		background: rgba(255, 255, 255, 0.68);
		padding: 0.9rem;
	}

	.checkbox-field input {
		width: 1rem;
		height: 1rem;
		margin-top: 0.2rem;
		padding: 0;
		accent-color: #1d6770;
	}

	button:disabled,
	input:disabled {
		cursor: wait;
		opacity: 0.72;
	}

	@media (max-width: 900px) {
		.auth-shell {
			grid-template-columns: 1fr;
		}

		h1 {
			max-width: 14ch;
		}
	}

	@media (max-width: 680px) {
		.auth-shell {
			padding: 0.75rem;
		}

		.brand-panel,
		.auth-panel {
			padding: 1.25rem;
		}

		.panel-header,
		.action-row {
			flex-direction: column;
		}

		.two-columns,
		.plan-options {
			grid-template-columns: 1fr;
		}

		.mode-switch,
		.action-row button {
			width: 100%;
		}
	}
</style>

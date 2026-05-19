<script lang="ts">
	import { onMount } from 'svelte';
	import type { AccountPermissions, StoredSession } from '$lib/auth';
	import {
		PERMISSION_FIELDS,
		PERMISSION_SCOPES,
		TEAM_ROLES,
		canManageTeam,
		createTeamUser,
		deactivateTeamUser,
		fetchTeamUsers,
		getPlanSeatSummary,
		updateTeamUser,
		updateTeamUserPermissions,
		type CreateTeamUserInput,
		type PermissionScope,
		type TeamUser,
		type UpdateTeamUserInput,
		type UserRole,
		type UserStatus
	} from '../infrastructure/team-api';

	let { session, onConfirmDeactivate } = $props<{
		session: StoredSession | null;
		onConfirmDeactivate: (user: TeamUser) => Promise<boolean>;
	}>();

	let users = $state<TeamUser[]>([]);
	let isLoading = $state(false);
	let isCreating = $state(false);
	let isCreateFormOpen = $state(false);
	let feedback = $state<{ tone: 'success' | 'error' | 'info'; text: string } | null>(null);
	let createInput = $state<CreateTeamUserInput>(createEmptyUserInput());
	let savingUserId = $state<string | null>(null);
	let savingPermissionsUserId = $state<string | null>(null);

	const token = $derived(session?.payload.token ?? '');
	const currentUser = $derived(session?.payload.user ?? null);
	const subscription = $derived(session?.payload.subscription ?? null);
	const canManage = $derived(canManageTeam(currentUser));
	const seatSummary = $derived(getPlanSeatSummary(subscription, users));

	onMount(() => {
		if (token && canManage) {
			void loadUsers();
		}
	});

	function createEmptyUserInput(): CreateTeamUserInput {
		return {
			name: '',
			email: '',
			password: '',
			role: 'professional',
			status: 'active'
		};
	}

	async function loadUsers() {
		if (!token || !canManage) return;

		isLoading = true;
		feedback = null;
		try {
			const result = await fetchTeamUsers(token);
			users = sortUsers(result.users);
		} catch (error) {
			feedback = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel carregar a equipe.'
			};
		} finally {
			isLoading = false;
		}
	}

	async function submitCreateUser(event: SubmitEvent) {
		event.preventDefault();
		if (!token || isCreating) return;
		if (seatSummary.isAtLimit && createInput.status === 'active') {
			feedback = {
				tone: 'error',
				text: 'O plano atual nao tem assentos ativos disponiveis para uma nova conta.'
			};
			return;
		}

		isCreating = true;
		feedback = null;
		try {
			const created = await createTeamUser(token, {
				...createInput,
				name: createInput.name.trim(),
				email: createInput.email.trim()
			});
			if (created) users = sortUsers([created, ...users]);
			createInput = createEmptyUserInput();
			isCreateFormOpen = false;
			feedback = {
				tone: 'success',
				text: 'Conta criada para a clinica. O funcionario ja pode acessar com e-mail e senha.'
			};
		} catch (error) {
			feedback = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel criar a conta.'
			};
		} finally {
			isCreating = false;
		}
	}

	async function updateUserField(user: TeamUser, patch: Partial<UpdateTeamUserInput>) {
		if (!token || savingUserId) return;

		savingUserId = user.user_id;
		feedback = null;
		try {
			const updated = await updateTeamUser(token, user.user_id, {
				name: patch.name ?? user.name,
				email: patch.email ?? user.email,
				role: patch.role ?? user.role,
				status: patch.status ?? user.status
			});
			if (updated) {
				users = sortUsers(users.map((item) => (item.user_id === updated.user_id ? updated : item)));
				feedback = { tone: 'success', text: 'Conta atualizada.' };
			}
		} catch (error) {
			await loadUsers();
			feedback = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel atualizar a conta.'
			};
		} finally {
			savingUserId = null;
		}
	}

	async function savePermissions(user: TeamUser) {
		if (!token || savingPermissionsUserId) return;

		savingPermissionsUserId = user.user_id;
		feedback = null;
		try {
			const permissions = await updateTeamUserPermissions(token, user.user_id, user.permissions);
			users = users.map((item) =>
				item.user_id === user.user_id ? { ...item, permissions } : item
			);
			feedback = { tone: 'success', text: 'Permissoes atualizadas.' };
		} catch (error) {
			await loadUsers();
			feedback = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel atualizar permissoes.'
			};
		} finally {
			savingPermissionsUserId = null;
		}
	}

	async function requestDeactivateUser(user: TeamUser) {
		if (!token || user.status === 'inactive') return;
		if (!(await onConfirmDeactivate(user))) return;

		savingUserId = user.user_id;
		feedback = null;
		try {
			await deactivateTeamUser(token, user.user_id);
			users = users.map((item) =>
				item.user_id === user.user_id ? { ...item, status: 'inactive' } : item
			);
			feedback = { tone: 'success', text: 'Conta desativada.' };
		} catch (error) {
			feedback = {
				tone: 'error',
				text: error instanceof Error ? error.message : 'Nao foi possivel desativar a conta.'
			};
		} finally {
			savingUserId = null;
		}
	}

	function updatePermissionDraft(
		user: TeamUser,
		key: keyof AccountPermissions,
		scope: PermissionScope
	) {
		users = users.map((item) =>
			item.user_id === user.user_id
				? {
						...item,
						permissions: {
							...item.permissions,
							[key]: scope
						}
					}
				: item
		);
	}

	function canChangeUser(user: TeamUser) {
		if (!currentUser) return false;
		if (currentUser.user_id === user.user_id) return false;
		if (user.role === 'owner' && currentUser.role !== 'owner') return false;
		return true;
	}

	function getRoleLabel(role: string) {
		return TEAM_ROLES.find((item) => item.value === role)?.label ?? role;
	}

	function getStatusLabel(status: string) {
		return status === 'active' ? 'Ativa' : 'Inativa';
	}

	function sortUsers(input: TeamUser[]) {
		return [...input].sort((left, right) => {
			if (left.role === 'owner' && right.role !== 'owner') return -1;
			if (right.role === 'owner' && left.role !== 'owner') return 1;
			if (left.status !== right.status) return left.status === 'active' ? -1 : 1;
			return left.name.localeCompare(right.name);
		});
	}
</script>

<section class="settings-workspace">
	<div class="settings-head">
		<div>
			<span class="section-kicker">Governanca</span>
			<h1>Contas da clinica</h1>
			<p>
				Cada conta pertence a esta clinica, usa as permissoes abaixo e consome assentos do plano.
			</p>
		</div>

		<button type="button" class="secondary-button" onclick={loadUsers} disabled={isLoading || !canManage}>
			Atualizar
		</button>
	</div>

	{#if !canManage}
		<article class="governance-empty">
			<strong>Acesso restrito</strong>
			<p>Somente responsaveis da clinica e administradores com permissao de equipe podem gerenciar contas.</p>
		</article>
	{:else}
		<div class="governance-metrics">
			<div>
				<span>Contas ativas</span>
				<strong>{seatSummary.activeUsers}</strong>
			</div>
			<div>
				<span>Limite do plano</span>
				<strong>{seatSummary.maxProfessionals || 'Sem limite'}</strong>
			</div>
			<div>
				<span>Assentos disponiveis</span>
				<strong>{seatSummary.maxProfessionals ? seatSummary.availableSeats : 'Livre'}</strong>
			</div>
		</div>

		{#if feedback}
			<div class={`governance-feedback ${feedback.tone}`}>{feedback.text}</div>
		{/if}

		<div class="governance-actions">
			<button
				type="button"
				class="primary-button"
				onclick={() => (isCreateFormOpen = !isCreateFormOpen)}
			>
				{isCreateFormOpen ? 'Fechar cadastro' : '+ Criar conta'}
			</button>
			<span>
				A criacao da conta funciona como autorizacao previa da clinica para o funcionario acessar.
			</span>
		</div>

		{#if isCreateFormOpen}
			<form class="governance-form" onsubmit={submitCreateUser}>
				<label>
					<span>Nome</span>
					<input type="text" bind:value={createInput.name} required />
				</label>
				<label>
					<span>E-mail</span>
					<input type="email" bind:value={createInput.email} required />
				</label>
				<label>
					<span>Senha inicial</span>
					<input type="password" bind:value={createInput.password} required minlength="6" />
				</label>
				<label>
					<span>Tipo da conta</span>
					<select bind:value={createInput.role}>
						{#each TEAM_ROLES.filter((role) => role.value !== 'owner') as role}
							<option value={role.value}>{role.label}</option>
						{/each}
					</select>
				</label>
				<label>
					<span>Status</span>
					<select bind:value={createInput.status}>
						<option value="active">Ativa</option>
						<option value="inactive">Inativa</option>
					</select>
				</label>
				<div class="governance-form-actions">
					<button type="button" class="secondary-button" onclick={() => (isCreateFormOpen = false)}>
						Cancelar
					</button>
					<button type="submit" class="primary-button" disabled={isCreating}>
						{isCreating ? 'Criando...' : 'Criar conta'}
					</button>
				</div>
			</form>
		{/if}

		<div class="governance-user-list">
			{#if isLoading}
				<p class="governance-empty">Carregando contas da clinica...</p>
			{/if}

			{#each users as user}
				<article class="governance-user-card">
					<header>
						<div>
							<strong>{user.name}</strong>
							<span>{user.email}</span>
						</div>
						<div class="governance-badges">
							<span>{getRoleLabel(user.role)}</span>
							<span class={user.status}>{getStatusLabel(user.status)}</span>
						</div>
					</header>

					<div class="governance-user-controls">
						<label>
							<span>Tipo</span>
							<select
								value={user.role}
								disabled={!canChangeUser(user) || savingUserId === user.user_id}
								onchange={(event) =>
									updateUserField(user, {
										role: (event.currentTarget as HTMLSelectElement).value as UserRole
									})}
							>
								{#each TEAM_ROLES as role}
									<option value={role.value}>{role.label}</option>
								{/each}
							</select>
						</label>

						<label>
							<span>Status</span>
							<select
								value={user.status}
								disabled={!canChangeUser(user) || savingUserId === user.user_id}
								onchange={(event) =>
									updateUserField(user, {
										status: (event.currentTarget as HTMLSelectElement).value as UserStatus
									})}
							>
								<option value="active">Ativa</option>
								<option value="inactive">Inativa</option>
							</select>
						</label>

						<button
							type="button"
							class="danger-button"
							disabled={!canChangeUser(user) || user.status === 'inactive' || savingUserId === user.user_id}
							onclick={() => requestDeactivateUser(user)}
						>
							Desativar
						</button>
					</div>

					<div class="permission-matrix">
						{#each PERMISSION_FIELDS as field}
							<label>
								<span>{field.label}</span>
								<select
									value={user.permissions[field.key]}
									disabled={!canChangeUser(user) || savingPermissionsUserId === user.user_id}
									onchange={(event) =>
										updatePermissionDraft(
											user,
											field.key,
											(event.currentTarget as HTMLSelectElement).value as PermissionScope
										)}
								>
									{#each PERMISSION_SCOPES as scope}
										<option value={scope.value}>{scope.label}</option>
									{/each}
								</select>
							</label>
						{/each}
					</div>

					<footer>
						<button
							type="button"
							class="secondary-button"
							disabled={!canChangeUser(user) || savingPermissionsUserId === user.user_id}
							onclick={() => savePermissions(user)}
						>
							{savingPermissionsUserId === user.user_id ? 'Salvando...' : 'Salvar permissoes'}
						</button>
					</footer>
				</article>
			{:else}
				{#if !isLoading}
					<p class="governance-empty">Nenhuma conta encontrada para esta clinica.</p>
				{/if}
			{/each}
		</div>
	{/if}
</section>

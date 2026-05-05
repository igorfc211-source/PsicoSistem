<script lang="ts">
	import { getInitials } from '$lib/shared/formatters';

	let {
		searchTerm,
		userName,
		theme,
		onToggleSidebar,
		onSearchTermChange,
		onEditProfile,
		onOpenSettings,
		onToggleTheme,
		onLogout
	} = $props<{
		searchTerm: string;
		userName: string;
		theme: 'light' | 'dark';
		onToggleSidebar: () => void;
		onSearchTermChange: (value: string) => void;
		onEditProfile: () => void;
		onOpenSettings: () => void;
		onToggleTheme: () => void;
		onLogout: () => void;
	}>();

	let isProfileMenuOpen = $state(false);

	// Executa uma acao do menu e fecha o dropdown para manter o fluxo responsivo.
	function handleProfileAction(action: () => void) {
		action();
		isProfileMenuOpen = false;
	}
</script>

<header class="topbar">
	<!-- Botao mobile: abre o menu lateral quando a tela esta estreita. -->
	<button type="button" class="mobile-menu-button" aria-label="Abrir menu" onclick={onToggleSidebar}>
		<span></span>
		<span></span>
		<span></span>
	</button>

	<!-- Busca global do painel: filtra aprendentes e futuramente pode incluir agendamentos. -->
	<label class="search-box">
		<span>?</span>
		<input
			value={searchTerm}
			placeholder="Buscar aprendente ou agendamento..."
			oninput={(event) => onSearchTermChange((event.currentTarget as HTMLInputElement).value)}
		/>
	</label>

	<!-- Area do usuario: notificacoes, avatar e acao de sair da sessao. -->
	<div class="user-area">
		<button type="button" class="bell-button" aria-label="Notificacoes">o</button>
		<div class="user-avatar">{getInitials(userName)}</div>
		<div class="profile-menu">
			<button
				type="button"
				class="user-menu"
				aria-haspopup="menu"
				aria-expanded={isProfileMenuOpen}
				onclick={() => (isProfileMenuOpen = !isProfileMenuOpen)}
			>
				{userName}
				<span>v</span>
			</button>

			{#if isProfileMenuOpen}
				<!-- Dropdown do perfil: concentrador de acoes pessoais e preferencias. -->
				<div class="profile-dropdown" role="menu">
					<button type="button" role="menuitem" onclick={() => handleProfileAction(onEditProfile)}>
						Editar meu perfil
					</button>
					<button type="button" role="menuitem" onclick={() => handleProfileAction(onOpenSettings)}>
						Configuracoes
					</button>
					<button type="button" role="menuitem" onclick={() => handleProfileAction(onToggleTheme)}>
						{theme === 'dark' ? 'Modo claro' : 'Modo escuro'}
					</button>
					<button type="button" role="menuitem" class="danger-item" onclick={() => handleProfileAction(onLogout)}>
						Sair
					</button>
				</div>
			{/if}
		</div>
	</div>
</header>

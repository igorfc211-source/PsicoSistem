<script lang="ts">
	import { getInitials } from '$lib/shared/formatters';

	let {
		searchTerm,
		userName,
		onSearchTermChange,
		onLogout
	} = $props<{
		searchTerm: string;
		userName: string;
		onSearchTermChange: (value: string) => void;
		onLogout: () => void;
	}>();
</script>

<header class="topbar">
	<label class="search-box">
		<span>?</span>
		<input
			value={searchTerm}
			placeholder="Buscar aprendente ou agendamento..."
			oninput={(event) => onSearchTermChange((event.currentTarget as HTMLInputElement).value)}
		/>
	</label>

	<div class="user-area">
		<button type="button" class="bell-button" aria-label="Notificacoes">o</button>
		<div class="user-avatar">{getInitials(userName)}</div>
		<button type="button" class="user-menu" onclick={onLogout}>
			{userName}
			<span>v</span>
		</button>
	</div>
</header>

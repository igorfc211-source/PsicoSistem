<script lang="ts">
	import { onDestroy } from 'svelte';
	import { fly } from 'svelte/transition';
	import type { Learner, VisitKind } from '$lib/modules/learners';
	import LearnerAvatar from '$lib/modules/learners/components/avatar/LearnerAvatar.svelte';
	import {
		getAgendaEventKindLabel,
		type AgendaEvent,
		type AgendaEventKind,
		type NewAgendaEventInput,
		type NewSessionAppointmentInput,
		type ScheduleItem
	} from '$lib/modules/scheduling';

	let {
		selectedDate,
		selectedDateLabel,
		learners,
		selectedLearnerId,
		dayItems,
		onCreateSession,
		onCreateEvent,
		onOpenLearner,
		onRemoveSession,
		onRemoveEvent
	} = $props<{
		selectedDate: string;
		selectedDateLabel: string;
		learners: Learner[];
		selectedLearnerId: string | null;
		dayItems: ScheduleItem[];
		onCreateSession: (input: NewSessionAppointmentInput) => boolean;
		onCreateEvent: (input: NewAgendaEventInput) => boolean;
		onOpenLearner: (id: string) => void;
		onRemoveSession: (learnerId: string, visitId: string) => void;
		onRemoveEvent: (event: AgendaEvent) => void;
	}>();

	let actionMode = $state<'session' | 'event'>('session');
	let sessionLearnerId = $state('');
	let sessionKind = $state<VisitKind>('session');
	let sessionStartTime = $state('09:00');
	let sessionEndTime = $state('09:50');
	let sessionTitle = $state('Sessao individual');
	let sessionLocation = $state('Consultorio');
	let sessionNotes = $state('');
	let eventTitle = $state('Reuniao de equipe');
	let eventKind = $state<AgendaEventKind>('meeting');
	let eventStartTime = $state('11:00');
	let eventEndTime = $state('12:00');
	let eventDescription = $state('');
	let confirmation = $state<{ id: number; text: string } | null>(null);
	let confirmationTimer: ReturnType<typeof setTimeout> | null = null;

	$effect(() => {
		const preferredLearnerId = selectedLearnerId ?? learners[0]?.id ?? '';
		if (!sessionLearnerId && preferredLearnerId) {
			sessionLearnerId = preferredLearnerId;
		}
	});

	const visitKindOptions: Array<{ value: VisitKind; label: string; defaultTitle: string }> = [
		{ value: 'session', label: 'Sessao', defaultTitle: 'Sessao individual' },
		{ value: 'assessment', label: 'Avaliacao', defaultTitle: 'Avaliacao inicial' },
		{ value: 'return', label: 'Retorno', defaultTitle: 'Retorno clinico' }
	];

	// Traduz o tipo de visita para exibicao humana nos cards da agenda.
	function getVisitKindLabel(kind: VisitKind) {
		return visitKindOptions.find((option) => option.value === kind)?.label ?? 'Sessao';
	}

	// Ao trocar o tipo, sugere um titulo coerente sem travar a edicao manual.
	function handleVisitKindChange(event: Event) {
		const value = (event.currentTarget as HTMLSelectElement).value as VisitKind;
		sessionKind = value;
		sessionTitle =
			visitKindOptions.find((option) => option.value === value)?.defaultTitle ?? sessionTitle;
	}

	// Cria uma sessao vinculada ao aprendente selecionado e limpa apenas campos descritivos.
	function handleCreateSession(event: SubmitEvent) {
		event.preventDefault();

		const created = onCreateSession({
			learnerId: sessionLearnerId,
			date: selectedDate,
			startTime: sessionStartTime,
			endTime: sessionEndTime,
			kind: sessionKind,
			title: sessionTitle,
			location: sessionLocation,
			notes: sessionNotes
		});

		if (created) {
			sessionNotes = '';
			showConfirmation('Sessao confirmada na agenda.');
		}
	}

	// Cria eventos livres da agenda, como reunioes ou bloqueios, sem prender a um aprendente.
	function handleCreateEvent(event: SubmitEvent) {
		event.preventDefault();

		const created = onCreateEvent({
			title: eventTitle,
			date: selectedDate,
			startTime: eventStartTime,
			endTime: eventEndTime,
			kind: eventKind,
			description: eventDescription
		});

		if (created) {
			eventDescription = '';
			showConfirmation('Evento confirmado na agenda.');
		}
	}

	function showConfirmation(text: string) {
		if (confirmationTimer) clearTimeout(confirmationTimer);

		confirmation = { id: Date.now(), text };
		confirmationTimer = setTimeout(() => {
			confirmation = null;
			confirmationTimer = null;
		}, 2200);
	}

	onDestroy(() => {
		if (confirmationTimer) clearTimeout(confirmationTimer);
	});
</script>

<section class="day-scheduler">
	<!-- Resumo do dia selecionado: data legivel e total de compromissos. -->
	<div class="day-scheduler-head">
		<div>
			<span>Dia selecionado</span>
			<h2>{selectedDateLabel}</h2>
			<p>Data NBR 5892: {selectedDate}</p>
		</div>
		<strong>{dayItems.length} compromisso{dayItems.length === 1 ? '' : 's'}</strong>
	</div>

	{#if confirmation}
		{#key confirmation.id}
			<div class="schedule-confirmation" transition:fly={{ y: -12, duration: 180 }}>
				{confirmation.text}
			</div>
		{/key}
	{/if}

	<!-- Linha do tempo: agrupa sessoes de aprendentes e eventos livres do mesmo dia. -->
	<div class="schedule-list" aria-label="Compromissos do dia">
		{#each dayItems as item (item.id)}
			<article class={`schedule-card ${item.kind} ${item.tone}`} in:fly={{ x: 24, duration: 220 }}>
				<!-- Coluna de horario: deixa claro o intervalo de cada compromisso. -->
				<div class="schedule-time">
					<strong>{item.startTime}</strong>
					<span>{item.endTime}</span>
				</div>

				<!-- Informacoes principais do compromisso, variando entre sessao e evento. -->
				<div class="schedule-info">
					<div class="schedule-title-row">
						{#if item.learner}
							<LearnerAvatar
								name={item.learner.name}
								photoUrl={item.learner.photoUrl}
								size="small"
							/>
						{/if}
						<div>
							<strong>{item.title}</strong>
							<p>{item.subtitle}</p>
						</div>
					</div>
					{#if item.kind === 'event' && item.event}
						<small>{getAgendaEventKindLabel(item.event.kind)}</small>
					{:else if item.visit}
						<small>
							{getVisitKindLabel(item.visit.kind)}
							- {item.visit.status === 'scheduled' ? 'Pendente' : item.visit.status}
						</small>
					{/if}
				</div>

				<!-- Acoes contextuais: abrir ficha/remover sessao ou excluir evento livre. -->
				<div class="schedule-actions">
					{#if item.kind === 'session' && item.learner}
						<button type="button" class="secondary-button" onclick={() => onOpenLearner(item.learner!.id)}>
							Abrir ficha
						</button>
						{#if item.visit}
							<button
								type="button"
								class="danger-button"
								onclick={() => onRemoveSession(item.learner!.id, item.visit!.id)}
							>
								Remover
							</button>
						{/if}
					{:else if item.event}
						<button type="button" class="danger-button" onclick={() => onRemoveEvent(item.event!)}>
							Excluir
						</button>
					{/if}
				</div>
			</article>
		{:else}
			<!-- Estado vazio: aparece quando o dia ainda nao possui compromissos. -->
			<div class="schedule-empty">
				<strong>Nenhum compromisso neste dia.</strong>
				<p>Esse horario esta livre para uma sessao, reuniao, supervisao ou bloqueio.</p>
			</div>
		{/each}
	</div>

	<!-- Escolha do fluxo de cadastro: sessao vinculada a aprendente ou evento livre. -->
	<div class="schedule-choice">
		<span>O que gostaria de fazer nesse dia?</span>
		<div>
			<button
				type="button"
				class:active={actionMode === 'session'}
				onclick={() => (actionMode = 'session')}
			>
				Agendar com aprendente
			</button>
			<button
				type="button"
				class:active={actionMode === 'event'}
				onclick={() => (actionMode = 'event')}
			>
				Adicionar evento
			</button>
		</div>
	</div>

	{#if actionMode === 'session'}
		<!-- Formulario de sessao: cria um compromisso dentro do prontuario do aprendente. -->
		<form class="schedule-form card" onsubmit={handleCreateSession} transition:fly={{ y: 12, duration: 170 }}>
			<div class="form-grid">
				<label>
					<span>Aprendente</span>
					<select bind:value={sessionLearnerId} required>
						<option value="" disabled>Selecione o aprendente</option>
						{#each learners as learner}
							<option value={learner.id}>{learner.name}</option>
						{/each}
					</select>
				</label>

				<label>
					<span>Tipo de visita</span>
					<select value={sessionKind} onchange={handleVisitKindChange}>
						{#each visitKindOptions as option}
							<option value={option.value}>{option.label}</option>
						{/each}
					</select>
				</label>

				<label>
					<span>Titulo do atendimento</span>
					<input bind:value={sessionTitle} placeholder="Sessao individual" required />
				</label>

				<label>
					<span>Inicio</span>
					<input type="time" bind:value={sessionStartTime} required />
				</label>

				<label>
					<span>Fim</span>
					<input type="time" bind:value={sessionEndTime} required />
				</label>

				<label>
					<span>Local / modalidade</span>
					<input bind:value={sessionLocation} placeholder="Consultorio, online, escola..." />
				</label>
			</div>

			<label>
				<span>Observacoes para a sessao</span>
				<textarea bind:value={sessionNotes} placeholder="Objetivo, combinados ou preparo necessario."></textarea>
			</label>

			<button type="submit" class="primary-button">Confirmar sessao</button>
		</form>
	{:else}
		<!-- Formulario de evento livre: usado para reuniao, supervisao ou bloqueio de agenda. -->
		<form class="schedule-form card" onsubmit={handleCreateEvent} transition:fly={{ y: 12, duration: 170 }}>
			<div class="form-grid">
				<label>
					<span>Titulo do evento</span>
					<input bind:value={eventTitle} required />
				</label>

				<label>
					<span>Categoria</span>
					<select bind:value={eventKind}>
						<option value="meeting">Reuniao</option>
						<option value="supervision">Supervisao</option>
						<option value="block">Bloqueio de horario</option>
						<option value="other">Outro evento</option>
					</select>
				</label>

				<label>
					<span>Inicio</span>
					<input type="time" bind:value={eventStartTime} required />
				</label>

				<label>
					<span>Fim</span>
					<input type="time" bind:value={eventEndTime} required />
				</label>
			</div>

			<label>
				<span>Descricao</span>
				<textarea bind:value={eventDescription} placeholder="Pauta, participantes ou contexto."></textarea>
			</label>

			<button type="submit" class="primary-button">Salvar evento</button>
		</form>
	{/if}
</section>

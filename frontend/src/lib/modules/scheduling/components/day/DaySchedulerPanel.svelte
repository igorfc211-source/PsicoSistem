<script lang="ts">
	import type { Learner } from '$lib/modules/learners';
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

	$effect(() => {
		const preferredLearnerId = selectedLearnerId ?? learners[0]?.id ?? '';
		if (!sessionLearnerId && preferredLearnerId) {
			sessionLearnerId = preferredLearnerId;
		}
	});

	// Cria uma sessao vinculada ao aprendente selecionado e limpa apenas campos descritivos.
	function handleCreateSession(event: SubmitEvent) {
		event.preventDefault();

		const created = onCreateSession({
			learnerId: sessionLearnerId,
			date: selectedDate,
			startTime: sessionStartTime,
			endTime: sessionEndTime,
			title: sessionTitle,
			location: sessionLocation,
			notes: sessionNotes
		});

		if (created) {
			sessionNotes = '';
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
		}
	}
</script>

<section class="day-scheduler">
	<div class="day-scheduler-head">
		<div>
			<span>Dia selecionado</span>
			<h2>{selectedDateLabel}</h2>
		</div>
		<strong>{dayItems.length} compromisso{dayItems.length === 1 ? '' : 's'}</strong>
	</div>

	<div class="schedule-list" aria-label="Compromissos do dia">
		{#each dayItems as item}
			<article class={`schedule-card ${item.kind} ${item.tone}`}>
				<div class="schedule-time">
					<strong>{item.startTime}</strong>
					<span>{item.endTime}</span>
				</div>

				<div class="schedule-info">
					<strong>{item.title}</strong>
					<p>{item.subtitle}</p>
					{#if item.kind === 'event' && item.event}
						<small>{getAgendaEventKindLabel(item.event.kind)}</small>
					{:else if item.visit}
						<small>{item.visit.status === 'scheduled' ? 'Agendada' : item.visit.status}</small>
					{/if}
				</div>

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
			<div class="schedule-empty">
				<strong>Nenhum compromisso neste dia.</strong>
				<p>Esse horario esta livre para uma sessao, reuniao, supervisao ou bloqueio.</p>
			</div>
		{/each}
	</div>

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
		<form class="schedule-form card" onsubmit={handleCreateSession}>
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
					<span>Tipo de atendimento</span>
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
		<form class="schedule-form card" onsubmit={handleCreateEvent}>
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

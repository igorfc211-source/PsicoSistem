<script lang="ts">
	import type { Visit, VisitKind, VisitStatus } from '../../domain/types';
	import handleCalendarSection from '../../../../../routes/+page.svelte'
	let {
		visit,
		onUpdateVisit,
		onRemoveVisit
	} = $props<{
		visit: Visit | null;
		onUpdateVisit: (visitId: string, patch: Partial<Visit>) => void;
		onRemoveVisit: (visitId: string) => void;
	}>();
</script>

<div class="visit-editor card">
	<!-- Cabecalho do editor: mostra a data ou orienta a selecionar um dia no calendario. -->
	<h3 class="flex-1 ">{visit ? visit.date : 'Selecione uma data'}</h3>
	<button onclick={ } class="w-4px bg-blue-600" >Definir agendamento</button>

	{#if visit}
		<!-- Dados principais da sessao selecionada. -->
		<label>
			<span>Titulo</span>
			<input
				value={visit.title}
				oninput={(event) =>
					onUpdateVisit(visit.id, {
						title: (event.currentTarget as HTMLInputElement).value
					})}
			/>
		</label>

		<!-- Tipo da visita: diferencia sessao, avaliacao e retorno no calendario. -->
		<label>
			<span>Tipo de visita</span>
			<select
				value={visit.kind}
				onchange={(event) =>
					onUpdateVisit(visit.id, {
						kind: (event.currentTarget as HTMLSelectElement).value as VisitKind
					})}
			>
				<option value="session">Sessao</option>
				<option value="assessment">Avaliacao</option>
				<option value="return">Retorno</option>
			</select>
		</label>

		<!-- Horarios precisos: inicio e fim do atendimento. -->
		<div class="form-grid compact">
			<label>
				<span>Inicio</span>
				<input
					type="time"
					value={visit.startTime}
					oninput={(event) =>
						onUpdateVisit(visit.id, {
							startTime: (event.currentTarget as HTMLInputElement).value
						})}
				/>
			</label>

			<label>
				<span>Fim</span>
				<input
					type="time"
					value={visit.endTime}
					oninput={(event) =>
						onUpdateVisit(visit.id, {
							endTime: (event.currentTarget as HTMLInputElement).value
						})}
				/>
			</label>
		</div>

		<!-- Local ou modalidade da sessao: consultorio, online, escola ou outro. -->
		<label>
			<span>Local / modalidade</span>
			<input
				value={visit.location}
				oninput={(event) =>
					onUpdateVisit(visit.id, {
						location: (event.currentTarget as HTMLInputElement).value
					})}
			/>
		</label>

		<!-- Status operacional usado para acompanhar frequencia e andamento. -->
		<label>
			<span>Status da visita</span>
			<select
				value={visit.status}
				onchange={(event) =>
					onUpdateVisit(visit.id, {
						status: (event.currentTarget as HTMLSelectElement).value as VisitStatus
					})}
			>
				<option value="scheduled">Agendada</option>
				<option value="completed">Realizada</option>
				<option value="missed">Faltou</option>
			</select>
		</label>

		<!-- Observacoes internas da visita selecionada. -->
		<label>
			<span>Observacoes</span>
			<textarea
				value={visit.notes}
				oninput={(event) =>
					onUpdateVisit(visit.id, {
						notes: (event.currentTarget as HTMLTextAreaElement).value
					})}
			></textarea>
		</label>

		<!-- Acao destrutiva isolada no fim do formulario para reduzir exclusoes acidentais. -->
		<button type="button" class="danger-button" onclick={() => onRemoveVisit(visit.id)}>
			Remover visita
		</button>
	{/if}
</div>

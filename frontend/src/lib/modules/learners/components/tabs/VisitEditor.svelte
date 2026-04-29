<script lang="ts">
	import type { Visit, VisitStatus } from '../../domain/types';

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
	<h3>{visit ? visit.date : 'Selecione uma data'}</h3>

	{#if visit}
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

		<button type="button" class="danger-button" onclick={() => onRemoveVisit(visit.id)}>
			Remover visita
		</button>
	{/if}
</div>
